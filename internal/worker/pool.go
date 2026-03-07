package worker

import (
	"context"
	"sync"

	"github.com/ak-repo/wim/internal/config"
	"github.com/ak-repo/wim/pkg/logger"
)

type Pool struct {
	size     int
	jobQueue chan *Job
	wg       sync.WaitGroup
	logger   logger.Logger
}

type Job struct {
	ID         string
	Topic      string
	Payload    []byte
	RetryCount int
}

type JobHandler func(ctx context.Context, job *Job) error

func NewPool(ctx context.Context, cfg config.WorkerConfig, log logger.Logger) *Pool {
	pool := &Pool{
		size:     cfg.PoolSize,
		jobQueue: make(chan *Job, cfg.QueueSize),
		logger:   log,
	}

	for i := 0; i < cfg.PoolSize; i++ {
		pool.wg.Add(1)
		go pool.worker(ctx, i)
	}

	return pool
}

func (p *Pool) Submit(ctx context.Context, job *Job) error {
	select {
	case p.jobQueue <- job:
		return nil
	default:
		return ErrJobQueueFull
	}
}

func (p *Pool) worker(ctx context.Context, id int) {
	defer p.wg.Done()
	p.logger.Info("worker started", "worker_id", id)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("worker stopping", "worker_id", id)
			return
		case job := <-p.jobQueue:
			if err := p.processJob(ctx, job); err != nil {
				p.logger.Error("job failed", "job_id", job.ID, "error", err)
			}
		}
	}
}

func (p *Pool) processJob(ctx context.Context, job *Job) error {
	p.logger.Info("processing job", "job_id", job.ID, "topic", job.Topic)

	switch job.Topic {
	case "inventory.stock_recalculation":
		return p.handleStockRecalculation(ctx, job)
	case "inventory.expiry_alert":
		return p.handleExpiryAlert(ctx, job)
	case "reports.generation":
		return p.handleReportGeneration(ctx, job)
	default:
		p.logger.Warn("unknown job topic", "topic", job.Topic)
	}

	return nil
}

func (p *Pool) handleStockRecalculation(ctx context.Context, job *Job) error {
	p.logger.Info("running stock recalculation")
	return nil
}

func (p *Pool) handleExpiryAlert(ctx context.Context, job *Job) error {
	p.logger.Info("running expiry alert check")
	return nil
}

func (p *Pool) handleReportGeneration(ctx context.Context, job *Job) error {
	p.logger.Info("generating report")
	return nil
}

func (p *Pool) Shutdown() {
	close(p.jobQueue)
	p.wg.Wait()
}
