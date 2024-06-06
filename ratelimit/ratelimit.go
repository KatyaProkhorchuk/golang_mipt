//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	maxCount int
	interval time.Duration
	stop chan struct{} // завершение работы
	timeout chan[]*time.Timer
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	timer := make([]*time.Timer, maxCount)
	timeout := make(chan[]*time.Timer, 1)
	for i := range timer {
		timer[i] = time.NewTimer(0)
	}
	timeout <- timer
	return &Limiter{maxCount: maxCount, 
		interval: interval, 
		stop: make(chan struct{}, 1), 
		timeout: timeout}
}
func (l *Limiter) Helper(timeout []*time.Timer, i int) bool {
	defer func() {l.timeout <- timeout}() 
	//если таймер сработал
	select {
	case <-timeout[i].C:
		timeout[i] =time.NewTimer(l.interval)
		return true
	default:
		return false
	}
}
func (l* Limiter) Stopped() bool {
	select {
	case <-l.stop:
		return true
	default:
		return false
	}
}
func (l *Limiter) Acquire(ctx context.Context) error {
	// проверим не закрыт ли канал
	if l.Stopped() {
		return ErrStopped
	}
	select{
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	i := 0
	for i < l.maxCount {
		select {
		case <-l.stop:
			return ErrStopped
		case <-ctx.Done():
			return ctx.Err()
		case timeout := <- l.timeout:
			if l.Helper(timeout, i) {
				return nil
			}
		default:
		}
	
		i = (i + 1) % l.maxCount
	}
	return nil
}

func (l *Limiter) Stop() {
	close(l.stop)
}