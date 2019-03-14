package watchdog;

import (
    "time"
)

type Watchdog struct {
    interval time.Duration
    timer *time.Timer
}

func New(interval time.Duration) *Watchdog {
    w := Watchdog{
        interval: interval,
        timer: time.NewTimer(interval),
    }
    return &w
}

func (w *Watchdog) Stop() {
    w.timer.Stop()
}

func (w *Watchdog) Reset() {
    w.timer.Stop()
    w.timer.Reset(w.interval)
}

func (w *Watchdog) TimeOverChannel() <-chan time.Time {
    return w.timer.C
}
