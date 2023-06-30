package main

import (
	"context"
	"time"

	"github.com/hoangtk0100/app-context/util/asyncjob"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func timer(name string) func() {
	start := time.Now()
	return func() {
		log.Info().Msgf("%s took %v", name, time.Since(start))
	}
}

func main() {
	defer timer("Job Group")()

	job1 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 3)
		log.Info().Msg("Done job 1")
		return nil
	}, asyncjob.WithName("Job 1"))

	job2 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 3)
		log.Info().Msg("Running job 2")
		return errors.New("err at job 2")
	}, asyncjob.WithName("Job 2"))

	job3 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 3)
		log.Info().Msg("Running job 3")
		return errors.New("err at job 3")
	}, asyncjob.WithName("Job 3"))

	job4 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 3)
		log.Info().Msg("Done job 4")
		return nil
	}, asyncjob.WithName("Job 4"))

	group := asyncjob.NewGroup(true, job1, job2, job3, job4)
	if err := group.Run(context.Background()); err != nil {
		log.Error().Err(err)
	}
}
