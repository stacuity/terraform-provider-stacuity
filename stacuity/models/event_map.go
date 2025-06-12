// Copyright (c) HashiCorp, Inc.

package models

type EventMapList struct {
	Success    bool               `json:"success"`
	Messages   []string           `json:"messages"`
	TotalItems int32              `json:"totalItems"`
	Limit      int32              `json:"limit"`
	Offset     int32              `json:"offset"`
	Data       []EventMapReadItem `json:"data"`
}

type EventMapSingle struct {
	Success    bool             `json:"success"`
	Messages   []string         `json:"messages"`
	TotalItems int32            `json:"totalItems"`
	Limit      int32            `json:"limit"`
	Offset     int32            `json:"offset"`
	Data       EventMapReadItem `json:"data"`
}

type EventMapResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
	Data     string   `json:"data"`
}

type EventMapReadItem struct {
	Id            string          `json:"id,omitempty"`
	Moniker       string          `json:"moniker"`
	Name          string          `json:"name"`
	EventScope    EventScope      `json:"eventScope"`
	Subscriptions *[]Subscription `json:"subscriptions"`
}

type Subscription struct {
	Id            string        `json:"id,omitempty"`
	EventMap      EventMap      `json:"eventMap"`
	EventEndpoint EventEndpoint `json:"eventEndpoint"`
	EventType     EventType     `json:"eventType"`
}

type SubscriptionResponseItem struct {
	EventSubscriptionId string       `json:"eventSubscriptionId"`
	EventEndpointId     string       `json:"eventEndpointId"`
	EventType           EventType    `json:"eventType"`
	EventHandler        EventHandler `json:"eventHandler"`
}

type EventHandler struct {
	Id      string `json:"id,omitempty"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
}

type EventType struct {
	Key     int32  `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type EventEndpoint struct {
	Id                 string `json:"id,omitempty"`
	Name               string `json:"name"`
	Moniker            string `json:"moniker"`
	Type               string `json:"type"`
	Active             bool   `json:"active,omitempty"`
	SummaryDescription string `json:"summaryDescription"`
}

type SubscriptionResponse struct {
	Success  bool                       `json:"success"`
	Messages []string                   `json:"messages"`
	Data     []SubscriptionResponseItem `json:"data"`
}

type SubscriptionList struct {
	Success    bool           `json:"success"`
	Messages   []string       `json:"messages"`
	TotalItems int32          `json:"totalItems"`
	Limit      int32          `json:"limit"`
	Offset     int32          `json:"offset"`
	Data       []Subscription `json:"data"`
}

type EventScope struct {
	Key     int32  `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type EventMapModifyItem struct {
	Id            string                            `json:"id,omitempty"`
	Name          string                            `json:"name"`
	Moniker       string                            `json:"moniker"`
	EventScope    string                            `json:"eventScope"`
	Subscriptions *[]EventMapSubscriptionModifyItem `json:"subscriptions"`
}

type EventMapSubscriptionModifyItem struct {
	EventEndpointId string `json:"eventEndpointId,omitempty"`
	EventTypeId     string `json:"eventTypeId"`
}
