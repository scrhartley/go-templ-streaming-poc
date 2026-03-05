package templates

import (
	"iter"
	"slices"
	"sync"
)

func Concurrent[R any](f func() R) func() R {
	c := make(chan R, 1)
	go func() {
		defer close(c)
		c <- f()
	}()
	return sync.OnceValue(func() R {
		return <-c
	})
}

func TryConcurrent[R any](f func() (R, error)) func() (R, error) {
	c := make(chan R, 1)
	e := make(chan error, 1)
	go func() {
		defer close(c)
		defer close(e)
		res, err := f()
		c <- res
		e <- err
	}()
	return sync.OnceValues(func() (R, error) {
		return <-c, <-e
	})
}

func ConcurrentSeq[R any](fs ...func() R) iter.Seq[func() R] {
	c := make(chan R, len(fs))
	for _, f := range fs {
		go func() {
			c <- f()
		}()
	}

	cachingResults := make([]func() R, len(fs))
	for i := range fs {
		cachingResults[i] = sync.OnceValue(func() R {
			return <-c
		})
	}
	return slices.Values(cachingResults)
}

func TryConcurrentSeq[R any](fs ...func() (R, error)) iter.Seq[func() (R, error)] {
	type result struct {
		value R
		err   error
	}

	c := make(chan result, len(fs))
	for _, f := range fs {
		go func() {
			val, err := f()
			c <- result{val, err}
		}()
	}

	cachingResults := make([]func() (R, error), len(fs))
	for i := range fs {
		cachingResults[i] = sync.OnceValues(func() (R, error) {
			res := <-c
			return res.value, res.err
		})
	}
	return slices.Values(cachingResults)
}
