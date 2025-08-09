package xui

// Default endpoints used by X-UI branches if not overridden in PanelSpec.Endpoints.
var defaultEndpoints = map[string]string{
    "login":        "/login",
    "listInbounds": "/panel/api/inbounds/list",
    "getInbound":   "/panel/api/inbounds/get/%d",
    "addInbound":   "/panel/api/inbounds/add",
    "updateInbound":"/panel/api/inbounds/update/%d",
    "deleteInbound":"/panel/api/inbounds/del/%d",
    "addClient":    "/panel/api/inbounds/addClient",
    "delClient":    "/panel/api/inbounds/delClient/%d",
}
