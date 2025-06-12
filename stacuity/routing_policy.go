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

// GetRoutingPolicies - Returns list of RoutingPolicies
func (c *Client) GetRoutingPolicies(pagingState models.PagingState) ([]models.RoutingPolicyReadItem, error) {
	querystring, _ := query.Values(pagingState)
	RoutingPolicyItems := []models.RoutingPolicyReadItem{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/RoutingPolicies?", c.HostURL)+querystring.Encode(), nil)
	if err != nil {
		return RoutingPolicyItems, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return RoutingPolicyItems, err
	}

	apiResponse := models.RoutingPolicyList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return RoutingPolicyItems, err
	}

	if !apiResponse.Success {
		return RoutingPolicyItems, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	RoutingPolicyItems = append(RoutingPolicyItems, apiResponse.Data...)

	return RoutingPolicyItems, nil
}

// GetRoutingPolicy - Returns a specific RoutingPolicy
func (c *Client) GetRoutingPolicy(RoutingPolicyId string) (models.RoutingPolicyReadItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/RoutingPolicies/%s", c.HostURL, RoutingPolicyId), nil)
	if err != nil {
		return models.RoutingPolicyReadItem{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return models.RoutingPolicyReadItem{}, err
	}

	apiResponse := models.RoutingPolicySingle{}
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return models.RoutingPolicyReadItem{}, err
	}

	if !apiResponse.Success {
		return models.RoutingPolicyReadItem{}, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return apiResponse.Data, nil
}

// CreateRoutingPolicy - Create a new Routing Policy
func (c *Client) CreateRoutingPolicy(RoutingPolicy models.RoutingPolicyModifyItem) (*models.RoutingPolicyResponse, error) {
	rb, err := json.Marshal(RoutingPolicy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/RoutingPolicies", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.RoutingPolicyResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// UpdateRoutingPolicy - Update a new Routing Policy
func (c *Client) UpdateRoutingPolicy(RoutingPolicyId string, RoutingPolicy models.RoutingPolicyModifyItem) (*models.RoutingPolicyResponse, error) {
	rb, err := json.Marshal(RoutingPolicy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/RoutingPolicies/%s", c.HostURL, RoutingPolicyId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.RoutingPolicyResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// DeleteRoutingPolicy - Delete a Routing Policy
func (c *Client) DeleteRoutingPolicy(RoutingPolicyId string) (*models.RoutingPolicyResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/RoutingPolicies/%s", c.HostURL, RoutingPolicyId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.RoutingPolicyResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}
