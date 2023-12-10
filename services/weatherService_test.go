package services

import (
	"Backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testValidateWeather(t *testing.T) {
	tests := []struct {
		name   string
		runner *models.Job
		want   *models.ResponseError
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			responseErr := ValidateJob(test.runner)
			assert.Equal(t, test.want, responseErr)
		})
	}
}
