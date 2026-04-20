package main

// Р С›Р С—Р С‘РЎРѓР В°Р Р…Р С‘Р Вµ Р СР В°РЎР‚РЎв‚¬РЎР‚РЎС“РЎвЂљР С•Р Р† РЎвЂљР С•Р В»РЎРЉР С”Р С• Р Т‘Р В»РЎРЏ Р С–Р ВµР Р…Р ВµРЎР‚Р В°РЎвЂ Р С‘Р С‘ OpenAPI (swag). Р В Р ВµР В°Р В»РЎРЉР Р…РЎвЂ№Р Вµ РЎвЂ¦Р ВµР Р…Р Т‘Р В»Р ВµРЎР‚РЎвЂ№ Р Р† internal/httpserver.
// Р С’Р Т‘Р СР С‘Р Р…-РЎР‚РЎС“РЎвЂЎР С”Р С‘ Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘Р в„–: РЎРѓР Р…Р В°РЎвЂЎР В°Р В»Р В° POST /auth/sign-in Р С—Р С•Р Т‘ Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»Р ВµР С РЎРѓ РЎР‚Р С•Р В»РЎРЉРЎР‹ admin РІР‚вЂќ cookie session_id РЎС“Р в„–Р Т‘РЎвЂРЎвЂљ Р Р† Try it out (credentials).

// --- system ---

// HealthCheck
// @Summary Р СџРЎР‚Р С•Р Р†Р ВµРЎР‚Р С”Р В° Р В¶Р С‘Р Р†Р С•РЎРѓРЎвЂљР С‘ РЎРѓР ВµРЎР‚Р Р†Р С‘РЎРѓР В°
// @Description Р вЂ™Р С•Р В·Р Р†РЎР‚Р В°РЎвЂ°Р В°Р ВµРЎвЂљ status ok
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func _swaggerHealth() {}

// --- log ---

// LogFindAll
// @Summary Р РЋР С—Р С‘РЎРѓР С•Р С” Р В·Р В°Р С—Р С‘РЎРѓР ВµР в„– Log
// @Tags log
// @Produce json
// @Success 200 {array} object
// @Router /log/find-all [get]
func _swaggerLogFindAll() {}

// --- category (Р С—РЎС“Р В±Р В»Р С‘РЎвЂЎР Р…Р С•) ---

// CategoryFindAll
// @Summary Р вЂќР ВµРЎР‚Р ВµР Р†Р С• Р Р†РЎРѓР ВµРЎвЂ¦ Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘Р в„–
// @Tags category
// @Produce json
// @Success 200 {array} object
// @Router /category/find-all [get]
func _swaggerCategoryFindAll() {}

// CategoryFindByID
// @Summary Р С™Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘РЎРЏ Р С—Р С• id (Р С—Р С•Р В»Р Р…Р С•Р Вµ Р Т‘Р ВµРЎР‚Р ВµР Р†Р С•)
// @Tags category
// @Produce json
// @Param id path int true "ID Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘Р С‘"
// @Success 200 {object} object
// @Failure 404 {object} map[string]interface{}
// @Router /category/find-by-id/{id} [get]
func _swaggerCategoryFindByID() {}

// CategoryFindBySlug
// @Summary Р С™Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘РЎРЏ Р С—Р С• slug
// @Tags category
// @Produce json
// @Param slug path string true "Slug"
// @Success 200 {object} object
// @Failure 404 {object} map[string]interface{}
// @Router /category/slug/{slug} [get]
func _swaggerCategoryFindBySlug() {}

// CategoryFindBySlugPath
// @Summary Р В Р В°Р В·РЎР‚Р ВµРЎв‚¬Р ВµР Р…Р С‘Р Вµ РЎвЂ Р ВµР С—Р С•РЎвЂЎР С”Р С‘ slug (Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘РЎРЏ / Р С—Р С•Р Т‘Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘РЎРЏ / РЎвЂљР С‘Р С—)
// @Description Р вЂ™ Swagger Р’В«Try it outР’В» Р Р†Р Р†Р ВµР Т‘Р С‘ РЎРѓР ВµР С–Р СР ВµР Р…РЎвЂљРЎвЂ№ РЎвЂЎР ВµРЎР‚Р ВµР В· %2F, Р Р…Р В°Р С—РЎР‚Р С‘Р СР ВµРЎР‚: elektronika%2Ftelefony
// @Tags category
// @Produce json
// @Param slugPath path string true "Р В¦Р ВµР С—Р С•РЎвЂЎР С”Р В° (Р С‘Р В»Р С‘ Р С•Р Т‘Р С‘Р Р… РЎРѓР ВµР С–Р СР ВµР Р…РЎвЂљ)"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /category/path/{slugPath} [get]
func _swaggerCategoryPath() {}

// --- category (Р В°Р Т‘Р СР С‘Р Р…; Р С—РЎР‚Р С‘ ADMIN_API_KEY РІР‚вЂќ Р В·Р В°Р С–Р С•Р В»Р С•Р Р†Р С•Р С” X-Admin-Key) ---

// CategoryCreate
// @Summary Р РЋР С•Р В·Р Т‘Р В°РЎвЂљРЎРЉ Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘РЎР‹
// @Tags category-admin
// @Accept json
// @Produce json
// @Param body body swaggerCreateCategory true "Р СћР ВµР В»Р С•"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /category/create-category [post]
func _swaggerCategoryCreate() {}

// CategoryUpdate
// @Summary Р С›Р В±Р Р…Р С•Р Р†Р С‘РЎвЂљРЎРЉ Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘РЎР‹
// @Tags category-admin
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body swaggerUpdateCategory true "Р СћР ВµР В»Р С•"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /category/update-category/{id} [put]
func _swaggerCategoryUpdate() {}

// CategoryDelete
// @Summary Р Р€Р Т‘Р В°Р В»Р С‘РЎвЂљРЎРЉ Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘РЎР‹
// @Tags category-admin
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]interface{}
// @Router /category/delete-category/{id} [delete]
func _swaggerCategoryDelete() {}

// Р СћР ВµР В»Р В° Р В·Р В°Р С—РЎР‚Р С•РЎРѓР С•Р Р† Р Т‘Р В»РЎРЏ Swagger UI
type swaggerCreateCategory struct {
	Name string  `json:"name" example:"Р С’Р Р†РЎвЂљР С•Р СР С•Р В±Р С‘Р В»Р С‘"`
	Slug *string `json:"slug,omitempty" example:"avtomobili"`
}

type swaggerUpdateCategory struct {
	Name string  `json:"name" example:"Р С’Р Р†РЎвЂљР С•Р СР С•Р В±Р С‘Р В»Р С‘"`
	Slug *string `json:"slug,omitempty"`
}

// --- auth ---

// AuthSignUp
// @Summary Р В Р ВµР С–Р С‘РЎРѓРЎвЂљРЎР‚Р В°РЎвЂ Р С‘РЎРЏ РІР‚вЂќ Р С•РЎвЂљР С—РЎР‚Р В°Р Р†Р С”Р В° Р С”Р С•Р Т‘Р В° (query where=telegram|sms)
// @Tags auth
// @Accept json
// @Produce json
// @Param where query string true "telegram Р С‘Р В»Р С‘ sms" Enums(telegram,sms)
// @Param body body swaggerSignUp true "Р вЂќР В°Р Р…Р Р…РЎвЂ№Р Вµ"
// @Success 200 {object} map[string]string
// @Router /auth/sign-up [post]
func _swaggerAuthSignUp() {}

// AuthVerifyMobile
// @Summary Р СџР С•Р Т‘РЎвЂљР Р†Р ВµРЎР‚Р В¶Р Т‘Р ВµР Р…Р С‘Р Вµ РЎвЂљР ВµР В»Р ВµРЎвЂћР С•Р Р…Р В° Р С—Р С• Р С”Р С•Р Т‘РЎС“
// @Tags auth
// @Produce json
// @Param code query string true "Р С™Р С•Р Т‘ Р С‘Р В· SMS/TG"
// @Success 200 {object} map[string]string
// @Router /auth/verify-mobile-code [post]
func _swaggerAuthVerifyMobile() {}

