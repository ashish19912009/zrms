package authz

default allow = false

allow {
    input.action == "read"
    startswith(input.resource, "/public/")
}

allow {
    some i
    policy := input.policies[i]
    policy.action == input.action
    glob.match(policy.resource, ["*"], input.resource)
}

# Helper function for glob-like matching
glob.match(pattern, delimiters, input) {
    pattern == input
}

glob.match(pattern, delimiters, input) {
    contains(pattern, "*")
    startswith(input, trim_suffix(pattern, "*"))
}