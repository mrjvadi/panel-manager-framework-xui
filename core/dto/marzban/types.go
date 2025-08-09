package marzban

type SystemInfo struct {
    Version   string                 `json:"version,omitempty"`
    Uptime    int64                  `json:"uptime,omitempty"`
    Meta      map[string]interface{} `json:"meta,omitempty"`
    Raw       map[string]interface{} `json:"raw,omitempty"`
}

type UserUsage struct {
    Username string `json:"username,omitempty"`
    Up       int64  `json:"up,omitempty"`
    Down     int64  `json:"down,omitempty"`
    Total    int64  `json:"total,omitempty"`
}

type Subscription struct {
    ID        string                 `json:"id,omitempty"`
    Username  string                 `json:"username,omitempty"`
    Link      string                 `json:"link,omitempty"`
    ExpiresAt *int64                 `json:"expires_at,omitempty"`
    Raw       map[string]interface{} `json:"raw,omitempty"`
}