// AuthSignIn
// @Summary Р вЂ™РЎвЂ¦Р С•Р Т‘ (РЎРѓРЎвЂљР В°Р Р†Р С‘РЎвЂљ cookie session_id)
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerSignIn true "Р вЂєР С•Р С–Р С‘Р Р…"
// @Success 200 {object} object
// @Router /auth/sign-in [post]
func _swaggerAuthSignIn() {}

// AuthMe
// @Summary Р СћР ВµР С”РЎС“РЎвЂ°Р С‘Р в„– Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»РЎРЉ
// @Tags auth
// @Produce json
// @Success 200 {object} object
// @Router /auth/me [get]
func _swaggerAuthMe() {}

// AuthIsAdmin
// @Summary Р СџРЎР‚Р С•Р Р†Р ВµРЎР‚Р С”Р В° РЎР‚Р С•Р В»Р С‘ admin
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]bool
// @Router /auth/isAdmin [get]
func _swaggerAuthIsAdmin() {}

// AuthLogout
// @Summary Р вЂ™РЎвЂ№РЎвЂ¦Р С•Р Т‘
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func _swaggerAuthLogout() {}

// AuthForgot
// @Summary Р вЂ”Р В°Р С—РЎР‚Р С•РЎРѓ Р С”Р С•Р Т‘Р В° РЎРѓР В±РЎР‚Р С•РЎРѓР В° Р Р…Р В° Р С—Р С•РЎвЂЎРЎвЂљРЎС“
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerForgotEmail true "email"
// @Success 200 {object} map[string]string
// @Router /auth/forgot-password [post]
func _swaggerAuthForgot() {}

// AuthVerifyForgot
// @Summary Р СџРЎР‚Р С•Р Р†Р ВµРЎР‚Р С”Р В° Р С”Р С•Р Т‘Р В° РЎРѓР В±РЎР‚Р С•РЎРѓР В°
// @Tags auth
// @Produce json
// @Param code query string true "Р С™Р С•Р Т‘ Р С‘Р В· Р С—Р С‘РЎРѓРЎРЉР СР В°"
// @Success 200 {object} map[string]int
// @Router /auth/verify-code [post]
func _swaggerAuthVerifyForgot() {}

// AuthChangePassword
// @Summary Р СњР С•Р Р†РЎвЂ№Р в„– Р С—Р В°РЎР‚Р С•Р В»РЎРЉ Р С—Р С•РЎРѓР В»Р Вµ verify-code
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerChangePassword true "Р СћР ВµР В»Р С•"
// @Success 200 {object} map[string]string
// @Router /auth/change-password [post]
func _swaggerAuthChangePassword() {}

type swaggerSignUp struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

type swaggerSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type swaggerChangePassword struct {
	UserID   int    `json:"userId"`
	Password string `json:"password"`
}

type swaggerForgotEmail struct {
	Email string `json:"email"`
}

// --- user ---

// UserFindAll
// @Summary Р РЋР С—Р С‘РЎРѓР С•Р С” Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»Р ВµР в„– (Р В°Р Т‘Р СР С‘Р Р…: cookie session_id + РЎР‚Р С•Р В»РЎРЉ admin)
// @Tags user-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /user/find-all [get]
func _swaggerUserFindAll() {}

// UserInfo
// @Summary Р С™Р В°РЎР‚РЎвЂљР С•РЎвЂЎР С”Р В° Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»РЎРЏ (РЎР‚Р ВµР в„–РЎвЂљР С‘Р Р…Р С–, Р В»Р С‘Р СР С‘РЎвЂљ Р С•Р В±РЎР‰РЎРЏР Р†Р В»Р ВµР Р…Р С‘Р в„–)
// @Tags user
// @Produce json
// @Param id path int true "User id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /user/info/{id} [get]
func _swaggerUserInfo() {}

// UserRemainingFreeAds
// @Summary Р С›РЎРѓРЎвЂљР В°РЎвЂљР С•Р С” Р В±Р ВµРЎРѓР С—Р В»Р В°РЎвЂљР Р…РЎвЂ№РЎвЂ¦ Р С•Р В±РЎР‰РЎРЏР Р†Р В»Р ВµР Р…Р С‘Р в„– (РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags user
// @Produce json
// @Success 200 {object} object
// @Router /user/remaining-free-ads [get]
func _swaggerUserRemainingFreeAds() {}

// UserShowNumber
// @Summary Р СџР С•Р С”Р В°Р В·Р В°РЎвЂљРЎРЉ Р Р…Р С•Р СР ВµРЎР‚ Р С—РЎР‚Р С•Р Т‘Р В°Р Р†РЎвЂ Р В° (РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags user
// @Produce json
// @Param userId path int true "Р СџРЎР‚Р С•Р Т‘Р В°Р Р†Р ВµРЎвЂ "
// @Success 200 {object} map[string]string
// @Router /user/show-number/{userId} [get]
func _swaggerUserShowNumber() {}

// UserUpdateSettings
// @Summary Р С›Р В±Р Р…Р С•Р Р†Р В»Р ВµР Р…Р С‘Р Вµ Р Р…Р В°РЎРѓРЎвЂљРЎР‚Р С•Р ВµР С” (multipart, РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags user
// @Accept mpfd
// @Produce json
// @Param fullName formData string false "Р В¤Р ВР С›"
// @Param phoneNumber formData string false "Р СћР ВµР В»Р ВµРЎвЂћР С•Р Р…"
// @Param isAnswersCall formData string false "true/false"
// @Param profileType formData string false "INDIVIDUAL|OOO|IP"
// @Param photo formData file false "Р С’Р Р†Р В°РЎвЂљР В°РЎР‚"
// @Success 200 {object} object
// @Router /user/update-settings [patch]
func _swaggerUserUpdateSettings() {}

// UserVerifyEmail
// @Summary Р С›РЎвЂљР С—РЎР‚Р В°Р Р†Р С‘РЎвЂљРЎРЉ Р С”Р С•Р Т‘ Р С—Р С•Р Т‘РЎвЂљР Р†Р ВµРЎР‚Р В¶Р Т‘Р ВµР Р…Р С‘РЎРЏ Р Р…Р В° Р С—Р С•РЎвЂЎРЎвЂљРЎС“ (РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags user
// @Produce json
// @Success 200 {object} map[string]string
// @Router /user/verify-email [post]
func _swaggerUserVerifyEmail() {}

// UserVerifyEmailCode
// @Summary Р СџР С•Р Т‘РЎвЂљР Р†Р ВµРЎР‚Р Т‘Р С‘РЎвЂљРЎРЉ Р С—Р С•РЎвЂЎРЎвЂљРЎС“ Р С—Р С• Р С”Р С•Р Т‘РЎС“ Р С‘Р В· Р С—Р С‘РЎРѓРЎРЉР СР В°
// @Tags user
// @Produce json
// @Param code query string true "Р С™Р С•Р Т‘"
// @Success 200 {object} map[string]string
// @Router /user/verify-code [post]
func _swaggerUserVerifyEmailCode() {}

// UserSetBalance
// @Summary Р Р€РЎРѓРЎвЂљР В°Р Р…Р С•Р Р†Р С‘РЎвЂљРЎРЉ bonusBalance (Р В°Р Т‘Р СР С‘Р Р…, РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags user-admin
// @Produce json
// @Param userId path int true "User id"
// @Param balance query string true "Р В§Р С‘РЎРѓР В»Р С•"
// @Success 200 {object} map[string]string
// @Router /user/set-balance/{userId} [put]
func _swaggerUserSetBalance() {}

