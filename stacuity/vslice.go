// Copyright (c) HashiCorp, Inc.

package stacuity

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"stacuity.com/go_client/models"

	"github.com/google/go-querystring/query"
)

// GetVSlices - Returns list of VSlices
func (c *Client) GetVSlices(pagingState models.PagingState) ([]models.VSliceReadItem, error) {
	querystring, _ := query.Values(pagingState)
	vSliceItems := []models.VSliceReadItem{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/vslices?", c.HostURL)+querystring.Encode(), nil)
	if err != nil {
		return vSliceItems, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return vSliceItems, err
	}

	vSlices := models.VSliceList{}
	err = json.Unmarshal(body, &vSlices)

	if err != nil {
		return vSliceItems, err
	}

	if !vSlices.Success {
		return vSliceItems, errors.New(strings.Join(vSlices.Messages, " "))
	}

	vSliceItems = append(vSliceItems, vSlices.Data...)

	return vSliceItems, nil
}

// GetVSlice - Returns a specific VSlice
func (c *Client) GetVSlice(vSliceId string) (models.VSliceReadItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/vslices/%s", c.HostURL, vSliceId), nil)
	if err != nil {
		return models.VSliceReadItem{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return models.VSliceReadItem{}, err
	}

	vSlices := models.VSliceSingle{}
	err = json.Unmarshal(body, &vSlices)

	if err != nil {
		return models.VSliceReadItem{}, err
	}

	if !vSlices.Success {
		return models.VSliceReadItem{}, errors.New(strings.Join(vSlices.Messages, " "))
	}

	return vSlices.Data, nil
}

// CreateVSlice - Create a new vSlice
func (c *Client) CreateVSlice(vSlice models.VSliceModifyItem) (*models.VSliceResponse, error) {
	rb, err := json.Marshal(vSlice)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/vslices", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.VSliceResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// UpdateVSlice - Update a new vSlice
func (c *Client) UpdateVSlice(vSliceId string, vSlice models.VSliceModifyItem) (*models.VSliceResponse, error) {
	rb, err := json.Marshal(vSlice)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/vslices/%s", c.HostURL, vSliceId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.VSliceResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// DeleteVSlice - Delete a vSlice
func (c *Client) DeleteVSlice(vSliceId string) (*models.VSliceResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/vslices/%s", c.HostURL, vSliceId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.VSliceResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}
