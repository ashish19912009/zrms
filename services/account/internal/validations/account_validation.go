package validations

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ashish19912009/zrms/services/account/internal/model"
)

var (
	allowedStatuses = []string{"active", "inactive", "suspended", "blocked", "limited"}
	allowedGenders  = []string{
		"male", "female",
		"transgender male", "transgender female",
		"non-binary", "genderqueer", "genderfluid",
		"agender", "bigender", "trigender", "pangender",
		"demiboy", "demigirl", "demigender",
		"two-spirit", "hijra", "fa'afafine",
		"neutrois", "maverique", "androgyne", "intergender",
		"questioning", "other", "prefer not to say"}
)

func SetAllowedStatuses(statuses []string) {
	allowedStatuses = statuses
}

func SetAllowedGenders(genders []string) {
	allowedGenders = genders
}

var (
	ErrMobileRequired           = errors.New("validation.error.mobile_required")
	ErrNameRequired             = errors.New("validation.error.name_required")
	ErrInvalidRole              = errors.New("validation.error.invalid_role")
	ErrInvalidStatus            = errors.New("validation.error.invalid_status")
	ErrAccountIDRequired        = errors.New("validation.error.account_id_required")
	ErrEmptyString              = errors.New("input cannot be empty")
	ErrLengthTooShort           = errors.New("input is too short")
	ErrLengthTooLong            = errors.New("input is too long")
	ErrRestrictedWordFound      = errors.New("input contains restricted word")
	ErrNotAlphanumeric          = errors.New("input must be alphanumeric")
	ErrDoesNotMatchPattern      = errors.New("input does not match required pattern")
	ErrMissingPrefix            = errors.New("input must start with required prefix")
	ErrMissingSuffix            = errors.New("input must end with required suffix")
	ErrNotNumeric               = errors.New("input must be numeric")
	ErrInvalidUTF8              = errors.New("input is not valid UTF-8")
	ErrMissingCharacter         = errors.New("input must contain required character")
	ErrInvalidOption            = errors.New("input is not an allowed option")
	ErrInvalidURL               = errors.New("input must be a valid URL")
	ErrNotLowercase             = errors.New("input must be all lowercase")
	ErrNotUppercase             = errors.New("input must be all uppercase")
	ErrPasswordEmpty            = errors.New("password is required")
	ErrPasswordTooShort         = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong          = errors.New("password must not exceed 64 characters")
	ErrPasswordMissingUpper     = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingLower     = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingDigit     = errors.New("password must contain at least one digit")
	ErrPasswordMissingSpecial   = errors.New("password must contain at least one special character")
	ErrPasswordContainsSpace    = errors.New("password must not contain spaces")
	ErrUUIDEmpty                = errors.New("UUID is required")
	ErrInvalidUUID              = errors.New("invalid UUID format")
	ErrEmailEmpty               = errors.New("email is required")
	ErrInvalidEmail             = errors.New("invalid email format")
	ErrTimeRequired             = errors.New("time value is required")
	ErrTimeInFuture             = errors.New("time must not be in the future")
	ErrTimeTooFarPast           = errors.New("time is too far in the past")
	ErrInvalidGender            = errors.New("Invalid gender selected")
	ErrDateFormatNotRecongnized = errors.New("Date format not recognized")
	ErrDateNotBeFuture          = errors.New("DOB cannot be in the future")
	ErrDateInNegative           = errors.New("DOB results in negative age")
	ErrDateUnrealistic          = errors.New("Age seems unrealistic (>150 years)")
	ErrAadhaarEmpty             = errors.New("aadhaar number is required")
	ErrAadhaarLength            = errors.New("aadhaar number must be exactly 12 digits")
	ErrAadhaarNotNumeric        = errors.New("aadhaar number must contain only digits")
	ErrAadhaarInvalidStart      = errors.New("aadhaar number cannot start with 0 or 1")
	ErrAadhaarInvalidChecksum   = errors.New("aadhaar number is invalid based on checksum (Verhoeff algorithm)")
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
	owner_id := TrimWhitespace(franchise.Franchise_Owner_id)

	// validate business name
	if err := ValidateNotEmpty(businessName); err != nil {
		return err
	}
	if err := ValidateLength(businessName, 5, 100); err != nil {
		return err
	}
	if err := ValidateNoRestrictedWords(businessName, []string{}); err != nil {
		return err
	}

	// validate LogoURL
	if err := ValidateNotEmpty(LogoURL); err != nil {
		return err
	}
	if err := ValidateLength(LogoURL, 5, 200); err != nil {
		return err
	}
	if err := ValidateURL(LogoURL); err != nil {
		return err
	}

	// validate domain
	if err := ValidateNotEmpty(SubDomain); err != nil {
		return err
	}
	if err := ValidateLength(SubDomain, 5, 50); err != nil {
		return err
	}
	// validate franchise_owner_id
	// validate domain
	if err := ValidateNotEmpty(owner_id); err != nil {
		return err
	}
	if err := ValidateUUID(owner_id); err != nil {
		return err
	}
	return nil
}

