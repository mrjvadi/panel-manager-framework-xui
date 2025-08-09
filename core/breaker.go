package core

import "time"
import "sync"

type Breaker struct {
    mu sync.Mutex
    state int
    fails int
    success int
    threshold int
    cooldown time.Duration
    openedAt time.Time
}

func NewBreaker(th int, cd time.Duration) *Breaker {
    if th <= 0 { th = 5 }
    if cd <= 0 { cd = 5 * time.Second }
    return &Breaker{ threshold: th, cooldown: cd }
}

func (b *Breaker) Allow() bool {
    b.mu.Lock(); defer b.mu.Unlock()
    switch b.state {
    case 0: return true
    case 1:
        if time.Since(b.openedAt) >= b.cooldown { b.state = 2; b.fails=0; b.success=0; return true }
        return false
    case 2: return true
    default: return true
    }
}
func (b *Breaker) OnSuccess() { b.mu.Lock(); if b.state==2 { b.success++; if b.success>=2 { b.state=0; b.success=0; b.fails=0 } }; b.mu.Unlock() }
func (b *Breaker) OnFailure() { b.mu.Lock(); if b.state==0 { b.fails++; if b.fails>=b.threshold { b.state=1; b.openedAt=time.Now() } } else if b.state==2 { b.state=1; b.openedAt=time.Now() }; b.mu.Unlock() }
