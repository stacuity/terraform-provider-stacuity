// Copyright (c) HashiCorp, Inc.

package models

type OperatorPolicyList struct {
	Success    bool                     `json:"success"`
	Messages   []string                 `json:"messages"`
	TotalItems int32                    `json:"totalItems"`
	Limit      int32                    `json:"limit"`
	Offset     int32                    `json:"offset"`
	Data       []OperatorPolicyReadItem `json:"data"`
}

type OperatorPolicySingle struct {
	Success    bool                   `json:"success"`
	Messages   []string               `json:"messages"`
	TotalItems int32                  `json:"totalItems"`
	Limit      int32                  `json:"limit"`
	Offset     int32                  `json:"offset"`
	Data       OperatorPolicyReadItem `json:"data"`
}

type OperatorPolicyResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
	Data     string   `json:"data"`
}

type OperatorPolicyReadItem struct {
	Id       string                 `json:"id,omitempty"`
	Moniker  string                 `json:"moniker"`
	Name     string                 `json:"name"`
	Allow2g  bool                   `json:"allow2g"`
	Allow3g  bool                   `json:"allow3g"`
	Allow45g bool                   `json:"allow45g"`
	Entries  *[]OperatorPolicyEntry `json:"entries"`
}

type OperatorPolicyModifyItem struct {
	Id      string                           `json:"id,omitempty"`
	Name    string                           `json:"name"`
	Moniker string                           `json:"moniker"`
	Entries *[]OperatorPolicyEntryModifyItem `json:"entries"`
}

type OperatorPolicyEntry struct {
	Id                         string                     `json:"id,omitempty"`
	OperatorId                 *int32                     `json:"operatorId"`
	Iso3                       *string                    `json:"iso3"`
	SteeringProfileEntryAction SteeringProfileEntryAction `json:"steeringProfileEntryAction"`
}

type SteeringProfileEntryAction struct {
	Key     int32  `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type OperatorPolicyEntryResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
}

type OperatorPolicyEntryList struct {
	Success    bool                  `json:"success"`
	Messages   []string              `json:"messages"`
	TotalItems int32                 `json:"totalItems"`
	Limit      int32                 `json:"limit"`
	Offset     int32                 `json:"offset"`
	Data       []OperatorPolicyEntry `json:"data"`
}

type OperatorPolicyEntryModifyItem struct {
	OperatorId                 *int32  `json:"operatorId"`
	Iso3                       *string `json:"iso3"`
	SteeringProfileEntryAction string  `json:"steeringProfileEntryAction"`
}
