package af

import (
    "time"
)

type Timer struct {
    startTicks int64
    pausedTicks int64
    paused bool
    started bool
}


func (t *Timer) Start() {
    t.started = true
    t.paused = false
    t.startTicks = time.Now().UnixNano()
}

func (t *Timer) Stop() {
    t.started = false
    t.paused = false
}

func (t *Timer) GetTicks() int64 {
    if t.started == true {
        if t.paused == true {
            return t.pausedTicks
        } else {
            return time.Now().UnixNano() - t.startTicks
        }
    }

    return 0
}

func (t* Timer) Pause() {
    if t.started == true && t.paused == false {
        t.paused = true
        t.pausedTicks = time.Now().UnixNano() - t.startTicks
    }
}

func (t* Timer) Unpause() {
    if t.paused == true {
        t.paused = false
        t.startTicks = time.Now().UnixNano() - t.pausedTicks
        t.pausedTicks = 0
    }
}

func (t* Timer) IsStarted() bool {
    return t.started
}

func (t* Timer) IsPaused() bool {
    return t.paused
}