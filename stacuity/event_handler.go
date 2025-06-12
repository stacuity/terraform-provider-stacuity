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

// GetEventHandlers - Returns list of EventHandlers
func (c *Client) GetEventHandlers(pagingState models.PagingState) ([]models.EventHandlerReadItem, error) {
	querystring, _ := query.Values(pagingState)
	EventHandlerItems := []models.EventHandlerReadItem{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/EventHandlers?", c.HostURL)+querystring.Encode(), nil)
	if err != nil {
		return EventHandlerItems, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return EventHandlerItems, err
	}

	apiResponse := models.EventHandlerList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return EventHandlerItems, err
	}

	if !apiResponse.Success {
		return EventHandlerItems, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	EventHandlerItems = append(EventHandlerItems, apiResponse.Data...)

	return EventHandlerItems, nil
}

// GetEventHandler - Returns a specific EventHandler
func (c *Client) GetEventHandler(EventHandlerId string) (models.EventHandlerReadItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/EventHandlers/%s", c.HostURL, EventHandlerId), nil)
	if err != nil {
		return models.EventHandlerReadItem{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return models.EventHandlerReadItem{}, err
	}

	apiResponse := models.EventHandlerSingle{}
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return models.EventHandlerReadItem{}, err
	}

	if !apiResponse.Success {
		return models.EventHandlerReadItem{}, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return apiResponse.Data, nil
}

// CreateEventHandler - Create a new Event Handler
func (c *Client) CreateEventHandler(EventHandler models.EventHandlerModifyItem) (*models.EventHandlerResponse, error) {
	rb, err := json.Marshal(EventHandler)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/EventHandlers", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.EventHandlerResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// UpdateEventHandler - Update a new Event Handler
func (c *Client) UpdateEventHandler(EventHandlerId string, EventHandler models.EventHandlerModifyItem) (*models.EventHandlerResponse, error) {
	rb, err := json.Marshal(EventHandler)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/EventHandlers/%s", c.HostURL, EventHandlerId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.EventHandlerResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// DeleteEventHandler - Delete a EventHandler
func (c *Client) DeleteEventHandler(EventHandlerId string) (*models.EventHandlerResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/EventHandlers/%s", c.HostURL, EventHandlerId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.EventHandlerResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}
