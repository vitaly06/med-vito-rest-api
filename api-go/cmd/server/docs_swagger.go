package main

// РћРїРёСЃР°РЅРёРµ РјР°СЂС€СЂСѓС‚РѕРІ С‚РѕР»СЊРєРѕ РґР»СЏ РіРµРЅРµСЂР°С†РёРё OpenAPI (swag). Р РµР°Р»СЊРЅС‹Рµ С…РµРЅРґР»РµСЂС‹ РІ internal/httpserver.
// РђРґРјРёРЅ-СЂСѓС‡РєРё РєР°С‚РµРіРѕСЂРёР№: СЃРЅР°С‡Р°Р»Р° POST /auth/sign-in РїРѕРґ РїРѕР»СЊР·РѕРІР°С‚РµР»РµРј СЃ СЂРѕР»СЊСЋ admin вЂ” cookie session_id СѓР№РґС‘С‚ РІ Try it out (credentials).

// --- system ---

// HealthCheck
// @Summary РџСЂРѕРІРµСЂРєР° Р¶РёРІРѕСЃС‚Рё СЃРµСЂРІРёСЃР°
// @Description Р’РѕР·РІСЂР°С‰Р°РµС‚ status ok
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func _swaggerHealth() {}

// --- log ---

// LogFindAll
// @Summary РЎРїРёСЃРѕРє Р·Р°РїРёСЃРµР№ Log
// @Tags log
// @Produce json
// @Success 200 {array} object
// @Router /log/find-all [get]
func _swaggerLogFindAll() {}

// --- category (РїСѓР±Р»РёС‡РЅРѕ) ---

// CategoryFindAll
// @Summary Р”РµСЂРµРІРѕ РІСЃРµС… РєР°С‚РµРіРѕСЂРёР№
// @Tags category
// @Produce json
// @Success 200 {array} object
// @Router /category/find-all [get]
func _swaggerCategoryFindAll() {}

// CategoryFindByID
// @Summary РљР°С‚РµРіРѕСЂРёСЏ РїРѕ id (РїРѕР»РЅРѕРµ РґРµСЂРµРІРѕ)
// @Tags category
// @Produce json
// @Param id path int true "ID РєР°С‚РµРіРѕСЂРёРё"
// @Success 200 {object} object
// @Failure 404 {object} map[string]interface{}
// @Router /category/find-by-id/{id} [get]
func _swaggerCategoryFindByID() {}

// CategoryFindBySlug
// @Summary РљР°С‚РµРіРѕСЂРёСЏ РїРѕ slug
// @Tags category
// @Produce json
// @Param slug path string true "Slug"
// @Success 200 {object} object
// @Failure 404 {object} map[string]interface{}
// @Router /category/slug/{slug} [get]
func _swaggerCategoryFindBySlug() {}

// CategoryFindBySlugPath
// @Summary Р Р°Р·СЂРµС€РµРЅРёРµ С†РµРїРѕС‡РєРё slug (РєР°С‚РµРіРѕСЂРёСЏ / РїРѕРґРєР°С‚РµРіРѕСЂРёСЏ / С‚РёРї)
// @Description Р’ Swagger В«Try it outВ» РІРІРµРґРё СЃРµРіРјРµРЅС‚С‹ С‡РµСЂРµР· %2F, РЅР°РїСЂРёРјРµСЂ: elektronika%2Ftelefony
// @Tags category
// @Produce json
// @Param slugPath path string true "Р¦РµРїРѕС‡РєР° (РёР»Рё РѕРґРёРЅ СЃРµРіРјРµРЅС‚)"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /category/path/{slugPath} [get]
func _swaggerCategoryPath() {}

// --- category (Р°РґРјРёРЅ; РїСЂРё ADMIN_API_KEY вЂ” Р·Р°РіРѕР»РѕРІРѕРє X-Admin-Key) ---

// CategoryCreate
// @Summary РЎРѕР·РґР°С‚СЊ РєР°С‚РµРіРѕСЂРёСЋ
// @Tags category-admin
// @Accept json
// @Produce json
// @Param body body swaggerCreateCategory true "РўРµР»Рѕ"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /category/create-category [post]
func _swaggerCategoryCreate() {}

// CategoryUpdate
// @Summary РћР±РЅРѕРІРёС‚СЊ РєР°С‚РµРіРѕСЂРёСЋ
// @Tags category-admin
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body swaggerUpdateCategory true "РўРµР»Рѕ"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /category/update-category/{id} [put]
func _swaggerCategoryUpdate() {}

// CategoryDelete
// @Summary РЈРґР°Р»РёС‚СЊ РєР°С‚РµРіРѕСЂРёСЋ
// @Tags category-admin
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]interface{}
// @Router /category/delete-category/{id} [delete]
func _swaggerCategoryDelete() {}

// РўРµР»Р° Р·Р°РїСЂРѕСЃРѕРІ РґР»СЏ Swagger UI
type swaggerCreateCategory struct {
	Name string  `json:"name" example:"РђРІС‚РѕРјРѕР±РёР»Рё"`
	Slug *string `json:"slug,omitempty" example:"avtomobili"`
}

