package core

// Generic helper functions (free functions) because Go methods cannot have type parameters.
func As[T any](m *Manager, id string) (T, bool) {
    s, ok := m.getSlot(id)
    if !ok || !s.enabled {
        var zero T
        return zero, false
    }
    v, ok := any(s.drv).(T)
    return v, ok
}
func SupportsExt[T any](m *Manager, id string) bool { _, ok := As[T](m, id); return ok }
