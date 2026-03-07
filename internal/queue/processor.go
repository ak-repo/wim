package queue

import (
	"context"

	"github.com/ak-repo/wim/internal/worker"
	"github.com/ak-repo/wim/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type Processor struct {
	logger     logger.Logger
	workerPool *worker.Pool
}

func NewProcessor(log logger.Logger, pool *worker.Pool) *Processor {
	return &Processor{
		logger:     log,
		workerPool: pool,
	}
}

func (p *Processor) Run(ctx context.Context, reader *kafka.Reader) error {
	p.logger.Info("starting message processor")

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("processor shutting down")
			return nil
		default:
			msg, err := reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				p.logger.Error("failed to fetch message", "error", err)
				continue
			}

			if err := p.processMessage(ctx, msg); err != nil {
				p.logger.Error("failed to process message", "error", err)
				if err := reader.CommitMessages(ctx, msg); err != nil {
					p.logger.Error("failed to commit message", "error", err)
				}
				continue
			}

			if err := reader.CommitMessages(ctx, msg); err != nil {
				p.logger.Error("failed to commit message", "error", err)
			}
		}
	}
}

func (p *Processor) processMessage(ctx context.Context, msg kafka.Message) error {
	p.logger.Info("processing message", "topic", msg.Topic, "partition", msg.Partition)

	job := &worker.Job{
		ID:         string(msg.Key),
		Topic:      msg.Topic,
		Payload:    msg.Value,
		RetryCount: 0,
	}

	return p.workerPool.Submit(ctx, job)
}
