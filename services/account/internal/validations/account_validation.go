package validations

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

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
	ErrMobileRequired      = errors.New("validation.error.mobile_required")
	ErrNameRequired        = errors.New("validation.error.name_required")
	ErrInvalidRole         = errors.New("validation.error.invalid_role")
	ErrInvalidStatus       = errors.New("validation.error.invalid_status")
	ErrAccountIDRequired   = errors.New("validation.error.account_id_required")
	ErrEmptyString         = errors.New("input cannot be empty")
	ErrLengthTooShort      = errors.New("input is too short")
	ErrLengthTooLong       = errors.New("input is too long")
	ErrRestrictedWordFound = errors.New("input contains restricted word")
	ErrNotAlphanumeric     = errors.New("input must be alphanumeric")
	ErrDoesNotMatchPattern = errors.New("input does not match required pattern")
	ErrMissingPrefix       = errors.New("input must start with required prefix")
	ErrMissingSuffix       = errors.New("input must end with required suffix")
	ErrNotNumeric          = errors.New("input must be numeric")
	ErrInvalidUTF8         = errors.New("input is not valid UTF-8")
	ErrMissingCharacter    = errors.New("input must contain required character")
	ErrInvalidOption       = errors.New("input is not an allowed option")
	ErrInvalidURL          = errors.New("input must be a valid URL")
	ErrNotLowercase        = errors.New("input must be all lowercase")
	ErrNotUppercase        = errors.New("input must be all uppercase")
)

// Length check (ValidateLength)

// Empty check (ValidateNotEmpty)

// Restricted words (ValidateNoRestrictedWords)

// Regex match (matchPattern)

// Alphanumeric check (isAlphanumeric)

// Prefix/suffix check (hasPrefix, hasSuffix)

// Numeric check (isNumeric)

// Trim spaces (trimWhitespace)

// UTF-8 validation (isValidUTF8)

// Contains specific character (containsCharacter)

// Check against list (isValidOption)

// URL validation (isValidURL)

// Lowercase check (isLowercase)

// Uppercase check (isUppercase)

// Remove non-alphanumeric characters (removeNonAlphanumeric)

