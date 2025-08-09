package xui

import (
	"github.com/mrjvadi/panel-manager-framework-xui/core"
)

const AlirezaName = "xui.alireza"

func init() { core.Register(AlirezaName, NewAlireza) }

type alireza struct{ generic }

func NewAlireza(sp core.PanelSpec, opts ...core.Option) (core.Driver, error) {
	if sp.Endpoints == nil {
		sp.Endpoints = map[string]string{}
	}
	if sp.Endpoints["login"] == "" {
		sp.Endpoints["login"] = "/api/auth/login"
	}
	if sp.Endpoints["listUsers"] == "" {
		sp.Endpoints["listUsers"] = "/api/users"
	}
	if sp.Endpoints["listInbounds"] == "" {
		sp.Endpoints["listInbounds"] = "/xui/inbound/list"
	}
	if sp.Version == "" {
		sp.Version = "alireza"
	}
	g := newGeneric(sp, opts...)
	return &alireza{*g}, nil
}

func (d *alireza) Name() string { return AlirezaName }
func (d *alireza) Version() string {
	if d.sp.Version != "" {
		return d.sp.Version
	}
	return "alireza"
}
