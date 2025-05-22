package validations_test

// import (
// 	"errors"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/ashish19912009/zrms/services/account/internal/model"
// 	"github.com/ashish19912009/zrms/services/account/internal/validations"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	ErrMobileRequired           = errors.New("validation.error.mobile_required")
// 	ErrNameRequired             = errors.New("validation.error.name_required")
// 	ErrInvalidRole              = errors.New("validation.error.invalid_role")
// 	ErrInvalidStatus            = errors.New("validation.error.invalid_status")
// 	ErrAccountIDRequired        = errors.New("validation.error.account_id_required")
// 	ErrEmptyString              = errors.New("input cannot be empty")
// 	ErrLengthTooShort           = errors.New("input is too short")
// 	ErrLengthTooLong            = errors.New("input is too long")
// 	ErrRestrictedWordFound      = errors.New("input contains restricted word")
// 	ErrNotAlphanumeric          = errors.New("input must be alphanumeric")
// 	ErrDoesNotMatchPattern      = errors.New("input does not match required pattern")
// 	ErrMissingPrefix            = errors.New("input must start with required prefix")
// 	ErrMissingSuffix            = errors.New("input must end with required suffix")
// 	ErrNotNumeric               = errors.New("input must be numeric")
// 	ErrInvalidUTF8              = errors.New("input is not valid UTF-8")
// 	ErrMissingCharacter         = errors.New("input must contain required character")
// 	ErrInvalidOption            = errors.New("input is not an allowed option")
// 	ErrInvalidURL               = errors.New("input must be a valid URL")
// 	ErrNotLowercase             = errors.New("input must be all lowercase")
// 	ErrNotUppercase             = errors.New("input must be all uppercase")
// 	ErrPasswordEmpty            = errors.New("password is required")
// 	ErrPasswordTooShort         = errors.New("password must be at least 8 characters long")
// 	ErrPasswordTooLong          = errors.New("password must not exceed 64 characters")
// 	ErrPasswordMissingUpper     = errors.New("password must contain at least one uppercase letter")
// 	ErrPasswordMissingLower     = errors.New("password must contain at least one lowercase letter")
// 	ErrPasswordMissingDigit     = errors.New("password must contain at least one digit")
// 	ErrPasswordMissingSpecial   = errors.New("password must contain at least one special character")
// 	ErrPasswordContainsSpace    = errors.New("password must not contain spaces")
// 	ErrUUIDEmpty                = errors.New("UUID is required")
// 	ErrInvalidUUID              = errors.New("invalid UUID format")
// 	ErrEmailEmpty               = errors.New("email is required")
// 	ErrInvalidEmail             = errors.New("invalid email format")
// 	ErrTimeRequired             = errors.New("time value is required")
// 	ErrTimeInFuture             = errors.New("time must not be in the future")
// 	ErrTimeTooFarPast           = errors.New("time is too far in the past")
// 	ErrInvalidGender            = errors.New("invalid gender selected")
// 	ErrDateFormatNotRecongnized = errors.New("date format not recognized")
// 	ErrDateNotBeFuture          = errors.New("dob cannot be in the future")
// 	ErrDateInNegative           = errors.New("dob results in negative age")
// 	ErrDateUnrealistic          = errors.New("age seems unrealistic (>150 years)")
// 	ErrAadhaarEmpty             = errors.New("aadhaar number is required")
// 	ErrAadhaarLength            = errors.New("aadhaar number must be exactly 12 digits")
// 	ErrAadhaarNotNumeric        = errors.New("aadhaar number must contain only digits")
// 	ErrAadhaarInvalidStart      = errors.New("aadhaar number cannot start with 0 or 1")
// 	ErrAadhaarInvalidChecksum   = errors.New("aadhaar number is invalid based on checksum (Verhoeff algorithm)")
// 	ErrInvalidPincode           = errors.New("invalid pincode: must be 6 digits")
// 	ErrInvalidLatitude          = errors.New("invalid latitude: must be between -90 and 90")
// 	ErrInvalidLongitude         = errors.New("invalid longitude: must be between -180 and 180")
// )

// func ptrTime(t time.Time) *time.Time {
// 	return &t
// }

