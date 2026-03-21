package worker

import "errors"

var ErrJobQueueFull = errors.New("job queue is full")
