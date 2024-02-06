package dadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	cleanUrl   = "https://cleaner.dadata.ru/api/v1"
	suggestUrl = "https://suggestions.dadata.ru/suggestions/api/4_1/rs"
)

type DaData struct {
	client    *http.Client
	token     string
	secretKey string
}

type DaDataClean []map[string]interface{}

type DaDataSuggest struct {
	Suggest []map[string]interface{} `json:"suggestions"`
}

type DaDataIpLocate struct {
	Location map[string]interface{} `json:"location"`
}

func New(token string, secretKey string) *DaData {
	return &DaData{
		client: &http.Client{
			Transport: http.DefaultTransport,
		},
		token:     token,
		secretKey: secretKey,
	}
}

func (d *DaData) GetCleanValue(path string, body string) (*DaDataClean, error) {
	req, err := http.NewRequest("POST", cleanUrl+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+d.token)
	req.Header.Set("X-Secret", d.secretKey)
	req.Header.Set("Content-Type", "application/json")

	response, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == 429 {
		return nil, errors.New("query limit reached, come back after midnight") //TODO move to const
	}

	if response.StatusCode != 200 {
		rBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("cannot get value from dadata api with error: %s", string(rBytes))
	}

	rBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res *DaDataClean
	err = json.Unmarshal(rBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *DaData) GetSuggestValue(path string, body string) (*DaDataSuggest, error) {
	req, err := http.NewRequest("POST", suggestUrl+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+d.token)
	req.Header.Set("Content-Type", "application/json")

	response, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == 429 {
		return nil, errors.New("query limit reached, come back after midnight") //TODO move to const
	}

	if response.StatusCode != 200 {
		rBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("cannot get value from dadata api with error: %s", string(rBytes))
	}

	rBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res DaDataSuggest
	err = json.Unmarshal(rBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (d *DaData) GetIpLocateValue(path string, body string) (*DaDataIpLocate, error) {
	req, err := http.NewRequest("POST", suggestUrl+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+d.token)
	req.Header.Set("Content-Type", "application/json")

	response, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == 429 {
		return nil, errors.New("query limit reached, come back after midnight") //TODO move to const
	}

	if response.StatusCode != 200 {
		rBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("cannot get value from dadata api with error: %s", string(rBytes))
	}

	rBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res DaDataIpLocate
	err = json.Unmarshal(rBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