// func TestValidateFranchise_Success(t *testing.T) {
// 	franchise := &model.Franchise{
// 		BusinessName:       "Zippy Foods",
// 		LogoURL:            "https://example.com/logo.png",
// 		SubDomain:          "zippyfoods",
// 		Franchise_Owner_id: uuid.New().String(),
// 	}

// 	if err := validations.ValidateFranchise(franchise); err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// func TestValidateFranchise_Failure(t *testing.T) {
// 	franchise := &model.Franchise{
// 		BusinessName:       "",
// 		LogoURL:            "",
// 		SubDomain:          "",
// 		Franchise_Owner_id: "invalid-uuid",
// 	}

// 	err := validations.ValidateFranchise(franchise)
// 	if err == nil {
// 		t.Errorf("expected validation error, got nil")
// 	}
// }

// func TestValidateFranchiseOwner_Success(t *testing.T) {
// 	owner := &model.FranchiseOwner{
// 		Name:     "Ashish Kumar",
// 		Gender:   "male",
// 		Dob:      "1990-01-01",
// 		MobileNo: "9876543210",
// 		Email:    "ashish@example.com",
// 		Address:  "123, Main Street, City",
// 		AadharNo: "234567891234",
// 		Status:   "active",
// 	}

// 	if err := validations.ValidateFranchiseOwner(owner); err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// }

// func TestValidateFranchiseOwner_Failure(t *testing.T) {
// 	owner := &model.FranchiseOwner{
// 		Name:     "",
// 		Gender:   "unknown",
// 		Dob:      time.Now().AddDate(1, 0, 0).Format("2006-01-02"), // Future date
// 		MobileNo: "",
// 		Email:    "invalid-email",
// 		Address:  "",
// 		AadharNo: "123",
// 		Status:   "ghost",
// 	}

// 	err := validations.ValidateFranchiseOwner(owner)
// 	if err == nil {
// 		t.Errorf("expected validation error, got nil")
// 	}
// }

// func TestValidateFranchise_ValidInput(t *testing.T) {
// 	franchise := &model.Franchise{
// 		BusinessName:       "Zippy Eats",
// 		LogoURL:            "https://example.com/logo.png",
// 		SubDomain:          "zippyeats",
// 		Franchise_Owner_id: uuid.New().String(),
// 	}

// 	err := validations.ValidateFranchise(franchise)
// 	assert.NoError(t, err)
// }

// func TestValidateFranchise_InvalidInputs(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		franchise *model.Franchise
// 	}{
// 		{"Empty BusinessName", &model.Franchise{BusinessName: "", LogoURL: "https://valid.url", SubDomain: "domain", Franchise_Owner_id: uuid.New().String()}},
// 		{"Short BusinessName", &model.Franchise{BusinessName: "Zip", LogoURL: "https://valid.url", SubDomain: "domain", Franchise_Owner_id: uuid.New().String()}},
// 		{"Empty LogoURL", &model.Franchise{BusinessName: "Valid Name", LogoURL: "", SubDomain: "domain", Franchise_Owner_id: uuid.New().String()}},
// 		{"Invalid LogoURL", &model.Franchise{BusinessName: "Valid Name", LogoURL: "not-a-url", SubDomain: "domain", Franchise_Owner_id: uuid.New().String()}},
// 		{"Empty SubDomain", &model.Franchise{BusinessName: "Valid Name", LogoURL: "https://valid.url", SubDomain: "", Franchise_Owner_id: uuid.New().String()}},
// 		{"Short SubDomain", &model.Franchise{BusinessName: "Valid Name", LogoURL: "https://valid.url", SubDomain: "abc", Franchise_Owner_id: uuid.New().String()}},
// 		{"Empty Owner ID", &model.Franchise{BusinessName: "Valid Name", LogoURL: "https://valid.url", SubDomain: "domain", Franchise_Owner_id: ""}},
// 		{"Invalid UUID", &model.Franchise{BusinessName: "Valid Name", LogoURL: "https://valid.url", SubDomain: "domain", Franchise_Owner_id: "invalid-uuid"}},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := validations.ValidateFranchise(tt.franchise)
// 			assert.Error(t, err)
// 		})
// 	}
// }

