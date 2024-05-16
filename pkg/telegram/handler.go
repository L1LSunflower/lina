package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/L1LSunflower/lina/internal/entities"
	"github.com/go-playground/validator/v10"
)

type Handler interface {
	MediaGroup(data *entities.MediaGroup) (*http.Response, error)
	SendMessage(data *entities.Message) (*http.Response, error)
	Send(data json.Marshaler, method string, headers map[string]string) (*http.Response, error)
}

type THandler struct {
	url       string
	apiToken  string
	client    *http.Client
	validator *validator.Validate
}

func NewTHandler(url string, apiToken string) *THandler {
	return &THandler{
		url:      url,
		apiToken: "bot" + apiToken,
	}
}

func (h *THandler) MediaGroup(data *entities.MediaGroup) (*http.Response, error) {
	resp, err := h.send(data, "sendMediaGroup", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram response status code is not 200: %d, response: %v", resp.StatusCode, resp)
	}
	return resp, nil
}

func (h *THandler) SendMessage(data *entities.Message) (*http.Response, error) {
	resp, err := h.send(data, "sendMessage", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram response status code is not 200: %d, response: %v", resp.StatusCode, resp)
	}
	return resp, nil
}

func (h *THandler) send(data json.Marshaler, method string, headers map[string]string) (*http.Response, error) {
	b, err := data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, h.url+"/bot"+h.apiToken+"/"+method, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
