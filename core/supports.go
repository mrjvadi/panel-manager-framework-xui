package core

// SupportsExt: آیا پنل، افزونه‌ی T را پیاده کرده است؟
func (m *Manager) SupportsExt[T any](id string) bool {
    _, ok := m.As[T](id)
    return ok
}

// SupportsFeature: آیا پنل، Feature خاصی را دارد؟
func (m *Manager) SupportsFeature(id string, f Feature) bool {
    if s, ok := m.getSlot(id); ok && s.enabled {
        return s.drv.Capabilities().Has(f)
    }
    return false
}
