package core

import (
	"context"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/pubsub"
	"github.com/hoangtk0100/app-context/util/asyncjob"
)

type SubJob struct {
	Name string
	Hdl  func(ctx context.Context, msg *pubsub.Message) error
}

type GroupSubJob interface {
	Run(ctx context.Context) error
}

type subscribeEngine struct {
	name   string
	jobs   []topicJobs
	ps     PubSubComponent
	ac     appctx.AppContext
	logger appctx.Logger
}

func NewSubscribeEngine(name string, ps PubSubComponent, ac appctx.AppContext) *subscribeEngine {
	return &subscribeEngine{
		name: name,
		ps:   ps,
		ac:   ac,
	}
}

type topicJobs struct {
	topic        pubsub.Topic
	isConcurrent bool
	jobs         []SubJob
}

func (engine *subscribeEngine) AddTopicJobs(topic pubsub.Topic, isConcurrent bool, jobs ...SubJob) {
	topicJobs := &topicJobs{
		topic:        topic,
		isConcurrent: isConcurrent,
		jobs:         jobs,
	}

	engine.jobs = append(engine.jobs, *topicJobs)
}

func (engine *subscribeEngine) Start() error {
	engine.logger = engine.ac.Logger(engine.name)
	for _, jobIndex := range engine.jobs {
		engine.startSubTopic(jobIndex.topic, jobIndex.isConcurrent, jobIndex.jobs...)
	}

	return nil
}

func (engine *subscribeEngine) startSubTopic(topic pubsub.Topic, isConcurrent bool, jobs ...SubJob) error {
	c, _ := engine.ps.Subscribe(context.Background(), topic)
	for _, item := range jobs {
		engine.logger.Info("Setup subscriber :", item.Name)
	}

	getJobHandler := func(job *SubJob, msg *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) error {
			engine.logger.Infof("Run job [%s] - Value: %v", job.Name, msg.Data())
			return job.Hdl(ctx, msg)
		}
	}

	go func() {
		for {
			msg := <-c

			jobHdls := make([]asyncjob.Job, len(jobs))
			for index := range jobs {
				jobHdlIdnex := getJobHandler(&jobs[index], msg)
				jobHdls[index] = asyncjob.NewJob(jobHdlIdnex, asyncjob.WithName(jobs[index].Name))
			}

			group := asyncjob.NewGroup(isConcurrent, jobHdls...)
			if err := group.Run(context.Background()); err != nil {
				engine.logger.Error(err)
			}
		}
	}()

	return nil
}
