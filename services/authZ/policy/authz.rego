package zrms.services.authz

default allow = false
default deny_reason = "default deny"
policy_version := "v1.0.0"

# --------------------------------------------------
# Input Validation
# --------------------------------------------------

input_valid {
    # Required fields
    input.resource != ""
    input.action != ""
    
    # Permissions must exist and be an object
    is_object(input.permissions)
}

# --------------------------------------------------
# Permission Evaluation
# --------------------------------------------------

permission_key = k {
    k := sprintf("%s:%s", [input.resource, input.action])
}

permission_exists {
    permission_key
    input.permissions[permission_key]
}

permission_allowed {
    permission_exists
    input.permissions[permission_key].allowed == true
}

# --------------------------------------------------
# Decision Logic
# --------------------------------------------------

allow {
    input_valid
    permission_allowed
}

# --------------------------------------------------
# Detailed Deny Reasons
# --------------------------------------------------

deny_reason = "invalid input: missing resource" {
    input.resource == ""
}

deny_reason = "invalid input: missing action" {
    input.action == ""
}

deny_reason = "invalid permissions structure" {
    not is_object(input.permissions)
}

deny_reason = "permission not found" {
    input_valid
    not permission_exists
}

deny_reason = "permission explicitly denied" {
    input_valid
    permission_exists
    not permission_allowed
}

# --------------------------------------------------
# Helper Functions
# --------------------------------------------------

is_object(x) {
    not is_string(x)
    not is_array(x)
    not is_set(x)
    not is_boolean(x)
    not is_null(x)
}

is_boolean(x) {
    x == true
} else = false {
    x == false
}