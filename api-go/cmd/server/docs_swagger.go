package main

// Описание маршрутов только для генерации OpenAPI (swag). Реальные хендлеры в internal/httpserver.
// Админ-ручки категорий: сначала POST /auth/sign-in под пользователем с ролью admin — cookie session_id уйдёт в Try it out (credentials).

// --- system ---

// HealthCheck
// @Summary Проверка живости сервиса
// @Description Возвращает status ok
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func _swaggerHealth() {}

// --- log ---

// LogFindAll
// @Summary Список записей Log
// @Tags log
// @Produce json
// @Success 200 {array} object
// @Router /log/find-all [get]
func _swaggerLogFindAll() {}

// --- category (публично) ---

// CategoryFindAll
// @Summary Дерево всех категорий
// @Tags category
// @Produce json
// @Success 200 {array} object
// @Router /category/find-all [get]
func _swaggerCategoryFindAll() {}

// CategoryFindByID
// @Summary Категория по id (полное дерево)
// @Tags category
// @Produce json
// @Param id path int true "ID категории"
// @Success 200 {object} object
// @Failure 404 {object} map[string]interface{}
// @Router /category/find-by-id/{id} [get]
func _swaggerCategoryFindByID() {}

// CategoryFindBySlug
// @Summary Категория по slug
// @Tags category
// @Produce json
// @Param slug path string true "Slug"
// @Success 200 {object} object
// @Failure 404 {object} map[string]interface{}
// @Router /category/slug/{slug} [get]
func _swaggerCategoryFindBySlug() {}

// CategoryFindBySlugPath
// @Summary Разрешение цепочки slug (категория / подкатегория / тип)
// @Description В Swagger «Try it out» введи сегменты через %2F, например: elektronika%2Ftelefony
// @Tags category
// @Produce json
// @Param slugPath path string true "Цепочка (или один сегмент)"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /category/path/{slugPath} [get]
func _swaggerCategoryPath() {}

// --- category (админ; при ADMIN_API_KEY — заголовок X-Admin-Key) ---

// CategoryCreate
// @Summary Создать категорию
// @Tags category-admin
// @Accept json
// @Produce json
// @Param body body swaggerCreateCategory true "Тело"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /category/create-category [post]
func _swaggerCategoryCreate() {}

// CategoryUpdate
// @Summary Обновить категорию
// @Tags category-admin
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param body body swaggerUpdateCategory true "Тело"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /category/update-category/{id} [put]
func _swaggerCategoryUpdate() {}

// CategoryDelete
// @Summary Удалить категорию
// @Tags category-admin
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]interface{}
// @Router /category/delete-category/{id} [delete]
func _swaggerCategoryDelete() {}

// Тела запросов для Swagger UI
type swaggerCreateCategory struct {
	Name string  `json:"name" example:"Автомобили"`
	Slug *string `json:"slug,omitempty" example:"avtomobili"`
}

type swaggerUpdateCategory struct {
	Name string  `json:"name" example:"Автомобили"`
	Slug *string `json:"slug,omitempty"`
}

// --- auth ---

// AuthSignUp
// @Summary Регистрация — отправка кода (query where=telegram|sms)
// @Tags auth
// @Accept json
// @Produce json
// @Param where query string true "telegram или sms" Enums(telegram,sms)
// @Param body body swaggerSignUp true "Данные"
// @Success 200 {object} map[string]string
// @Router /auth/sign-up [post]
func _swaggerAuthSignUp() {}

// AuthVerifyMobile
// @Summary Подтверждение телефона по коду
// @Tags auth
// @Produce json
// @Param code query string true "Код из SMS/TG"
// @Success 200 {object} map[string]string
// @Router /auth/verify-mobile-code [post]
func _swaggerAuthVerifyMobile() {}

// AuthSignIn
// @Summary Вход (ставит cookie session_id)
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerSignIn true "Логин"
// @Success 200 {object} object
// @Router /auth/sign-in [post]
func _swaggerAuthSignIn() {}

