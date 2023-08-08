package time

import "time"

type Time struct {
	start time.Time
	end   time.Time
}

func NewTime() *Time {
	return &Time{
		start: time.Now(),
		end:   time.Now(),
	}
}

func (t *Time) Start() {
	t.start = time.Now()
}

func (t *Time) Finish() *Time {
	t.end = time.Now()

	return t
}

func (t *Time) Diff() int64 {
	return t.end.Sub(t.start).Nanoseconds()
}

func (t *Time) Nano() int64 {
	t.Finish()
	return t.end.Sub(t.start).Nanoseconds()
}

func (t *Time) Mili() int64 {
	t.Finish()
	return t.end.Sub(t.start).Milliseconds()
}

func (t *Time) Micro() int64 {
	t.Finish()
	return t.end.Sub(t.start).Microseconds()
}
