package xui

type Inbound struct {
    ID     int                    `json:"id,omitempty"`
    Remark string                 `json:"remark,omitempty"`
    Protocol string               `json:"protocol,omitempty"`
    Port   int                    `json:"port,omitempty"`
    Settings map[string]any       `json:"settings,omitempty"`
    Sniffing map[string]any       `json:"sniffing,omitempty"`
    StreamSettings map[string]any `json:"streamSettings,omitempty"`
    Raw    map[string]any         `json:"raw,omitempty"`
}

type InboundCreate struct {
    Remark string                 `json:"remark,omitempty"`
    Protocol string               `json:"protocol,omitempty"`
    Port   int                    `json:"port,omitempty"`
    Settings map[string]any       `json:"settings,omitempty"`
    Sniffing map[string]any       `json:"sniffing,omitempty"`
    StreamSettings map[string]any `json:"streamSettings,omitempty"`
}

type InboundUpdate = InboundCreate

type TrafficRecord struct {
    Email string `json:"email,omitempty"`
    Up    int64  `json:"up,omitempty"`
    Down  int64  `json:"down,omitempty"`
    Total int64  `json:"total,omitempty"`
}

type CloneInboundOptions struct {
    Port   *int   `json:"port,omitempty"`   // اگر ست باشد، همین پورت استفاده می‌شود
    Remark *string `json:"remark,omitempty"` // اگر ست باشد، همین اسم استفاده می‌شود
}