type swaggerUpdateCategory struct {
	Name string  `json:"name" example:"РђРІС‚РѕРјРѕР±РёР»Рё"`
	Slug *string `json:"slug,omitempty"`
}

// --- auth ---

// AuthSignUp
// @Summary Р РµРіРёСЃС‚СЂР°С†РёСЏ вЂ” РѕС‚РїСЂР°РІРєР° РєРѕРґР° (query where=telegram|sms)
// @Tags auth
// @Accept json
// @Produce json
// @Param where query string true "telegram РёР»Рё sms" Enums(telegram,sms)
// @Param body body swaggerSignUp true "Р”Р°РЅРЅС‹Рµ"
// @Success 200 {object} map[string]string
// @Router /auth/sign-up [post]
func _swaggerAuthSignUp() {}

// AuthVerifyMobile
// @Summary РџРѕРґС‚РІРµСЂР¶РґРµРЅРёРµ С‚РµР»РµС„РѕРЅР° РїРѕ РєРѕРґСѓ
// @Tags auth
// @Produce json
// @Param code query string true "РљРѕРґ РёР· SMS/TG"
// @Success 200 {object} map[string]string
// @Router /auth/verify-mobile-code [post]
func _swaggerAuthVerifyMobile() {}

// AuthSignIn
// @Summary Р’С…РѕРґ (СЃС‚Р°РІРёС‚ cookie session_id)
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerSignIn true "Р›РѕРіРёРЅ"
// @Success 200 {object} object
// @Router /auth/sign-in [post]
func _swaggerAuthSignIn() {}

// AuthMe
// @Summary РўРµРєСѓС‰РёР№ РїРѕР»СЊР·РѕРІР°С‚РµР»СЊ
// @Tags auth
// @Produce json
// @Success 200 {object} object
// @Router /auth/me [get]
func _swaggerAuthMe() {}

// AuthIsAdmin
// @Summary РџСЂРѕРІРµСЂРєР° СЂРѕР»Рё admin
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]bool
// @Router /auth/isAdmin [get]
func _swaggerAuthIsAdmin() {}

// AuthLogout
// @Summary Р’С‹С…РѕРґ
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func _swaggerAuthLogout() {}

// AuthForgot
// @Summary Р—Р°РїСЂРѕСЃ РєРѕРґР° СЃР±СЂРѕСЃР° РЅР° РїРѕС‡С‚Сѓ
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerForgotEmail true "email"
// @Success 200 {object} map[string]string
// @Router /auth/forgot-password [post]
func _swaggerAuthForgot() {}

// AuthVerifyForgot
// @Summary РџСЂРѕРІРµСЂРєР° РєРѕРґР° СЃР±СЂРѕСЃР°
// @Tags auth
// @Produce json
// @Param code query string true "РљРѕРґ РёР· РїРёСЃСЊРјР°"
// @Success 200 {object} map[string]int
// @Router /auth/verify-code [post]
func _swaggerAuthVerifyForgot() {}

// AuthChangePassword
// @Summary РќРѕРІС‹Р№ РїР°СЂРѕР»СЊ РїРѕСЃР»Рµ verify-code
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerChangePassword true "РўРµР»Рѕ"
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
// @Summary РЎРїРёСЃРѕРє РїРѕР»СЊР·РѕРІР°С‚РµР»РµР№ (Р°РґРјРёРЅ: cookie session_id + СЂРѕР»СЊ admin)
// @Tags user-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /user/find-all [get]
func _swaggerUserFindAll() {}

// UserInfo
// @Summary РљР°СЂС‚РѕС‡РєР° РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ (СЂРµР№С‚РёРЅРі, Р»РёРјРёС‚ РѕР±СЉСЏРІР»РµРЅРёР№)
// @Tags user
// @Produce json
// @Param id path int true "User id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /user/info/{id} [get]
func _swaggerUserInfo() {}

// UserRemainingFreeAds
// @Summary РћСЃС‚Р°С‚РѕРє Р±РµСЃРїР»Р°С‚РЅС‹С… РѕР±СЉСЏРІР»РµРЅРёР№ (СЃРµСЃСЃРёСЏ)
// @Tags user
// @Produce json
// @Success 200 {object} object
// @Router /user/remaining-free-ads [get]
func _swaggerUserRemainingFreeAds() {}

// UserShowNumber
// @Summary РџРѕРєР°Р·Р°С‚СЊ РЅРѕРјРµСЂ РїСЂРѕРґР°РІС†Р° (СЃРµСЃСЃРёСЏ)
// @Tags user
// @Produce json
// @Param userId path int true "РџСЂРѕРґР°РІРµС†"
// @Success 200 {object} map[string]string
// @Router /user/show-number/{userId} [get]
func _swaggerUserShowNumber() {}

