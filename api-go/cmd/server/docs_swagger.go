package main

// OpenAPI (swag). Р РµРіРµРЅРµСЂР°С†РёСЏ: swag init -g main.go -o ./docs -d ./cmd/server --parseInternal --outputTypes go,json
// (swagger.yaml РЅРµ РіРµРЅРµСЂРёРј: go-yaml РЅРµ РїСЂРёРЅРёРјР°РµС‚ U+0080..U+009F РёР· СЃС‚Р°СЂС‹С… Р±РёС‚С‹С… UTF-8 СЃС‚СЂРѕРє РІ @Summary/@Description.)
// РђРЅРЅРѕС‚Р°С†РёРё С‚РѕР»СЊРєРѕ Р·РґРµСЃСЊ (internal/httpserver Р±РµР· swag). РЎРµСЃСЃРёСЏ: POST /auth/sign-in РёР»Рё Authorize в†’ session_id.

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

// swaggerCreateDraftJSON С‚РµР»Рѕ С‡РµСЂРЅРѕРІРёРєР°: РїРѕР»СЏ РѕРїС†РёРѕРЅР°Р»СЊРЅС‹; price/quantity РІ JSON РјРѕРіСѓС‚ Р±С‹С‚СЊ СЃС‚СЂРѕРєРѕР№ РёР»Рё С‡РёСЃР»РѕРј.
type swaggerCreateDraftJSON struct {
	Name          string            `json:"name" example:""`
	Price         string            `json:"price" example:""`
	Quantity      string            `json:"quantity" example:""`
	State         string            `json:"state" example:"NEW"`
	Description   string            `json:"description" example:""`
	Address       string            `json:"address" example:""`
	CategoryID    string            `json:"categoryId" example:""`
	SubcategoryID string            `json:"subcategoryId" example:""`
	SubCategoryID string            `json:"subCategoryId" example:""`
	TypeID        string            `json:"typeId" example:""`
	FieldValues   map[string]string `json:"fieldValues"`
	VideoURL      string            `json:"videoUrl" example:""`
}

// ProductCreate
// @Summary РЎРѕР·РґР°С‚СЊ РѕР±СЉСЏРІР»РµРЅРёРµ РёР»Рё С‡РµСЂРЅРѕРІРёРє
// @Description РџРѕР»РЅР°СЏ С„РѕСЂРјР° Рё РІР°Р»РёРґР°С†РёСЏ - РѕР±СЉСЏРІР»РµРЅРёРµ РЅР° РјРѕРґРµСЂР°С†РёСЋ (isDraft=false, moderateState MODERATE). РќРµРїРѕР»РЅС‹Рµ РґР°РЅРЅС‹Рµ РёР»Рё РѕС€РёР±РєР° РІР°Р»РёРґР°С†РёРё - СЃРѕС…СЂР°РЅСЏРµС‚СЃСЏ С‡РµСЂРЅРѕРІРёРє (isDraft=true, DRAFT), РєР°Рє create-draft.
// @Description multipart/form-data (РґРѕ 8 images) РёР»Рё application/json. Р§РёСЃР»Р° РІ JSON РґРѕРїСѓСЃС‚РёРјС‹ РєР°Рє number. РќСѓР¶РЅР° cookie session_id.
// @Security SessionId
// @Tags product,product-draft
// @Accept json
// @Accept mpfd
// @Produce json
// @Param body body swaggerCreateDraftJSON false "РўРѕР»СЊРєРѕ РґР»СЏ JSON; РїСЂРё multipart РїРѕР»СЏ С„РѕСЂРјС‹ СЃРј. Р±СЌРєРµРЅРґ"
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
// @Security SessionId
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Router /product/{id} [delete]
func _swaggerProductDelete() {}

// ProductPatch
// @Summary РћР±РЅРѕРІРёС‚СЊ РѕР±СЉСЏРІР»РµРЅРёРµ РёР»Рё С‡РµСЂРЅРѕРІРёРє
// @Description Р”Р»СЏ DRAFT - РјСЏРіРєР°СЏ РІР°Р»РёРґР°С†РёСЏ; РґР»СЏ РѕСЃС‚Р°Р»СЊРЅС‹С… - РїРѕР»РЅР°СЏ. multipart/form-data РёР»Рё application/json.
// @Security SessionId
// @Tags product,product-draft
// @Accept json
// @Accept mpfd
// @Produce json
// @Param id path int true "Product id"
// @Param body body swaggerCreateDraftJSON false "Р§Р°СЃС‚РёС‡РЅРѕРµ С‚РµР»Рѕ РґР»СЏ JSON"
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

