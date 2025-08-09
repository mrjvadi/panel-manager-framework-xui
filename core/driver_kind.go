package core

type DriverKind string
const (
    DriverMarzban    DriverKind = "marzban"
    DriverXUISanaei  DriverKind = "xui.sanaei"
    DriverXUIAlireza DriverKind = "xui.alireza"
)
func (k DriverKind) String() string { return string(k) }
