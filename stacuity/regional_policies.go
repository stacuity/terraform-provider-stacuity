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

// GetRegionalPolicies - Returns list of RegionalPolicies
func (c *Client) GetRegionalPolicies(pagingState models.PagingState) ([]models.RegionalPolicyReadItem, error) {
	querystring, _ := query.Values(pagingState)
	RegionalPolicyItems := []models.RegionalPolicyReadItem{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/RegionalPolicies?", c.HostURL)+querystring.Encode(), nil)
	if err != nil {
		return RegionalPolicyItems, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return RegionalPolicyItems, err
	}

	apiResponse := models.RegionalPolicyList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return RegionalPolicyItems, err
	}

	if !apiResponse.Success {
		return RegionalPolicyItems, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	RegionalPolicyItems = append(RegionalPolicyItems, apiResponse.Data...)

	return RegionalPolicyItems, nil
}

// GetRegionalPolicy - Returns a specific RegionalPolicy
func (c *Client) GetRegionalPolicy(RegionalPolicyId string) (models.RegionalPolicyReadItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/RegionalPolicies/%s", c.HostURL, RegionalPolicyId), nil)
	if err != nil {
		return models.RegionalPolicyReadItem{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return models.RegionalPolicyReadItem{}, err
	}

	apiResponse := models.RegionalPolicySingle{}
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return models.RegionalPolicyReadItem{}, err
	}

	if !apiResponse.Success {
		return models.RegionalPolicyReadItem{}, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return apiResponse.Data, nil
}

// GetRegionalPolicyEntries - Returns a specific RegionalPolicy entries
func (c *Client) GetRegionalPolicyEntries(RegionalPolicyId string) ([]models.RegionalPolicyEntry, error) {
	RegionalPolicyEntries := []models.RegionalPolicyEntry{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/RegionalPolicies/%s/entries", c.HostURL, RegionalPolicyId), nil)
	if err != nil {
		return RegionalPolicyEntries, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return RegionalPolicyEntries, err
	}

	apiResponse := models.RegionalPolicyEntryList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return RegionalPolicyEntries, err
	}

	if !apiResponse.Success {
		return RegionalPolicyEntries, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	RegionalPolicyEntries = append(RegionalPolicyEntries, apiResponse.Data...)

	return RegionalPolicyEntries, nil
}

// CreateRegionalPolicy - Create a new Regional Policy
func (c *Client) CreateRegionalPolicy(RegionalPolicy models.RegionalPolicyModifyItem) (*models.RegionalPolicyResponse, error) {
	rb, err := json.Marshal(RegionalPolicy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/RegionalPolicies", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.RegionalPolicyResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// AddRegionalPolicyEntries - add new Regional Policy entries
func (c *Client) AddRegionalPolicyEntries(RegionalPolicyEntries []models.RegionalPolicyEntryModifyItem, RegionalPolicyId string) (*models.RegionalPolicyEntryResponse, error) {
	apiResponse := models.RegionalPolicyEntryResponse{}

	for _, elem := range RegionalPolicyEntries {
		rb, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/RegionalPolicies/%s/entries", c.HostURL, RegionalPolicyId), strings.NewReader(string(rb)))
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return nil, err
		}

		if !apiResponse.Success {
			return nil, errors.New(strings.Join(apiResponse.Messages, " "))
		}
	}

	return &apiResponse, nil
}

// UpdateRegionalPolicy - Update a new Regional Policy
func (c *Client) UpdateRegionalPolicy(RegionalPolicyId string, RegionalPolicy models.RegionalPolicyModifyItem) (*models.RegionalPolicyResponse, error) {
	rb, err := json.Marshal(RegionalPolicy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/RegionalPolicies/%s", c.HostURL, RegionalPolicyId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.RegionalPolicyResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// DeleteRegionalPolicy - Delete a Regional Policy
func (c *Client) DeleteRegionalPolicy(RegionalPolicyId string) (*models.RegionalPolicyResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/RegionalPolicies/%s", c.HostURL, RegionalPolicyId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.RegionalPolicyResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}
