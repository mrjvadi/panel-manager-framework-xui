package core

import "context"

type HealthChecker interface {
    Health(ctx context.Context) error
}

// Health: health-check یک پنل اگر پیاده‌سازی شده باشد
func (m *Manager) Health(ctx context.Context, id string) error {
    if s, ok := m.getSlot(id); ok && s.enabled {
        if h, ok := s.drv.(HealthChecker); ok {
            return h.Health(ctx)
        }
    }
    return nil
}
