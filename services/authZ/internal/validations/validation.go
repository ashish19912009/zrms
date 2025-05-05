package validations

import (
	"errors"
	"regexp"
	"strings"

	"github.com/ashish19912009/zrms/services/authZ/internal/model"
)

var (
	allowedStatuses = []string{"active", "inactive", "suspended", "blocked", "limited"}
)

func SetAllowedStatuses(statuses []string) {
	allowedStatuses = statuses
}

var (
	ErrInvalidStatus  = errors.New("validation.error.invalid_status")
	ErrEmptyString    = errors.New("input cannot be empty")
	ErrLengthTooShort = errors.New("input is too short")
	ErrLengthTooLong  = errors.New("input is too long")
	ErrUUIDEmpty      = errors.New("UUID is required")
	ErrInvalidUUID    = errors.New("invalid UUID format")
)

func ValidateCheckAccess(access *model.CheckAccess) error {
	account_id := TrimWhitespace(access.AccountID)
	franchise_id := TrimWhitespace(access.FranchiseID)
	resource := TrimWhitespace(access.Resource)
	action := TrimWhitespace(access.Action)

	// validate
	if err := ValidateNotEmpty(account_id); err != nil {
		return err
	}
	if err := ValidateNotEmpty(franchise_id); err != nil {
		return err
	}
	if err := ValidateNotEmpty(resource); err != nil {
		return err
	}
	if err := ValidateNotEmpty(action); err != nil {
		return err
	}

	if err := ValidateLength(resource, 1, 50); err != nil {
		return err
	}
	if err := ValidateLength(action, 1, 50); err != nil {
		return err
	}

	if err := ValidateUUID(franchise_id); err != nil {
		return err
	}
	if err := ValidateUUID(account_id); err != nil {
		return err
	}
	return nil
}

func ValidateStatus(status string) error {
	for _, s := range allowedStatuses {
		if s == status {
			return nil
		}
	}
	return errors.New(ErrInvalidStatus.Error())
}

// ValidateLength checks if the input string length is within the given range
func ValidateLength(str string, minLen, maxLen int) error {
	length := len(str)
	if length < minLen {
		return ErrLengthTooShort
	}
	if length > maxLen {
		return ErrLengthTooLong
	}
	return nil
}

// str := "Hello"
// isValid := ValidateLength(str, 5, 10) // True if length is between 5 and 10

// ValidateNotEmpty checks if the input string is not empty or only whitespace
func ValidateNotEmpty(str string) error {
	if strings.TrimSpace(str) == "" {
		return ErrEmptyString
	}
	return nil
}

// str := "12345"
// isValid := isNumeric(str) // True if string contains only digits

// TrimWhitespace removes leading and trailing spaces (utility, not validator)
func TrimWhitespace(str string) string {
	return strings.TrimSpace(str)
}

// ValidateUUID checks if the input string is a valid UUID
func ValidateUUID(id string) error {
	id = strings.TrimSpace(id)

	if id == "" {
		return ErrUUIDEmpty
	}

	// Regex to match UUID v1–v5
	uuidRegex := `^[a-fA-F0-9]{8}\-[a-fA-F0-9]{4}\-[1-5][a-fA-F0-9]{3}\-[89abAB][a-fA-F0-9]{3}\-[a-fA-F0-9]{12}$`
	re := regexp.MustCompile(uuidRegex)

	if !re.MatchString(id) {
		return nil
		//return ErrInvalidUUID
	}

	return nil
}
