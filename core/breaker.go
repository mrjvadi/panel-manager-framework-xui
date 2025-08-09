package core

import (
    "sync"
    "time"
)

// Breaker: مدارشکن سبک
type Breaker struct {
    mu        sync.Mutex
    state     int // 0 closed, 1 open, 2 half-open
    fails     int
    success   int
    threshold int
    cooldown  time.Duration
    openedAt  time.Time
}

func NewBreaker(threshold int, cooldown time.Duration) *Breaker {
    if threshold <= 0 { threshold = 5 }
    if cooldown <= 0 { cooldown = 5 * time.Second }
    return &Breaker{ threshold: threshold, cooldown: cooldown }
}

func (b *Breaker) Allow() bool {
    b.mu.Lock(); defer b.mu.Unlock()
    switch b.state {
    case 0: // closed
        return True
    case 1: // open
        if time.Since(b.openedAt) >= b.cooldown {
            b.state = 2 // half-open
            b.success = 0; b.fails = 0
            return True
        }
        return False
    case 2: // half-open
        return True
    default:
        return True
    }
}

func (b *Breaker) OnSuccess() {
    b.mu.Lock(); defer b.mu.Unlock()
    switch b.state {
    case 0:
        b.success++
    case 2: // half-open
        b.success++
        if b.success >= 2 { // چند موفقیت متوالی برای بستن کامل
            b.state = 0
            b.success = 0; b.fails = 0
        }
    }
}

func (b *Breaker) OnFailure() {
    b.mu.Lock(); defer b.mu.Unlock()
    b.fails++
    switch b.state {
    case 0:
        if b.fails >= b.threshold {
            b.state = 1
            b.openedAt = time.Now()
        }
    case 2:
        // در حالت half-open با یک شکست باز می‌ماند
        b.state = 1
        b.openedAt = time.Now()
    }
}

const True = true