// UserUpdateSettings
// @Summary РћР±РЅРѕРІР»РµРЅРёРµ РЅР°СЃС‚СЂРѕРµРє (multipart, СЃРµСЃСЃРёСЏ)
// @Tags user
// @Accept mpfd
// @Produce json
// @Param fullName formData string false "Р¤РРћ"
// @Param phoneNumber formData string false "РўРµР»РµС„РѕРЅ"
// @Param isAnswersCall formData string false "true/false"
// @Param profileType formData string false "INDIVIDUAL|OOO|IP"
// @Param photo formData file false "РђРІР°С‚Р°СЂ"
// @Success 200 {object} object
// @Router /user/update-settings [patch]
func _swaggerUserUpdateSettings() {}

// UserVerifyEmail
// @Summary РћС‚РїСЂР°РІРёС‚СЊ РєРѕРґ РїРѕРґС‚РІРµСЂР¶РґРµРЅРёСЏ РЅР° РїРѕС‡С‚Сѓ (СЃРµСЃСЃРёСЏ)
// @Tags user
// @Produce json
// @Success 200 {object} map[string]string
// @Router /user/verify-email [post]
func _swaggerUserVerifyEmail() {}

// UserVerifyEmailCode
// @Summary РџРѕРґС‚РІРµСЂРґРёС‚СЊ РїРѕС‡С‚Сѓ РїРѕ РєРѕРґСѓ РёР· РїРёСЃСЊРјР°
// @Tags user
// @Produce json
// @Param code query string true "РљРѕРґ"
// @Success 200 {object} map[string]string
// @Router /user/verify-code [post]
func _swaggerUserVerifyEmailCode() {}

// UserSetBalance
// @Summary РЈСЃС‚Р°РЅРѕРІРёС‚СЊ bonusBalance (Р°РґРјРёРЅ, СЃРµСЃСЃРёСЏ)
// @Tags user-admin
// @Produce json
// @Param userId path int true "User id"
// @Param balance query string true "Р§РёСЃР»Рѕ"
// @Success 200 {object} map[string]string
// @Router /user/set-balance/{userId} [put]
func _swaggerUserSetBalance() {}

// UserToggleBanned
// @Summary Р‘Р°РЅ / СЂР°Р·Р±Р°РЅ (Р°РґРјРёРЅ, СЃРµСЃСЃРёСЏ)
// @Tags user-admin
// @Produce json
// @Param id path int true "User id"
// @Success 200 {object} map[string]string
// @Router /user/toggle-banned/{id} [put]
func _swaggerUserToggleBanned() {}

// UserAdminPatch
// @Summary РћР±РЅРѕРІРёС‚СЊ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ (Р°РґРјРёРЅ, СЃРµСЃСЃРёСЏ)
// @Tags user-admin
// @Accept json
// @Produce json
// @Param id path int true "User id"
// @Param body body swaggerAdminUpdateUser true "РџРѕР»СЏ"
// @Success 200 {object} map[string]string
// @Router /user/{id} [patch]
func _swaggerUserAdminPatch() {}

// UserAdminDelete
// @Summary РЈРґР°Р»РёС‚СЊ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ (Р°РґРјРёРЅ, СЃРµСЃСЃРёСЏ)
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
// @Summary РЎРѕР·РґР°С‚СЊ С‚РѕРІР°СЂ (multipart, РґРѕ 8 images, СЃРµСЃСЃРёСЏ)
// @Tags product
// @Accept mpfd
// @Produce json
// @Success 201 {object} object
// @Router /product/create [post]
func _swaggerProductCreate() {}

// ProductAll
// @Summary РЎРїРёСЃРѕРє С‚РѕРІР°СЂРѕРІ / РїРѕРёСЃРє (query; optional СЃРµСЃСЃРёСЏ РґР»СЏ РёР·Р±СЂР°РЅРЅРѕРіРѕ)
// @Tags product
// @Produce json
// @Router /product/all-products [get]
func _swaggerProductAll() {}

// ProductCard
// @Summary РљР°СЂС‚РѕС‡РєР° С‚РѕРІР°СЂР°
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Router /product/product-card/{id} [get]
func _swaggerProductCard() {}

// ProductDelete
// @Summary РЈРґР°Р»РёС‚СЊ СЃРІРѕР№ С‚РѕРІР°СЂ
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Router /product/{id} [delete]
func _swaggerProductDelete() {}

// ProductPatch
// @Summary РћР±РЅРѕРІРёС‚СЊ С‚РѕРІР°СЂ (multipart)
// @Tags product
// @Accept mpfd
// @Produce json
// @Param id path int true "Product id"
// @Router /product/{id} [patch]
func _swaggerProductPatch() {}

// ProductModerate
// @Summary РњРѕРґРµСЂР°С†РёСЏ (admin)
// @Tags product-admin
// @Produce json
// @Param id path int true "Product id"
// @Param status query string true "APPROVED|DENIDED"
// @Param reason query string false "РџСЂРёС‡РёРЅР° РїСЂРё DENIDED"
// @Router /product/moderate-product/{id} [put]
func _swaggerProductModerate() {}

// --- review ---