// func TestValidateFranchiseOwner_ValidInput(t *testing.T) {
// 	owner := &model.FranchiseOwner{
// 		Name:     "Ashish Kumar",
// 		Gender:   "male",
// 		Dob:      "1990-01-01",
// 		MobileNo: "9876543210",
// 		Email:    "ashish@example.com",
// 		Address:  "123 Baker Street",
// 		AadharNo: "234567891234",
// 		Status:   "active",
// 	}

// 	err := validations.ValidateFranchiseOwner(owner)
// 	assert.NoError(t, err)
// }

// func TestValidateFranchiseOwner_InvalidInputs(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		owner *model.FranchiseOwner
// 	}{
// 		{"Empty Name", &model.FranchiseOwner{Name: "", Gender: "male", Dob: "1990-01-01", MobileNo: "9876543210", Email: "ashish@example.com", Address: "123", AadharNo: "234567891234", Status: "active"}},
// 		{"Short Name", &model.FranchiseOwner{Name: "Ash", Gender: "male", Dob: "1990-01-01", MobileNo: "9876543210", Email: "ashish@example.com", Address: "123 Baker", AadharNo: "234567891234", Status: "active"}},
// 		{"Invalid Gender", &model.FranchiseOwner{Name: "Ashish", Gender: "invalid", Dob: "1990-01-01", MobileNo: "9876543210", Email: "ashish@example.com", Address: "123 Baker", AadharNo: "234567891234", Status: "active"}},
// 		{"Future DOB", &model.FranchiseOwner{Name: "Ashish", Gender: "male", Dob: time.Now().AddDate(1, 0, 0).Format("2006-01-02"), MobileNo: "9876543210", Email: "ashish@example.com", Address: "123 Baker", AadharNo: "234567891234", Status: "active"}},
// 		{"Invalid MobileNo", &model.FranchiseOwner{Name: "Ashish", Gender: "male", Dob: "1990-01-01", MobileNo: "invalid", Email: "ashish@example.com", Address: "123 Baker", AadharNo: "234567891234", Status: "active"}},
// 		{"Invalid Email", &model.FranchiseOwner{Name: "Ashish", Gender: "male", Dob: "1990-01-01", MobileNo: "9876543210", Email: "invalid-email", Address: "123 Baker", AadharNo: "234567891234", Status: "active"}},
// 		{"Short Address", &model.FranchiseOwner{Name: "Ashish", Gender: "male", Dob: "1990-01-01", MobileNo: "9876543210", Email: "ashish@example.com", Address: "12", AadharNo: "234567891234", Status: "active"}},
// 		{"Invalid Aadhaar", &model.FranchiseOwner{Name: "Ashish", Gender: "male", Dob: "1990-01-01", MobileNo: "9876543210", Email: "ashish@example.com", Address: "123 Baker", AadharNo: "123", Status: "active"}},
// 		{"Invalid Status", &model.FranchiseOwner{Name: "Ashish", Gender: "male", Dob: "1990-01-01", MobileNo: "9876543210", Email: "ashish@example.com", Address: "123 Baker", AadharNo: "234567891234", Status: "invalid"}},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := validations.ValidateFranchiseOwner(tt.owner)
// 			assert.Error(t, err)
// 		})
// 	}
// }

// func TestValidateFranchiseAccounts_Valid(t *testing.T) {
// 	acc := &model.FranchiseAccount{
// 		FranchiseID: uuid.NewString(),
// 		EmployeeID:  "EMP123",
// 		LoginID:     "login123",
// 		Password:    "StrongP@ssw0rd",
// 		AccountType: "packer",
// 		Name:        "John Doe",
// 		MobileNo:    "9876543210",
// 		Email:       "john@example.com",
// 		RoleID:      "role-abc",
// 		Status:      "active",
// 		CreatedAt:   ptrTime(time.Now()),
// 		UpdatedAt:   ptrTime(time.Now()),
// 		DeletedAt:   nil,
// 	}
// 	err := validations.ValidateFranchiseAccounts(acc)
// 	assert.NoError(t, err)
// }

