package authenticationv2

import (
	"errors"
	"time"
	"tunn/utils/timer"
)

var ErrInBusy = errors.New("in busy")

type Bucket struct {
	timeout time.Duration
	sig     chan struct{}
}

//
// NewBucket
// @Description:
// @param timeout
// @return *Bucket
//
func NewBucket(timeout time.Duration) *Bucket {
	b := &Bucket{
		timeout: timeout,
		sig:     make(chan struct{}, 1),
	}
	go func() {
		b.sig <- struct{}{}
	}()
	return b
}

//
// Occupy
// @Description:
// @receiver b
// @return error
//
func (b *Bucket) Occupy() error {
	err := timer.TimeoutTask(func() error {
		<-b.sig
		return nil
	}, b.timeout)
	if err == timer.Timeout {
		return ErrInBusy
	}
	return err
}

//
// Leave
// @Description:
// @receiver b
//
func (b *Bucket) Leave() {
	go func() {
		b.sig <- struct{}{}
	}()
}
