// Copyright (c) HashiCorp, Inc.

package models

type RegionalPolicyList struct {
	Success    bool                     `json:"success"`
	Messages   []string                 `json:"messages"`
	TotalItems int32                    `json:"totalItems"`
	Limit      int32                    `json:"limit"`
	Offset     int32                    `json:"offset"`
	Data       []RegionalPolicyReadItem `json:"data"`
}

type RegionalPolicySingle struct {
	Success    bool                   `json:"success"`
	Messages   []string               `json:"messages"`
	TotalItems int32                  `json:"totalItems"`
	Limit      int32                  `json:"limit"`
	Offset     int32                  `json:"offset"`
	Data       RegionalPolicyReadItem `json:"data"`
}

type RegionalPolicyResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
	Data     string   `json:"data"`
}

type RegionalPolicyReadItem struct {
	Id      string                 `json:"id,omitempty"`
	Moniker string                 `json:"moniker"`
	Name    string                 `json:"name"`
	Active  bool                   `json:"active"`
	IsFixed bool                   `json:"isFixed"`
	Entries *[]RegionalPolicyEntry `json:"entries"`
}

type RegionalPolicyModifyItem struct {
	Id      string                           `json:"id,omitempty"`
	Name    string                           `json:"name"`
	Moniker string                           `json:"moniker"`
	Entries *[]RegionalPolicyEntryModifyItem `json:"entries"`
}

type RegionalPolicyEntry struct {
	Id                      string          `json:"id,omitempty"`
	OperatorId              *int32          `json:"operatorId"`
	RegionalGatewayPolicyId *string         `json:"regionalGatewayPolicyId"`
	Iso3                    *string         `json:"iso3"`
	RegionalGateway         regionalGateway `json:"regionalGateway"`
}

type regionalGateway struct {
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
}

type RegionalPolicyEntryResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
}

type RegionalPolicyEntryList struct {
	Success    bool                  `json:"success"`
	Messages   []string              `json:"messages"`
	TotalItems int32                 `json:"totalItems"`
	Limit      int32                 `json:"limit"`
	Offset     int32                 `json:"offset"`
	Data       []RegionalPolicyEntry `json:"data"`
}

type RegionalPolicyEntryModifyItem struct {
	RegionalGatewayId *string `json:"regionalGatewayId"`
	OperatorId        *int32  `json:"operatorId"`
	Iso3              *string `json:"iso3"`
}