// ProductCreateDraft
// @Summary РЇРІРЅРѕ СЃРѕР·РґР°С‚СЊ С‡РµСЂРЅРѕРІРёРє
// @Description Р’СЃРµ РїРѕР»СЏ РѕРїС†РёРѕРЅР°Р»СЊРЅС‹, РјРѕР¶РЅРѕ РїСѓСЃС‚РѕР№ JSON-РѕР±СЉРµРєС‚. JSON РёР»Рё multipart (images). РћС‚РІРµС‚: product, isDraft=true, moderateState=DRAFT. РЎРЅР°С‡Р°Р»Р° Authorize (session_id) РёР»Рё POST /auth/sign-in.
// @Security SessionId
// @Tags product-draft
// @Accept json
// @Accept mpfd
// @Produce json
// @Param body body swaggerCreateDraftJSON false "РўРµР»Рѕ (РѕРїС†РёРѕРЅР°Р»СЊРЅРѕ)"
// @Success 201 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /product/create-draft [post]
func _swaggerProductCreateDraft() {}

// ProductMyDrafts
// @Summary РЎРїРёСЃРѕРє РјРѕРёС… С‡РµСЂРЅРѕРІРёРєРѕРІ
// @Description РўРѕР»СЊРєРѕ Р°РІС‚РѕСЂРёР·РѕРІР°РЅРЅС‹Р№ РїРѕР»СЊР·РѕРІР°С‚РµР»СЊ. Cookie session_id.
// @Security SessionId
// @Tags product-draft
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /product/my-drafts [get]
func _swaggerProductMyDrafts() {}

// ProductPublishDraft
// @Summary Р’С‹Р»РѕР¶РёС‚СЊ С‡РµСЂРЅРѕРІРёРє (РЅР° РјРѕРґРµСЂР°С†РёСЋ)
// @Description РџСЂРѕРІРµСЂРєР°: РёРјСЏ РЅРµ РїСѓСЃС‚РѕРµ Рё РЅРµ РґРµС„РѕР»С‚РЅРѕРµ Р§РµСЂРЅРѕРІРёРє, Р°РґСЂРµСЃ РІР°Р»РёРґРЅС‹Р№, С†РµРЅР° >= 1, state NEW|USED. РЈСЃРїРµС…: isDraft=false, moderateState=MODERATE.
// @Security SessionId
// @Tags product-draft
// @Produce json
// @Param id path int true "ID РїСЂРѕРґСѓРєС‚Р°-С‡РµСЂРЅРѕРІРёРєР°"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /product/publish-draft/{id} [post]
func _swaggerProductPublishDraft() {}

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

// swaggerTinkoffNotification вЂ” С‚РµР»Рѕ webhook РўРёРЅСЊРєРѕС„С„ (РєР°Рє Nest PaymentNotificationDto).
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