// ReviewSend
// @Summary РћСЃС‚Р°РІРёС‚СЊ РѕС‚Р·С‹РІ РїСЂРѕРґР°РІС†Сѓ (СЃРµСЃСЃРёСЏ)
// @Tags review
// @Accept json
// @Produce json
// @Param body body swaggerSendReview true "РўРµР»Рѕ"
// @Success 200 {object} map[string]string
// @Router /review/send-review [post]
func _swaggerReviewSend() {}

// ReviewUserReviews
// @Summary РћРґРѕР±СЂРµРЅРЅС‹Рµ РѕС‚Р·С‹РІС‹ Рѕ РїРѕР»СЊР·РѕРІР°С‚РµР»Рµ
// @Tags review
// @Produce json
// @Param id path int true "User id"
// @Router /review/user-reviews/{id} [get]
func _swaggerReviewUserReviews() {}

// ReviewModerate
// @Summary РњРѕРґРµСЂР°С†РёСЏ РѕС‚Р·С‹РІР° (admin)
// @Tags review-admin
// @Produce json
// @Param id path int true "Review id"
// @Param status query string true "APPROVED|DENIDED"
// @Router /review/moderate-review/{id} [put]
func _swaggerReviewModerate() {}

// ReviewModerateList
// @Summary РћС‡РµСЂРµРґСЊ РѕС‚Р·С‹РІРѕРІ РЅР° РјРѕРґРµСЂР°С†РёСЋ (admin)
// @Tags review-admin
// @Produce json
// @Router /review/all-reviews-to-moderate [get]
func _swaggerReviewModerateList() {}

type swaggerSendReview struct {
	Text           *string `json:"text,omitempty"`
	Rating         float64 `json:"rating" example:"5"`
	ReviewedUserID int32   `json:"reviewedUserId"`
}

// --- chat (СЃРµСЃСЃРёСЏ cookie session_id) ---

// ChatStart
// @Summary РќР°С‡Р°С‚СЊ С‡Р°С‚ РїРѕ С‚РѕРІР°СЂСѓ
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
// @Summary РЎРїРёСЃРѕРє С‡Р°С‚РѕРІ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ
// @Tags chat
// @Produce json
// @Success 200 {array} object
// @Router /chat [get]
func _swaggerChatList() {}

// ChatMessages
// @Summary РЎРѕРѕР±С‰РµРЅРёСЏ С‡Р°С‚Р° (РїР°РіРёРЅР°С†РёСЏ)
// @Tags chat
// @Produce json
// @Param id path int true "Chat id"
// @Param page query int false "РЎС‚СЂР°РЅРёС†Р°" default(1)
// @Param limit query int false "Р›РёРјРёС‚" default(50)
// @Success 200 {object} object
// @Router /chat/{id}/messages [get]
func _swaggerChatMessages() {}

// ChatInfo
// @Summary РРЅС„РѕСЂРјР°С†РёСЏ Рѕ С‡Р°С‚Рµ
// @Tags chat
// @Produce json
// @Param id path int true "Chat id"
// @Success 200 {object} object
// @Router /chat/{id} [get]
func _swaggerChatInfo() {}

type swaggerStartChat struct {
	ProductID int32 `json:"productId" example:"1"`
}

// --- payment (Рў-Р‘Р°РЅРє / Tinkoff; СЃРµСЃСЃРёСЏ вЂ” cookie session_id) ---

// PaymentCreate
// @Summary РЎРѕР·РґР°РЅРёРµ РїР»Р°С‚РµР¶Р° РґР»СЏ РїРѕРїРѕР»РЅРµРЅРёСЏ Р±Р°Р»Р°РЅСЃР°
// @Description Init РІ Рў-Р‘Р°РЅРє. РќСѓР¶РЅС‹ TINKOFF_TERMINAL_KEY Рё TINKOFF_SECRET_KEY. РђРІС‚РѕСЂРёР·Р°С†РёСЏ: cookie session_id РїРѕСЃР»Рµ POST /auth/sign-in.
// @Tags payment
// @Accept json
// @Produce json
// @Param body body swaggerCreatePayment true "РЎСѓРјРјР° РІ СЂСѓР±Р»СЏС… (РјРёРЅ. 1)"
// @Success 201 {object} swaggerPaymentCreateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /payment/create [post]
func _swaggerPaymentCreate() {}

// PaymentNotification
// @Summary Webhook СѓРІРµРґРѕРјР»РµРЅРёР№ Рў-Р‘Р°РЅРєР° Рѕ СЃС‚Р°С‚СѓСЃРµ РїР»Р°С‚РµР¶Р°
// @Description Р‘РµР· СЃРµСЃСЃРёРё. РџРѕРґРїРёСЃСЊ Token РїСЂРѕРІРµСЂСЏРµС‚СЃСЏ РїРѕ РїРѕР»СЏРј С‚РµР»Р°. РўРµР»Рѕ вЂ” РєР°Рє РїСЂРёС…РѕРґРёС‚ РѕС‚ Р±Р°РЅРєР°; РїСЂРёРјРµСЂ РЅРёР¶Рµ.
// @Tags payment
// @Accept json
// @Produce json
// @Param body body swaggerTinkoffNotification true "РЈРІРµРґРѕРјР»РµРЅРёРµ"
// @Success 200 {object} swaggerPaymentNotifyResponse
// @Failure 400 {object} map[string]interface{}
// @Router /payment/notification [post]
func _swaggerPaymentNotification() {}

