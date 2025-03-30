package validations_test

import (
	"testing"

	"github.com/ashish19912009/zrms/services/account/internal/model"
	"github.com/ashish19912009/zrms/services/account/internal/validations"
	"github.com/stretchr/testify/assert"
)

func TestValidateMobileNo(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expects error
	}{
		{"empty", "", validations.ErrMobileRequired},
		{"spaces only", "   ", validations.ErrMobileRequired},
		{"less digits", "12345", validations.ErrMobileRequired},
		{"alphabets", "abcdefghij", validations.ErrMobileRequired},
		{"valid", "9876543210", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validations.ValidateMobileNo(tt.input)
			assert.Equal(t, tt.expects, err)
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expects error
	}{
		{"empty", "", validations.ErrNameRequired},
		{"spaces only", "    ", validations.ErrNameRequired},
		{"valid", "John Doe", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validations.ValidateName(tt.input)
			assert.Equal(t, tt.expects, err)
		})
	}
}

func TestValidateRole(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expects error
	}{
		{"invalid", "ceo", validations.ErrInvalidRole},
		{"empty", "", validations.ErrInvalidRole},
		{"valid", "admin", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validations.ValidateRole(tt.input)
			assert.Equal(t, tt.expects, err)
		})
	}
}

func TestValidateStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		expects error
	}{
		{"invalid", "paused", validations.ErrInvalidStatus},
		{"empty", "", validations.ErrInvalidStatus},
		{"valid", "active", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validations.ValidateStatus(tt.input)
			assert.Equal(t, tt.expects, err)
		})
	}
}

func TestValidateAccount(t *testing.T) {
	invalidAcc := &model.Account{MobileNo: "", Name: "", Role: "ceo", Status: "paused"}
	err := validations.ValidateAccount(invalidAcc)
	assert.Error(t, err)

	validAcc := &model.Account{MobileNo: "9876543210", Name: "Alice", Role: "admin", Status: "active"}
	err = validations.ValidateAccount(validAcc)
	assert.NoError(t, err)
}

func TestValidateAccountUpdate(t *testing.T) {
	accMissingID := &model.Account{MobileNo: "9876543210", Name: "Bob", Role: "manager", Status: "inactive"}
	err := validations.ValidateAccountUpdate(accMissingID)
	assert.Equal(t, validations.ErrAccountIDRequired, err)

	valid := &model.Account{ID: "acc-001", MobileNo: "9876543210", Name: "Bob", Role: "manager", Status: "inactive"}
	err = validations.ValidateAccountUpdate(valid)
	assert.NoError(t, err)
}
