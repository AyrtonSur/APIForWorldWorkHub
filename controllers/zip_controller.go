package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

func CheckZip(context *gin.Context) {
	var input struct {
		Zipcode string `json:"zipcode" validate:"required,zipcode"`
	}

	if err := context.BindJSON(&input); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON"})
		return
	}

	url := fmt.Sprintf("http://api.zippopotam.us/us/%s", input.Zipcode)
	resp, err := http.Get(url)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error in http"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "ZipCode does not exist."})
		return
	}

	var zipCodeInfo ZipCodeInfo
	if err := json.NewDecoder(resp.Body).Decode(&zipCodeInfo); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error while decoding json."})
		return
	}

	context.IndentedJSON(http.StatusOK, zipCodeInfo)
}
