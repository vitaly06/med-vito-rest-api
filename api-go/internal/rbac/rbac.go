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

const (
	PermAdsManageOwn         = "ads.manage.own"
	PermAdsReadPublic        = "ads.read.public"
	PermPaymentsUse          = "payments.use"
	PermPromotionUse         = "promotion.use"
	PermBannerRequest        = "banner.request"
	PermModerationQueueRead  = "moderation.queue.read"
	PermModerationDecision   = "moderation.decision"
	PermModerationAILogsRead = "moderation.ai_logs.read"
	PermUserBlock            = "user.block"
	PermUserUnblock          = "user.unblock"
	PermDealsReadAll         = "deals.read.all"
	PermReservationsReadAll  = "reservations.read.all"
	PermBalancesRead         = "balances.read"
	PermBalancesWrite        = "balances.write"
	PermUsersRead            = "users.read"
	PermUsersEdit            = "users.edit"
	PermRolesAssign          = "roles.assign"
	PermTariffsManage        = "tariffs.manage"
	PermAnalyticsRead        = "analytics.read"
	PermAIFiltersManage      = "ai.filters.manage"
	PermAITrain              = "ai.train"
	PermIntegrationsGov      = "integrations.gov"
	PermBackupsManage        = "backups.manage"
	PermUpdatesManage        = "updates.manage"
	PermAppealsCreate        = "appeals.create"
	PermAppealsReview        = "appeals.review"
)

var rolePermissions = map[string]map[string]struct{}{
	RoleGuest: {
		PermAdsReadPublic: {},
	},
	RoleUser: {
		PermAdsManageOwn:  {},
		PermAdsReadPublic: {},
		PermPaymentsUse:   {},
		PermPromotionUse:  {},
		PermBannerRequest: {},
		PermAppealsCreate: {},
	},
	RoleUserVerified: {
		PermAdsManageOwn:  {},
		PermAdsReadPublic: {},
		PermPaymentsUse:   {},
		PermPromotionUse:  {},
		PermBannerRequest: {},
		PermAppealsCreate: {},
	},
	RoleSeniorModerator: {
		PermModerationQueueRead:  {},
		PermModerationDecision:   {},
		PermModerationAILogsRead: {},
		PermUserBlock:            {},
		PermUserUnblock:          {},
		PermDealsReadAll:         {},
		PermReservationsReadAll:  {},
		PermBalancesRead:         {},
		PermAppealsReview:        {},
	},
	RoleAdmin: {
		PermModerationQueueRead:  {},
		PermModerationDecision:   {},
		PermModerationAILogsRead: {},
		PermUserBlock:            {},
		PermUserUnblock:          {},
		PermDealsReadAll:         {},
		PermReservationsReadAll:  {},
		PermBalancesRead:         {},
		PermBalancesWrite:        {},
		PermUsersRead:            {},
		PermUsersEdit:            {},
		PermRolesAssign:          {},
		PermTariffsManage:        {},
		PermAnalyticsRead:        {},
		PermAIFiltersManage:      {},
		PermAppealsReview:        {},
	},
	RoleSuperAdmin: {
		PermModerationQueueRead:  {},
		PermModerationDecision:   {},
		PermModerationAILogsRead: {},
		PermUserBlock:            {},
		PermUserUnblock:          {},
		PermDealsReadAll:         {},
		PermReservationsReadAll:  {},
		PermBalancesRead:         {},
		PermBalancesWrite:        {},
		PermUsersRead:            {},
		PermUsersEdit:            {},
		PermRolesAssign:          {},
		PermTariffsManage:        {},
		PermAnalyticsRead:        {},
		PermAIFiltersManage:      {},
		PermAITrain:              {},
		PermIntegrationsGov:      {},
		PermBackupsManage:        {},
		PermUpdatesManage:        {},
		PermAppealsReview:        {},
	},
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

func HasPermission(roleName *string, permission string) bool {
	role := NormalizeRole(roleName)
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	_, ok = perms[permission]
	return ok
}