// PaymentHistory
// @Summary РСЃС‚РѕСЂРёСЏ РїР»Р°С‚РµР¶РµР№ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ
// @Description Р”Рѕ 50 Р·Р°РїРёСЃРµР№, РЅРѕРІС‹Рµ СЃРІРµСЂС…Сѓ. РЎРµСЃСЃРёСЏ: cookie session_id.
// @Tags payment
// @Produce json
// @Success 200 {array} swaggerPaymentHistoryItem
// @Failure 401 {object} map[string]interface{}
// @Router /payment/history [get]
func _swaggerPaymentHistory() {}

// PaymentCheckStatus
// @Summary РџСЂРѕРІРµСЂРєР° СЃС‚Р°С‚СѓСЃР° РїР»Р°С‚РµР¶Р° РІ Рў-Р‘Р°РЅРєРµ (GetState)
// @Description РЎРµСЃСЃРёСЏ: cookie session_id. Р’ С‚РµР»Рµ вЂ” paymentId РёР· РѕС‚РІРµС‚Р° Init РёР»Рё СѓРІРµРґРѕРјР»РµРЅРёСЏ.
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
	Description *string `json:"description,omitempty" example:"РџРѕРїРѕР»РЅРµРЅРёРµ Р±Р°Р»Р°РЅСЃР°"`
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
	Message string `json:"message,omitempty" example:"Р‘Р°Р»Р°РЅСЃ СѓСЃРїРµС€РЅРѕ РїРѕРїРѕР»РЅРµРЅ"`
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

// swaggerTinkoffNotification вЂ” РїРѕР»СЏ РєР°Рє РІ Nest PaymentNotificationDto (СЂРµР°Р»СЊРЅС‹Р№ webhook РјРѕР¶РµС‚ РґРѕР±Р°РІР»СЏС‚СЊ РїРѕР»СЏ).
type swaggerTinkoffNotification struct {
	TerminalKey string `json:"TerminalKey" example:"1766153689307DEMO"`
	OrderID     string `json:"OrderId" example:"123-1735123456789"`
	Success     bool   `json:"Success" example:"true"`
	Status      string `json:"Status" example:"CONFIRMED"`
	PaymentID   string `json:"PaymentId" example:"2673412345"`
	Amount      int64  `json:"Amount" example:"100000"`
	Token       string `json:"Token" example:"РїРѕРґРїРёСЃСЊ_РѕС‚_Р±Р°РЅРєР°"`
	ErrorCode   string `json:"ErrorCode,omitempty" example:"0"`
	Pan         string `json:"Pan,omitempty" example:"430000******0777"`
}

// --- promotion ---

// PromotionAll
// @Summary Р’СЃРµ С‚РёРїС‹ РїСЂРѕРґРІРёР¶РµРЅРёСЏ (С‚Р°СЂРёС„С‹)
// @Tags promotion
// @Produce json
// @Success 200 {array} object
// @Router /promotion/all-promotions [get]
func _swaggerPromotionAll() {}

// PromotionAdd
// @Summary РџРѕРґРєР»СЋС‡РёС‚СЊ РїСЂРѕРґРІРёР¶РµРЅРёРµ Рє С‚РѕРІР°СЂСѓ (СЃРµСЃСЃРёСЏ)
// @Tags promotion
// @Accept json
// @Produce json
// @Param body body swaggerAddPromotion true "РўРµР»Рѕ"
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
// @Summary РЎС‚Р°С‚РёСЃС‚РёРєР° РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ (РїСЂРѕСЃРјРѕС‚СЂС‹, С‚РµР»РµС„РѕРЅ, РёР·Р±СЂР°РЅРЅРѕРµ)
// @Tags statistics
// @Produce json
// @Param period query string false "day week month quarter half-year year"
// @Param categoryId query int false "Р¤РёР»СЊС‚СЂ РїРѕ РєР°С‚РµРіРѕСЂРёРё (С‡РµСЂРµР· SubCategory)"
// @Param region query string false "РџРѕРґСЃС‚СЂРѕРєР° РІ address (ILIKE)"
// @Param productId query int false "РљРѕРЅРєСЂРµС‚РЅС‹Р№ С‚РѕРІР°СЂ"
// @Success 200 {object} object
// @Router /statistics/analytic [get]
func _swaggerStatisticsAnalytic() {}

// StatisticsProducts
// @Summary РђРЅР°Р»РёС‚РёРєР° РїРѕ РєР°Р¶РґРѕРјСѓ С‚РѕРІР°СЂСѓ РїСЂРѕРґР°РІС†Р°
// @Tags statistics
// @Produce json
// @Success 200 {array} object
// @Router /statistics/products-analytic [get]
func _swaggerStatisticsProducts() {}