func ValidateFranchiseOwner(f_owner *model.FranchiseOwner) error {
	name := f_owner.Name
	gender := f_owner.Gender
	dob := f_owner.Dob
	mobile_no := f_owner.MobileNo
	email := f_owner.Email
	address := f_owner.Address
	aadhar_no := f_owner.AadharNo
	status := f_owner.Status

	// validate name
	if err := ValidateNotEmpty(name); err != nil {
		return err
	}
	if err := ValidateLength(name, 5, 100); err != nil {
		return err
	}
	if err := ValidateNoRestrictedWords(name, []string{}); err != nil {
		return err
	}
	// validate gender
	if err := ValidateNotEmpty(gender); err != nil {
		return err
	}
	if err := ValidateGenders(gender); err != nil {
		return err
	}
	// validate dob
	if err := ValidateNotEmpty(dob); err != nil {
		return err
	}
	if err := validateDOB(dob); err != nil {
		return err
	}
	// validate mobile no
	if err := ValidateNotEmpty(mobile_no); err != nil {
		return err
	}
	if err := ValidateMobileNo(dob); err != nil {
		return err
	}
	// validate email
	if err := ValidateNotEmpty(email); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	// validate address
	if err := ValidateNotEmpty(address); err != nil {
		return err
	}
	if err := ValidateLength(address, 5, 500); err != nil {
		return err
	}
	// validate aadhar
	if err := ValidateNotEmpty(aadhar_no); err != nil {
		return err
	}
	if err := ValidateAadhaarNumber(aadhar_no); err != nil {
		return err
	}
	// validate status
	if err := ValidateNotEmpty(status); err != nil {
		return err
	}
	if err := ValidateStatus(status); err != nil {
		return err
	}
	return nil
}

