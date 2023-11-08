package common

import (
	"context"
	"sync"
	"time"
)

var timerPool sync.Pool

func getTimer(d time.Duration) *time.Timer {
	if t, _ := timerPool.Get().(*time.Timer); t != nil {
		t.Reset(d)
		return t
	}
	return time.NewTimer(d)
}

func putTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	timerPool.Put(t)
}

// Sleep awaits for provided interval.
// Can be interrupted by context cancelation.
func Sleep(ctx context.Context, interval time.Duration) (err error) {
	timer := getTimer(interval)
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-timer.C:
	}
	putTimer(timer)
	return err
}
