package utils

import (
	"time"
)

type Debouncer struct {
	cancelUntil           time.Time
	intervalInMiliseconds int
	taskToRun             func()
}

func NewDebouncer(interval int, action func()) *Debouncer {
	return &Debouncer{
		cancelUntil:           time.Now(),
		intervalInMiliseconds: interval,
		taskToRun:             action,
	}
}

func (d *Debouncer) Invoke() {
	d.delayedInvoke()
}

func (d *Debouncer) delayedInvoke() {
	now := time.Now()
	if now.Before(d.cancelUntil) {
		return
	}

	d.cancelUntil = now.Add(time.Duration(d.intervalInMiliseconds) * time.Millisecond)

	timer := time.NewTimer(time.Duration(d.intervalInMiliseconds) * time.Millisecond)
	go func() {
		<-timer.C
		d.taskToRun()
	}()
}