// UserToggleBanned
// @Summary Р вЂР В°Р Р… / РЎР‚Р В°Р В·Р В±Р В°Р Р… (Р В°Р Т‘Р СР С‘Р Р…, РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags user-admin
// @Produce json
// @Param id path int true "User id"
// @Success 200 {object} map[string]string
// @Router /user/toggle-banned/{id} [put]
func _swaggerUserToggleBanned() {}

// UserAdminPatch
// @Summary Р С›Р В±Р Р…Р С•Р Р†Р С‘РЎвЂљРЎРЉ Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»РЎРЏ (Р В°Р Т‘Р СР С‘Р Р…, РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags user-admin
// @Accept json
// @Produce json
// @Param id path int true "User id"
// @Param body body swaggerAdminUpdateUser true "Р СџР С•Р В»РЎРЏ"
// @Success 200 {object} map[string]string
// @Router /user/{id} [patch]
func _swaggerUserAdminPatch() {}

// UserAdminDelete
// @Summary Р Р€Р Т‘Р В°Р В»Р С‘РЎвЂљРЎРЉ Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»РЎРЏ (Р В°Р Т‘Р СР С‘Р Р…, РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags user-admin
// @Produce json
// @Param id path int true "User id"
// @Success 200 {object} map[string]string
// @Router /user/{id} [delete]
func _swaggerUserAdminDelete() {}

type swaggerAdminUpdateUser struct {
	FullName     *string  `json:"fullName,omitempty"`
	Email        *string  `json:"email,omitempty"`
	PhoneNumber  *string  `json:"phoneNumber,omitempty"`
	ProfileType  *string  `json:"profileType,omitempty"`
	BonusBalance *float64 `json:"bonusBalance,omitempty"`
}

// --- product ---

// ProductCreate
// @Summary Р РЋР С•Р В·Р Т‘Р В°РЎвЂљРЎРЉ РЎвЂљР С•Р Р†Р В°РЎР‚ (multipart, Р Т‘Р С• 8 images, РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags product
// @Accept mpfd
// @Produce json
// @Success 201 {object} object
// @Router /product/create [post]
func _swaggerProductCreate() {}

// ProductAll
// @Summary Р РЋР С—Р С‘РЎРѓР С•Р С” РЎвЂљР С•Р Р†Р В°РЎР‚Р С•Р Р† / Р С—Р С•Р С‘РЎРѓР С” (query; optional РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ Р Т‘Р В»РЎРЏ Р С‘Р В·Р В±РЎР‚Р В°Р Р…Р Р…Р С•Р С–Р С•)
// @Tags product
// @Produce json
// @Router /product/all-products [get]
func _swaggerProductAll() {}

// ProductCard
// @Summary Р С™Р В°РЎР‚РЎвЂљР С•РЎвЂЎР С”Р В° РЎвЂљР С•Р Р†Р В°РЎР‚Р В°
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Router /product/product-card/{id} [get]
func _swaggerProductCard() {}

// ProductDelete
// @Summary Р Р€Р Т‘Р В°Р В»Р С‘РЎвЂљРЎРЉ РЎРѓР Р†Р С•Р в„– РЎвЂљР С•Р Р†Р В°РЎР‚
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Router /product/{id} [delete]
func _swaggerProductDelete() {}

// ProductPatch
// @Summary Р С›Р В±Р Р…Р С•Р Р†Р С‘РЎвЂљРЎРЉ РЎвЂљР С•Р Р†Р В°РЎР‚ (multipart)
// @Tags product
// @Accept mpfd
// @Produce json
// @Param id path int true "Product id"
// @Router /product/{id} [patch]
func _swaggerProductPatch() {}

// ProductModerate
// @Summary Р СљР С•Р Т‘Р ВµРЎР‚Р В°РЎвЂ Р С‘РЎРЏ (admin)
// @Tags product-admin
// @Produce json
// @Param id path int true "Product id"
// @Param status query string true "APPROVED|DENIDED"
// @Param reason query string false "Р СџРЎР‚Р С‘РЎвЂЎР С‘Р Р…Р В° Р С—РЎР‚Р С‘ DENIDED"
// @Router /product/moderate-product/{id} [put]
func _swaggerProductModerate() {}

// --- review ---

// ReviewSend
// @Summary Р С›РЎРѓРЎвЂљР В°Р Р†Р С‘РЎвЂљРЎРЉ Р С•РЎвЂљР В·РЎвЂ№Р Р† Р С—РЎР‚Р С•Р Т‘Р В°Р Р†РЎвЂ РЎС“ (РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags review
// @Accept json
// @Produce json
// @Param body body swaggerSendReview true "Р СћР ВµР В»Р С•"
// @Success 200 {object} map[string]string
// @Router /review/send-review [post]
func _swaggerReviewSend() {}

// ReviewUserReviews
// @Summary Р С›Р Т‘Р С•Р В±РЎР‚Р ВµР Р…Р Р…РЎвЂ№Р Вµ Р С•РЎвЂљР В·РЎвЂ№Р Р†РЎвЂ№ Р С• Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»Р Вµ
// @Tags review
// @Produce json
// @Param id path int true "User id"
// @Router /review/user-reviews/{id} [get]
func _swaggerReviewUserReviews() {}

// ReviewModerate
// @Summary Р СљР С•Р Т‘Р ВµРЎР‚Р В°РЎвЂ Р С‘РЎРЏ Р С•РЎвЂљР В·РЎвЂ№Р Р†Р В° (admin)
// @Tags review-admin
// @Produce json
// @Param id path int true "Review id"
// @Param status query string true "APPROVED|DENIDED"
// @Router /review/moderate-review/{id} [put]
func _swaggerReviewModerate() {}

// ReviewModerateList
// @Summary Р С›РЎвЂЎР ВµРЎР‚Р ВµР Т‘РЎРЉ Р С•РЎвЂљР В·РЎвЂ№Р Р†Р С•Р Р† Р Р…Р В° Р СР С•Р Т‘Р ВµРЎР‚Р В°РЎвЂ Р С‘РЎР‹ (admin)
// @Tags review-admin
// @Produce json
// @Router /review/all-reviews-to-moderate [get]
func _swaggerReviewModerateList() {}

type swaggerSendReview struct {
	Text           *string `json:"text,omitempty"`
	Rating         float64 `json:"rating" example:"5"`
	ReviewedUserID int32   `json:"reviewedUserId"`
}

// --- chat (РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ cookie session_id) ---

// ChatStart
// @Summary Р СњР В°РЎвЂЎР В°РЎвЂљРЎРЉ РЎвЂЎР В°РЎвЂљ Р С—Р С• РЎвЂљР С•Р Р†Р В°РЎР‚РЎС“
// @Tags chat
// @Accept json
// @Produce json
// @Param body body swaggerStartChat true "productId"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /chat/start [post]
func _swaggerChatStart() {}

// ChatList
// @Summary Р РЋР С—Р С‘РЎРѓР С•Р С” РЎвЂЎР В°РЎвЂљР С•Р Р† Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»РЎРЏ
// @Tags chat
// @Produce json
// @Success 200 {array} object
// @Router /chat [get]
func _swaggerChatList() {}

// ChatMessages
// @Summary Р РЋР С•Р С•Р В±РЎвЂ°Р ВµР Р…Р С‘РЎРЏ РЎвЂЎР В°РЎвЂљР В° (Р С—Р В°Р С–Р С‘Р Р…Р В°РЎвЂ Р С‘РЎРЏ)
// @Tags chat
// @Produce json
// @Param id path int true "Chat id"
// @Param page query int false "Р РЋРЎвЂљРЎР‚Р В°Р Р…Р С‘РЎвЂ Р В°" default(1)
// @Param limit query int false "Р вЂєР С‘Р СР С‘РЎвЂљ" default(50)
// @Success 200 {object} object
// @Router /chat/{id}/messages [get]
func _swaggerChatMessages() {}

