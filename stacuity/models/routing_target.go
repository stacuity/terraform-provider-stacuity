// Copyright (c) HashiCorp, Inc.

package models

type RoutingTargetList struct {
	Success    bool                    `json:"success"`
	Messages   []string                `json:"messages"`
	TotalItems int32                   `json:"totalItems"`
	Limit      int32                   `json:"limit"`
	Offset     int32                   `json:"offset"`
	Data       []RoutingTargetReadItem `json:"data"`
}

type RoutingTargetSingle struct {
	Success    bool                  `json:"success"`
	Messages   []string              `json:"messages"`
	TotalItems int32                 `json:"totalItems"`
	Limit      int32                 `json:"limit"`
	Offset     int32                 `json:"offset"`
	Data       RoutingTargetReadItem `json:"data"`
}

type RoutingTargetResponse struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
	Data     string   `json:"data"`
}

type RoutingTargetReadItem struct {
	Id                           string                    `json:"id,omitempty"`
	Name                         string                    `json:"name"`
	Moniker                      string                    `json:"moniker"`
	RoutingTargetType            RoutingTargetType         `json:"routingTargetType"`
	RoutingTargetStatus          RoutingTargetStatus       `json:"routingTargetStatus"`
	VSlice                       VSlice                    `json:"vSlice"`
	ConfigurationData            *ConfigurationData        `json:"configurationData,omitempty"`
	PublicInstanceConfiguration  string                    `json:"publicInstanceConfiguration"`
	NetworkKey                   int32                     `json:"networkKey"`
	RoutingTargetTypeInstance    RoutingTargetTypeInstance `json:"routingTargetTypeInstance"`
	RoutingRedundancyZoneMoniker string                    `json:"routingRedundancyZoneMoniker"`
	RoutingRedundancyZoneName    string                    `json:"routingRedundancyZoneName"`
	RegionalGatewayMoniker       string                    `json:"regionalGatewayMoniker"`
	RegionalGatewayName          string                    `json:"regionalGatewayName"`
}

type RoutingTargetTypeInstance struct {
	Id      int32  `json:"id,omitempty"`
	Name    string `json:"name"`
	Moniker string `json:"moniker"`
}

type ConfigurationData struct {
	WireGuardConfig *WireGuardConfig `json:"wireGuardConfig,omitempty"`
	VpnConfig       *VpnConfig       `json:"vpnConfig,omitempty"`
}

type WireGuardConfig struct {
	LocalSubnets         string `json:"localSubnets"`
	LocalPublicKey       string `json:"localPublicKey"`
	RemotePublicKey      string `json:"remotePublicKey"`
	RemoteSubnets        string `json:"remoteSubnets"`
	RemotePeerIPAddress  string `json:"remotePeerIPAddress"`
	RemotePeerPortNumber int32  `json:"remotePeerPortNumber"`
}

type VpnConfig struct {
	LocalSubnets           string `json:"localSubnets"`
	RemoteSubnets          string `json:"remoteSubnets"`
	RemotePeerAddress      string `json:"remotePeerAddress"`
	RemoteEncryptionDomain string `json:"remoteEncryptionDomain"`
	LocalEncryptionDomain  string `json:"localEncryptionDomain"`
	PresharedKey           string `json:"presharedKey"`
	KeyExchangeType        string `json:"keyExchangeType"`
	VpnIkeOption           string `json:"vpnIkeOption"`
	VpnEspOption           string `json:"vpnEspOption"`
	Phase1Lifetime         int32  `json:"phase1Lifetime"`
	Phase2Lifetime         int32  `json:"phase2Lifetime"`
}

type RoutingTargetType struct {
	Key     int32  `json:"key,omitempty"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type RoutingTargetStatus struct {
	Key     int32  `json:"key,omitempty"`
	Moniker string `json:"moniker"`
	Name    string `json:"name"`
	Active  bool   `json:"active"`
}

type RoutingTargetModifyItem struct {
	Id                           string             `json:"id,omitempty"`
	Name                         string             `json:"name"`
	Moniker                      string             `json:"moniker"`
	RoutingTargetType            string             `json:"routingTargetType"`
	RoutingRedundancyZoneMoniker string             `json:"routingRedundancyZoneMoniker"`
	ConfigurationData            *ConfigurationData `json:"configurationData,omitempty"`
	SubnetAddress                string             `json:"subnetAddress"`
	VSlice                       string             `json:"vslice"`
	RoutingTargetTypeInstanceId  string             `json:"routingTargetTypeInstanceId"`
}
