package core

func (m *Manager) SupportsExt[T any](id string) bool { _, ok := m.As[T](id); return ok }
