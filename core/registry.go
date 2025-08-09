package core

import "sync"

type FactoryFn func(PanelSpec, ...Option) (Driver, error)

var (
    regMu sync.RWMutex
    reg = map[string]FactoryFn{}
)

func Register(name string, fn FactoryFn) { regMu.Lock(); defer regMu.Unlock(); reg[name] = fn }
func Factory(name string) (FactoryFn, bool) { regMu.RLock(); defer regMu.RUnlock(); fn, ok := reg[name]; return fn, ok }
