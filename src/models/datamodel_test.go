package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDataModel struct {
	DataModel
	Test string `json:"test"`
}

func TestDataModel_UnmarshalJson(t *testing.T) {
	data := []byte(`{
		"test": "we shall see if this works"
	}`)
	dataModel := &testDataModel{}
	err := dataModel.UnmarshalJson(data)
	assert.NoError(t, err)
	assert.Equal(t, &testDataModel{Test: "we shall see if this works"}, dataModel)
}
