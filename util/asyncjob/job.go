package asyncjob

import (
	"context"
	"log"
	"time"
)

type Job interface {
	Name() string
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	RetryIndex() int
	SetRetryDurations(times []time.Duration)
	LastError() error
}

const (
	defaultMaxTimeout = time.Second * 10
)

var (
	defaultRetryTimes = []time.Duration{time.Second, time.Second * 2, time.Second * 4}
)

type JobState int

const (
	StateInit JobState = iota
	StateRunning
	StateFailed
	StateTimeout
	StateCompleted
	StateRetryFailed
)

type JobHandler func(ctx context.Context) error

type jobConfig struct {
	Name           string
	MaxTimeout     time.Duration
	RetryDurations []time.Duration
}

type OptionHandler func(*jobConfig)

func WithName(name string) OptionHandler {
	return func(jc *jobConfig) {
		jc.Name = name
	}
}

func WithRetryDurations(times []time.Duration) OptionHandler {
	return func(jc *jobConfig) {
		jc.RetryDurations = times
	}
}

type job struct {
	config     jobConfig
	handler    JobHandler
	state      JobState
	retryIndex int
	stopChan   chan bool
	lastErr    error
}

func NewJob(handler JobHandler, options ...OptionHandler) *job {
	newJob := job{
		config: jobConfig{
			MaxTimeout:     defaultMaxTimeout,
			RetryDurations: defaultRetryTimes,
		},
		handler:    handler,
		retryIndex: -1,
		state:      StateInit,
		stopChan:   make(chan bool),
		lastErr:    nil,
	}

	for index := range options {
		options[index](&newJob.config)
	}

	return &newJob
}

func (j *job) Name() string {
	return j.config.Name
}

func (j *job) Execute(ctx context.Context) error {
	log.Println("Execute", j.config.Name)
	j.state = StateRunning

	if err := j.handler(ctx); err != nil {
		j.state = StateFailed
		j.lastErr = err
		return err
	}

	j.state = StateCompleted

	return nil
}

func (j *job) Retry(ctx context.Context) error {
	if j.retryIndex == len(j.config.RetryDurations)-1 {
		return nil
	}

	j.retryIndex++
	time.Sleep(j.config.RetryDurations[j.retryIndex])

	err := j.Execute(ctx)
	if err == nil {
		j.state = StateCompleted
		return nil
	}

	if j.retryIndex == len(j.config.RetryDurations)-1 {
		j.state = StateRetryFailed
		j.lastErr = err
		return err
	}

	j.state = StateFailed
	return err
}

func (j *job) State() JobState {
	return j.state
}

func (j *job) RetryIndex() int {
	return j.retryIndex
}

func (j *job) SetRetryDurations(times []time.Duration) {
	if len(times) == 0 {
		return
	}

	j.config.RetryDurations = times
}

func (j *job) LastError() error {
	return j.lastErr
}
