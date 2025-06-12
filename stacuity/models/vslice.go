// Copyright (c) HashiCorp, Inc.

package models

type VSliceList struct {
	Success    bool             `json:"success"`
	Messages   []string         `json:"messages"`
	TotalItems int32            `json:"totalItems"`
	Limit      int32            `json:"limit"`
	Offset     int32            `json:"offset"`
	Data       []VSliceReadItem `json:"data"`
}

type VSliceSingle struct {
	Success    bool           `json:"success"`
	Messages   []string       `json:"messages"`
	TotalItems int32          `json:"totalItems"`
	Limit      int32          `json:"limit"`
	Offset     int32          `json:"offset"`
	Data       VSliceReadItem `json:"data"`
}

type VSliceResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
	Data     string   `json:"data"`
}

type VSlice struct {
	Id      string `json:"id,omitempty"`
	Moniker string `json:"moniker,"`
	Name    string `json:"name,"`
}

type VSliceReadItem struct {
	Id                 string          `json:"id,omitempty"`
	Name               string          `json:"name"`
	Moniker            string          `json:"moniker"`
	Subnets            []string        `json:"subnets"`
	DNSServers         []string        `json:"dnsServers"`
	EventMap           EventMap        `json:"eventMap"`
	DNSMode            DNSMode         `json:"dnsMode"`
	IpAddressFamily    IpAddressFamily `json:"ipAddressFamily"`
	EndpointGroupCount int32           `json:"endpointGroupCount"`
	EndpointCount      int32           `json:"endpointCount"`
}

type VSliceModifyItem struct {
	Id               string   `json:"id,omitempty"`
	Name             string   `json:"name"`
	Moniker          string   `json:"moniker"`
	DNSServers       []string `json:"dnsServers"`
	DNSMode          string   `json:"dnsMode"`
	IpAddressFamily  string   `json:"ipAddressFamily"`
	SubnetAddress    string   `json:"subnetAddress"`
	IpAllocationType string   `json:"ipAllocationType,omitempty"`
	EventMap         string   `json:"eventMap"`
}

type IpAddressFamily struct {
	Key     int32  `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type DNSMode struct {
	Key     int32  `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type EventMap struct {
	Id      string `json:"id"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
}
