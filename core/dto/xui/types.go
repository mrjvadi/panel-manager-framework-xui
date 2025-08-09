package xui

type Inbound struct {
	ID             int            `json:"id,omitempty"`
	Remark         string         `json:"remark,omitempty"`
	Protocol       string         `json:"protocol,omitempty"`
	Port           int            `json:"port,omitempty"`
	Settings       map[string]any `json:"settings,omitempty"`
	Sniffing       map[string]any `json:"sniffing,omitempty"`
	StreamSettings map[string]any `json:"streamSettings,omitempty"`
	Raw            map[string]any `json:"raw,omitempty"`
}
type InboundCreate struct {
	Remark         string         `json:"remark,omitempty"`
	Protocol       string         `json:"protocol,omitempty"`
	Port           int            `json:"port,omitempty"`
	Settings       map[string]any `json:"settings,omitempty"`
	Sniffing       map[string]any `json:"sniffing,omitempty"`
	StreamSettings map[string]any `json:"streamSettings,omitempty"`
}
type InboundUpdate = InboundCreate
type TrafficRecord struct {
	Email           string `json:"email,omitempty"`
	Up, Down, Total int64
}

type Client struct {
	ID         int    `json:"id,omitempty"`
	Email      string `json:"email,omitempty"`
	Enable     bool   `json:"enable,omitempty"`
	ExpiryTime int64  `json:"expiryTime,omitempty"`
	TotalGB    int64  `json:"totalGB,omitempty"`
	LimitIP    int    `json:"limitIp,omitempty"`
	Flow       string `json:"flow,omitempty"`
	SubID      string `json:"subId,omitempty"`
	Comment    string `json:"comment,omitempty"`
}

type ClientCreate struct {
	Email      string
	Enable     bool
	ExpiryTime int64
	TotalGB    int64
	LimitIP    int
	Flow       string
	SubID      string
	Comment    string
}
