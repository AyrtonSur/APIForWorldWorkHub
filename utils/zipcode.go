package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ZipCodeInfo struct {
	PostCode string `json:"post code"`
	Country  string `json:"country"`
	Places   []struct {
		PlaceName string `json:"place name"`
		State     string `json:"state"`
	} `json:"places"`
}

func ValidateZipCode(zipCode string) (bool, error) {
	url := fmt.Sprintf("http://api.zippopotam.us/us/%s", zipCode)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var zipCodeInfo ZipCodeInfo
	if err := json.NewDecoder(resp.Body).Decode(&zipCodeInfo); err != nil {
		return false, err
	}

	return true, nil
}