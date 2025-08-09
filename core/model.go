package core

// مدل‌های متحد

type User struct {
    ID          string `json:"id"`
    Username    string `json:"username"`
    ExpireAt    *int64 `json:"expireAt,omitempty"`
    TrafficUsed *int64 `json:"trafficUsed,omitempty"`
    Raw         any    `json:"raw,omitempty"`
}

type Inbound struct {
    ID     string `json:"id"`
    Type   string `json:"type"`
    Remark string `json:"remark,omitempty"`
    Port   *int   `json:"port,omitempty"`
    Raw    any    `json:"raw,omitempty"`
}

// قابلیت‌ها و Featureها

type Feature string

const (
    // Marzban
    FeatureSubscriptions Feature = "subscriptions"
    FeatureUsersUsage    Feature = "users_usage"
    FeatureSystemInfo    Feature = "system_info"
    // X-UI
    FeatureXUIClients    Feature = "xui_clients"
)

type Capabilities struct {
    UsersCRUD     bool
    InboundsCRUD  bool
    TrafficStats  bool
    UserSuspend   bool
    UserReset     bool
    Extra         map[Feature]bool
}

func (c Capabilities) Has(f Feature) bool { return c.Extra != nil && c.Extra[f] }
