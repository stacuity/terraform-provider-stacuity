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

// GetEndpointGroups - Returns list of EndpointGroups
func (c *Client) GetEndpointGroups(pagingState models.PagingState) ([]models.EndpointGroupReadItem, error) {
	querystring, _ := query.Values(pagingState)
	EndpointGroupItems := []models.EndpointGroupReadItem{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/EndpointGroups?", c.HostURL)+querystring.Encode(), nil)
	if err != nil {
		return EndpointGroupItems, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return EndpointGroupItems, err
	}

	apiResponse := models.EndpointGroupList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return EndpointGroupItems, err
	}

	if !apiResponse.Success {
		return EndpointGroupItems, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	EndpointGroupItems = append(EndpointGroupItems, apiResponse.Data...)

	return EndpointGroupItems, nil
}

// GetEndpointGroup - Returns a specific EndpointGroup
func (c *Client) GetEndpointGroup(EndpointGroupId string) (models.EndpointGroupReadItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/EndpointGroups/%s", c.HostURL, EndpointGroupId), nil)
	if err != nil {
		return models.EndpointGroupReadItem{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return models.EndpointGroupReadItem{}, err
	}

	apiResponse := models.EndpointGroupSingle{}
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return models.EndpointGroupReadItem{}, err
	}

	if !apiResponse.Success {
		return models.EndpointGroupReadItem{}, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return apiResponse.Data, nil
}

// CreateEndpointGroup - Create a new Routing Policy
func (c *Client) CreateEndpointGroup(EndpointGroup models.EndpointGroupModifyItem) (*models.EndpointGroupResponse, error) {
	rb, err := json.Marshal(EndpointGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/EndpointGroups", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.EndpointGroupResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// UpdateEndpointGroup - Update a new Routing Policy
func (c *Client) UpdateEndpointGroup(EndpointGroupId string, EndpointGroup models.EndpointGroupModifyItem) (*models.EndpointGroupResponse, error) {
	rb, err := json.Marshal(EndpointGroup)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/EndpointGroups/%s", c.HostURL, EndpointGroupId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.EndpointGroupResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// DeleteEndpointGroup - Delete a EndpointGroup
func (c *Client) DeleteEndpointGroup(EndpointGroupId string) (*models.EndpointGroupResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/EndpointGroups/%s", c.HostURL, EndpointGroupId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.EndpointGroupResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}
