package core

type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Up       int64  `json:"up,omitempty"`
	Down     int64  `json:"down,omitempty"`
}

type Inbound struct {
	ID       string `json:"id,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Port     int    `json:"port,omitempty"`
}

// NOTE: TLS و BasicAuth در core/auth.go تعریف شده‌اند.
// اینجا فقط ازشان استفاده می‌کنیم.
type PanelSpec struct {
	ID        string
	BaseURL   string
	Auth      BasicAuth
	TLS       TLS
	Version   string
	Endpoints map[string]string
}

type Capabilities struct{ bits uint64 }

func (c Capabilities) Has(x Capabilities) bool { return (c.bits & x.bits) != 0 }
