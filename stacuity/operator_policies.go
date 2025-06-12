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

// GetOperatorPolicies - Returns list of OperatorPolicies
func (c *Client) GetOperatorPolicies(pagingState models.PagingState) ([]models.OperatorPolicyReadItem, error) {
	querystring, _ := query.Values(pagingState)
	OperatorPolicyItems := []models.OperatorPolicyReadItem{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/OperatorPolicies?", c.HostURL)+querystring.Encode(), nil)
	if err != nil {
		return OperatorPolicyItems, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return OperatorPolicyItems, err
	}

	apiResponse := models.OperatorPolicyList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return OperatorPolicyItems, err
	}

	if !apiResponse.Success {
		return OperatorPolicyItems, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	OperatorPolicyItems = append(OperatorPolicyItems, apiResponse.Data...)

	return OperatorPolicyItems, nil
}

// GetOperatorPolicy - Returns a specific OperatorPolicy
func (c *Client) GetOperatorPolicy(OperatorPolicyId string) (models.OperatorPolicyReadItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/OperatorPolicies/%s", c.HostURL, OperatorPolicyId), nil)
	if err != nil {
		return models.OperatorPolicyReadItem{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return models.OperatorPolicyReadItem{}, err
	}

	apiResponse := models.OperatorPolicySingle{}
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return models.OperatorPolicyReadItem{}, err
	}

	if !apiResponse.Success {
		return models.OperatorPolicyReadItem{}, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return apiResponse.Data, nil
}

// GetOperatorPolicyEntries - Returns a specific OperatorPolicy entries
func (c *Client) GetOperatorPolicyEntries(OperatorPolicyId string) ([]models.OperatorPolicyEntry, error) {
	OperatorPolicyEntries := []models.OperatorPolicyEntry{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/OperatorPolicies/%s/entries", c.HostURL, OperatorPolicyId), nil)
	if err != nil {
		return OperatorPolicyEntries, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return OperatorPolicyEntries, err
	}

	apiResponse := models.OperatorPolicyEntryList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return OperatorPolicyEntries, err
	}

	if !apiResponse.Success {
		return OperatorPolicyEntries, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	OperatorPolicyEntries = append(OperatorPolicyEntries, apiResponse.Data...)

	return OperatorPolicyEntries, nil
}

// CreateOperatorPolicy - Create a new Operator Policy
func (c *Client) CreateOperatorPolicy(OperatorPolicy models.OperatorPolicyModifyItem) (*models.OperatorPolicyResponse, error) {
	rb, err := json.Marshal(OperatorPolicy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/OperatorPolicies", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.OperatorPolicyResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// AddOperatorPolicyEntries - add new Operator Policy entries
func (c *Client) AddOperatorPolicyEntries(OperatorPolicyEntries []models.OperatorPolicyEntryModifyItem, OperatorPolicyId string) (*models.OperatorPolicyEntryResponse, error) {
	rb, err := json.Marshal(OperatorPolicyEntries)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/OperatorPolicies/%s/entries", c.HostURL, OperatorPolicyId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.OperatorPolicyEntryResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// UpdateOperatorPolicy - Update a new Operator Policy
func (c *Client) UpdateOperatorPolicy(OperatorPolicyId string, OperatorPolicy models.OperatorPolicyModifyItem) (*models.OperatorPolicyResponse, error) {
	rb, err := json.Marshal(OperatorPolicy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/OperatorPolicies/%s", c.HostURL, OperatorPolicyId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.OperatorPolicyResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// DeleteOperatorPolicy - Delete a Operator Policy
func (c *Client) DeleteOperatorPolicy(OperatorPolicyId string) (*models.OperatorPolicyResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/OperatorPolicies/%s", c.HostURL, OperatorPolicyId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.OperatorPolicyResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}