// AuthMe
// @Summary Текущий пользователь
// @Tags auth
// @Produce json
// @Success 200 {object} object
// @Router /auth/me [get]
func _swaggerAuthMe() {}

// AuthIsAdmin
// @Summary Проверка роли admin
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]bool
// @Router /auth/isAdmin [get]
func _swaggerAuthIsAdmin() {}

// AuthLogout
// @Summary Выход
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func _swaggerAuthLogout() {}

// AuthForgot
// @Summary Запрос кода сброса на почту
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerForgotEmail true "email"
// @Success 200 {object} map[string]string
// @Router /auth/forgot-password [post]
func _swaggerAuthForgot() {}

// AuthVerifyForgot
// @Summary Проверка кода сброса
// @Tags auth
// @Produce json
// @Param code query string true "Код из письма"
// @Success 200 {object} map[string]int
// @Router /auth/verify-code [post]
func _swaggerAuthVerifyForgot() {}

// AuthChangePassword
// @Summary Новый пароль после verify-code
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerChangePassword true "Тело"
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
// @Summary Список пользователей (админ: cookie session_id + роль admin)
// @Tags user-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /user/find-all [get]
func _swaggerUserFindAll() {}

// UserInfo
// @Summary Карточка пользователя (рейтинг, лимит объявлений)
// @Tags user
// @Produce json
// @Param id path int true "User id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /user/info/{id} [get]
func _swaggerUserInfo() {}

// UserRemainingFreeAds
// @Summary Остаток бесплатных объявлений (сессия)
// @Tags user
// @Produce json
// @Success 200 {object} object
// @Router /user/remaining-free-ads [get]
func _swaggerUserRemainingFreeAds() {}

// UserShowNumber
// @Summary Показать номер продавца (сессия)
// @Tags user
// @Produce json
// @Param userId path int true "Продавец"
// @Success 200 {object} map[string]string
// @Router /user/show-number/{userId} [get]
func _swaggerUserShowNumber() {}

// UserUpdateSettings
// @Summary Обновление настроек (multipart, сессия)
// @Tags user
// @Accept mpfd
// @Produce json
// @Param fullName formData string false "ФИО"
// @Param phoneNumber formData string false "Телефон"
// @Param isAnswersCall formData string false "true/false"
// @Param profileType formData string false "INDIVIDUAL|OOO|IP"
// @Param photo formData file false "Аватар"
// @Success 200 {object} object
// @Router /user/update-settings [patch]
func _swaggerUserUpdateSettings() {}

// UserVerifyEmail
// @Summary Отправить код подтверждения на почту (сессия)
// @Tags user
// @Produce json
// @Success 200 {object} map[string]string
// @Router /user/verify-email [post]
func _swaggerUserVerifyEmail() {}

// UserVerifyEmailCode
// @Summary Подтвердить почту по коду из письма
// @Tags user
// @Produce json
// @Param code query string true "Код"
// @Success 200 {object} map[string]string
// @Router /user/verify-code [post]
func _swaggerUserVerifyEmailCode() {}

// UserSetBalance
// @Summary Установить bonusBalance (админ, сессия)
// @Tags user-admin
// @Produce json
// @Param userId path int true "User id"
// @Param balance query string true "Число"
// @Success 200 {object} map[string]string
// @Router /user/set-balance/{userId} [put]
func _swaggerUserSetBalance() {}

// UserToggleBanned
// @Summary Бан / разбан (админ, сессия)
// @Tags user-admin
// @Produce json
// @Param id path int true "User id"
// @Success 200 {object} map[string]string
// @Router /user/toggle-banned/{id} [put]
func _swaggerUserToggleBanned() {}

// UserAdminPatch
// @Summary Обновить пользователя (админ, сессия)
// @Tags user-admin
// @Accept json
// @Produce json
// @Param id path int true "User id"
// @Param body body swaggerAdminUpdateUser true "Поля"
// @Success 200 {object} map[string]string
// @Router /user/{id} [patch]
func _swaggerUserAdminPatch() {}

