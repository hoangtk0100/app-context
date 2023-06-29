package asyncjob

import (
	"context"
	"sync"

	"github.com/hoangtk0100/app-context/util"
)

type group struct {
	jobs         []Job
	isConcurrent bool
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	return &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}
}

func (g *group) Run(ctx context.Context) error {
	if g.isConcurrent {
		numJobs := len(g.jobs)
		g.wg.Add(numJobs)

		errChan := make(chan error, numJobs)
		for _, jobIndex := range g.jobs {
			go func(job Job) {
				defer util.Recovery()

				errChan <- g.runJob(ctx, job)
				g.wg.Done()
			}(jobIndex)
		}

		g.wg.Wait()

		// Return the 1st error in errChan
		for index := 1; index <= numJobs; index++ {
			if err := <-errChan; err != nil {
				return err
			}
		}
	} else {
		for _, jobIndex := range g.jobs {
			if err := g.runJob(ctx, jobIndex); err != nil {
				return err
			}
		}
	}

	return nil
}

// Retry if needed
func (g *group) runJob(ctx context.Context, job Job) error {
	if err := job.Execute(ctx); err != nil {
		for {
			if job.State() == StateRetryFailed {
				return err
			}

			if job.Retry(ctx) == nil {
				return nil
			}
		}
	}

	return nil
}