// registerSwaggerDocSymbols вЂ” СЃСЃС‹Р»РєРё РЅР° СЃРёРјРІРѕР»С‹ РґР»СЏ swag; РёРЅР°С‡Рµ staticcheck СЂСѓРіР°РµС‚СЃСЏ РЅР° В«unusedВ».
func init() {
	_ = []any{
		_swaggerAddressSuggestions, _swaggerAddressValidate, _swaggerCDEKCities, _swaggerCDEKDeliveryPoints, _swaggerCDEKCalculate,
		_swaggerAuthChangePassword, _swaggerAuthForgot, _swaggerAuthIsAdmin, _swaggerAuthLogout, _swaggerAuthMe,
		_swaggerDealGetByID, _swaggerDealGetCDEKQR, _swaggerDealMarkShipped,
		_swaggerAuthSignIn, _swaggerAuthSignUp, _swaggerAuthVerifyForgot, _swaggerAuthVerifyMobile,
		_swaggerBannerAllModerate, _swaggerBannerCreate, _swaggerBannerList, _swaggerBannerModerate, _swaggerBannerRandom,
		_swaggerCategoryCreate, _swaggerCategoryDelete, _swaggerCategoryFindAll, _swaggerCategoryFindByID,
		_swaggerCategoryFindBySlug, _swaggerCategoryPath, _swaggerCategoryUpdate,
		_swaggerChatInfo, _swaggerChatList, _swaggerChatMessages, _swaggerChatStart,
		_swaggerHealth, _swaggerKnowledgeBaseCreate, _swaggerKnowledgeBaseDelete, _swaggerKnowledgeBaseGetByID,
		_swaggerKnowledgeBaseList, _swaggerKnowledgeBaseUpdate, _swaggerLogFindAll,
		_swaggerPaymentCheckStatus, _swaggerPaymentCreate, _swaggerPaymentHistory, _swaggerPaymentNotification,
		_swaggerProductAll, _swaggerProductCard, _swaggerProductCreate, _swaggerProductCreateDraft,
		_swaggerProductDelete, _swaggerProductModerate, _swaggerProductMyDrafts, _swaggerProductPatch,
		_swaggerProductPublishDraft,
		_swaggerModerationGetProduct, _swaggerModerationList,
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
		swaggerCreateCategory{}, swaggerCreateDraftJSON{}, swaggerCreatePayment{}, swaggerCreateSupportTicket{}, swaggerForgotEmail{},
		swaggerDealCDEKQRResponse{}, swaggerDealMarkShippedRequest{},
		swaggerPaymentCheckStateResponse{}, swaggerPaymentCreateResponse{}, swaggerPaymentHistoryItem{},
		swaggerKnowledgeBaseArticle{}, swaggerKnowledgeBaseArticleRequest{}, swaggerKnowledgeBaseCreateResponse{},
		swaggerKnowledgeBaseDeleteResponse{}, swaggerKnowledgeBaseUpdateResponse{},
		swaggerModerationListResponse{}, swaggerModerationProductDetail{},
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

// --- deals (Р±РµР·РѕРїР°СЃРЅР°СЏ СЃРґРµР»РєР°; СЃРµСЃСЃРёСЏ cookie session_id) ---

// DealGetByID
// @Summary РџРѕР»СѓС‡РёС‚СЊ СЃРґРµР»РєСѓ РїРѕ ID
// @Description Р’РѕР·РІСЂР°С‰Р°РµС‚ РїРѕР»РЅСѓСЋ РєР°СЂС‚РѕС‡РєСѓ СЃРґРµР»РєРё, РІРєР»СЋС‡Р°СЏ Р±Р»РѕРє cdek (track, trackingUrl, trackPending).
// @Security SessionId
// @Tags deals
// @Produce json
// @Param id path int true "ID СЃРґРµР»РєРё"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /deals/{id} [get]
func _swaggerDealGetByID() {}

// DealMarkShipped
// @Summary РџРѕРґС‚РІРµСЂРґРёС‚СЊ РѕС‚РїСЂР°РІРєСѓ (РїСЂРѕРґР°РІРµС†)
// @Description РњРѕР¶РЅРѕ РїРµСЂРµРґР°С‚СЊ С‚РѕР»СЊРєРѕ orderUuid - С‚СЂРµРє РїРѕРґС‚СЏРЅРµС‚СЃСЏ Р°РІС‚РѕРјР°С‚РёС‡РµСЃРєРё РёР· CDEK, РєРѕРіРґР° Р±СѓРґРµС‚ РїСЂРёСЃРІРѕРµРЅ.
// @Security SessionId
// @Tags deals
// @Accept json
// @Produce json
// @Param id path int true "ID СЃРґРµР»РєРё"
// @Param body body swaggerDealMarkShippedRequest false "Р”Р°РЅРЅС‹Рµ РѕС‚РіСЂСѓР·РєРё CDEK"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /deals/{id}/mark-shipped [post]
func _swaggerDealMarkShipped() {}

// DealGetCDEKQR
// @Summary РџРѕР»СѓС‡РёС‚СЊ QR Рё С‚СЂРµРє CDEK РґР»СЏ СЃРґРµР»РєРё
// @Description Р’РѕР·РІСЂР°С‰Р°РµС‚ qrCodeData/qrCodeUrl, trackNumber, trackingUrl Рё orderUuid. QR Р±РµСЂРµС‚СЃСЏ РЅР°РїСЂСЏРјСѓСЋ РёР· РѕС‚РІРµС‚Р° CDEK API.
// @Security SessionId
// @Tags deals
// @Produce json
// @Param id path int true "ID СЃРґРµР»РєРё"
// @Success 200 {object} swaggerDealCDEKQRResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /deals/{id}/cdek-qr [get]
func _swaggerDealGetCDEKQR() {}

type swaggerDealMarkShippedRequest struct {
	CDEKOrderUUID   *string `json:"cdekOrderUuid,omitempty" example:"6f61a0f8-9260-4e6d-8d17-43a988ab86b8"`
	CDEKTrackNumber *string `json:"cdekTrackNumber,omitempty" example:"1401262037"`
}

type swaggerDealCDEKQRResponse struct {
	QRCodeData   *string `json:"qrCodeData,omitempty" example:"iVBORw0KGgoAAAANSUhEUgAA..."`
	QRCodeURL    *string `json:"qrCodeUrl,omitempty" example:"https://api.cdek.ru/v2/.../barcode.pdf"`
	TrackNumber  *string `json:"trackNumber,omitempty" example:"1401262037"`
	TrackingURL  *string `json:"trackingUrl,omitempty" example:"https://www.cdek.ru/ru/tracking?order_id=1401262037"`
	OrderUUID    *string `json:"orderUuid,omitempty" example:"6f61a0f8-9260-4e6d-8d17-43a988ab86b8"`
	TrackPending bool    `json:"trackPending" example:"true"`
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

// --- auth (VK) ---

// AuthVKURL
// @Summary РџРѕР»СѓС‡РёС‚СЊ СЃСЃС‹Р»РєСѓ VK OAuth
// @Tags auth
// @Produce json
// @Param state query string false "РџСЂРѕРёР·РІРѕР»СЊРЅС‹Р№ state"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /auth/vk/url [get]
func _swaggerAuthVKURL() {}

// AuthVKSignIn
// @Summary Р’С…РѕРґ С‡РµСЂРµР· VK
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerVKSignIn true "РљРѕРґ VK OAuth"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/sign-in [post]
func _swaggerAuthVKSignIn() {}

// AuthVKOnboardingStatus
// @Summary РЎС‚Р°С‚СѓСЃ VK onboarding
// @Security SessionId
// @Tags auth
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/onboarding/status [get]
func _swaggerAuthVKOnboardingStatus() {}

// AuthVKOnboardingStartEmail
// @Summary РќР°С‡Р°С‚СЊ РїРѕРґС‚РІРµСЂР¶РґРµРЅРёРµ email РґР»СЏ VK onboarding
// @Security SessionId
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerVKEmail true "Email"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/onboarding/start-email [post]
func _swaggerAuthVKOnboardingStartEmail() {}

// AuthVKOnboardingVerifyEmail
// @Summary РџРѕРґС‚РІРµСЂРґРёС‚СЊ email РєРѕРґРѕРј РґР»СЏ VK onboarding
// @Security SessionId
// @Tags auth
// @Produce json
// @Param code query string true "РљРѕРґ РёР· РїРёСЃСЊРјР°"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/onboarding/verify-email [post]
func _swaggerAuthVKOnboardingVerifyEmail() {}

// AuthVKOnboardingStartPhone
// @Summary РќР°С‡Р°С‚СЊ РїРѕРґС‚РІРµСЂР¶РґРµРЅРёРµ С‚РµР»РµС„РѕРЅР° РґР»СЏ VK onboarding
// @Security SessionId
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerVKPhone true "РўРµР»РµС„РѕРЅ"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/onboarding/start-phone [post]
func _swaggerAuthVKOnboardingStartPhone() {}

// AuthVKOnboardingVerifyPhone
// @Summary РџРѕРґС‚РІРµСЂРґРёС‚СЊ С‚РµР»РµС„РѕРЅ РєРѕРґРѕРј РґР»СЏ VK onboarding
// @Security SessionId
// @Tags auth
// @Produce json
// @Param code query string true "РљРѕРґ РёР· SMS"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/onboarding/verify-phone [post]
func _swaggerAuthVKOnboardingVerifyPhone() {}

type swaggerVKSignIn struct {
	Code     string `json:"code" example:"vk_oauth_code"`
	State    string `json:"state" example:"state123"`
	DeviceID string `json:"device_id" example:"device-abc-123"`
}

type swaggerVKEmail struct {
	Email string `json:"email" example:"user@example.com"`
}

type swaggerVKPhone struct {
	PhoneNumber string `json:"phoneNumber" example:"+79991234567"`
}

// --- user extra ---

// UserChangeRole
// @Summary РЎРјРµРЅРёС‚СЊ СЂРѕР»СЊ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ (admin)
// @Tags user-admin
// @Accept json
// @Produce json
// @Param id path int true "User id"
// @Param body body swaggerUserRoleChange true "РќРѕРІР°СЏ СЂРѕР»СЊ"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /user/{id}/role [put]
func _swaggerUserChangeRole() {}

type swaggerUserRoleChange struct {
	Role string `json:"role" example:"admin"`
}

// --- product extra ---

// ProductAvailableFilters
// @Summary Р”РѕСЃС‚СѓРїРЅС‹Рµ С„РёР»СЊС‚СЂС‹ РєР°С‚Р°Р»РѕРіР°
// @Tags product
// @Produce json
// @Param categoryId query int false "РљР°С‚РµРіРѕСЂРёСЏ"
// @Param subCategoryId query int false "РџРѕРґРєР°С‚РµРіРѕСЂРёСЏ"
// @Param typeId query int false "РўРёРї"
// @Success 200 {object} object
// @Router /product/available-filters [get]
func _swaggerProductAvailableFilters() {}

// ProductRandom
// @Summary РЎР»СѓС‡Р°Р№РЅС‹Рµ С‚РѕРІР°СЂС‹
// @Tags product
// @Produce json
// @Success 200 {array} object
// @Router /product/random-products [get]
func _swaggerProductRandom() {}

// ProductRecommended
// @Summary Р РµРєРѕРјРµРЅРґРѕРІР°РЅРЅС‹Рµ С‚РѕРІР°СЂС‹
// @Tags product
// @Produce json
// @Success 200 {array} object
// @Router /product/recommended [get]
func _swaggerProductRecommended() {}

// ProductUserProducts
// @Summary РўРѕРІР°СЂС‹ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ
// @Tags product
// @Produce json
// @Param id path int true "User id"
// @Success 200 {array} object
// @Router /product/user-products/{id} [get]
func _swaggerProductUserProducts() {}

// ProductAddToFavorites
// @Summary Р”РѕР±Р°РІРёС‚СЊ С‚РѕРІР°СЂ РІ РёР·Р±СЂР°РЅРЅРѕРµ
// @Security SessionId
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]interface{}
// @Router /product/add-to-favorites/{id} [post]
func _swaggerProductAddToFavorites() {}

// ProductRemoveFromFavorites
// @Summary РЈРґР°Р»РёС‚СЊ С‚РѕРІР°СЂ РёР· РёР·Р±СЂР°РЅРЅРѕРіРѕ
// @Security SessionId
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]interface{}
// @Router /product/remove-from-favorites/{id} [delete]
func _swaggerProductRemoveFromFavorites() {}

// ProductMyFavorites
// @Summary РњРѕРµ РёР·Р±СЂР°РЅРЅРѕРµ
// @Security SessionId
// @Tags product
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /product/my-favorites [get]
func _swaggerProductMyFavorites() {}

// ProductToggle
// @Summary РЎРєСЂС‹С‚СЊ РёР»Рё РѕРїСѓР±Р»РёРєРѕРІР°С‚СЊ СЃРІРѕР№ С‚РѕРІР°СЂ
// @Security SessionId
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /product/toggle-product/{id} [put]
func _swaggerProductToggle() {}

// ProductAllToModerate
// @Summary РўРѕРІР°СЂС‹ РЅР° РјРѕРґРµСЂР°С†РёРё (admin)
// @Tags product-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /product/all-products-to-moderate [get]
func _swaggerProductAllToModerate() {}

// ProductPromoted
// @Summary РџСЂРѕРґРІРёРіР°РµРјС‹Рµ С‚РѕРІР°СЂС‹ (admin)
// @Tags product-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /product/promoted-products [get]
func _swaggerProductPromoted() {}

// ProductTogglePromotion
// @Summary Р’РєР»СЋС‡РёС‚СЊ РёР»Рё РІС‹РєР»СЋС‡РёС‚СЊ РїСЂРѕРґРІРёР¶РµРЅРёРµ (admin)
// @Tags product-admin
// @Produce json
// @Param promotionId path int true "Promotion id"
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /product/toggle-promotion/{promotionId} [put]
func _swaggerProductTogglePromotion() {}

// --- review appeals ---

// ReviewCreateAppeal
// @Summary РџРѕРґР°С‚СЊ Р°РїРµР»Р»СЏС†РёСЋ РЅР° РѕС‚Р·С‹РІ
// @Security SessionId
// @Tags review
// @Accept json
// @Produce json
// @Param body body swaggerReviewAppealCreate true "РўРµР»Рѕ"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /review/appeals [post]
func _swaggerReviewCreateAppeal() {}

// ReviewMyAppeals
// @Summary РњРѕРё Р°РїРµР»Р»СЏС†РёРё РЅР° РѕС‚Р·С‹РІС‹
// @Security SessionId
// @Tags review
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /review/my-appeals [get]
func _swaggerReviewMyAppeals() {}

// ReviewAllAppeals
// @Summary Р’СЃРµ Р°РїРµР»Р»СЏС†РёРё РЅР° РѕС‚Р·С‹РІС‹ (moderator/admin)
// @Tags review-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /review/all-appeals [get]
func _swaggerReviewAllAppeals() {}

// ReviewResolveAppeal
// @Summary Р Р°Р·СЂРµС€РёС‚СЊ Р°РїРµР»Р»СЏС†РёСЋ РЅР° РѕС‚Р·С‹РІ (moderator/admin)
// @Tags review-admin
// @Accept json
// @Produce json
// @Param id path int true "Appeal id"
// @Param body body swaggerResolveAppeal true "Р РµС€РµРЅРёРµ"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /review/resolve-appeal/{id} [put]
func _swaggerReviewResolveAppeal() {}

type swaggerReviewAppealCreate struct {
	ReviewID int    `json:"reviewId" example:"12"`
	Reason   string `json:"reason" example:"РћС‚Р·С‹РІ СЃРѕРґРµСЂР¶РёС‚ РЅРµРґРѕСЃС‚РѕРІРµСЂРЅС‹Рµ СЃРІРµРґРµРЅРёСЏ"`
}

type swaggerResolveAppeal struct {
	Status        string `json:"status" example:"RESOLVED"`
	ModeratorNote string `json:"moderatorNote" example:"РђРїРµР»Р»СЏС†РёСЏ СЂР°СЃСЃРјРѕС‚СЂРµРЅР°"`
}

// --- reservation ---

// ReservationCreate
// @Summary РЎРѕР·РґР°С‚СЊ СЂРµР·РµСЂРІ С‚РѕРІР°СЂР°
// @Security SessionId
// @Tags reservation
// @Accept json
// @Produce json
// @Param body body swaggerReservationCreate true "РўРµР»Рѕ"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/ [post]
func _swaggerReservationCreate() {}

// ReservationMy
// @Summary РњРѕРё СЂРµР·РµСЂРІС‹
// @Security SessionId
// @Tags reservation
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/my [get]
func _swaggerReservationMy() {}

// ReservationProductInfo
// @Summary РРЅС„РѕСЂРјР°С†РёСЏ Рѕ СЂРµР·РµСЂРІРёСЂРѕРІР°РЅРёРё С‚РѕРІР°СЂР°
// @Tags reservation
// @Produce json
// @Param productId path int true "Product id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /reservations/product/{productId} [get]
func _swaggerReservationProductInfo() {}

// ReservationCancelByBuyer
// @Summary РћС‚РјРµРЅР° СЂРµР·РµСЂРІР° РїРѕРєСѓРїР°С‚РµР»РµРј
// @Security SessionId
// @Tags reservation
// @Produce json
// @Param id path int true "Reservation id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/{id}/cancel-by-buyer [post]
func _swaggerReservationCancelByBuyer() {}

// ReservationCancelBySeller
// @Summary РћС‚РјРµРЅР° СЂРµР·РµСЂРІР° РїСЂРѕРґР°РІС†РѕРј
// @Security SessionId
// @Tags reservation
// @Accept json
// @Produce json
// @Param id path int true "Reservation id"
// @Param body body swaggerReservationCancelReason true "РџСЂРёС‡РёРЅР°"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/{id}/cancel-by-seller [post]
func _swaggerReservationCancelBySeller() {}

// ReservationCancel
// @Summary РЈРЅРёРІРµСЂСЃР°Р»СЊРЅР°СЏ РѕС‚РјРµРЅР° СЂРµР·РµСЂРІР°
// @Security SessionId
// @Tags reservation
// @Accept json
// @Produce json
// @Param id path int true "Reservation id"
// @Param body body swaggerReservationCancelReason false "РџСЂРёС‡РёРЅР°"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/{id}/cancel [post]
func _swaggerReservationCancel() {}

// ReservationExtend
// @Summary РџСЂРѕРґР»РёС‚СЊ СЂРµР·РµСЂРІ
// @Security SessionId
// @Tags reservation
// @Produce json
// @Param id path int true "Reservation id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/{id}/extend [post]
func _swaggerReservationExtend() {}

// ReservationUpdateProductSettings
// @Summary РћР±РЅРѕРІРёС‚СЊ РЅР°СЃС‚СЂРѕР№РєРё СЂРµР·РµСЂРІРёСЂРѕРІР°РЅРёСЏ С‚РѕРІР°СЂР°
// @Security SessionId
// @Tags reservation
// @Accept json
// @Produce json
// @Param body body swaggerReservationProductSettings true "РўРµР»Рѕ"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/product-settings [put]
func _swaggerReservationUpdateProductSettings() {}

type swaggerReservationCreate struct {
	ProductID int     `json:"productId" example:"6157119"`
	Hours     *int    `json:"hours,omitempty" example:"24"`
	Note      *string `json:"note,omitempty" example:"РџСЂРѕС€Сѓ РїСЂРёРґРµСЂР¶Р°С‚СЊ С‚РѕРІР°СЂ РґРѕ РІРµС‡РµСЂР°"`
}

type swaggerReservationCancelReason struct {
	Reason *string `json:"reason,omitempty" example:"РџРѕРєСѓРїР°С‚РµР»СЊ РЅРµ РІС‹С€РµР» РЅР° СЃРІСЏР·СЊ"`
}

type swaggerReservationProductSettings struct {
	ProductID         int  `json:"productId" example:"6157119"`
	AllowReservations bool `json:"allowReservations" example:"true"`
}

// --- statistics extra ---

// StatisticsSearchQueries
// @Summary РЎС‚Р°С‚РёСЃС‚РёРєР° РїРѕРёСЃРєРѕРІС‹С… Р·Р°РїСЂРѕСЃРѕРІ
// @Security SessionId
// @Tags statistics
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /statistics/search-queries [get]
func _swaggerStatisticsSearchQueries() {}

// StatisticsCabinetDashboard
// @Summary Р”Р°С€Р±РѕСЂРґ Р»РёС‡РЅРѕРіРѕ РєР°Р±РёРЅРµС‚Р°
// @Security SessionId
// @Tags statistics
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /statistics/cabinet-dashboard [get]
func _swaggerStatisticsCabinetDashboard() {}

// --- banner extra ---

// BannerMyStats
// @Summary РњРѕСЏ СЃС‚Р°С‚РёСЃС‚РёРєР° РїРѕ Р±Р°РЅРЅРµСЂР°Рј
// @Security SessionId
// @Tags banner
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /banner/my-stats/all [get]
func _swaggerBannerMyStats() {}

// BannerViewTrack
// @Summary Р—Р°С„РёРєСЃРёСЂРѕРІР°С‚СЊ РїСЂРѕСЃРјРѕС‚СЂ Р±Р°РЅРЅРµСЂР°
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Success 200 {object} map[string]string
// @Router /banner/{id}/view [post]
func _swaggerBannerViewTrack() {}

// BannerStats
// @Summary РЎС‚Р°С‚РёСЃС‚РёРєР° РєРѕРЅРєСЂРµС‚РЅРѕРіРѕ Р±Р°РЅРЅРµСЂР°
// @Security SessionId
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /banner/{id}/stats [get]
func _swaggerBannerStats() {}

// BannerUpdate
// @Summary РћР±РЅРѕРІРёС‚СЊ Р±Р°РЅРЅРµСЂ (admin)
// @Tags banner
// @Accept json
// @Produce json
// @Param id path int true "Banner id"
// @Param body body swaggerBannerUpdate true "РўРµР»Рѕ"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /banner/{id} [put]
func _swaggerBannerUpdate() {}

// BannerDelete
// @Summary РЈРґР°Р»РёС‚СЊ Р±Р°РЅРЅРµСЂ (admin)
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /banner/{id} [delete]
func _swaggerBannerDelete() {}

// BannerGetByID
// @Summary РџРѕР»СѓС‡РёС‚СЊ Р±Р°РЅРЅРµСЂ РїРѕ id
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Success 200 {object} object
// @Failure 404 {object} map[string]interface{}
// @Router /banner/{id} [get]
func _swaggerBannerGetByID() {}

type swaggerBannerUpdate struct {
	Name          *string `json:"name,omitempty" example:"Р›РµС‚РЅСЏСЏ Р°РєС†РёСЏ"`
	PhotoURL      *string `json:"photoUrl,omitempty" example:"https://cdn.example.com/banner.jpg"`
	Place         *string `json:"place,omitempty" example:"PRODUCT_FEED"`
	NavigateToURL *string `json:"navigateToUrl,omitempty" example:"https://torguisam.ru/product/6157119"`
}

// --- moderation extra ---

// ModerationSummary
// @Summary РЎРІРѕРґРєР° РїРѕ РјРѕРґРµСЂР°С†РёРё
// @Tags moderation-admin
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/summary [get]
func _swaggerModerationSummary() {}

// ModerationAuditLogs
// @Summary Р–СѓСЂРЅР°Р» Р°СѓРґРёС‚Р° РјРѕРґРµСЂР°С†РёРё
// @Tags moderation-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/audit-logs [get]
func _swaggerModerationAuditLogs() {}

// ModerationAppeals
// @Summary РЎРїРёСЃРѕРє Р°РїРµР»Р»СЏС†РёР№ РјРѕРґРµСЂР°С†РёРё
// @Tags moderation-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/appeals [get]
func _swaggerModerationAppeals() {}

// ModerationReviewAppeal
// @Summary Р Р°СЃСЃРјРѕС‚СЂРµС‚СЊ Р°РїРµР»Р»СЏС†РёСЋ РјРѕРґРµСЂР°С†РёРё
// @Tags moderation-admin
// @Accept json
// @Produce json
// @Param id path int true "Appeal id"
// @Param body body swaggerModerationAppealReview true "Р РµС€РµРЅРёРµ"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/appeals/{id}/review [put]
func _swaggerModerationReviewAppeal() {}

type swaggerModerationAppealReview struct {
	Status        string `json:"status" example:"APPROVED"`
	ReviewComment string `json:"reviewComment" example:"РђРїРµР»Р»СЏС†РёСЏ СЂР°СЃСЃРјРѕС‚СЂРµРЅР° РјРѕРґРµСЂР°С‚РѕСЂРѕРј"`
}

// --- deals extra ---

// DealCreate
// @Summary РЎРѕР·РґР°С‚СЊ Р±РµР·РѕРїР°СЃРЅСѓСЋ СЃРґРµР»РєСѓ
// @Security SessionId
// @Tags deals
// @Accept json
// @Produce json
// @Param body body swaggerDealCreate true "РўРµР»Рѕ"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /deals/ [post]
func _swaggerDealCreate() {}

// DealMyPurchases
// @Summary РњРѕРё РїРѕРєСѓРїРєРё
// @Security SessionId
// @Tags deals
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /deals/my-purchases [get]
func _swaggerDealMyPurchases() {}

// DealMySales
// @Summary РњРѕРё РїСЂРѕРґР°Р¶Рё
// @Security SessionId
// @Tags deals
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /deals/my-sales [get]
func _swaggerDealMySales() {}

// DealMyAll
// @Summary Р’СЃРµ РјРѕРё СЃРґРµР»РєРё
// @Security SessionId
// @Tags deals
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /deals/my [get]
func _swaggerDealMyAll() {}

// DealPay
// @Summary РћРїР»Р°С‚РёС‚СЊ СЃРґРµР»РєСѓ
// @Security SessionId
// @Tags deals
// @Produce json
// @Param id path int true "Deal id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /deals/{id}/pay [post]
func _swaggerDealPay() {}

// DealSyncPayment
// @Summary РЎРёРЅС…СЂРѕРЅРёР·РёСЂРѕРІР°С‚СЊ СЃС‚Р°С‚СѓСЃ РѕРїР»Р°С‚С‹ СЃРґРµР»РєРё
// @Security SessionId
// @Tags deals
// @Produce json
// @Param id path int true "Deal id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /deals/{id}/sync-payment [post]
func _swaggerDealSyncPayment() {}

// DealSetCDEKHandoff
// @Summary РЈРєР°Р·Р°С‚СЊ СЃРїРѕСЃРѕР± РїРµСЂРµРґР°С‡Рё С‚РѕРІР°СЂР° РІ CDEK
// @Security SessionId
// @Tags deals
// @Accept json
// @Produce json
// @Param id path int true "Deal id"
// @Param body body swaggerDealHandoff true "РўРµР»Рѕ"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /deals/{id}/cdek-handoff [post]
func _swaggerDealSetCDEKHandoff() {}

// DealConfirmDelivery
// @Summary РџРѕРґС‚РІРµСЂРґРёС‚СЊ РїРѕР»СѓС‡РµРЅРёРµ С‚РѕРІР°СЂР°
// @Security SessionId
// @Tags deals
// @Produce json
// @Param id path int true "Deal id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /deals/{id}/confirm-delivery [post]
func _swaggerDealConfirmDelivery() {}

// DealOpenDispute
// @Summary РћС‚РєСЂС‹С‚СЊ СЃРїРѕСЂ РїРѕ СЃРґРµР»РєРµ
// @Security SessionId
// @Tags deals
// @Accept json
// @Produce json
// @Param id path int true "Deal id"
// @Param body body swaggerDealDispute true "РџСЂРёС‡РёРЅР°"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /deals/{id}/open-dispute [post]
func _swaggerDealOpenDispute() {}

// DealCancel
// @Summary Отменить сделку
// @Security SessionId
// @Tags deals
// @Produce json
// @Param id path int true "Deal id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /deals/{id}/cancel [post]
func _swaggerDealCancel() {}

type swaggerDealCreate struct {
	ProductID         int     `json:"productId" example:"6157119"`
	DeliveryCost      int     `json:"deliveryCost" example:"185"`
	CDEKTariffCode    *int    `json:"cdekTariffCode,omitempty" example:"136"`
    CDEKTariffName    *string `json:"cdekTariffName,omitempty" example:"Посылка склад-склад"`
	CDEKFromCityCode  *int    `json:"cdekFromCityCode,omitempty" example:"261"`
	CDEKToCityCode    *int    `json:"cdekToCityCode,omitempty" example:"44"`
	CDEKFromPvzCode   *string `json:"cdekFromPvzCode,omitempty" example:"ORN24"`
	CDEKToPvzCode     *string `json:"cdekToPvzCode,omitempty" example:"MSK12"`
    CDEKToAddress     *string `json:"cdekToAddress,omitempty" example:"ул. Пример, д. 1"`
	CDEKPackageWeight *int    `json:"cdekPackageWeight,omitempty" example:"500"`
	CDEKPackageLength *int    `json:"cdekPackageLength,omitempty" example:"17"`
	CDEKPackageWidth  *int    `json:"cdekPackageWidth,omitempty" example:"12"`
	CDEKPackageHeight *int    `json:"cdekPackageHeight,omitempty" example:"9"`
	CDEKRecipientMode *string `json:"cdekRecipientMode,omitempty" example:"pvz"`
}

type swaggerDealHandoff struct {
	Mode            string  `json:"mode" example:"pvz"`
	CDEKFromPvzCode *string `json:"cdekFromPvzCode,omitempty" example:"ORN24"`
    CDEKFromAddress *string `json:"cdekFromAddress,omitempty" example:"г. Оренбург, ул. Чкалова, 59"`
}

type swaggerDealDispute struct {
    Reason string `json:"reason" example:"Получен товар в ненадлежащем состоянии"`
}

// AdminDealList
// @Summary Список сделок для администратора и модератора
// @Security SessionId
// @Tags admin-deals
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /admin/deals/list [get]
func _swaggerAdminDealList() {}

// AdminDealGet
// @Summary Получить сделку по id для администратора
// @Security SessionId
// @Tags admin-deals
// @Produce json
// @Param id path int true "Deal id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/deals/{id} [get]
func _swaggerAdminDealGet() {}

// AdminDealSetStatus
// @Summary Изменить статус сделки
// @Security SessionId
// @Tags admin-deals
// @Accept json
// @Produce json
// @Param id path int true "Deal id"
// @Param body body swaggerAdminDealStatus true "Тело"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/deals/{id}/status [patch]
func _swaggerAdminDealSetStatus() {}

// AdminDealLogs
// @Summary Получить логи по сделке
// @Security SessionId
// @Tags admin-deals
// @Produce json
// @Param id path int true "Deal id"
// @Success 200 {array} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /admin/deals/{id}/logs [get]
func _swaggerAdminDealLogs() {}

type swaggerAdminDealStatus struct {
	Status string `json:"status" example:"CANCELLED"`
}

// --- cdek extra ---

// CDEKTariffs
// @Summary РџРѕР»СѓС‡РёС‚СЊ РґРѕСЃС‚СѓРїРЅС‹Рµ С‚Р°СЂРёС„С‹ CDEK
// @Tags cdek
// @Accept json
// @Produce json
// @Param body body swaggerCDEKCalculateRequest true "Body"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /cdek/tariffs [post]
func _swaggerCDEKTariffs() {}
