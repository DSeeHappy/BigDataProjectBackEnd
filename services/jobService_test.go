package services

import (
	"Backend/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestJobService_ValidateJob(t *testing.T) {
	tests := []struct {
		name string
		job  *models.Job
		want *models.ResponseError
	}{
		{
			name: "Invalid_Name",
			job: &models.Job{
				Name: "",
			},
			want: &models.ResponseError{
				Message: "Invalid name: ",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_State",
			job: &models.Job{
				Name:  "John",
				State: "",
			},
			want: &models.ResponseError{
				Message: "Invalid state: ",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_Zip_Code",
			job: &models.Job{
				Name:    "John",
				State:   "CA",
				ZipCode: "",
			},
			want: &models.ResponseError{
				Message: "Invalid Zip Code: ",
				Status:  http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseErr := ValidateJob(tt.job)
			assert.NotEmpty(t, responseErr)
			assert.Equal(t, tt.want.Message, responseErr.Message)
			assert.Equal(t, tt.want.Status, responseErr.Status)
		})
	}
}