// UserAdminDelete
// @Summary Удалить пользователя (админ, сессия)
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
// @Summary Создать товар (multipart, до 8 images, сессия)
// @Tags product
// @Accept mpfd
// @Produce json
// @Success 201 {object} object
// @Router /product/create [post]
func _swaggerProductCreate() {}

// ProductAll
// @Summary Список товаров / поиск (query; optional сессия для избранного)
// @Tags product
// @Produce json
// @Router /product/all-products [get]
func _swaggerProductAll() {}

// ProductCard
// @Summary Карточка товара
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Router /product/product-card/{id} [get]
func _swaggerProductCard() {}

// ProductDelete
// @Summary Удалить свой товар
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Router /product/{id} [delete]
func _swaggerProductDelete() {}

// ProductPatch
// @Summary Обновить товар (multipart)
// @Tags product
// @Accept mpfd
// @Produce json
// @Param id path int true "Product id"
// @Router /product/{id} [patch]
func _swaggerProductPatch() {}

// ProductModerate
// @Summary Модерация (admin)
// @Tags product-admin
// @Produce json
// @Param id path int true "Product id"
// @Param status query string true "APPROVED|DENIDED"
// @Param reason query string false "Причина при DENIDED"
// @Router /product/moderate-product/{id} [put]
func _swaggerProductModerate() {}

// --- review ---

// ReviewSend
// @Summary Оставить отзыв продавцу (сессия)
// @Tags review
// @Accept json
// @Produce json
// @Param body body swaggerSendReview true "Тело"
// @Success 200 {object} map[string]string
// @Router /review/send-review [post]
func _swaggerReviewSend() {}

// ReviewUserReviews
// @Summary Одобренные отзывы о пользователе
// @Tags review
// @Produce json
// @Param id path int true "User id"
// @Router /review/user-reviews/{id} [get]
func _swaggerReviewUserReviews() {}

// ReviewModerate
// @Summary Модерация отзыва (admin)
// @Tags review-admin
// @Produce json
// @Param id path int true "Review id"
// @Param status query string true "APPROVED|DENIDED"
// @Router /review/moderate-review/{id} [put]
func _swaggerReviewModerate() {}

// ReviewModerateList
// @Summary Очередь отзывов на модерацию (admin)
// @Tags review-admin
// @Produce json
// @Router /review/all-reviews-to-moderate [get]
func _swaggerReviewModerateList() {}

type swaggerSendReview struct {
	Text           *string `json:"text,omitempty"`
	Rating         float64 `json:"rating" example:"5"`
	ReviewedUserID int32   `json:"reviewedUserId"`
}

// --- chat (сессия cookie session_id) ---

// ChatStart
// @Summary Начать чат по товару
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
// @Summary Список чатов пользователя
// @Tags chat
// @Produce json
// @Success 200 {array} object
// @Router /chat [get]
func _swaggerChatList() {}

// ChatMessages
// @Summary Сообщения чата (пагинация)
// @Tags chat
// @Produce json
// @Param id path int true "Chat id"
// @Param page query int false "Страница" default(1)
// @Param limit query int false "Лимит" default(50)
// @Success 200 {object} object
// @Router /chat/{id}/messages [get]
func _swaggerChatMessages() {}

// ChatInfo
// @Summary Информация о чате
// @Tags chat
// @Produce json
// @Param id path int true "Chat id"
// @Success 200 {object} object
// @Router /chat/{id} [get]
func _swaggerChatInfo() {}

type swaggerStartChat struct {
	ProductID int32 `json:"productId" example:"1"`
}

// --- payment (Т-Банк / Tinkoff; сессия — cookie session_id) ---

// PaymentCreate
// @Summary Создание платежа для пополнения баланса
// @Description Init в Т-Банк. Нужны TINKOFF_TERMINAL_KEY и TINKOFF_SECRET_KEY. Авторизация: cookie session_id после POST /auth/sign-in.
// @Tags payment
// @Accept json
// @Produce json
// @Param body body swaggerCreatePayment true "Сумма в рублях (мин. 1)"
// @Success 201 {object} swaggerPaymentCreateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /payment/create [post]
func _swaggerPaymentCreate() {}

