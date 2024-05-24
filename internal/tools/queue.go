package tools

import (
	"context"
	"log"
	"sync"
	"time"
)

const defaultTimeout = 5 * time.Second

type Jobs interface {
	Go(f func(ctx context.Context))
	GoWithContext(ctx context.Context, f func(ctx context.Context))
	GoWithTimeout(timeout time.Duration, f func(ctx context.Context))
	GoWithTimeoutWithContext(ctx context.Context, timeout time.Duration, f func(ctx context.Context))
	Wait()
	Stop()
	Sleep(timeout time.Duration)
}

// TODO: next iteration chnage default int into atomic
type Queue struct {
	// workerPool is size of pull workers
	workerPoolSize *int
	// workersCount local free workers counter
	workersCount *int
	// wg as *sync.WaitGroup
	wg *sync.WaitGroup
	// qChan chanel for runners
	qChan chan task
}

type task struct {
	ctx    context.Context
	cancel context.CancelFunc
	task   func(ctx context.Context)
}

func NewQueue(workerPoolSize int) *Queue {
	workersCount := 0
	return &Queue{
		workerPoolSize: &workerPoolSize,
		workersCount:   &workersCount,
		wg:             &sync.WaitGroup{},
		qChan:          make(chan task, workerPoolSize),
	}
}

func (q *Queue) Go(f func(ctx context.Context)) {
	if len(q.qChan) >= *q.workerPoolSize {
		q.Wait()
	}
	q.qChan <- task{ctx: context.Background(), cancel: q.cancelFunc(), task: f}
}

func (q *Queue) GoWithContext(ctx context.Context, f func(ctx context.Context)) {
	if len(q.qChan) >= *q.workerPoolSize {
		q.Wait()
	}
	q.qChan <- task{ctx: ctx, cancel: q.cancelFunc(), task: f}
}

func (q *Queue) GoWithTimeout(timeout time.Duration, f func(ctx context.Context)) {
	if len(q.qChan) >= *q.workerPoolSize {
		q.Wait()
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	q.qChan <- task{ctx: ctx, cancel: q.cancelFuncWithCancel(cancel), task: f}
}

func (q *Queue) GoWithTimeoutWithContext(ctx context.Context, timeout time.Duration, f func(ctx context.Context)) {
	if len(q.qChan) >= *q.workerPoolSize {
		q.Wait()
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	q.qChan <- task{ctx: ctx, cancel: q.cancelFuncWithCancel(cancel), task: f}
}

func (q *Queue) Wait() {
	log.Println("active jobs: ", *q.workersCount)
	q.Sleep(defaultTimeout)
	if *q.workerPoolSize <= *q.workersCount {
		q.Wait()
	}
}

func (q *Queue) Stop() {
	close(q.qChan)
	for *q.workersCount == 0 {
		q.Sleep(defaultTimeout)
	}
}

func (q *Queue) Sleep(timeout time.Duration) {
	time.Sleep(timeout)
}

func (q *Queue) cancelFunc() func() {
	return func() {
		*q.workersCount--
	}
}

func (q *Queue) cancelFuncWithCancel(cancelFunc func()) func() {
	return func() {
		cancelFunc()
		*q.workersCount--
	}
}

func (q *Queue) addJob(job *task) {
	if *q.workerPoolSize <= *q.workersCount {
		q.wg.Wait()
		//q.Wait()
	}
	q.wg.Add(1)
	*q.workersCount++
	go func() {
		defer q.wg.Done()
		defer job.cancel()
		job.task(job.ctx)
	}()
}

func (q *Queue) Run() {
	for {
		select {
		case job := <-q.qChan:
			if *q.workerPoolSize <= *q.workersCount {
				q.Wait()
			}
			go q.addJob(&job)
		default:
			q.Wait()
		}
	}
}
