// Copyright (c) HashiCorp, Inc.

package stacuity

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
	"stacuity.com/go_client/models"
)

// GetRoutingTargets - Returns list of RoutingTargets
func (c *Client) GetRoutingTargets(pagingState models.PagingState) ([]models.RoutingTargetReadItem, error) {
	querystring, _ := query.Values(pagingState)
	routingTargetItems := []models.RoutingTargetReadItem{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/routingtargets?", c.HostURL)+querystring.Encode(), nil)
	if err != nil {
		return routingTargetItems, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return routingTargetItems, err
	}

	apiResponse := models.RoutingTargetList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return routingTargetItems, err
	}

	if !apiResponse.Success {
		return routingTargetItems, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	routingTargetItems = append(routingTargetItems, apiResponse.Data...)

	return routingTargetItems, nil
}

// GetRoutingTarget - Returns a specific RoutingTarget
func (c *Client) GetRoutingTarget(routingTargetId string) (models.RoutingTargetReadItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/routingtargets/%s", c.HostURL, routingTargetId), nil)
	if err != nil {
		return models.RoutingTargetReadItem{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return models.RoutingTargetReadItem{}, err
	}

	apiResponse := models.RoutingTargetSingle{}
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return models.RoutingTargetReadItem{}, err
	}

	if !apiResponse.Success {
		return models.RoutingTargetReadItem{}, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return apiResponse.Data, nil
}

// CreateRoutingTarget - Create a new Routing Target
func (c *Client) CreateRoutingTarget(routingTarget models.RoutingTargetModifyItem) (*models.RoutingTargetResponse, error) {
	rb, err := json.Marshal(routingTarget)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/routingtargets", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.RoutingTargetResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// UpdateRoutingTarget - Update a new Routing Target
func (c *Client) UpdateRoutingTarget(routingTargetId string, routingTarget models.RoutingTargetModifyItem) (*models.RoutingTargetResponse, error) {
	rb, err := json.Marshal(routingTarget)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/routingtargets/%s", c.HostURL, routingTargetId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.RoutingTargetResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// DeleteRoutingTarget - Delete a Routing Target
func (c *Client) DeleteRoutingTarget(routingTargetId string) (*models.RoutingTargetResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/routingtargets/%s", c.HostURL, routingTargetId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.RoutingTargetResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}