// ChatInfo
// @Summary Р ВР Р…РЎвЂћР С•РЎР‚Р СР В°РЎвЂ Р С‘РЎРЏ Р С• РЎвЂЎР В°РЎвЂљР Вµ
// @Tags chat
// @Produce json
// @Param id path int true "Chat id"
// @Success 200 {object} object
// @Router /chat/{id} [get]
func _swaggerChatInfo() {}

type swaggerStartChat struct {
	ProductID int32 `json:"productId" example:"1"`
}

// --- payment (Р Сћ-Р вЂР В°Р Р…Р С” / Tinkoff; РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ РІР‚вЂќ cookie session_id) ---

// PaymentCreate
// @Summary Р РЋР С•Р В·Р Т‘Р В°Р Р…Р С‘Р Вµ Р С—Р В»Р В°РЎвЂљР ВµР В¶Р В° Р Т‘Р В»РЎРЏ Р С—Р С•Р С—Р С•Р В»Р Р…Р ВµР Р…Р С‘РЎРЏ Р В±Р В°Р В»Р В°Р Р…РЎРѓР В°
// @Description Init Р Р† Р Сћ-Р вЂР В°Р Р…Р С”. Р СњРЎС“Р В¶Р Р…РЎвЂ№ TINKOFF_TERMINAL_KEY Р С‘ TINKOFF_SECRET_KEY. Р С’Р Р†РЎвЂљР С•РЎР‚Р С‘Р В·Р В°РЎвЂ Р С‘РЎРЏ: cookie session_id Р С—Р С•РЎРѓР В»Р Вµ POST /auth/sign-in.
// @Tags payment
// @Accept json
// @Produce json
// @Param body body swaggerCreatePayment true "Р РЋРЎС“Р СР СР В° Р Р† РЎР‚РЎС“Р В±Р В»РЎРЏРЎвЂ¦ (Р СР С‘Р Р…. 1)"
// @Success 201 {object} swaggerPaymentCreateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /payment/create [post]
func _swaggerPaymentCreate() {}

// PaymentNotification
// @Summary Webhook РЎС“Р Р†Р ВµР Т‘Р С•Р СР В»Р ВµР Р…Р С‘Р в„– Р Сћ-Р вЂР В°Р Р…Р С”Р В° Р С• РЎРѓРЎвЂљР В°РЎвЂљРЎС“РЎРѓР Вµ Р С—Р В»Р В°РЎвЂљР ВµР В¶Р В°
// @Description Р вЂР ВµР В· РЎРѓР ВµРЎРѓРЎРѓР С‘Р С‘. Р СџР С•Р Т‘Р С—Р С‘РЎРѓРЎРЉ Token Р С—РЎР‚Р С•Р Р†Р ВµРЎР‚РЎРЏР ВµРЎвЂљРЎРѓРЎРЏ Р С—Р С• Р С—Р С•Р В»РЎРЏР С РЎвЂљР ВµР В»Р В°. Р СћР ВµР В»Р С• РІР‚вЂќ Р С”Р В°Р С” Р С—РЎР‚Р С‘РЎвЂ¦Р С•Р Т‘Р С‘РЎвЂљ Р С•РЎвЂљ Р В±Р В°Р Р…Р С”Р В°; Р С—РЎР‚Р С‘Р СР ВµРЎР‚ Р Р…Р С‘Р В¶Р Вµ.
// @Tags payment
// @Accept json
// @Produce json
// @Param body body swaggerTinkoffNotification true "Р Р€Р Р†Р ВµР Т‘Р С•Р СР В»Р ВµР Р…Р С‘Р Вµ"
// @Success 200 {object} swaggerPaymentNotifyResponse
// @Failure 400 {object} map[string]interface{}
// @Router /payment/notification [post]
func _swaggerPaymentNotification() {}

// PaymentHistory
// @Summary Р ВРЎРѓРЎвЂљР С•РЎР‚Р С‘РЎРЏ Р С—Р В»Р В°РЎвЂљР ВµР В¶Р ВµР в„– Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»РЎРЏ
// @Description Р вЂќР С• 50 Р В·Р В°Р С—Р С‘РЎРѓР ВµР в„–, Р Р…Р С•Р Р†РЎвЂ№Р Вµ РЎРѓР Р†Р ВµРЎР‚РЎвЂ¦РЎС“. Р РЋР ВµРЎРѓРЎРѓР С‘РЎРЏ: cookie session_id.
// @Tags payment
// @Produce json
// @Success 200 {array} swaggerPaymentHistoryItem
// @Failure 401 {object} map[string]interface{}
// @Router /payment/history [get]
func _swaggerPaymentHistory() {}

// PaymentCheckStatus
// @Summary Р СџРЎР‚Р С•Р Р†Р ВµРЎР‚Р С”Р В° РЎРѓРЎвЂљР В°РЎвЂљРЎС“РЎРѓР В° Р С—Р В»Р В°РЎвЂљР ВµР В¶Р В° Р Р† Р Сћ-Р вЂР В°Р Р…Р С”Р Вµ (GetState)
// @Description Р РЋР ВµРЎРѓРЎРѓР С‘РЎРЏ: cookie session_id. Р вЂ™ РЎвЂљР ВµР В»Р Вµ РІР‚вЂќ paymentId Р С‘Р В· Р С•РЎвЂљР Р†Р ВµРЎвЂљР В° Init Р С‘Р В»Р С‘ РЎС“Р Р†Р ВµР Т‘Р С•Р СР В»Р ВµР Р…Р С‘РЎРЏ.
// @Tags payment
// @Accept json
// @Produce json
// @Param body body swaggerCheckPayment true "paymentId"
// @Success 200 {object} swaggerPaymentCheckStateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /payment/check-status [post]
func _swaggerPaymentCheckStatus() {}

type swaggerCreatePayment struct {
	Amount      float64 `json:"amount" example:"1000"`
	Description *string `json:"description,omitempty" example:"Р СџР С•Р С—Р С•Р В»Р Р…Р ВµР Р…Р С‘Р Вµ Р В±Р В°Р В»Р В°Р Р…РЎРѓР В°"`
}

type swaggerCheckPayment struct {
	PaymentID string `json:"paymentId" example:"2673412345"`
}

type swaggerPaymentCreateResponse struct {
	PaymentID  string  `json:"paymentId" example:"2673412345"`
	PaymentURL string  `json:"paymentUrl" example:"https://securepay.tinkoff.ru/..."`
	OrderID    string  `json:"orderId" example:"123-1735123456789"`
	Amount     float64 `json:"amount" example:"1000"`
}

type swaggerPaymentNotifyResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"Р вЂР В°Р В»Р В°Р Р…РЎРѓ РЎС“РЎРѓР С—Р ВµРЎв‚¬Р Р…Р С• Р С—Р С•Р С—Р С•Р В»Р Р…Р ВµР Р…"`
}

type swaggerPaymentHistoryItem struct {
	ID         int32   `json:"id"`
	OrderID    string  `json:"orderId"`
	PaymentID  string  `json:"paymentId"`
	UserID     int32   `json:"userId"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status" example:"COMPLETED"`
	PaymentURL *string `json:"paymentUrl"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

type swaggerPaymentCheckStateResponse struct {
	Status  string  `json:"status" example:"CONFIRMED"`
	Amount  float64 `json:"amount" example:"1000"`
	OrderID string  `json:"orderId"`
}

// swaggerTinkoffNotification РІР‚вЂќ Р С—Р С•Р В»РЎРЏ Р С”Р В°Р С” Р Р† Nest PaymentNotificationDto (РЎР‚Р ВµР В°Р В»РЎРЉР Р…РЎвЂ№Р в„– webhook Р СР С•Р В¶Р ВµРЎвЂљ Р Т‘Р С•Р В±Р В°Р Р†Р В»РЎРЏРЎвЂљРЎРЉ Р С—Р С•Р В»РЎРЏ).
type swaggerTinkoffNotification struct {
	TerminalKey string `json:"TerminalKey" example:"1766153689307DEMO"`
	OrderID     string `json:"OrderId" example:"123-1735123456789"`
	Success     bool   `json:"Success" example:"true"`
	Status      string `json:"Status" example:"CONFIRMED"`
	PaymentID   string `json:"PaymentId" example:"2673412345"`
	Amount      int64  `json:"Amount" example:"100000"`
	Token       string `json:"Token" example:"Р С—Р С•Р Т‘Р С—Р С‘РЎРѓРЎРЉ_Р С•РЎвЂљ_Р В±Р В°Р Р…Р С”Р В°"`
	ErrorCode   string `json:"ErrorCode,omitempty" example:"0"`
	Pan         string `json:"Pan,omitempty" example:"430000******0777"`
}

// --- promotion ---

// PromotionAll
// @Summary Р вЂ™РЎРѓР Вµ РЎвЂљР С‘Р С—РЎвЂ№ Р С—РЎР‚Р С•Р Т‘Р Р†Р С‘Р В¶Р ВµР Р…Р С‘РЎРЏ (РЎвЂљР В°РЎР‚Р С‘РЎвЂћРЎвЂ№)
// @Tags promotion
// @Produce json
// @Success 200 {array} object
// @Router /promotion/all-promotions [get]
func _swaggerPromotionAll() {}

// PromotionAdd
// @Summary Р СџР С•Р Т‘Р С”Р В»РЎР‹РЎвЂЎР С‘РЎвЂљРЎРЉ Р С—РЎР‚Р С•Р Т‘Р Р†Р С‘Р В¶Р ВµР Р…Р С‘Р Вµ Р С” РЎвЂљР С•Р Р†Р В°РЎР‚РЎС“ (РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags promotion
// @Accept json
// @Produce json
// @Param body body swaggerAddPromotion true "Р СћР ВµР В»Р С•"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /promotion/add-promotion [post]
func _swaggerPromotionAdd() {}

type swaggerAddPromotion struct {
	ProductID   int32 `json:"productId" example:"123"`
	PromotionID int32 `json:"promotionId" example:"1"`
	Days        int32 `json:"days" example:"7"`
}

// --- statistics ---

// StatisticsAnalytic
// @Summary Р РЋРЎвЂљР В°РЎвЂљР С‘РЎРѓРЎвЂљР С‘Р С”Р В° Р С—Р С•Р В»РЎРЉР В·Р С•Р Р†Р В°РЎвЂљР ВµР В»РЎРЏ (Р С—РЎР‚Р С•РЎРѓР СР С•РЎвЂљРЎР‚РЎвЂ№, РЎвЂљР ВµР В»Р ВµРЎвЂћР С•Р Р…, Р С‘Р В·Р В±РЎР‚Р В°Р Р…Р Р…Р С•Р Вµ)
// @Tags statistics
// @Produce json
// @Param period query string false "day week month quarter half-year year"
// @Param categoryId query int false "Р В¤Р С‘Р В»РЎРЉРЎвЂљРЎР‚ Р С—Р С• Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘Р С‘ (РЎвЂЎР ВµРЎР‚Р ВµР В· SubCategory)"
// @Param region query string false "Р СџР С•Р Т‘РЎРѓРЎвЂљРЎР‚Р С•Р С”Р В° Р Р† address (ILIKE)"
// @Param productId query int false "Р С™Р С•Р Р…Р С”РЎР‚Р ВµРЎвЂљР Р…РЎвЂ№Р в„– РЎвЂљР С•Р Р†Р В°РЎР‚"
// @Success 200 {object} object
// @Router /statistics/analytic [get]
func _swaggerStatisticsAnalytic() {}

// StatisticsProducts
// @Summary Р С’Р Р…Р В°Р В»Р С‘РЎвЂљР С‘Р С”Р В° Р С—Р С• Р С”Р В°Р В¶Р Т‘Р С•Р СРЎС“ РЎвЂљР С•Р Р†Р В°РЎР‚РЎС“ Р С—РЎР‚Р С•Р Т‘Р В°Р Р†РЎвЂ Р В°
// @Tags statistics
// @Produce json
// @Success 200 {array} object
// @Router /statistics/products-analytic [get]
func _swaggerStatisticsProducts() {}

// --- support (РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ) ---

// SupportCreateTicket
// @Summary Р РЋР С•Р В·Р Т‘Р В°РЎвЂљРЎРЉ РЎвЂљР С‘Р С”Р ВµРЎвЂљ Р С—Р С•Р Т‘Р Т‘Р ВµРЎР‚Р В¶Р С”Р С‘
// @Tags support
// @Accept json
// @Produce json
// @Param body body swaggerCreateSupportTicket true "Р СћР ВµР В»Р С•"
// @Success 201 {object} object
// @Router /support/tickets [post]
func _swaggerSupportCreateTicket() {}

// SupportMyTickets
// @Summary Р СљР С•Р С‘ РЎвЂљР С‘Р С”Р ВµРЎвЂљРЎвЂ№ (Р С—Р В°Р С–Р С‘Р Р…Р В°РЎвЂ Р С‘РЎРЏ, РЎвЂћР С‘Р В»РЎРЉРЎвЂљРЎР‚РЎвЂ№ query)
// @Tags support
// @Produce json
// @Router /support/tickets/my [get]
func _swaggerSupportMyTickets() {}

// SupportAllTickets
// @Summary Р вЂ™РЎРѓР Вµ РЎвЂљР С‘Р С”Р ВµРЎвЂљРЎвЂ№ (Р СР С•Р Т‘Р ВµРЎР‚Р В°РЎвЂљР С•РЎР‚/admin)
// @Tags support
// @Produce json
// @Router /support/tickets/all [get]
func _swaggerSupportAllTickets() {}

// SupportStats
// @Summary Р РЋРЎвЂљР В°РЎвЂљР С‘РЎРѓРЎвЂљР С‘Р С”Р В° РЎвЂљР С‘Р С”Р ВµРЎвЂљР С•Р Р† (РЎвЂљР С•Р В»РЎРЉР С”Р С• admin)
// @Tags support
// @Produce json
// @Router /support/stats [get]
func _swaggerSupportStats() {}

// SupportGetTicket
// @Summary Р СћР С‘Р С”Р ВµРЎвЂљ РЎРѓ РЎРѓР С•Р С•Р В±РЎвЂ°Р ВµР Р…Р С‘РЎРЏР СР С‘
// @Tags support
// @Produce json
// @Param id path int true "Ticket id"
// @Router /support/tickets/{id} [get]
func _swaggerSupportGetTicket() {}

// SupportSendMessage
// @Summary Р РЋР С•Р С•Р В±РЎвЂ°Р ВµР Р…Р С‘Р Вµ Р Р† РЎвЂљР С‘Р С”Р ВµРЎвЂљ
// @Tags support
// @Accept json
// @Param id path int true "Ticket id"
// @Param body body swaggerSupportMessage true "Р СћР ВµР С”РЎРѓРЎвЂљ"
// @Router /support/tickets/{id}/messages [post]
func _swaggerSupportSendMessage() {}

// SupportUpdateTicket
// @Summary Р С›Р В±Р Р…Р С•Р Р†Р С‘РЎвЂљРЎРЉ РЎвЂљР С‘Р С”Р ВµРЎвЂљ (Р СР С•Р Т‘Р ВµРЎР‚Р В°РЎвЂљР С•РЎР‚/admin)
// @Tags support
// @Accept json
// @Param id path int true "Ticket id"
// @Param body body swaggerUpdateSupportTicket true "Р СџР С•Р В»РЎРЏ"
// @Router /support/tickets/{id} [put]
func _swaggerSupportUpdateTicket() {}

// SupportAssignTicket
// @Summary Р СњР В°Р В·Р Р…Р В°РЎвЂЎР С‘РЎвЂљРЎРЉ РЎвЂљР С‘Р С”Р ВµРЎвЂљ Р Р…Р В° РЎРѓР ВµР В±РЎРЏ
// @Tags support
// @Param id path int true "Ticket id"
// @Router /support/tickets/{id}/assign [put]
func _swaggerSupportAssignTicket() {}

type swaggerCreateSupportTicket struct {
	Theme    string  `json:"theme" example:"TECHNICAL_ISSUE"`
	Subject  string  `json:"subject"`
	Message  string  `json:"message"`
	Priority *string `json:"priority,omitempty" example:"MEDIUM"`
}

type swaggerSupportMessage struct {
	Text string `json:"text"`
}

type swaggerUpdateSupportTicket struct {
	Status   *string `json:"status,omitempty" example:"IN_PROGRESS"`
	Priority *string `json:"priority,omitempty" example:"HIGH"`
}

// --- address (DaData, Р В±Р ВµР В· РЎРѓР ВµРЎРѓРЎРѓР С‘Р С‘) ---

// AddressSuggestions
// @Summary Р СџР С•Р Т‘РЎРѓР С”Р В°Р В·Р С”Р С‘ Р В°Р Т‘РЎР‚Р ВµРЎРѓР В° (DaData)
// @Tags address
// @Produce json
// @Param query query string true "Р РЋРЎвЂљРЎР‚Р С•Р С”Р В° Р С—Р С•Р С‘РЎРѓР С”Р В°"
// @Param limit query int false "Р вЂєР С‘Р СР С‘РЎвЂљ" default(5)
// @Success 200 {array} object
// @Router /address/suggestions [get]
func _swaggerAddressSuggestions() {}

// AddressValidate
// @Summary Р СџРЎР‚Р С•Р Р†Р ВµРЎР‚Р С”Р В° Р В°Р Т‘РЎР‚Р ВµРЎРѓР В° Р С—Р С• Р С—Р ВµРЎР‚Р Р†Р С•Р в„– Р С—Р С•Р Т‘РЎРѓР С”Р В°Р В·Р С”Р Вµ DaData
// @Tags address
// @Accept json
// @Produce json
// @Param body body swaggerValidateAddress true "Р С’Р Т‘РЎР‚Р ВµРЎРѓ"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /address/validate [post]
func _swaggerAddressValidate() {}

type swaggerValidateAddress struct {
	Address        string `json:"address"`
	AddressDetails any    `json:"addressDetails,omitempty"`
}

// --- banner ---

// BannerCreate
// @Summary Р РЋР С•Р В·Р Т‘Р В°РЎвЂљРЎРЉ Р В±Р В°Р Р…Р Р…Р ВµРЎР‚ (multipart: image, name, place, navigateToUrl; РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ)
// @Tags banner
// @Accept mpfd
// @Produce json
// @Success 201 {object} object
// @Router /banner [post]
func _swaggerBannerCreate() {}

// BannerRandom
// @Summary Р РЋР В»РЎС“РЎвЂЎР В°Р в„–Р Р…РЎвЂ№Р Вµ Р С•Р Т‘Р С•Р В±РЎР‚Р ВµР Р…Р Р…РЎвЂ№Р Вµ Р В±Р В°Р Р…Р Р…Р ВµРЎР‚РЎвЂ№
// @Tags banner
// @Produce json
// @Router /banner/random [get]
func _swaggerBannerRandom() {}

// BannerList
// @Summary Р РЋР С—Р С‘РЎРѓР С•Р С” Р С•Р Т‘Р С•Р В±РЎР‚Р ВµР Р…Р Р…РЎвЂ№РЎвЂ¦ Р В±Р В°Р Р…Р Р…Р ВµРЎР‚Р С•Р Р† (query place Р С•Р С—РЎвЂ Р С‘Р С•Р Р…Р В°Р В»РЎРЉР Р…Р С•)
// @Tags banner
// @Produce json
// @Router /banner [get]
func _swaggerBannerList() {}

// BannerModerate
// @Summary Р СљР С•Р Т‘Р ВµРЎР‚Р В°РЎвЂ Р С‘РЎРЏ Р В±Р В°Р Р…Р Р…Р ВµРЎР‚Р В° (admin, query status)
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Param status query string true "MODERATE|APPROVED|DENIDED"
// @Router /banner/moderate/{id} [put]
func _swaggerBannerModerate() {}

// BannerAllModerate
// @Summary Р С›РЎвЂЎР ВµРЎР‚Р ВµР Т‘РЎРЉ Р В±Р В°Р Р…Р Р…Р ВµРЎР‚Р С•Р Р† Р Р…Р В° Р СР С•Р Т‘Р ВµРЎР‚Р В°РЎвЂ Р С‘РЎР‹ (admin)
// @Tags banner
// @Produce json
// @Router /banner/all-banners-to-moderate [get]
func _swaggerBannerAllModerate() {}

// --- subcategory / subcategory-type / type-field (Р В°Р Т‘Р СР С‘Р Р… РІР‚вЂќ РЎРѓР ВµРЎРѓРЎРѓР С‘РЎРЏ + РЎР‚Р С•Р В»РЎРЉ admin) ---

// SubcategoryFindAll
// @Summary Р РЋР С—Р С‘РЎРѓР С•Р С” Р С—Р С•Р Т‘Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘Р в„–
// @Tags subcategory
// @Produce json
// @Success 200 {array} object
// @Router /subcategory/find-all [get]
func _swaggerSubcategoryFindAll() {}

// SubcategoryFindByID
// @Summary Р СџР С•Р Т‘Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘РЎРЏ Р С—Р С• id
// @Tags subcategory
// @Produce json
// @Param id path int true "ID"
// @Router /subcategory/find-by-id/{id} [get]
func _swaggerSubcategoryFindByID() {}

// SubcategoryTypeFindAll
// @Summary Р вЂ™РЎРѓР Вµ РЎвЂљР С‘Р С—РЎвЂ№ Р С—Р С•Р Т‘Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘Р в„–
// @Tags subcategory-type
// @Produce json
// @Router /subcategory-type/find-all [get]
func _swaggerSubcategoryTypeFindAll() {}

// SubcategoryTypeFindByID
// @Summary Р СћР С‘Р С— Р С—Р С•Р Т‘Р С”Р В°РЎвЂљР ВµР С–Р С•РЎР‚Р С‘Р С‘ Р С—Р С• id
// @Tags subcategory-type
// @Produce json
// @Param id path int true "ID"
// @Router /subcategory-type/find-by-id/{id} [get]
func _swaggerSubcategoryTypeFindByID() {}

// TypeFieldFindAll
// @Summary Р вЂ™РЎРѓР Вµ РЎвЂ¦Р В°РЎР‚Р В°Р С”РЎвЂљР ВµРЎР‚Р С‘РЎРѓРЎвЂљР С‘Р С”Р С‘ (Р С—Р С•Р В»РЎРЏ РЎвЂљР С‘Р С—Р В°)
// @Tags type-field
// @Produce json
// @Router /type-field/find-all [get]
func _swaggerTypeFieldFindAll() {}

// TypeFieldFindByID
// @Summary Р ТђР В°РЎР‚Р В°Р С”РЎвЂљР ВµРЎР‚Р С‘РЎРѓРЎвЂљР С‘Р С”Р В° Р С—Р С• id
// @Tags type-field
// @Produce json
// @Param id path int true "ID"
// @Router /type-field/find-by-id/{id} [get]
func _swaggerTypeFieldFindByID() {}

// registerSwaggerDocSymbols РІР‚вЂќ РЎРѓРЎРѓРЎвЂ№Р В»Р С”Р С‘ Р Р…Р В° Р В·Р В°Р С–Р В»РЎС“РЎв‚¬Р С”Р С‘ Р С‘ РЎвЂљР С‘Р С—РЎвЂ№ РЎвЂљР С•Р В»РЎРЉР С”Р С• Р Т‘Р В»РЎРЏ swag; Р В±Р ВµР В· РЎРЊРЎвЂљР С•Р С–Р С• gopls/staticcheck Р Р†Р С‘Р Т‘РЎРЏРЎвЂљ Р’В«unusedР’В».
func init() {
	_ = []any{
		_swaggerAddressSuggestions, _swaggerAddressValidate, _swaggerCDEKCities, _swaggerCDEKDeliveryPoints, _swaggerCDEKCalculate,
		_swaggerAuthChangePassword, _swaggerAuthForgot, _swaggerAuthIsAdmin, _swaggerAuthLogout, _swaggerAuthMe,
		_swaggerAuthSignIn, _swaggerAuthSignUp, _swaggerAuthVerifyForgot, _swaggerAuthVerifyMobile,
		_swaggerBannerAllModerate, _swaggerBannerCreate, _swaggerBannerList, _swaggerBannerModerate, _swaggerBannerRandom,
		_swaggerCategoryCreate, _swaggerCategoryDelete, _swaggerCategoryFindAll, _swaggerCategoryFindByID,
		_swaggerCategoryFindBySlug, _swaggerCategoryPath, _swaggerCategoryUpdate,
		_swaggerChatInfo, _swaggerChatList, _swaggerChatMessages, _swaggerChatStart,
		_swaggerHealth, _swaggerLogFindAll,
		_swaggerPaymentCheckStatus, _swaggerPaymentCreate, _swaggerPaymentHistory, _swaggerPaymentNotification,
		_swaggerProductAll, _swaggerProductCard, _swaggerProductCreate, _swaggerProductDelete,
		_swaggerProductModerate, _swaggerProductPatch,
		_swaggerPromotionAdd, _swaggerPromotionAll,
		_swaggerReviewModerate, _swaggerReviewModerateList, _swaggerReviewSend, _swaggerReviewUserReviews,
		_swaggerStatisticsAnalytic, _swaggerStatisticsProducts,
		_swaggerSubcategoryFindAll, _swaggerSubcategoryFindByID, _swaggerSubcategoryTypeFindAll, _swaggerSubcategoryTypeFindByID,
		_swaggerSupportAllTickets, _swaggerSupportAssignTicket, _swaggerSupportCreateTicket, _swaggerSupportGetTicket,
		_swaggerSupportMyTickets, _swaggerSupportSendMessage, _swaggerSupportStats, _swaggerSupportUpdateTicket,
		_swaggerTypeFieldFindAll, _swaggerTypeFieldFindByID,
		_swaggerUserAdminDelete, _swaggerUserAdminPatch, _swaggerUserFindAll, _swaggerUserInfo, _swaggerUserRemainingFreeAds,
		_swaggerUserSetBalance, _swaggerUserShowNumber, _swaggerUserToggleBanned, _swaggerUserUpdateSettings,
		_swaggerUserVerifyEmail, _swaggerUserVerifyEmailCode,
		swaggerAddPromotion{}, swaggerAdminUpdateUser{}, swaggerChangePassword{}, swaggerCheckPayment{},
		swaggerCreateCategory{}, swaggerCreatePayment{}, swaggerCreateSupportTicket{}, swaggerForgotEmail{},
		swaggerPaymentCheckStateResponse{}, swaggerPaymentCreateResponse{}, swaggerPaymentHistoryItem{},
		swaggerPaymentNotifyResponse{}, swaggerSendReview{}, swaggerSignIn{}, swaggerSignUp{}, swaggerStartChat{},
		swaggerSupportMessage{}, swaggerTinkoffNotification{}, swaggerUpdateCategory{}, swaggerUpdateSupportTicket{},
		swaggerValidateAddress{},
	}
}

// --- knowledge-base ---

// KnowledgeBaseList
// @Summary Р РЋР С—Р С‘РЎРѓР С•Р С” РЎРѓРЎвЂљР В°РЎвЂљР ВµР в„– Р В±Р В°Р В·РЎвЂ№ Р В·Р Р…Р В°Р Р…Р С‘Р в„–
// @Tags knowledge-base
// @Produce json
// @Success 200 {array} swaggerKnowledgeBaseArticle
// @Router /knowledge-base/ [get]
func _swaggerKnowledgeBaseList() {}

// KnowledgeBaseGetByID
// @Summary Р РЋРЎвЂљР В°РЎвЂљРЎРЉРЎРЏ Р В±Р В°Р В·РЎвЂ№ Р В·Р Р…Р В°Р Р…Р С‘Р в„– Р С—Р С• id
// @Tags knowledge-base
// @Produce json
// @Param id path int true "ID РЎРѓРЎвЂљР В°РЎвЂљРЎРЉР С‘"
// @Success 200 {object} swaggerKnowledgeBaseArticle
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /knowledge-base/{id} [get]
func _swaggerKnowledgeBaseGetByID() {}

// KnowledgeBaseCreate
// @Summary Р РЋР С•Р В·Р Т‘Р В°РЎвЂљРЎРЉ РЎРѓРЎвЂљР В°РЎвЂљРЎРЉРЎР‹ Р В±Р В°Р В·РЎвЂ№ Р В·Р Р…Р В°Р Р…Р С‘Р в„–
// @Tags knowledge-base-admin
// @Accept json
// @Produce json
// @Param body body swaggerKnowledgeBaseArticleRequest true "Р СћР ВµР В»Р С•"
// @Success 201 {object} swaggerKnowledgeBaseCreateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /knowledge-base/ [post]
func _swaggerKnowledgeBaseCreate() {}

// KnowledgeBaseUpdate
// @Summary Р С›Р В±Р Р…Р С•Р Р†Р С‘РЎвЂљРЎРЉ РЎРѓРЎвЂљР В°РЎвЂљРЎРЉРЎР‹ Р В±Р В°Р В·РЎвЂ№ Р В·Р Р…Р В°Р Р…Р С‘Р в„–
// @Tags knowledge-base-admin
// @Accept json
// @Produce json
// @Param id path int true "ID РЎРѓРЎвЂљР В°РЎвЂљРЎРЉР С‘"
// @Param body body swaggerKnowledgeBaseArticleRequest true "Р СћР ВµР В»Р С•"
// @Success 200 {object} swaggerKnowledgeBaseUpdateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /knowledge-base/{id} [put]
func _swaggerKnowledgeBaseUpdate() {}

// KnowledgeBaseDelete
// @Summary Р Р€Р Т‘Р В°Р В»Р С‘РЎвЂљРЎРЉ РЎРѓРЎвЂљР В°РЎвЂљРЎРЉРЎР‹ Р В±Р В°Р В·РЎвЂ№ Р В·Р Р…Р В°Р Р…Р С‘Р в„–
// @Tags knowledge-base-admin
// @Produce json
// @Param id path int true "ID РЎРѓРЎвЂљР В°РЎвЂљРЎРЉР С‘"
// @Success 200 {object} swaggerKnowledgeBaseDeleteResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /knowledge-base/{id} [delete]
func _swaggerKnowledgeBaseDelete() {}

type swaggerKnowledgeBaseArticle struct {
	ID        int32  `json:"id" example:"1"`
	Title     string `json:"title" example:"Р С™Р В°Р С” Р С•РЎвЂћР С•РЎР‚Р СР С‘РЎвЂљРЎРЉ Р В·Р В°Р С”Р В°Р В·"`
	Content   string `json:"content" example:"Р СћР ВµР С”РЎРѓРЎвЂљ РЎРѓРЎвЂљР В°РЎвЂљРЎРЉР С‘..."`
	CreatedAt string `json:"createdAt" example:"2026-03-31T10:00:00Z"`
	UpdatedAt string `json:"updatedAt" example:"2026-03-31T10:00:00Z"`
}

type swaggerKnowledgeBaseArticleRequest struct {
	Title   string `json:"title" example:"Р С™Р В°Р С” Р С•РЎвЂћР С•РЎР‚Р СР С‘РЎвЂљРЎРЉ Р В·Р В°Р С”Р В°Р В·"`
	Content string `json:"content" example:"Р СћР ВµР С”РЎРѓРЎвЂљ РЎРѓРЎвЂљР В°РЎвЂљРЎРЉР С‘..."`
}

type swaggerKnowledgeBaseCreateResponse struct {
	Message string                      `json:"message" example:"Р РЋРЎвЂљР В°РЎвЂљРЎРЉРЎРЏ РЎС“РЎРѓР С—Р ВµРЎв‚¬Р Р…Р С• РЎРѓР С•Р В·Р Т‘Р В°Р Р…Р В°"`
	Article swaggerKnowledgeBaseArticle `json:"article"`
}

type swaggerKnowledgeBaseUpdateResponse struct {
	Message string                      `json:"message" example:"Р РЋРЎвЂљР В°РЎвЂљРЎРЉРЎРЏ РЎС“РЎРѓР С—Р ВµРЎв‚¬Р Р…Р С• Р С•Р В±Р Р…Р С•Р Р†Р В»Р ВµР Р…Р В°"`
	Article swaggerKnowledgeBaseArticle `json:"article"`
}

type swaggerKnowledgeBaseDeleteResponse struct {
	Message string `json:"message" example:"Р РЋРЎвЂљР В°РЎвЂљРЎРЉРЎРЏ РЎС“РЎРѓР С—Р ВµРЎв‚¬Р Р…Р С• РЎС“Р Т‘Р В°Р В»Р ВµР Р…Р В°"`
}

// ModerationList
// @Summary РЎРїРёСЃРѕРє С‚РѕРІР°СЂРѕРІ AI-РјРѕРґРµСЂР°С†РёРё
// @Tags moderation-admin
// @Produce json
// @Param filter query string false "ALL|DENIED|MANUAL|APPROVED_AI"
// @Param page query int false "РќРѕРјРµСЂ СЃС‚СЂР°РЅРёС†С‹"
// @Success 200 {object} swaggerModerationListResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/products [get]
func _swaggerModerationList() {}

// ModerationGetProduct
// @Summary Р”РµС‚Р°Р»Рё С‚РѕРІР°СЂР° РёР· AI-РјРѕРґРµСЂР°С†РёРё
// @Tags moderation-admin
// @Produce json
// @Param id path int true "ID С‚РѕРІР°СЂР°"
// @Success 200 {object} swaggerModerationProductDetail
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /admin/moderation/products/{id} [get]
func _swaggerModerationGetProduct() {}

type swaggerModerationListResponse struct {
	Items []swaggerModerationListItem `json:"items"`
	Total int                         `json:"total" example:"12"`
	Page  int                         `json:"page" example:"1"`
	Pages int                         `json:"pages" example:"1"`
}

type swaggerModerationListItem struct {
	ID                        int32                    `json:"id" example:"1000001"`
	Name                      string                   `json:"name" example:"РўРѕРЅРѕРјРµС‚СЂ"`
	Price                     int32                    `json:"price" example:"3500"`
	Images                    []string                 `json:"images"`
	ModerateState             string                   `json:"moderateState" example:"AI_REVIEWED"`
	ModerationRejectionReason *string                  `json:"moderationRejectionReason,omitempty" example:"РўРµРєСЃС‚: РµСЃС‚СЊ РєРѕРЅС‚Р°РєС‚С‹"`
	CreatedAt                 string                   `json:"createdAt" example:"2026-04-06T10:00:00Z"`
	UpdatedAt                 string                   `json:"updatedAt" example:"2026-04-06T10:00:00Z"`
	Category                  swaggerModerationRefItem `json:"category"`
	SubCategory               swaggerModerationRefItem `json:"subCategory"`
	User                      swaggerModerationUser    `json:"user"`
}

type swaggerModerationProductDetail struct {
	ID                        int32                         `json:"id" example:"1000001"`
	Name                      string                        `json:"name" example:"РўРѕРЅРѕРјРµС‚СЂ"`
	Price                     int32                         `json:"price" example:"3500"`
	Description               string                        `json:"description" example:"РћРїРёСЃР°РЅРёРµ С‚РѕРІР°СЂР°"`
	Images                    []string                      `json:"images"`
	VideoURL                  *string                       `json:"videoUrl,omitempty" example:"https://example.com/video.mp4"`
	ModerateState             string                        `json:"moderateState" example:"AI_REVIEWED"`
	ModerationRejectionReason *string                       `json:"moderationRejectionReason,omitempty" example:"Р¤РѕС‚Рѕ: С‚СЂРµР±СѓРµС‚СЃСЏ СЂСѓС‡РЅР°СЏ РїСЂРѕРІРµСЂРєР°"`
	CreatedAt                 string                        `json:"createdAt" example:"2026-04-06T10:00:00Z"`
	UpdatedAt                 string                        `json:"updatedAt" example:"2026-04-06T10:00:00Z"`
	Category                  swaggerModerationRefItem      `json:"category"`
	SubCategory               swaggerModerationRefItem      `json:"subCategory"`
	Type                      *swaggerModerationTypeRefItem `json:"type,omitempty"`
	User                      swaggerModerationUserDetail   `json:"user"`
	FieldValues               []swaggerModerationFieldValue `json:"fieldValues"`
}

type swaggerModerationRefItem struct {
	ID   int32  `json:"id" example:"1"`
	Name string `json:"name" example:"РњРµРґС‚РµС…РЅРёРєР°"`
}

type swaggerModerationTypeRefItem struct {
	ID   int32   `json:"id" example:"1"`
	Name *string `json:"name" example:"РўРѕРЅРѕРјРµС‚СЂС‹"`
}

type swaggerModerationUser struct {
	ID          int32  `json:"id" example:"1"`
	FullName    string `json:"fullName" example:"РРІР°РЅ РРІР°РЅРѕРІ"`
	Email       string `json:"email" example:"ivan@example.com"`
	PhoneNumber string `json:"phoneNumber" example:"+79990000000"`
}

type swaggerModerationUserDetail struct {
	ID          int32  `json:"id" example:"1"`
	FullName    string `json:"fullName" example:"РРІР°РЅ РРІР°РЅРѕРІ"`
	Email       string `json:"email" example:"ivan@example.com"`
	PhoneNumber string `json:"phoneNumber" example:"+79990000000"`
	ProfileType string `json:"profileType" example:"INDIVIDUAL"`
}

type swaggerModerationFieldValue struct {
	Value string                    `json:"value" example:"Omron"`
	Field swaggerModerationFieldRef `json:"field"`
}

type swaggerModerationFieldRef struct {
	ID   int32  `json:"id" example:"1"`
	Name string `json:"name" example:"РџСЂРѕРёР·РІРѕРґРёС‚РµР»СЊ"`
}

// --- cdek ---

// CDEKCities
// @Summary Search CDEK cities
// @Tags cdek
// @Produce json
// @Param city query string true "City name"
// @Param limit query int false "Limit" default(20)
// @Success 200 {array} object
// @Failure 400 {object} map[string]interface{}
// @Router /cdek/cities [get]
func _swaggerCDEKCities() {}

// CDEKDeliveryPoints
// @Summary Get CDEK delivery points (PVZ) by city code
// @Tags cdek
// @Produce json
// @Param cityCode query int true "CDEK city code"
// @Success 200 {array} object
// @Failure 400 {object} map[string]interface{}
// @Router /cdek/delivery-points [get]
func _swaggerCDEKDeliveryPoints() {}

// CDEKCalculate
// @Summary Calculate CDEK delivery tariff
// @Tags cdek
// @Accept json
// @Produce json
// @Param body body swaggerCDEKCalculateRequest true "Body"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /cdek/calculate [post]
func _swaggerCDEKCalculate() {}

type swaggerCDEKCalculateRequest struct {
	TariffCode   int `json:"tariffCode" example:"136"`
	FromCityCode int `json:"fromCityCode" example:"44"`
	ToCityCode   int `json:"toCityCode" example:"270"`
	Weight       int `json:"weight" example:"1000"`
	Length       int `json:"length" example:"20"`
	Width        int `json:"width" example:"20"`
	Height       int `json:"height" example:"20"`
}
