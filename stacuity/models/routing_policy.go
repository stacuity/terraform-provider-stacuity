// Copyright (c) HashiCorp, Inc.

package models

type RoutingPolicyList struct {
	Success    bool                    `json:"success"`
	Messages   []string                `json:"messages"`
	TotalItems int32                   `json:"totalItems"`
	Limit      int32                   `json:"limit"`
	Offset     int32                   `json:"offset"`
	Data       []RoutingPolicyReadItem `json:"data"`
}

type RoutingPolicySingle struct {
	Success    bool                  `json:"success"`
	Messages   []string              `json:"messages"`
	TotalItems int32                 `json:"totalItems"`
	Limit      int32                 `json:"limit"`
	Offset     int32                 `json:"offset"`
	Data       RoutingPolicyReadItem `json:"data"`
}

type RoutingPolicyResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
	Data     string   `json:"data"`
}

type RoutingPolicyReadItem struct {
	Id                              string         `json:"id,omitempty"`
	Name                            string         `json:"name"`
	Moniker                         string         `json:"moniker"`
	VSlice                          VSlice         `json:"vslice"`
	RateLimitUplink                 RateLimit      `json:"rateLimitUplink"`
	RateLimitDownlink               RateLimit      `json:"rateLimitDownlink"`
	PacketDiscardUplinkPercentage   int32          `json:"packetDiscardUplinkPercentage"`
	PacketDiscardDownlinkPercentage int32          `json:"packetDiscardDownlinkPercentage"`
	RoutingPolicyStatus             RoutingPolicy  `json:"routingPolicyStatus"`
	RoutingPolicyRules              *[]Rule        `json:"routingPolicyRules"`
	RoutingPolicyEdgeServices       *[]EdgeService `json:"routingPolicyEdgeServices"`
}

type RoutingPolicyModifyItem struct {
	Id                              string        `json:"id,omitempty"`
	Name                            string        `json:"name"`
	Moniker                         string        `json:"moniker"`
	VSlice                          string        `json:"vslice"`
	RoutingPolicyStatus             string        `json:"routingPolicyStatus"`
	RoutingPolicyRules              []ModifyRule  `json:"routingPolicyRules"`
	RoutingPolicyEdgeServices       []EdgeService `json:"routingPolicyEdgeServices"`
	RateLimitUplinkMoniker          string        `json:"rateLimitUplinkMoniker"`
	RateLimitDownlinkMoniker        string        `json:"rateLimitDownlinkMoniker"`
	PacketDiscardUplinkPercentage   int32         `json:"packetDiscardUplinkPercentage"`
	PacketDiscardDownlinkPercentage int32         `json:"packetDiscardDownlinkPercentage"`
}

type RateLimit struct {
	Key     int    `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type RoutingPolicy struct {
	Key     int    `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type ModifyRule struct {
	Description            string  `json:"description"`
	RuleAction             string  `json:"ruleAction"`
	RuleDirection          string  `json:"ruleDirection"`
	SourceIpPattern        *string `json:"sourceIpPattern"`
	DestinationIpPattern   *string `json:"destinationIpPattern"`
	DivertIp               *string `json:"divertIp"`
	DivertPort             *string `json:"divertPort"`
	TransportProtocol      *string `json:"transportProtocol"`
	SourcePortPattern      *string `json:"sourcePortPattern"`
	DestinationPortPattern *string `json:"destinationPortPattern"`
	RoutingTarget          *string `json:"routingTarget"`
	Reflexive              bool    `json:"reflexive"`
	RegionalGateway        *string `json:"regionalGateway"`
	Enabled                bool    `json:"enabled"`
}

type Rule struct {
	Id                     string             `json:"id"`
	RoutingPolicyId        string             `json:"routingPolicyId"`
	Description            string             `json:"description"`
	RuleAction             RuleAction         `json:"ruleAction"`
	RuleDirection          RuleDirection      `json:"ruleDirection"`
	Precedence             int32              `json:"precedence"`
	SourceIpPattern        *string            `json:"sourceIpPattern"`
	DestinationIpPattern   *string            `json:"destinationIpPattern"`
	DivertIp               *string            `json:"divertIp"`
	DivertPort             *string            `json:"divertPort"`
	TransportProtocol      *TransportProtocol `json:"transportProtocol"`
	SourcePortPattern      *string            `json:"sourcePortPattern"`
	DestinationPortPattern *string            `json:"destinationPortPattern"`
	RoutingTarget          *RoutingTarget     `json:"routingTarget"`
	Reflexive              bool               `json:"reflexive"`
	Enabled                bool               `json:"enabled"`
	RegionalGateway        *RegionalGateway   `json:"regionalGateway"`
}

type RuleAction struct {
	Key     int    `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type RuleDirection struct {
	Key     int    `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type TransportProtocol struct {
	Key     int    `json:"key"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type RoutingTarget struct {
	Id      string `json:"id"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
}

type RegionalGateway struct {
	Id      string `json:"id"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
}

type EdgeService struct {
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	IconShape              string   `json:"iconShape"`
	Moniker                string   `json:"moniker"`
	Available              bool     `json:"available"`
	Enabled                bool     `json:"enabled"`
	HasInstance            bool     `json:"hasInstance"`
	EdgeServiceInstanceIds []string `json:"edgeServiceInstanceIds"`
}
