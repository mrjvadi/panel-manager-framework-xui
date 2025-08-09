package xui

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
