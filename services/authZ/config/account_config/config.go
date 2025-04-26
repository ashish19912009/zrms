package config

import (
	"os"
	"strings"
)

var (
	AllowedRoles    []string
	AllowedStatuses []string
)

func LoadConfig() {
	roles := os.Getenv("ALLOWED_ROLES")
	if roles == "" {
		roles = "admin,manager,delivery"
	}
	AllowedRoles = strings.Split(roles, ",")

	statuses := os.Getenv("ALLOWED_STATUSES")
	if statuses == "" {
		statuses = "active,inactive,suspended"
	}
	AllowedStatuses = strings.Split(statuses, ",")
}
