// Copyright (c) HashiCorp, Inc.

package models

type EventHandlerList struct {
	Success    bool                   `json:"success"`
	Messages   []string               `json:"messages"`
	TotalItems int32                  `json:"totalItems"`
	Limit      int32                  `json:"limit"`
	Offset     int32                  `json:"offset"`
	Data       []EventHandlerReadItem `json:"data"`
}

type EventHandlerSingle struct {
	Success    bool                 `json:"success"`
	Messages   []string             `json:"messages"`
	TotalItems int32                `json:"totalItems"`
	Limit      int32                `json:"limit"`
	Offset     int32                `json:"offset"`
	Data       EventHandlerReadItem `json:"data"`
}

type EventHandlerResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
	Data     string   `json:"data"`
}

type EventHandlerReadItem struct {
	Id                 string        `json:"id,omitempty"`
	Moniker            string        `json:"moniker"`
	Name               string        `json:"name"`
	ConfigurationData  Configuration `json:"configurationData"`
	SummaryDescription string        `json:"summaryDescription"`
	EventEndpointType  EventEndpoint `json:"eventEndpointType"`
}

type Configuration struct {
	WebhookConfig *WebhookConfig `json:"webhookConfig"`
}

type WebhookConfig struct {
	BearerToken *string `json:"bearerToken"`
	Username    *string `json:"username"`
	Password    *string `json:"password"`
	Timeout     string  `json:"timeout"`
	Url         string  `json:"url"`
}

type EventHandlerModifyItem struct {
	Id                string        `json:"id,omitempty"`
	Name              string        `json:"name"`
	Moniker           string        `json:"moniker"`
	ConfigurationData Configuration `json:"configurationData"`
	EventEndpointType string        `json:"eventEndpointType"`
}