func ValidateFranchise(franchise *model.Franchise) error {
	businessName := TrimWhitespace(franchise.BusinessName)
	LogoURL := TrimWhitespace(franchise.LogoURL)
	SubDomain := TrimWhitespace(franchise.SubDomain)

	// validate business name
	if err := ValidateName(businessName); err != nil {
		return err
	}
	if err := ValidateLength(businessName, 5, 100); err != nil {
		return err
	}
	if err := ValidateNotEmpty(businessName); err != nil {
		return err
	}
	if err := ValidateNoRestrictedWords(businessName, []string{}); err != nil {
		return err
	}

	// validate LogoURL
	if err := ValidateName(LogoURL); err != nil {
		return err
	}
	if err := ValidateLength(LogoURL, 5, 150); err != nil {
		return err
	}
	if err := ValidateNotEmpty(LogoURL); err != nil {
		return err
	}
	if err := ValidateURL(LogoURL); err != nil {
		return err
	}
	return nil

	// validate domain
	if err := ValidateName(SubDomain); err != nil {
		return err
	}
	if err := ValidateLength(businessName, 5, 100); err != nil {
		return err
	}
	if err := ValidateNotEmpty(businessName); err != nil {
		return err
	}
	if err := ValidateNoRestrictedWords(businessName, []string{}); err != nil {
		return err
	}
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

// validateDOB validates the given date of birth string against multiple patterns
func validateDOB(dobStr string) (bool, string) {
	// Define possible date formats
	dateFormats := []string{
		"2006-01-02", "02-01-2006", "01-02-2006", // YYYY-MM-DD, DD-MM-YYYY, MM-DD-YYYY
		"02/01/2006", "01/02/2006", "2006/01/02", // DD/MM/YYYY, MM/DD/YYYY, YYYY/MM/DD
		"02.01.2006", "01.02.2006", "2006.01.02", // DD.MM.YYYY, MM.DD.YYYY, YYYY.MM.DD
	}

	// Try to parse the date with each format
	var dob time.Time
	var err error
	for _, format := range dateFormats {
		dob, err = time.Parse(format, dobStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return false, "Date format not recognized"
	}

	// Check if the date is in the future
	if dob.After(time.Now()) {
		return false, "DOB cannot be in the future"
	}

	// Calculate the age
	age := time.Now().Year() - dob.Year()
	if time.Now().YearDay() < dob.YearDay() {
		age--
	}

	// Check for unrealistic age
	if age < 0 {
		return false, "DOB results in negative age"
	} else if age > 150 {
		return false, "Age seems unrealistic (>150 years)"
	}

	return true, fmt.Sprintf("Valid DOB. Age: %d years", age)
}

// ValidateAadhaarNumber validates the given Aadhaar number
func ValidateAadhaarNumber(aadhaar string) (bool, string) {
	// Check if the Aadhaar number is exactly 12 digits
	if len(aadhaar) != 12 {
		return false, "Aadhaar number must be exactly 12 digits."
	}

	// Check if the Aadhaar number contains only digits
	match, _ := regexp.MatchString("^[0-9]+$", aadhaar)
	if !match {
		return false, "Aadhaar number must contain only digits."
	}

	// Check if the first digit is 0 or 1
	firstDigit := aadhaar[0]
	if firstDigit == '0' || firstDigit == '1' {
		return false, "Aadhaar number cannot start with 0 or 1."
	}

	// Validate using the Luhn Algorithm
	if !isValidAadhaarLuhn(aadhaar) {
		return false, "Aadhaar number is invalid based on checksum (Luhn algorithm)."
	}

	return true, "Valid Aadhaar number."
}

// isValidAadhaarLuhn applies the Luhn algorithm to validate the Aadhaar number checksum
func isValidAadhaarLuhn(aadhaar string) bool {
	var sum int
	shouldDouble := false

	// Iterate over the digits from right to left
	for i := len(aadhaar) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(aadhaar[i]))
		if shouldDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		shouldDouble = !shouldDouble
	}

	// The sum should be divisible by 10 for the number to be valid
	return sum%10 == 0
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

// str := "Non-empty string"
// isValid := ValidateNotEmpty(str) // True if string is not empty

// ValidateNoRestrictedWords checks if the input contains any restricted words
func ValidateNoRestrictedWords(str string, restrictedWords []string) error {
	commonRestrictedWord := []string{
		"admin", "root", "system", "support", "moderator",
		"password", "token", "key", "credentials", "secret",
		"free", "click here", "subscribe", "buy now", "discount",
		"idiot", "stupid", "dumb", "hate", "ugly", "null", "undefined",
		"fuck",
		"f***", "s***", "b****", "a******", "idiot", "moron", "jerk", "stupid", // Profanity
		"sex", "porn", "nude", "rape", "masturbate", // Sexual terms
		"kill", "murder", "shoot", "bomb", "attack", // Violent terms
		"terrorist", "k***", "n****", "racist", "misogynistic", // Hate speech
	}
	if len(restrictedWords) == 0 {
		restrictedWords = commonRestrictedWord
	} else {
		restrictedWords = append(restrictedWords, commonRestrictedWord...)
	}
	for _, word := range restrictedWords {
		if strings.Contains(strings.ToLower(str), strings.ToLower(word)) {
			return ErrRestrictedWordFound
		}
	}
	return nil
}

// restrictedWords := []string{"admin", "password", "root"}
// str := "User password is weak"
// containsRestricted := ValidateNoRestrictedWords(str, restrictedWords) // True if restricted word is found

// ValidatePattern checks if the input matches the given regex pattern
func ValidatePattern(str string, pattern string) error {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(str) {
		return ErrDoesNotMatchPattern
	}
	return nil
}

// emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
// email := "test@example.com"
// isValid := matchPattern(email, emailPattern) // True if email is valid

// ValidateAlphanumeric checks if the input contains only letters and numbers
func ValidateAlphanumeric(str string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !re.MatchString(str) {
		return ErrNotAlphanumeric
	}
	return nil
}

// str := "Hello123"
// isValid := isAlphanumeric(str) // True if the string is alphanumeric

// ValidateHasPrefix checks if the input starts with the given prefix
func ValidateHasPrefix(str, prefix string) error {
	if !strings.HasPrefix(str, prefix) {
		return ErrMissingPrefix
	}
	return nil
}

// str := "http://example.com"
// isValidPrefix := hasPrefix(str, "http") // True if it starts with 'http'

// ValidateHasSuffix checks if the input ends with the given suffix
func ValidateHasSuffix(str, suffix string) error {
	if !strings.HasSuffix(str, suffix) {
		return ErrMissingSuffix
	}
	return nil
}

// str := "http://example.com"
// isValidSuffix := hasSuffix(str, ".com") // True if it ends with '.com'

// ValidateNumeric checks if the input is a valid numeric string
func ValidateNumeric(str string) error {
	if _, err := strconv.Atoi(str); err != nil {
		return ErrNotNumeric
	}
	return nil
}

// str := "12345"
// isValid := isNumeric(str) // True if string contains only digits

// TrimWhitespace removes leading and trailing spaces (utility, not validator)
func TrimWhitespace(str string) string {
	return strings.TrimSpace(str)
}

// str := "   Hello World!   "
// trimmed := trimWhitespace(str) // "Hello World!"

// ValidateUTF8 checks if the string is valid UTF-8
func ValidateUTF8(str string) error {
	if !utf8.ValidString(str) {
		return ErrInvalidUTF8
	}
	return nil
}

// str := "Hello World"
// isValid := isValidUTF8(str) // True if the string is valid UTF-8

// ValidateContainsCharacter checks if the input contains a specific character
func ValidateContainsCharacter(str, char string) error {
	if !strings.Contains(str, char) {
		return ErrMissingCharacter
	}
	return nil
}

// str := "Hello, World!"
// isValid := containsCharacter(str, ",") // True if it contains a comma

// ValidateInOptions checks if the input is in a predefined list of valid options
func ValidateInOptions(str string, validOptions []string) error {
	for _, option := range validOptions {
		if str == option {
			return nil
		}
	}
	return ErrInvalidOption
}

// validOptions := []string{"admin", "user", "guest"}
// str := "user"
// isValid := isValidOption(str, validOptions) // True if 'str' is in the list of valid options

// ValidateURL checks if the input is a valid URL
func ValidateURL(str string) error {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return ErrInvalidURL
	}
	return nil
}

// url := "https://www.example.com"
// isValid := isValidURL(url) // True if the string is a valid URL

// ValidateLowercase checks if the input is all lowercase
func ValidateLowercase(str string) error {
	if str != strings.ToLower(str) {
		return ErrNotLowercase
	}
	return nil
}

// str := "hello"
// isValid := isLowercase(str) // True if the string is lowercase

// ValidateUppercase checks if the input is all uppercase
func ValidateUppercase(str string) error {
	if str != strings.ToUpper(str) {
		return ErrNotUppercase
	}
	return nil
}

// str := "HELLO"
// isValid := isUppercase(str) // True if the string is uppercase

// RemoveNonAlphanumeric removes all non-alphanumeric characters from the input (utility)
func RemoveNonAlphanumeric(str string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return re.ReplaceAllString(str, "")
}

// str := "Hello, World!123"
// cleaned := removeNonAlphanumeric(str) // "HelloWorld123"
