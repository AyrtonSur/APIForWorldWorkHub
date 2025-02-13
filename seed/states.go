package seed

import (
	"example/APIForWorldWorkHub/models"
	"gorm.io/gorm"
)

var States = []models.Region{
	{Name: "Alabama", Abbreviation: "AL"},
	{Name: "Alaska", Abbreviation: "AK"},
	{Name: "Arizona", Abbreviation: "AZ"},
	{Name: "Arkansas", Abbreviation: "AR"},
	{Name: "California", Abbreviation: "CA"},
	{Name: "Colorado", Abbreviation: "CO"},
	{Name: "Connecticut", Abbreviation: "CT"},
	{Name: "Delaware", Abbreviation: "DE"},
	{Name: "Florida", Abbreviation: "FL"},
	{Name: "Georgia", Abbreviation: "GA"},
	{Name: "Hawaii", Abbreviation: "HI"},
	{Name: "Idaho", Abbreviation: "ID"},
	{Name: "Illinois", Abbreviation: "IL"},
	{Name: "Indiana", Abbreviation: "IN"},
	{Name: "Iowa", Abbreviation: "IA"},
	{Name: "Kansas", Abbreviation: "KS"},
	{Name: "Kentucky", Abbreviation: "KY"},
	{Name: "Louisiana", Abbreviation: "LA"},
	{Name: "Maine", Abbreviation: "ME"},
	{Name: "Maryland", Abbreviation: "MD"},
	{Name: "Massachusetts", Abbreviation: "MA"},
	{Name: "Michigan", Abbreviation: "MI"},
	{Name: "Minnesota", Abbreviation: "MN"},
	{Name: "Mississippi", Abbreviation: "MS"},
	{Name: "Missouri", Abbreviation: "MO"},
	{Name: "Montana", Abbreviation: "MT"},
	{Name: "Nebraska", Abbreviation: "NE"},
	{Name: "Nevada", Abbreviation: "NV"},
	{Name: "New Hampshire", Abbreviation: "NH"},
	{Name: "New Jersey", Abbreviation: "NJ"},
	{Name: "New Mexico", Abbreviation: "NM"},
	{Name: "New York", Abbreviation: "NY"},
	{Name: "North Carolina", Abbreviation: "NC"},
	{Name: "North Dakota", Abbreviation: "ND"},
	{Name: "Ohio", Abbreviation: "OH"},
	{Name: "Oklahoma", Abbreviation: "OK"},
	{Name: "Oregon", Abbreviation: "OR"},
	{Name: "Pennsylvania", Abbreviation: "PA"},
	{Name: "Rhode Island", Abbreviation: "RI"},
	{Name: "South Carolina", Abbreviation: "SC"},
	{Name: "South Dakota", Abbreviation: "SD"},
	{Name: "Tennessee", Abbreviation: "TN"},
	{Name: "Texas", Abbreviation: "TX"},
	{Name: "Utah", Abbreviation: "UT"},
	{Name: "Vermont", Abbreviation: "VT"},
	{Name: "Virginia", Abbreviation: "VA"},
	{Name: "Washington", Abbreviation: "WA"},
	{Name: "West Virginia", Abbreviation: "WV"},
	{Name: "Wisconsin", Abbreviation: "WI"},
	{Name: "Wyoming", Abbreviation: "WY"},
}

func InitializeStates(db *gorm.DB) {
	for _, state := range States {
		db.FirstOrCreate(&state, models.Region{Name: state.Name, Abbreviation: state.Abbreviation})
	}
}