// PaymentNotification
// @Summary Webhook уведомлений Т-Банка о статусе платежа
// @Description Без сессии. Подпись Token проверяется по полям тела. Тело — как приходит от банка; пример ниже.
// @Tags payment
// @Accept json
// @Produce json
// @Param body body swaggerTinkoffNotification true "Уведомление"
// @Success 200 {object} swaggerPaymentNotifyResponse
// @Failure 400 {object} map[string]interface{}
// @Router /payment/notification [post]
func _swaggerPaymentNotification() {}

// PaymentHistory
// @Summary История платежей пользователя
// @Description До 50 записей, новые сверху. Сессия: cookie session_id.
// @Tags payment
// @Produce json
// @Success 200 {array} swaggerPaymentHistoryItem
// @Failure 401 {object} map[string]interface{}
// @Router /payment/history [get]
func _swaggerPaymentHistory() {}

// PaymentCheckStatus
// @Summary Проверка статуса платежа в Т-Банке (GetState)
// @Description Сессия: cookie session_id. В теле — paymentId из ответа Init или уведомления.
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
	Description *string `json:"description,omitempty" example:"Пополнение баланса"`
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
	Message string `json:"message,omitempty" example:"Баланс успешно пополнен"`
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

// swaggerTinkoffNotification — поля как в Nest PaymentNotificationDto (реальный webhook может добавлять поля).
type swaggerTinkoffNotification struct {
	TerminalKey string `json:"TerminalKey" example:"1766153689307DEMO"`
	OrderID     string `json:"OrderId" example:"123-1735123456789"`
	Success     bool   `json:"Success" example:"true"`
	Status      string `json:"Status" example:"CONFIRMED"`
	PaymentID   string `json:"PaymentId" example:"2673412345"`
	Amount      int64  `json:"Amount" example:"100000"`
	Token       string `json:"Token" example:"подпись_от_банка"`
	ErrorCode   string `json:"ErrorCode,omitempty" example:"0"`
	Pan         string `json:"Pan,omitempty" example:"430000******0777"`
}

// --- promotion ---

// PromotionAll
// @Summary Все типы продвижения (тарифы)
// @Tags promotion
// @Produce json
// @Success 200 {array} object
// @Router /promotion/all-promotions [get]
func _swaggerPromotionAll() {}

// PromotionAdd
// @Summary Подключить продвижение к товару (сессия)
// @Tags promotion
// @Accept json
// @Produce json
// @Param body body swaggerAddPromotion true "Тело"
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
// @Summary Статистика пользователя (просмотры, телефон, избранное)
// @Tags statistics
// @Produce json
// @Param period query string false "day week month quarter half-year year"
// @Param categoryId query int false "Фильтр по категории (через SubCategory)"
// @Param region query string false "Подстрока в address (ILIKE)"
// @Param productId query int false "Конкретный товар"
// @Success 200 {object} object
// @Router /statistics/analytic [get]
func _swaggerStatisticsAnalytic() {}

// StatisticsProducts
// @Summary Аналитика по каждому товару продавца
// @Tags statistics
// @Produce json
// @Success 200 {array} object
// @Router /statistics/products-analytic [get]
func _swaggerStatisticsProducts() {}

// --- support (сессия) ---

// SupportCreateTicket
// @Summary Создать тикет поддержки
// @Tags support
// @Accept json
// @Produce json
// @Param body body swaggerCreateSupportTicket true "Тело"
// @Success 201 {object} object
// @Router /support/tickets [post]
func _swaggerSupportCreateTicket() {}

// SupportMyTickets
// @Summary Мои тикеты (пагинация, фильтры query)
// @Tags support
// @Produce json
// @Router /support/tickets/my [get]
func _swaggerSupportMyTickets() {}

// SupportAllTickets
// @Summary Все тикеты (модератор/admin)
// @Tags support
// @Produce json
// @Router /support/tickets/all [get]
func _swaggerSupportAllTickets() {}

// SupportStats
// @Summary Статистика тикетов (только admin)
// @Tags support
// @Produce json
// @Router /support/stats [get]
func _swaggerSupportStats() {}