func ValidateFranchiseAccounts(acc *model.FranchiseAccount) error {
	franchise_id := acc.FranchiseID
	emp_id := acc.EmployeeID
	login_id := acc.LoginID
	password := acc.Password
	account_type := acc.AccountType
	name := acc.Name
	mobileNo := acc.MobileNo
	email := acc.Email
	roleID := acc.RoleID
	status := acc.Status
	created_at := acc.CreatedAt
	updated_at := acc.UpdatedAt
	deleted_at := acc.DeletedAt

	// validate franchise_id
	if err := ValidateNotEmpty(franchise_id); err != nil {
		return err
	}
	if err := ValidateUUID(franchise_id); err != nil {
		return err
	}
	// validate emp_id
	if err := ValidateNotEmpty(emp_id); err != nil {
		return err
	}
	// validate login
	if err := ValidateNotEmpty(login_id); err != nil {
		return err
	}
	if err := ValidateLength(login_id, 5, 50); err != nil {
		return err
	}
	// validate password
	if err := ValidateNotEmpty(password); err != nil {
		return err
	}
	if err := ValidatePassword(password); err != nil {
		return err
	}
	// validate account type
	if err := ValidateNotEmpty(account_type); err != nil {
		return err
	}
	// validate name
	if err := ValidateNotEmpty(name); err != nil {
		return err
	}
	if err := ValidateLength(name, 5, 100); err != nil {
		return err
	}
	// validate mobile no
	if err := ValidateNotEmpty(mobileNo); err != nil {
		return err
	}
	if err := ValidateMobileNo(mobileNo); err != nil {
		return err
	}
	// validate email
	if err := ValidateNotEmpty(email); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	// validate role
	if err := ValidateNotEmpty(roleID); err != nil {
		return err
	}
	// validate status
	if err := ValidateNotEmpty(status); err != nil {
		return err
	}
	if err := ValidateStatus(status); err != nil {
		return err
	}
	// validate created_at
	if !created_at.IsZero() {
		if err := ValidateTime(created_at); err != nil {
			return err
		}
	}
	// validate updated_at
	if !updated_at.IsZero() {
		if err := ValidateTime(updated_at); err != nil {
			return err
		}
	}
	// validate deleted_at
	if !deleted_at.IsZero() {
		if err := ValidateTime(deleted_at); err != nil {
			return err
		}
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

func ValidateStatus(status string) error {
	for _, s := range allowedStatuses {
		if s == status {
			return nil
		}
	}
	return errors.New(ErrInvalidStatus.Error())
}

func ValidateGenders(gender string) error {
	for _, s := range allowedGenders {
		if s == gender {
			return nil
		}
	}
	return errors.New(ErrInvalidGender.Error())
}

// validateDOB validates the given date of birth string against multiple patterns
func validateDOB(dobStr string) error {
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
		return errors.New(ErrDateFormatNotRecongnized.Error())
	}

	// Check if the date is in the future
	if dob.After(time.Now()) {
		return errors.New(ErrDateNotBeFuture.Error())
	}

	// Calculate the age
	age := time.Now().Year() - dob.Year()
	if time.Now().YearDay() < dob.YearDay() {
		age--
	}

	// Check for unrealistic age
	if age < 0 {
		return errors.New(ErrDateInNegative.Error())
	} else if age > 150 {
		return errors.New(ErrDateUnrealistic.Error())
	}
	return errors.New("something wrong in date")
}

// ValidateAadhaarNumber validates the format and structure of an Aadhaar number (India)
func ValidateAadhaarNumber(aadhaar string) error {
	// Trim any whitespace
	aadhaar = strings.TrimSpace(aadhaar)

	if aadhaar == "" {
		return ErrAadhaarEmpty
	}

	if len(aadhaar) != 12 {
		return ErrAadhaarLength
	}

	if matched, _ := regexp.MatchString(`^[0-9]{12}$`, aadhaar); !matched {
		return ErrAadhaarNotNumeric
	}

	if aadhaar[0] == '0' || aadhaar[0] == '1' {
		return ErrAadhaarInvalidStart
	}

	if !isValidAadhaarLuhn(aadhaar) {
		return ErrAadhaarInvalidChecksum
	}

	return nil
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

// ValidatePassword checks password strength based on common security rules
func ValidatePassword(password string) error {
	password = strings.TrimSpace(password)

	if password == "" {
		return ErrPasswordEmpty
	}

	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	if len(password) > 64 {
		return ErrPasswordTooLong
	}

	if strings.Contains(password, " ") {
		return ErrPasswordContainsSpace
	}

	if match, _ := regexp.MatchString(`[A-Z]`, password); !match {
		return ErrPasswordMissingUpper
	}

	if match, _ := regexp.MatchString(`[a-z]`, password); !match {
		return ErrPasswordMissingLower
	}

	if match, _ := regexp.MatchString(`[0-9]`, password); !match {
		return ErrPasswordMissingDigit
	}

	if match, _ := regexp.MatchString(`[\W_]`, password); !match {
		return ErrPasswordMissingSpecial
	}

	return nil
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
		return ErrInvalidUUID
	}

	return nil
}

// ValidateEmail checks if the input string is a valid email address
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return ErrEmailEmpty
	}

	// Basic email regex (simplified but practical)
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)

	if !re.MatchString(email) {
		return ErrInvalidEmail
	}

	return nil
}

// ValidateTime checks if the time is non-zero, not in the future, and not before a minimum threshold
func ValidateTime(t *time.Time) error {
	if t.IsZero() {
		return ErrTimeRequired
	}

	if t.After(time.Now()) {
		return ErrTimeInFuture
	}

	// Optional: Reject dates before 1900
	minAllowed := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	if t.Before(minAllowed) {
		return ErrTimeTooFarPast
	}

	return nil
}