// --- support (СЃРµСЃСЃРёСЏ) ---

// SupportCreateTicket
// @Summary РЎРѕР·РґР°С‚СЊ С‚РёРєРµС‚ РїРѕРґРґРµСЂР¶РєРё
// @Tags support
// @Accept json
// @Produce json
// @Param body body swaggerCreateSupportTicket true "РўРµР»Рѕ"
// @Success 201 {object} object
// @Router /support/tickets [post]
func _swaggerSupportCreateTicket() {}

// SupportMyTickets
// @Summary РњРѕРё С‚РёРєРµС‚С‹ (РїР°РіРёРЅР°С†РёСЏ, С„РёР»СЊС‚СЂС‹ query)
// @Tags support
// @Produce json
// @Router /support/tickets/my [get]
func _swaggerSupportMyTickets() {}

// SupportAllTickets
// @Summary Р’СЃРµ С‚РёРєРµС‚С‹ (РјРѕРґРµСЂР°С‚РѕСЂ/admin)
// @Tags support
// @Produce json
// @Router /support/tickets/all [get]
func _swaggerSupportAllTickets() {}

// SupportStats
// @Summary РЎС‚Р°С‚РёСЃС‚РёРєР° С‚РёРєРµС‚РѕРІ (С‚РѕР»СЊРєРѕ admin)
// @Tags support
// @Produce json
// @Router /support/stats [get]
func _swaggerSupportStats() {}

// SupportGetTicket
// @Summary РўРёРєРµС‚ СЃ СЃРѕРѕР±С‰РµРЅРёСЏРјРё
// @Tags support
// @Produce json
// @Param id path int true "Ticket id"
// @Router /support/tickets/{id} [get]
func _swaggerSupportGetTicket() {}

// SupportSendMessage
// @Summary РЎРѕРѕР±С‰РµРЅРёРµ РІ С‚РёРєРµС‚
// @Tags support
// @Accept json
// @Param id path int true "Ticket id"
// @Param body body swaggerSupportMessage true "РўРµРєСЃС‚"
// @Router /support/tickets/{id}/messages [post]
func _swaggerSupportSendMessage() {}

// SupportUpdateTicket
// @Summary РћР±РЅРѕРІРёС‚СЊ С‚РёРєРµС‚ (РјРѕРґРµСЂР°С‚РѕСЂ/admin)
// @Tags support
// @Accept json
// @Param id path int true "Ticket id"
// @Param body body swaggerUpdateSupportTicket true "РџРѕР»СЏ"
// @Router /support/tickets/{id} [put]
func _swaggerSupportUpdateTicket() {}

// SupportAssignTicket
// @Summary РќР°Р·РЅР°С‡РёС‚СЊ С‚РёРєРµС‚ РЅР° СЃРµР±СЏ
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

// --- address (DaData, Р±РµР· СЃРµСЃСЃРёРё) ---

// AddressSuggestions
// @Summary РџРѕРґСЃРєР°Р·РєРё Р°РґСЂРµСЃР° (DaData)
// @Tags address
// @Produce json
// @Param query query string true "РЎС‚СЂРѕРєР° РїРѕРёСЃРєР°"
// @Param limit query int false "Р›РёРјРёС‚" default(5)
// @Success 200 {array} object
// @Router /address/suggestions [get]
func _swaggerAddressSuggestions() {}

// AddressValidate
// @Summary РџСЂРѕРІРµСЂРєР° Р°РґСЂРµСЃР° РїРѕ РїРµСЂРІРѕР№ РїРѕРґСЃРєР°Р·РєРµ DaData
// @Tags address
// @Accept json
// @Produce json
// @Param body body swaggerValidateAddress true "РђРґСЂРµСЃ"
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
// @Summary РЎРѕР·РґР°С‚СЊ Р±Р°РЅРЅРµСЂ (multipart: image, name, place, navigateToUrl; СЃРµСЃСЃРёСЏ)
// @Tags banner
// @Accept mpfd
// @Produce json
// @Success 201 {object} object
// @Router /banner [post]
func _swaggerBannerCreate() {}

// BannerRandom
// @Summary РЎР»СѓС‡Р°Р№РЅС‹Рµ РѕРґРѕР±СЂРµРЅРЅС‹Рµ Р±Р°РЅРЅРµСЂС‹
// @Tags banner
// @Produce json
// @Router /banner/random [get]
func _swaggerBannerRandom() {}

// BannerList
// @Summary РЎРїРёСЃРѕРє РѕРґРѕР±СЂРµРЅРЅС‹С… Р±Р°РЅРЅРµСЂРѕРІ (query place РѕРїС†РёРѕРЅР°Р»СЊРЅРѕ)
// @Tags banner
// @Produce json
// @Router /banner [get]
func _swaggerBannerList() {}

// BannerModerate
// @Summary РњРѕРґРµСЂР°С†РёСЏ Р±Р°РЅРЅРµСЂР° (admin, query status)
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Param status query string true "MODERATE|APPROVED|DENIDED"
// @Router /banner/moderate/{id} [put]
func _swaggerBannerModerate() {}

