package validations

import (
	"errors"
	"regexp"
	"strings"

	"github.com/ashish19912009/zrms/services/account/internal/model"
)

var (
	allowedRoles    = []string{"admin", "manager", "delivery"}
	allowedStatuses = []string{"active", "inactive", "suspended"}
)

func SetAllowedRoles(roles []string) {
	allowedRoles = roles
}

func SetAllowedStatuses(statuses []string) {
	allowedStatuses = statuses
}

var (
	ErrMobileRequired    = errors.New("validation.error.mobile_required")
	ErrNameRequired      = errors.New("validation.error.name_required")
	ErrInvalidRole       = errors.New("validation.error.invalid_role")
	ErrInvalidStatus     = errors.New("validation.error.invalid_status")
	ErrAccountIDRequired = errors.New("validation.error.account_id_required")
)

func ValidateAccount(acc *model.Account) error {
	if err := ValidateMobileNo(acc.MobileNo); err != nil {
		return err
	}
	if err := ValidateName(acc.Name); err != nil {
		return err
	}
	if err := ValidateRole(acc.Role); err != nil {
		return err
	}
	if err := ValidateStatus(acc.Status); err != nil {
		return err
	}
	return nil
}

func ValidateAccountUpdate(acc *model.Account) error {
	if acc.ID == "" {
		return errors.New(ErrAccountIDRequired.Error())
	}
	return nil
}

func ValidateMobileNo(mobile string) error {
	if strings.TrimSpace(mobile) == "" {
		return errors.New(ErrMobileRequired.Error())
	}
	pattern := `^\d{10}$`
	if matched, _ := regexp.MatchString(pattern, mobile); !matched {
		return errors.New(ErrMobileRequired.Error())
	}
	return nil
}

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New(ErrNameRequired.Error())
	}
	return nil
}

func ValidateRole(role string) error {
	for _, r := range allowedRoles {
		if r == role {
			return nil
		}
	}
	return errors.New(ErrInvalidRole.Error())
}

func ValidateStatus(status string) error {
	for _, s := range allowedStatuses {
		if s == status {
			return nil
		}
	}
	return errors.New(ErrInvalidStatus.Error())
}
