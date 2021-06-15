package concurrency

import (
	"errors"
	"strings"
	"sync"
)

type MultiErrorGroup struct {
	wg     *sync.WaitGroup
	errors []error
	m      *sync.Mutex
}

func NewMultiErrorGroup() *MultiErrorGroup {
	return &MultiErrorGroup{
		errors: make([]error, 0),
		wg:     &sync.WaitGroup{},
		m:      &sync.Mutex{},
	}
}

func (a *MultiErrorGroup) Go(f func() error) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := f(); err != nil {
			a.m.Lock()
			a.errors = append(a.errors, err)
			a.m.Unlock()
		}
	}()
}

func (a *MultiErrorGroup) Wait() error {
	a.wg.Wait()
	return mergeErrors(a.errors)
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	errStrs := make([]string, len(errs))
	for i, e := range errs {
		errStrs[i] = e.Error()
	}

	return errors.New(strings.Join(errStrs, "\n"))
}