// BannerAllModerate
// @Summary РћС‡РµСЂРµРґСЊ Р±Р°РЅРЅРµСЂРѕРІ РЅР° РјРѕРґРµСЂР°С†РёСЋ (admin)
// @Tags banner
// @Produce json
// @Router /banner/all-banners-to-moderate [get]
func _swaggerBannerAllModerate() {}

// --- subcategory / subcategory-type / type-field (Р°РґРјРёРЅ вЂ” СЃРµСЃСЃРёСЏ + СЂРѕР»СЊ admin) ---

// SubcategoryFindAll
// @Summary РЎРїРёСЃРѕРє РїРѕРґРєР°С‚РµРіРѕСЂРёР№
// @Tags subcategory
// @Produce json
// @Success 200 {array} object
// @Router /subcategory/find-all [get]
func _swaggerSubcategoryFindAll() {}

// SubcategoryFindByID
// @Summary РџРѕРґРєР°С‚РµРіРѕСЂРёСЏ РїРѕ id
// @Tags subcategory
// @Produce json
// @Param id path int true "ID"
// @Router /subcategory/find-by-id/{id} [get]
func _swaggerSubcategoryFindByID() {}

// SubcategoryTypeFindAll
// @Summary Р’СЃРµ С‚РёРїС‹ РїРѕРґРєР°С‚РµРіРѕСЂРёР№
// @Tags subcategory-type
// @Produce json
// @Router /subcategory-type/find-all [get]
func _swaggerSubcategoryTypeFindAll() {}

// SubcategoryTypeFindByID
// @Summary РўРёРї РїРѕРґРєР°С‚РµРіРѕСЂРёРё РїРѕ id
// @Tags subcategory-type
// @Produce json
// @Param id path int true "ID"
// @Router /subcategory-type/find-by-id/{id} [get]
func _swaggerSubcategoryTypeFindByID() {}

// TypeFieldFindAll
// @Summary Р’СЃРµ С…Р°СЂР°РєС‚РµСЂРёСЃС‚РёРєРё (РїРѕР»СЏ С‚РёРїР°)
// @Tags type-field
// @Produce json
// @Router /type-field/find-all [get]
func _swaggerTypeFieldFindAll() {}

// TypeFieldFindByID
// @Summary РҐР°СЂР°РєС‚РµСЂРёСЃС‚РёРєР° РїРѕ id
// @Tags type-field
// @Produce json
// @Param id path int true "ID"
// @Router /type-field/find-by-id/{id} [get]
func _swaggerTypeFieldFindByID() {}