// func TestValidateFranchiseAccounts_InvalidMobile(t *testing.T) {
// 	now := time.Now()
// 	acc := &model.FranchiseAccount{
// 		FranchiseID: uuid.NewString(),
// 		EmployeeID:  "EMP123",
// 		LoginID:     "login123",
// 		Password:    "StrongP@ssw0rd",
// 		AccountType: "packer",
// 		Name:        "John Doe",
// 		MobileNo:    "123", // Invalid mobile
// 		Email:       "john@example.com",
// 		RoleID:      "role-abc",
// 		Status:      "active",
// 		CreatedAt:   &now,
// 		UpdatedAt:   &now,
// 	}
// 	err := validations.ValidateFranchiseAccounts(acc)
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "mobile")
// }

// func TestValidateFranchiseDocument_Valid(t *testing.T) {
// 	doc := &model.FranchiseDocument{
// 		FranchiseID:    uuid.NewString(),
// 		DocumentTypeID: uuid.NewString(),
// 		DocumentURL:    "https://example.com/doc.pdf",
// 		UploadedBy:     "uploader-123",
// 		Status:         "pending",
// 		Remark:         "Valid document",
// 	}
// 	err := validations.ValidateFranchiseDocument(doc)
// 	assert.NoError(t, err)
// }

// func TestValidateFranchiseDocument_InvalidURL(t *testing.T) {
// 	doc := &model.FranchiseDocument{
// 		FranchiseID:    uuid.NewString(),
// 		DocumentTypeID: uuid.NewString(),
// 		DocumentURL:    "not_a_url",
// 		UploadedBy:     "uploader-123",
// 		Status:         "pending",
// 		Remark:         "Valid document",
// 	}
// 	err := validations.ValidateFranchiseDocument(doc)
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "url")
// }

// func TestValidateFranchiseDocument_RestrictedRemark(t *testing.T) {
// 	doc := &model.FranchiseDocument{
// 		FranchiseID:    uuid.NewString(),
// 		DocumentTypeID: uuid.NewString(),
// 		DocumentURL:    "https://example.com/doc.pdf",
// 		UploadedBy:     "uploader-123",
// 		Status:         "pending",
// 		Remark:         "This is admin only", // restricted word: admin
// 	}
// 	err := validations.ValidateFranchiseDocument(doc)
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "restricted")
// }

// func TestValidateUTF8(t *testing.T) {
// 	err := validations.ValidateUTF8("Hello World")
// 	if err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}

// 	invalid := string([]byte{0xff, 0xfe, 0xfd})
// 	err = validations.ValidateUTF8(invalid)
// 	if !errors.Is(err, ErrInvalidUTF8) {
// 		t.Errorf("expected ErrInvalidUTF8, got %v", err)
// 	}
// }

// func TestValidateContainsCharacter(t *testing.T) {
// 	if err := validations.ValidateContainsCharacter("Hello, World!", ","); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateContainsCharacter("Hello World", ","); !errors.Is(err, ErrMissingCharacter) {
// 		t.Errorf("expected ErrMissingCharacter, got %v", err)
// 	}
// }

// func TestValidateInOptions(t *testing.T) {
// 	options := []string{"admin", "user", "guest"}
// 	if err := validations.ValidateInOptions("user", options); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateInOptions("unknown", options); !errors.Is(err, ErrInvalidOption) {
// 		t.Errorf("expected ErrInvalidOption, got %v", err)
// 	}
// }

// func TestValidateURL(t *testing.T) {
// 	if err := validations.ValidateURL("https://www.example.com"); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateURL("invalid-url"); !errors.Is(err, ErrInvalidURL) {
// 		t.Errorf("expected ErrInvalidURL, got %v", err)
// 	}
// }

// func TestValidateLowercase(t *testing.T) {
// 	if err := validations.ValidateLowercase("hello"); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateLowercase("Hello"); !errors.Is(err, ErrNotLowercase) {
// 		t.Errorf("expected ErrNotLowercase, got %v", err)
// 	}
// }

// func TestValidateUppercase(t *testing.T) {
// 	if err := validations.ValidateUppercase("HELLO"); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateUppercase("Hello"); !errors.Is(err, ErrNotUppercase) {
// 		t.Errorf("expected ErrNotUppercase, got %v", err)
// 	}
// }

