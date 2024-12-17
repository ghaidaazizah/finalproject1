package service

import (
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AIService struct {
	Client HTTPClient
}

func (s *AIService) AnalyzeData(table map[string][]string, query, token string) (string, error) {
	if table == nil || len(table) == 0 {
		return "", fmt.Errorf("invalid input: table is empty")
	}

	payload := model.Inputs{
		Table: table,
		Query: query,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/google/tapas-large-finetuned-wtq", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(body))
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseData, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	cells, ok := result["cells"].([]interface{})
	if !ok || len(cells) == 0 {
		return "", fmt.Errorf("response missing or invalid 'cells' key")
	}

	res, ok := cells[0].(string)
	if !ok {
		return "", fmt.Errorf("first element in 'cells' is not a string")
	}

	return res, nil
}

func (s *AIService) ChatWithAI(context, query, token string) (model.ChatResponse, error) {
	type Payload struct {
		Inputs string `json:"inputs"`
	}

	bodyPayload := Payload{
		Inputs: query,
	}

	fmt.Println(bodyPayload)

	jsonData, err := json.Marshal(&bodyPayload)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api-inference.huggingface.co/models/microsoft/Phi-3.5-mini-instruct", bytes.NewBuffer(jsonData))
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Println(string(responseData))

	var chatResponses []model.ChatResponse
	if err := json.Unmarshal(responseData, &chatResponses); err != nil {
		return model.ChatResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(chatResponses) == 0 {
		return model.ChatResponse{}, fmt.Errorf("no responses received from API")
	}

	return chatResponses[0], nil
}