// SupportGetTicket
// @Summary Тикет с сообщениями
// @Tags support
// @Produce json
// @Param id path int true "Ticket id"
// @Router /support/tickets/{id} [get]
func _swaggerSupportGetTicket() {}

// SupportSendMessage
// @Summary Сообщение в тикет
// @Tags support
// @Accept json
// @Param id path int true "Ticket id"
// @Param body body swaggerSupportMessage true "Текст"
// @Router /support/tickets/{id}/messages [post]
func _swaggerSupportSendMessage() {}

// SupportUpdateTicket
// @Summary Обновить тикет (модератор/admin)
// @Tags support
// @Accept json
// @Param id path int true "Ticket id"
// @Param body body swaggerUpdateSupportTicket true "Поля"
// @Router /support/tickets/{id} [put]
func _swaggerSupportUpdateTicket() {}

// SupportAssignTicket
// @Summary Назначить тикет на себя
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

// --- address (DaData, без сессии) ---

// AddressSuggestions
// @Summary Подсказки адреса (DaData)
// @Tags address
// @Produce json
// @Param query query string true "Строка поиска"
// @Param limit query int false "Лимит" default(5)
// @Success 200 {array} object
// @Router /address/suggestions [get]
func _swaggerAddressSuggestions() {}

// AddressValidate
// @Summary Проверка адреса по первой подсказке DaData
// @Tags address
// @Accept json
// @Produce json
// @Param body body swaggerValidateAddress true "Адрес"
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
// @Summary Создать баннер (multipart: image, name, place, navigateToUrl; сессия)
// @Tags banner
// @Accept mpfd
// @Produce json
// @Success 201 {object} object
// @Router /banner [post]
func _swaggerBannerCreate() {}

// BannerRandom
// @Summary Случайные одобренные баннеры
// @Tags banner
// @Produce json
// @Router /banner/random [get]
func _swaggerBannerRandom() {}

// BannerList
// @Summary Список одобренных баннеров (query place опционально)
// @Tags banner
// @Produce json
// @Router /banner [get]
func _swaggerBannerList() {}

// BannerModerate
// @Summary Модерация баннера (admin, query status)
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Param status query string true "MODERATE|APPROVED|DENIDED"
// @Router /banner/moderate/{id} [put]
func _swaggerBannerModerate() {}

// BannerAllModerate
// @Summary Очередь баннеров на модерацию (admin)
// @Tags banner
// @Produce json
// @Router /banner/all-banners-to-moderate [get]
func _swaggerBannerAllModerate() {}

// --- subcategory / subcategory-type / type-field (админ — сессия + роль admin) ---

// SubcategoryFindAll
// @Summary Список подкатегорий
// @Tags subcategory
// @Produce json
// @Success 200 {array} object
// @Router /subcategory/find-all [get]
func _swaggerSubcategoryFindAll() {}

// SubcategoryFindByID
// @Summary Подкатегория по id
// @Tags subcategory
// @Produce json
// @Param id path int true "ID"
// @Router /subcategory/find-by-id/{id} [get]
func _swaggerSubcategoryFindByID() {}

// SubcategoryTypeFindAll
// @Summary Все типы подкатегорий
// @Tags subcategory-type
// @Produce json
// @Router /subcategory-type/find-all [get]
func _swaggerSubcategoryTypeFindAll() {}

// SubcategoryTypeFindByID
// @Summary Тип подкатегории по id
// @Tags subcategory-type
// @Produce json
// @Param id path int true "ID"
// @Router /subcategory-type/find-by-id/{id} [get]
func _swaggerSubcategoryTypeFindByID() {}

// TypeFieldFindAll
// @Summary Все характеристики (поля типа)
// @Tags type-field
// @Produce json
// @Router /type-field/find-all [get]
func _swaggerTypeFieldFindAll() {}

// TypeFieldFindByID
// @Summary Характеристика по id
// @Tags type-field
// @Produce json
// @Param id path int true "ID"
// @Router /type-field/find-by-id/{id} [get]
func _swaggerTypeFieldFindByID() {}
