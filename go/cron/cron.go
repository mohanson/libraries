// Package cron provides a simple way to schedule tasks at fixed intervals.
package cron

import (
	"log"
	"time"
)

// Cron schedules tasks at regular intervals. This function receives two parameters: the elapsed time (duration) between
// consecutive task executions (period), and the initial delay before the first task execution. This delay is applied
// after aligning to the next period boundary.
func Cron(e time.Duration, d time.Duration) <-chan struct{} {
	if d >= e {
		log.Panicln("cron: the delay should be less than the time elapsed between events")
	}
	r := make(chan struct{})
	go func() {
		for {
			n := time.Now()
			// Wait for the next scheduled event by adding the elapsed time and then waiting for the delay.
			time.Sleep(n.Add(e).Truncate(e).Sub(n) + d)
			r <- struct{}{}
		}
	}()
	return r
}
