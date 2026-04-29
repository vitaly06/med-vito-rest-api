package rbac

import "strings"

const (
	RoleGuest           = "GUEST"
	RoleUser            = "USER"
	RoleUserVerified    = "USER_VERIFIED"
	RoleSeniorModerator = "SENIOR_MODERATOR"
	RoleAdmin           = "ADMIN"
	RoleSuperAdmin      = "SUPERADMIN"
)

var roleLevels = map[string]int{
	RoleGuest:           0,
	RoleUser:            10,
	RoleUserVerified:    15,
	RoleSeniorModerator: 70,
	RoleAdmin:           90,
	RoleSuperAdmin:      100,
}

// NormalizeRole поддерживает legacy-имена ролей из текущей БД.
func NormalizeRole(roleName *string) string {
	if roleName == nil {
		return RoleGuest
	}
	v := strings.ToUpper(strings.TrimSpace(*roleName))
	switch v {
	case "DEFAULT":
		return RoleUser
	case "MODERATOR", "CONTENT_MANAGER", "TECH_SUPPORT":
		return RoleSeniorModerator
	case "ADMIN":
		return RoleAdmin
	case "SUPERADMIN", "ROOT":
		return RoleSuperAdmin
	case RoleGuest, RoleUser, RoleUserVerified, RoleSeniorModerator:
		return v
	default:
		return RoleUser
	}
}

func RoleLevel(roleName *string) int {
	return roleLevels[NormalizeRole(roleName)]
}

func HasMinRole(roleName *string, minLevel int) bool {
	return RoleLevel(roleName) >= minLevel
}
