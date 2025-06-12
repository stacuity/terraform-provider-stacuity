// Copyright (c) HashiCorp, Inc.

package models

type EndpointGroupList struct {
	Success    bool                    `json:"success"`
	Messages   []string                `json:"messages"`
	TotalItems int32                   `json:"totalItems"`
	Limit      int32                   `json:"limit"`
	Offset     int32                   `json:"offset"`
	Data       []EndpointGroupReadItem `json:"data"`
}

type EndpointGroupSingle struct {
	Success    bool                  `json:"success"`
	Messages   []string              `json:"messages"`
	TotalItems int32                 `json:"totalItems"`
	Limit      int32                 `json:"limit"`
	Offset     int32                 `json:"offset"`
	Data       EndpointGroupReadItem `json:"data"`
}

type EndpointGroupResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
	Data     string   `json:"data"`
}

type EndpointGroupReadItem struct {
	Id                    string                `json:"id,omitempty"`
	Moniker               string                `json:"moniker"`
	Name                  string                `json:"name"`
	EndpointsAssigned     int32                 `json:"endpointsAssigned"`
	VSlice                VSlice                `json:"vslice"`
	EventMap              *EventMap             `json:"eventMap"`
	RoutingPolicy         *RoutingPolicy        `json:"routingPolicy"`
	SteeringProfile       *SteeringProfile      `json:"steeringProfile"`
	RegionalGatewayPolicy RegionalGatewayPolicy `json:"regionalGatewayPolicy"`
	IPAllocationType      IPAllocationType      `json:"ipAllocationType"`
	CustomerId            string                `json:"customerId"`
}

type SteeringProfile struct {
	Id      string `json:"id"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
}

type RegionalGatewayPolicy struct {
	Id      string `json:"id"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
}

type IPAllocationType struct {
	Key     int    `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type EndpointGroupModifyItem struct {
	Id                    string  `json:"id,omitempty"`
	Name                  string  `json:"name"`
	Moniker               string  `json:"moniker"`
	EventMap              *string `json:"eventMap"`
	RoutingPolicy         *string `json:"routingPolicy"`
	SteeringProfile       *string `json:"steeringProfile"`
	RegionalGatewayPolicy string  `json:"regionalGatewayPolicy"`
	IPAllocationType      string  `json:"ipAllocationType"`
	VSlice                string  `json:"vslice"`
}