// func TestRemoveNonAlphanumeric(t *testing.T) {
// 	input := "Hello, World!123"
// 	expected := "HelloWorld123"
// 	result := validations.RemoveNonAlphanumeric(input)
// 	if result != expected {
// 		t.Errorf("expected %s, got %s", expected, result)
// 	}
// }

// func TestValidatePassword(t *testing.T) {
// 	valid := "Aa1@validpassword"
// 	if err := validations.ValidatePassword(valid); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}

// 	tests := map[string]string{
// 		"":                      ErrPasswordEmpty.Error(),
// 		"short":                 ErrPasswordTooShort.Error(),
// 		strings.Repeat("a", 65): ErrPasswordTooLong.Error(),
// 		"no upper1@":            ErrPasswordMissingUpper.Error(),
// 		"NOLOWER1@":             ErrPasswordMissingLower.Error(),
// 		"NoNumber@":             ErrPasswordMissingDigit.Error(),
// 		"NoSpecial1":            ErrPasswordMissingSpecial.Error(),
// 		"Contains space1@A":     ErrPasswordContainsSpace.Error(),
// 	}
// 	for input, expected := range tests {
// 		err := validations.ValidatePassword(input)
// 		if err == nil || err.Error() != expected {
// 			t.Errorf("expected %v, got %v for input %q", expected, err, input)
// 		}
// 	}
// }

// func TestValidateUUID(t *testing.T) {
// 	valid := "550e8400-e29b-41d4-a716-446655440000"
// 	if err := validations.ValidateUUID(valid); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	invalid := "not-a-uuid"
// 	if err := validations.ValidateUUID(invalid); !errors.Is(err, ErrInvalidUUID) {
// 		t.Errorf("expected ErrInvalidUUID, got %v", err)
// 	}
// }

// func TestValidateEmail(t *testing.T) {
// 	if err := validations.ValidateEmail("user@example.com"); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateEmail("invalid-email"); !errors.Is(err, ErrInvalidEmail) {
// 		t.Errorf("expected ErrInvalidEmail, got %v", err)
// 	}
// }

// func TestValidateTime(t *testing.T) {
// 	now := time.Now()
// 	past := time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC)
// 	future := now.Add(24 * time.Hour)

// 	if err := validations.ValidateTime(&now); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateTime(&time.Time{}); !errors.Is(err, ErrTimeRequired) {
// 		t.Errorf("expected ErrTimeRequired, got %v", err)
// 	}
// 	if err := validations.ValidateTime(&future); !errors.Is(err, ErrTimeInFuture) {
// 		t.Errorf("expected ErrTimeInFuture, got %v", err)
// 	}
// 	if err := validations.ValidateTime(&past); !errors.Is(err, ErrTimeTooFarPast) {
// 		t.Errorf("expected ErrTimeTooFarPast, got %v", err)
// 	}
// }

// func TestValidatePincode(t *testing.T) {
// 	if err := validations.ValidatePincode("560001"); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidatePincode("12345"); !errors.Is(err, ErrInvalidPincode) {
// 		t.Errorf("expected ErrInvalidPincode, got %v", err)
// 	}
// }

// func TestValidateLatitude(t *testing.T) {
// 	if err := validations.ValidateLatitude(45); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateLatitude(100); !errors.Is(err, ErrInvalidLatitude) {
// 		t.Errorf("expected ErrInvalidLatitude, got %v", err)
// 	}
// }

// func TestValidateLongitude(t *testing.T) {
// 	if err := validations.ValidateLongitude(80); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateLongitude(-200); !errors.Is(err, ErrInvalidLongitude) {
// 		t.Errorf("expected ErrInvalidLongitude, got %v", err)
// 	}
// }

// func TestValidateCoordinates(t *testing.T) {
// 	if err := validations.ValidateCoordinates(45, 80); err != nil {
// 		t.Errorf("expected nil, got %v", err)
// 	}
// 	if err := validations.ValidateCoordinates(100, 80); err == nil || !strings.Contains(err.Error(), "latitude") {
// 		t.Errorf("expected latitude error, got %v", err)
// 	}
// 	if err := validations.ValidateCoordinates(45, 200); err == nil || !strings.Contains(err.Error(), "longitude") {
// 		t.Errorf("expected longitude error, got %v", err)
// 	}
// }
