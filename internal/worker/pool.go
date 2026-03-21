package worker

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/ak-repo/wim/internal/config"
	batchSvc "github.com/ak-repo/wim/internal/service/batch"
	inventorySvc "github.com/ak-repo/wim/internal/service/inventory"
	reportSvc "github.com/ak-repo/wim/internal/service/report"
	"github.com/ak-repo/wim/pkg/logger"
)

type Dependencies struct {
	Inventory *inventorySvc.Service
	Batch     *batchSvc.Service
	Report    *reportSvc.Service
}

type Pool struct {
	size     int
	jobQueue chan *Job
	wg       sync.WaitGroup
	logger   logger.Logger
	deps     Dependencies
	maxRetry int
}

type Job struct {
	ID         string
	Topic      string
	Payload    []byte
	RetryCount int
}

type JobHandler func(ctx context.Context, job *Job) error

func NewPool(ctx context.Context, cfg config.WorkerConfig, log logger.Logger, deps Dependencies) *Pool {
	pool := &Pool{
		size:     cfg.PoolSize,
		jobQueue: make(chan *Job, cfg.QueueSize),
		logger:   log,
		deps:     deps,
		maxRetry: cfg.RetryCount,
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
				if job.RetryCount < p.maxRetry {
					job.RetryCount++
					delay := time.Duration(job.RetryCount) * time.Second
					p.logger.Warn("job retry scheduled", "job_id", job.ID, "retry_count", job.RetryCount, "delay", delay, "error", err)
					go func(j *Job, d time.Duration) {
						timer := time.NewTimer(d)
						defer timer.Stop()
						select {
						case <-ctx.Done():
							return
						case <-timer.C:
							_ = p.Submit(ctx, j)
						}
					}(job, delay)
					continue
				}

				p.logger.Error("job moved to dead-letter", "job_id", job.ID, "error", err)
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
	if p.deps.Inventory == nil {
		p.logger.Warn("inventory service not configured for stock recalculation")
		return nil
	}

	var payload struct {
		ProductID   string `json:"product_id"`
		WarehouseID string `json:"warehouse_id"`
	}
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return err
	}

	total, err := p.deps.Inventory.GetTotalQuantity(ctx, payload.ProductID, payload.WarehouseID)
	if err != nil {
		return err
	}

	p.logger.Info("stock recalculation completed", "product_id", payload.ProductID, "warehouse_id", payload.WarehouseID, "available_qty", total)
	return nil
}

func (p *Pool) handleExpiryAlert(ctx context.Context, job *Job) error {
	if p.deps.Batch == nil {
		p.logger.Warn("batch service not configured for expiry alerts")
		return nil
	}

	var payload struct {
		Days int `json:"days"`
	}
	if len(job.Payload) > 0 {
		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			return err
		}
	}

	batches, err := p.deps.Batch.GetExpiringSoon(ctx, payload.Days)
	if err != nil {
		return err
	}

	p.logger.Info("expiry alert check completed", "days", payload.Days, "expiring_count", len(batches))
	return nil
}

func (p *Pool) handleReportGeneration(ctx context.Context, job *Job) error {
	if p.deps.Report == nil {
		p.logger.Warn("report service not configured for report generation")
		return nil
	}

	var payload struct {
		Type string `json:"type"`
		Days int    `json:"days"`
	}
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return err
	}

	switch payload.Type {
	case "expiry":
		rows, err := p.deps.Report.ExpiryReport(ctx, payload.Days)
		if err != nil {
			return err
		}
		p.logger.Info("expiry report generated", "rows", len(rows), "days", payload.Days)
	default:
		p.logger.Info("report type not supported", "type", payload.Type)
	}

	return nil
}

func (p *Pool) Shutdown() {
	close(p.jobQueue)
	p.wg.Wait()
}
