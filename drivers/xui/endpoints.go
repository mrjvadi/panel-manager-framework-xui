package xui

var defaultEndpoints = map[string]string{
	"login":         "/login",
	"listInbounds":  "/panel/inbound/list",
	"getInbound":    "/panel/api/inbounds/get/%d",
	"addInbound":    "/panel/api/inbounds/add",
	"updateInbound": "/panel/inbound/update/%d",
	"deleteInbound": "/panel/inbound/del/%d",
	"addClient":     "/panel/inbound/addClient",
	"delClient":     "/panel/inbounds/delClient/%d",
}
