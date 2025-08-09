package core

import "strings"

type DriverKind int

const (
    DriverMarzban DriverKind = iota
    DriverXUIGeneric
    DriverXUISanaei
    DriverXUIAlireza
)

func (k DriverKind) String() string {
    switch k {
    case DriverMarzban:
        return "marzban"
    case DriverXUIGeneric:
        return "xui.generic"
    case DriverXUISanaei:
        return "xui.sanaei"
    case DriverXUIAlireza:
        return "xui.alireza"
    default:
        return ""
    }
}

// FamilyPrefix: پیشوند خانواده برای یک نوع درایور (برای گروه‌بندی)
func (k DriverKind) FamilyPrefix() string {
    name := k.String()
    if i := strings.IndexByte(name, '.'); i > 0 { return name[:i+1] }
    return name
}
