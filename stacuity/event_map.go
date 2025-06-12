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

// GetEventMaps - Returns list of EventMaps
func (c *Client) GetEventMaps(pagingState models.PagingState) ([]models.EventMapReadItem, error) {
	querystring, _ := query.Values(pagingState)
	EventMapItems := []models.EventMapReadItem{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/EventMaps?", c.HostURL)+querystring.Encode(), nil)
	if err != nil {
		return EventMapItems, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return EventMapItems, err
	}

	apiResponse := models.EventMapList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return EventMapItems, err
	}

	if !apiResponse.Success {
		return EventMapItems, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	EventMapItems = append(EventMapItems, apiResponse.Data...)

	return EventMapItems, nil
}

// GetEventMap - Returns a specific EventMap
func (c *Client) GetEventMap(EventMapId string) (models.EventMapReadItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/EventMaps/%s", c.HostURL, EventMapId), nil)
	if err != nil {
		return models.EventMapReadItem{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return models.EventMapReadItem{}, err
	}

	apiResponse := models.EventMapSingle{}
	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		return models.EventMapReadItem{}, err
	}

	if !apiResponse.Success {
		return models.EventMapReadItem{}, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return apiResponse.Data, nil
}

// GetEventMapSubscriptions - Returns a specific EventMap subscriptions
func (c *Client) GetEventMapSubscriptions(EventMapId string) ([]models.Subscription, error) {
	eventMapSubscriptionItems := []models.Subscription{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/EventMaps/%s/subscriptions", c.HostURL, EventMapId), nil)
	if err != nil {
		return eventMapSubscriptionItems, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return eventMapSubscriptionItems, err
	}

	apiResponse := models.SubscriptionList{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return eventMapSubscriptionItems, err
	}

	if !apiResponse.Success {
		return eventMapSubscriptionItems, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	eventMapSubscriptionItems = append(eventMapSubscriptionItems, apiResponse.Data...)

	return eventMapSubscriptionItems, nil
}

// CreateEventMap - Create a new Event Map
func (c *Client) CreateEventMap(EventMap models.EventMapModifyItem) (*models.EventMapResponse, error) {
	rb, err := json.Marshal(EventMap)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/EventMaps", c.HostURL), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.EventMapResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// AddEventMapSubscriptions - add new Event Map subscriptions
func (c *Client) AddEventMapSubscriptions(EventMapSubscription []models.EventMapSubscriptionModifyItem, EventMapId string) (*models.SubscriptionResponse, error) {
	rb, err := json.Marshal(EventMapSubscription)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/EventMaps/%s/subscriptions", c.HostURL, EventMapId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.SubscriptionResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// UpdateEventMap - Update a new Event Map
func (c *Client) UpdateEventMap(EventMapId string, EventMap models.EventMapModifyItem) (*models.EventMapResponse, error) {
	rb, err := json.Marshal(EventMap)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/EventMaps/%s", c.HostURL, EventMapId), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.EventMapResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}

// DeleteEventMap - Delete a Event Map
func (c *Client) DeleteEventMap(EventMapId string) (*models.EventMapResponse, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/EventMaps/%s", c.HostURL, EventMapId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	apiResponse := models.EventMapResponse{}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	if !apiResponse.Success {
		return nil, errors.New(strings.Join(apiResponse.Messages, " "))
	}

	return &apiResponse, nil
}
