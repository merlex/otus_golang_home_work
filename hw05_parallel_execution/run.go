package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(w *sync.WaitGroup, m *sync.Mutex, in chan Task, countErrors *int, limitErrors int) {
	defer w.Done()

	for task := range in {
		if err := task(); err != nil {
			isLimitReached := func() bool {
				m.Lock()
				defer m.Unlock()

				*countErrors++
				return *countErrors >= limitErrors
			}()

			if isLimitReached && limitErrors > 0 {
				return
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	mut := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	wg.Add(n)

	in := make(chan Task, len(tasks))
	for _, v := range tasks {
		in <- v
	}
	close(in)

	countErrors := 0
	limitErrors := m
	if limitErrors < 0 {
		limitErrors = 0
	}

	for i := 0; i < n; i++ {
		go worker(wg, mut, in, &countErrors, limitErrors)
	}
	wg.Wait()

	if countErrors >= limitErrors && limitErrors > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