// registerSwaggerDocSymbols вЂ” СЃСЃС‹Р»РєРё РЅР° Р·Р°РіР»СѓС€РєРё Рё С‚РёРїС‹ С‚РѕР»СЊРєРѕ РґР»СЏ swag; Р±РµР· СЌС‚РѕРіРѕ gopls/staticcheck РІРёРґСЏС‚ В«unusedВ».
func init() {
	_ = []any{
		_swaggerAddressSuggestions, _swaggerAddressValidate,
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
// @Summary РЎРїРёСЃРѕРє СЃС‚Р°С‚РµР№ Р±Р°Р·С‹ Р·РЅР°РЅРёР№
// @Tags knowledge-base
// @Produce json
// @Success 200 {array} swaggerKnowledgeBaseArticle
// @Router /knowledge-base/ [get]
func _swaggerKnowledgeBaseList() {}

// KnowledgeBaseGetByID
// @Summary РЎС‚Р°С‚СЊСЏ Р±Р°Р·С‹ Р·РЅР°РЅРёР№ РїРѕ id
// @Tags knowledge-base
// @Produce json
// @Param id path int true "ID СЃС‚Р°С‚СЊРё"
// @Success 200 {object} swaggerKnowledgeBaseArticle
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /knowledge-base/{id} [get]
func _swaggerKnowledgeBaseGetByID() {}

// KnowledgeBaseCreate
// @Summary РЎРѕР·РґР°С‚СЊ СЃС‚Р°С‚СЊСЋ Р±Р°Р·С‹ Р·РЅР°РЅРёР№
// @Tags knowledge-base-admin
// @Accept json
// @Produce json
// @Param body body swaggerKnowledgeBaseArticleRequest true "РўРµР»Рѕ"
// @Success 201 {object} swaggerKnowledgeBaseCreateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /knowledge-base/ [post]
func _swaggerKnowledgeBaseCreate() {}

// KnowledgeBaseUpdate
// @Summary РћР±РЅРѕРІРёС‚СЊ СЃС‚Р°С‚СЊСЋ Р±Р°Р·С‹ Р·РЅР°РЅРёР№
// @Tags knowledge-base-admin
// @Accept json
// @Produce json
// @Param id path int true "ID СЃС‚Р°С‚СЊРё"
// @Param body body swaggerKnowledgeBaseArticleRequest true "РўРµР»Рѕ"
// @Success 200 {object} swaggerKnowledgeBaseUpdateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /knowledge-base/{id} [put]
func _swaggerKnowledgeBaseUpdate() {}

// KnowledgeBaseDelete
// @Summary РЈРґР°Р»РёС‚СЊ СЃС‚Р°С‚СЊСЋ Р±Р°Р·С‹ Р·РЅР°РЅРёР№
// @Tags knowledge-base-admin
// @Produce json
// @Param id path int true "ID СЃС‚Р°С‚СЊРё"
// @Success 200 {object} swaggerKnowledgeBaseDeleteResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /knowledge-base/{id} [delete]
func _swaggerKnowledgeBaseDelete() {}

type swaggerKnowledgeBaseArticle struct {
	ID        int32  `json:"id" example:"1"`
	Title     string `json:"title" example:"РљР°Рє РѕС„РѕСЂРјРёС‚СЊ Р·Р°РєР°Р·"`
	Content   string `json:"content" example:"РўРµРєСЃС‚ СЃС‚Р°С‚СЊРё..."`
	CreatedAt string `json:"createdAt" example:"2026-03-31T10:00:00Z"`
	UpdatedAt string `json:"updatedAt" example:"2026-03-31T10:00:00Z"`
}

type swaggerKnowledgeBaseArticleRequest struct {
	Title   string `json:"title" example:"РљР°Рє РѕС„РѕСЂРјРёС‚СЊ Р·Р°РєР°Р·"`
	Content string `json:"content" example:"РўРµРєСЃС‚ СЃС‚Р°С‚СЊРё..."`
}

type swaggerKnowledgeBaseCreateResponse struct {
	Message string                      `json:"message" example:"РЎС‚Р°С‚СЊСЏ СѓСЃРїРµС€РЅРѕ СЃРѕР·РґР°РЅР°"`
	Article swaggerKnowledgeBaseArticle `json:"article"`
}

type swaggerKnowledgeBaseUpdateResponse struct {
	Message string                      `json:"message" example:"РЎС‚Р°С‚СЊСЏ СѓСЃРїРµС€РЅРѕ РѕР±РЅРѕРІР»РµРЅР°"`
	Article swaggerKnowledgeBaseArticle `json:"article"`
}

type swaggerKnowledgeBaseDeleteResponse struct {
	Message string `json:"message" example:"РЎС‚Р°С‚СЊСЏ СѓСЃРїРµС€РЅРѕ СѓРґР°Р»РµРЅР°"`
}

// ModerationList
// @Summary Список товаров AI-модерации
// @Tags moderation-admin
// @Produce json
// @Param filter query string false "ALL|DENIED|MANUAL|APPROVED_AI"
// @Param page query int false "Номер страницы"
// @Success 200 {object} swaggerModerationListResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/products [get]
func _swaggerModerationList() {}

// ModerationGetProduct
// @Summary Детали товара из AI-модерации
// @Tags moderation-admin
// @Produce json
// @Param id path int true "ID товара"
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
	Name                      string                   `json:"name" example:"Тонометр"`
	Price                     int32                    `json:"price" example:"3500"`
	Images                    []string                 `json:"images"`
	ModerateState             string                   `json:"moderateState" example:"AI_REVIEWED"`
	ModerationRejectionReason *string                  `json:"moderationRejectionReason,omitempty" example:"Текст: есть контакты"`
	CreatedAt                 string                   `json:"createdAt" example:"2026-04-06T10:00:00Z"`
	UpdatedAt                 string                   `json:"updatedAt" example:"2026-04-06T10:00:00Z"`
	Category                  swaggerModerationRefItem `json:"category"`
	SubCategory               swaggerModerationRefItem `json:"subCategory"`
	User                      swaggerModerationUser    `json:"user"`
}

type swaggerModerationProductDetail struct {
	ID                        int32                         `json:"id" example:"1000001"`
	Name                      string                        `json:"name" example:"Тонометр"`
	Price                     int32                         `json:"price" example:"3500"`
	Description               string                        `json:"description" example:"Описание товара"`
	Images                    []string                      `json:"images"`
	VideoURL                  *string                       `json:"videoUrl,omitempty" example:"https://example.com/video.mp4"`
	ModerateState             string                        `json:"moderateState" example:"AI_REVIEWED"`
	ModerationRejectionReason *string                       `json:"moderationRejectionReason,omitempty" example:"Фото: требуется ручная проверка"`
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
	Name string `json:"name" example:"Медтехника"`
}

type swaggerModerationTypeRefItem struct {
	ID   int32   `json:"id" example:"1"`
	Name *string `json:"name" example:"Тонометры"`
}

type swaggerModerationUser struct {
	ID          int32  `json:"id" example:"1"`
	FullName    string `json:"fullName" example:"Иван Иванов"`
	Email       string `json:"email" example:"ivan@example.com"`
	PhoneNumber string `json:"phoneNumber" example:"+79990000000"`
}

type swaggerModerationUserDetail struct {
	ID          int32  `json:"id" example:"1"`
	FullName    string `json:"fullName" example:"Иван Иванов"`
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
	Name string `json:"name" example:"Производитель"`
}
