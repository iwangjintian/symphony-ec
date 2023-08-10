/*
   MIT License

   Copyright (c) Microsoft Corporation.

   Permission is hereby granted, free of charge, to any person obtaining a copy
   of this software and associated documentation files (the "Software"), to deal
   in the Software without restriction, including without limitation the rights
   to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
   copies of the Software, and to permit persons to whom the Software is
   furnished to do so, subject to the following conditions:

   The above copyright notice and this permission notice shall be included in all
   copies or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
   IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
   FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
   AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
   LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
   OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
   SOFTWARE

*/

package vendors

import (
	"github.com/azure/symphony/api/pkg/apis/v1alpha1/managers/configs"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/managers"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/providers"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/providers/pubsub"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/vendors"
	"github.com/azure/symphony/coa/pkg/logger"
	"github.com/valyala/fasthttp"
)

var zLog = logger.NewLogger("coa.runtime")

type ConfigsVendor struct {
	vendors.Vendor
	ConfigsManager *configs.ConfigsManager
}

func (e *ConfigsVendor) Init(config vendors.VendorConfig, factories []managers.IManagerFactroy, providers map[string]map[string]providers.IProvider, pubsubProvider pubsub.IPubSubProvider) error {
	err := e.Vendor.Init(config, factories, providers, pubsubProvider)
	if err != nil {
		return err
	}
	for _, m := range e.Managers {
		if c, ok := m.(*configs.ConfigsManager); ok {
			e.ConfigsManager = c
		}
	}
	if e.ConfigsManager == nil {
		return v1alpha2.NewCOAError(nil, "configs manager is not supplied", v1alpha2.MissingConfig)
	}
	return nil
}

func (e *ConfigsVendor) GetInfo() vendors.VendorInfo {
	return vendors.VendorInfo{
		Version:  e.Vendor.Version,
		Name:     "Configs",
		Producer: "Microsoft",
	}
}

func (e *ConfigsVendor) GetEndpoints() []v1alpha2.Endpoint {
	route := "configs"
	if e.Route != "" {
		route = e.Route
	}
	return []v1alpha2.Endpoint{
		{
			Methods:    []string{fasthttp.MethodGet, fasthttp.MethodPost, fasthttp.MethodDelete},
			Route:      route,
			Version:    e.Version,
			Handler:    e.onConfigs,
			Parameters: []string{"name?"},
		},
	}
}
func (e *ConfigsVendor) onConfigs(request v1alpha2.COARequest) v1alpha2.COAResponse {
	resp := v1alpha2.COAResponse{
		State:       v1alpha2.MethodNotAllowed,
		Body:        []byte("{\"result\":\"405 - method not allowed\"}"),
		ContentType: "application/json",
	}
	return resp
}
