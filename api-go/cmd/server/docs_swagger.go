package main

// OpenAPI (swag). Регенерация: swag init -g main.go -o ./docs -d ./cmd/server --parseInternal --outputTypes go,json
// (swagger.yaml не генерим: go-yaml не принимает U+0080..U+009F из старых битых UTF-8 строк в @Summary/@Description.)
// Аннотации только здесь (internal/httpserver без swag). Сессия: POST /auth/sign-in или Authorize → session_id.

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

// swaggerCreateDraftJSON тело черновика: поля опциональны; price/quantity в JSON могут быть строкой или числом.
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
// @Summary Создать объявление или черновик
// @Description Полная форма и валидация - объявление на модерацию (isDraft=false, moderateState MODERATE). Неполные данные или ошибка валидации - сохраняется черновик (isDraft=true, DRAFT), как create-draft.
// @Description multipart/form-data (до 8 images) или application/json. Числа в JSON допустимы как number. Нужна cookie session_id.
// @Security SessionId
// @Tags product,product-draft
// @Accept json
// @Accept mpfd
// @Produce json
// @Param body body swaggerCreateDraftJSON false "Только для JSON; при multipart поля формы см. бэкенд"
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
// @Security SessionId
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Router /product/{id} [delete]
func _swaggerProductDelete() {}

// ProductPatch
// @Summary Обновить объявление или черновик
// @Description Для DRAFT - мягкая валидация; для остальных - полная. multipart/form-data или application/json.
// @Security SessionId
// @Tags product,product-draft
// @Accept json
// @Accept mpfd
// @Produce json
// @Param id path int true "Product id"
// @Param body body swaggerCreateDraftJSON false "Частичное тело для JSON"
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

// ProductCreateDraft
// @Summary Явно создать черновик
// @Description Все поля опциональны, можно пустой JSON-объект. JSON или multipart (images). Ответ: product, isDraft=true, moderateState=DRAFT. Сначала Authorize (session_id) или POST /auth/sign-in.
// @Security SessionId
// @Tags product-draft
// @Accept json
// @Accept mpfd
// @Produce json
// @Param body body swaggerCreateDraftJSON false "Тело (опционально)"
// @Success 201 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /product/create-draft [post]
func _swaggerProductCreateDraft() {}

// ProductMyDrafts
// @Summary Список моих черновиков
// @Description Только авторизованный пользователь. Cookie session_id.
// @Security SessionId
// @Tags product-draft
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /product/my-drafts [get]
func _swaggerProductMyDrafts() {}

// ProductPublishDraft
// @Summary Выложить черновик (на модерацию)
// @Description Проверка: имя не пустое и не дефолтное Черновик, адрес валидный, цена >= 1, state NEW|USED. Успех: isDraft=false, moderateState=MODERATE.
// @Security SessionId
// @Tags product-draft
// @Produce json
// @Param id path int true "ID продукта-черновика"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /product/publish-draft/{id} [post]
func _swaggerProductPublishDraft() {}

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

// swaggerTinkoffNotification — тело webhook Тинькофф (как Nest PaymentNotificationDto).
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

// registerSwaggerDocSymbols — ссылки на символы для swag; иначе staticcheck ругается на «unused».
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
// @Summary Список статей базы знаний
// @Tags knowledge-base
// @Produce json
// @Success 200 {array} swaggerKnowledgeBaseArticle
// @Router /knowledge-base/ [get]
func _swaggerKnowledgeBaseList() {}

// KnowledgeBaseGetByID
// @Summary Статья базы знаний по id
// @Tags knowledge-base
// @Produce json
// @Param id path int true "ID статьи"
// @Success 200 {object} swaggerKnowledgeBaseArticle
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /knowledge-base/{id} [get]
func _swaggerKnowledgeBaseGetByID() {}

// KnowledgeBaseCreate
// @Summary Создать статью базы знаний
// @Tags knowledge-base-admin
// @Accept json
// @Produce json
// @Param body body swaggerKnowledgeBaseArticleRequest true "Тело"
// @Success 201 {object} swaggerKnowledgeBaseCreateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /knowledge-base/ [post]
func _swaggerKnowledgeBaseCreate() {}

// KnowledgeBaseUpdate
// @Summary Обновить статью базы знаний
// @Tags knowledge-base-admin
// @Accept json
// @Produce json
// @Param id path int true "ID статьи"
// @Param body body swaggerKnowledgeBaseArticleRequest true "Тело"
// @Success 200 {object} swaggerKnowledgeBaseUpdateResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /knowledge-base/{id} [put]
func _swaggerKnowledgeBaseUpdate() {}

// KnowledgeBaseDelete
// @Summary Удалить статью базы знаний
// @Tags knowledge-base-admin
// @Produce json
// @Param id path int true "ID статьи"
// @Success 200 {object} swaggerKnowledgeBaseDeleteResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /knowledge-base/{id} [delete]
func _swaggerKnowledgeBaseDelete() {}

type swaggerKnowledgeBaseArticle struct {
	ID        int32  `json:"id" example:"1"`
	Title     string `json:"title" example:"Как оформить заказ"`
	Content   string `json:"content" example:"Текст статьи..."`
	CreatedAt string `json:"createdAt" example:"2026-03-31T10:00:00Z"`
	UpdatedAt string `json:"updatedAt" example:"2026-03-31T10:00:00Z"`
}

type swaggerKnowledgeBaseArticleRequest struct {
	Title   string `json:"title" example:"Как оформить заказ"`
	Content string `json:"content" example:"Текст статьи..."`
}

type swaggerKnowledgeBaseCreateResponse struct {
	Message string                      `json:"message" example:"Статья успешно создана"`
	Article swaggerKnowledgeBaseArticle `json:"article"`
}

type swaggerKnowledgeBaseUpdateResponse struct {
	Message string                      `json:"message" example:"Статья успешно обновлена"`
	Article swaggerKnowledgeBaseArticle `json:"article"`
}

type swaggerKnowledgeBaseDeleteResponse struct {
	Message string `json:"message" example:"Статья успешно удалена"`
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

// --- deals (безопасная сделка; сессия cookie session_id) ---

// DealGetByID
// @Summary Получить сделку по ID
// @Description Возвращает полную карточку сделки, включая блок cdek (track, trackingUrl, trackPending).
// @Security SessionId
// @Tags deals
// @Produce json
// @Param id path int true "ID сделки"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /deals/{id} [get]
func _swaggerDealGetByID() {}

// DealMarkShipped
// @Summary Подтвердить отправку (продавец)
// @Description Можно передать только orderUuid - трек подтянется автоматически из CDEK, когда будет присвоен.
// @Security SessionId
// @Tags deals
// @Accept json
// @Produce json
// @Param id path int true "ID сделки"
// @Param body body swaggerDealMarkShippedRequest false "Данные отгрузки CDEK"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /deals/{id}/mark-shipped [post]
func _swaggerDealMarkShipped() {}

// DealGetCDEKQR
// @Summary Получить QR и трек CDEK для сделки
// @Description Возвращает qrCodeData/qrCodeUrl, trackNumber, trackingUrl и orderUuid. QR берется напрямую из ответа CDEK API.
// @Security SessionId
// @Tags deals
// @Produce json
// @Param id path int true "ID сделки"
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
// @Summary Получить ссылку VK OAuth
// @Tags auth
// @Produce json
// @Param state query string false "Произвольный state"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Router /auth/vk/url [get]
func _swaggerAuthVKURL() {}

// AuthVKSignIn
// @Summary Вход через VK
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerVKSignIn true "Код VK OAuth"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/sign-in [post]
func _swaggerAuthVKSignIn() {}

// AuthVKOnboardingStatus
// @Summary Статус VK onboarding
// @Security SessionId
// @Tags auth
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/onboarding/status [get]
func _swaggerAuthVKOnboardingStatus() {}

// AuthVKOnboardingStartEmail
// @Summary Начать подтверждение email для VK onboarding
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
// @Summary Подтвердить email кодом для VK onboarding
// @Security SessionId
// @Tags auth
// @Produce json
// @Param code query string true "Код из письма"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/onboarding/verify-email [post]
func _swaggerAuthVKOnboardingVerifyEmail() {}

// AuthVKOnboardingStartPhone
// @Summary Начать подтверждение телефона для VK onboarding
// @Security SessionId
// @Tags auth
// @Accept json
// @Produce json
// @Param body body swaggerVKPhone true "Телефон"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/vk/onboarding/start-phone [post]
func _swaggerAuthVKOnboardingStartPhone() {}

// AuthVKOnboardingVerifyPhone
// @Summary Подтвердить телефон кодом для VK onboarding
// @Security SessionId
// @Tags auth
// @Produce json
// @Param code query string true "Код из SMS"
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
// @Summary Сменить роль пользователя (admin)
// @Tags user-admin
// @Accept json
// @Produce json
// @Param id path int true "User id"
// @Param body body swaggerUserRoleChange true "Новая роль"
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
// @Summary Доступные фильтры каталога
// @Tags product
// @Produce json
// @Param categoryId query int false "Категория"
// @Param subCategoryId query int false "Подкатегория"
// @Param typeId query int false "Тип"
// @Success 200 {object} object
// @Router /product/available-filters [get]
func _swaggerProductAvailableFilters() {}

// ProductRandom
// @Summary Случайные товары
// @Tags product
// @Produce json
// @Success 200 {array} object
// @Router /product/random-products [get]
func _swaggerProductRandom() {}

// ProductRecommended
// @Summary Рекомендованные товары
// @Tags product
// @Produce json
// @Success 200 {array} object
// @Router /product/recommended [get]
func _swaggerProductRecommended() {}

// ProductUserProducts
// @Summary Товары пользователя
// @Tags product
// @Produce json
// @Param id path int true "User id"
// @Success 200 {array} object
// @Router /product/user-products/{id} [get]
func _swaggerProductUserProducts() {}

// ProductAddToFavorites
// @Summary Добавить товар в избранное
// @Security SessionId
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]interface{}
// @Router /product/add-to-favorites/{id} [post]
func _swaggerProductAddToFavorites() {}

// ProductRemoveFromFavorites
// @Summary Удалить товар из избранного
// @Security SessionId
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]interface{}
// @Router /product/remove-from-favorites/{id} [delete]
func _swaggerProductRemoveFromFavorites() {}

// ProductMyFavorites
// @Summary Мое избранное
// @Security SessionId
// @Tags product
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /product/my-favorites [get]
func _swaggerProductMyFavorites() {}

// ProductToggle
// @Summary Скрыть или опубликовать свой товар
// @Security SessionId
// @Tags product
// @Produce json
// @Param id path int true "Product id"
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /product/toggle-product/{id} [put]
func _swaggerProductToggle() {}

// ProductAllToModerate
// @Summary Товары на модерации (admin)
// @Tags product-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /product/all-products-to-moderate [get]
func _swaggerProductAllToModerate() {}

// ProductPromoted
// @Summary Продвигаемые товары (admin)
// @Tags product-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /product/promoted-products [get]
func _swaggerProductPromoted() {}

// ProductTogglePromotion
// @Summary Включить или выключить продвижение (admin)
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
// @Summary Подать апелляцию на отзыв
// @Security SessionId
// @Tags review
// @Accept json
// @Produce json
// @Param body body swaggerReviewAppealCreate true "Тело"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /review/appeals [post]
func _swaggerReviewCreateAppeal() {}

// ReviewMyAppeals
// @Summary Мои апелляции на отзывы
// @Security SessionId
// @Tags review
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /review/my-appeals [get]
func _swaggerReviewMyAppeals() {}

// ReviewAllAppeals
// @Summary Все апелляции на отзывы (moderator/admin)
// @Tags review-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /review/all-appeals [get]
func _swaggerReviewAllAppeals() {}

// ReviewResolveAppeal
// @Summary Разрешить апелляцию на отзыв (moderator/admin)
// @Tags review-admin
// @Accept json
// @Produce json
// @Param id path int true "Appeal id"
// @Param body body swaggerResolveAppeal true "Решение"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /review/resolve-appeal/{id} [put]
func _swaggerReviewResolveAppeal() {}

type swaggerReviewAppealCreate struct {
	ReviewID int    `json:"reviewId" example:"12"`
	Reason   string `json:"reason" example:"Отзыв содержит недостоверные сведения"`
}

type swaggerResolveAppeal struct {
	Status        string `json:"status" example:"RESOLVED"`
	ModeratorNote string `json:"moderatorNote" example:"Апелляция рассмотрена"`
}

// --- reservation ---

// ReservationCreate
// @Summary Создать резерв товара
// @Security SessionId
// @Tags reservation
// @Accept json
// @Produce json
// @Param body body swaggerReservationCreate true "Тело"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/ [post]
func _swaggerReservationCreate() {}

// ReservationMy
// @Summary Мои резервы
// @Security SessionId
// @Tags reservation
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/my [get]
func _swaggerReservationMy() {}

// ReservationProductInfo
// @Summary Информация о резервировании товара
// @Tags reservation
// @Produce json
// @Param productId path int true "Product id"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /reservations/product/{productId} [get]
func _swaggerReservationProductInfo() {}

// ReservationCancelByBuyer
// @Summary Отмена резерва покупателем
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
// @Summary Отмена резерва продавцом
// @Security SessionId
// @Tags reservation
// @Accept json
// @Produce json
// @Param id path int true "Reservation id"
// @Param body body swaggerReservationCancelReason true "Причина"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/{id}/cancel-by-seller [post]
func _swaggerReservationCancelBySeller() {}

// ReservationCancel
// @Summary Универсальная отмена резерва
// @Security SessionId
// @Tags reservation
// @Accept json
// @Produce json
// @Param id path int true "Reservation id"
// @Param body body swaggerReservationCancelReason false "Причина"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/{id}/cancel [post]
func _swaggerReservationCancel() {}

// ReservationExtend
// @Summary Продлить резерв
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
// @Summary Обновить настройки резервирования товара
// @Security SessionId
// @Tags reservation
// @Accept json
// @Produce json
// @Param body body swaggerReservationProductSettings true "Тело"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /reservations/product-settings [put]
func _swaggerReservationUpdateProductSettings() {}

type swaggerReservationCreate struct {
	ProductID int     `json:"productId" example:"6157119"`
	Hours     *int    `json:"hours,omitempty" example:"24"`
	Note      *string `json:"note,omitempty" example:"Прошу придержать товар до вечера"`
}

type swaggerReservationCancelReason struct {
	Reason *string `json:"reason,omitempty" example:"Покупатель не вышел на связь"`
}

type swaggerReservationProductSettings struct {
	ProductID         int  `json:"productId" example:"6157119"`
	AllowReservations bool `json:"allowReservations" example:"true"`
}

// --- statistics extra ---

// StatisticsSearchQueries
// @Summary Статистика поисковых запросов
// @Security SessionId
// @Tags statistics
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /statistics/search-queries [get]
func _swaggerStatisticsSearchQueries() {}

// StatisticsCabinetDashboard
// @Summary Дашборд личного кабинета
// @Security SessionId
// @Tags statistics
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /statistics/cabinet-dashboard [get]
func _swaggerStatisticsCabinetDashboard() {}

// --- banner extra ---

// BannerMyStats
// @Summary Моя статистика по баннерам
// @Security SessionId
// @Tags banner
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /banner/my-stats/all [get]
func _swaggerBannerMyStats() {}

// BannerViewTrack
// @Summary Зафиксировать просмотр баннера
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Success 200 {object} map[string]string
// @Router /banner/{id}/view [post]
func _swaggerBannerViewTrack() {}

// BannerStats
// @Summary Статистика конкретного баннера
// @Security SessionId
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Router /banner/{id}/stats [get]
func _swaggerBannerStats() {}

// BannerUpdate
// @Summary Обновить баннер (admin)
// @Tags banner
// @Accept json
// @Produce json
// @Param id path int true "Banner id"
// @Param body body swaggerBannerUpdate true "Тело"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /banner/{id} [put]
func _swaggerBannerUpdate() {}

// BannerDelete
// @Summary Удалить баннер (admin)
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /banner/{id} [delete]
func _swaggerBannerDelete() {}

// BannerGetByID
// @Summary Получить баннер по id
// @Tags banner
// @Produce json
// @Param id path int true "Banner id"
// @Success 200 {object} object
// @Failure 404 {object} map[string]interface{}
// @Router /banner/{id} [get]
func _swaggerBannerGetByID() {}

type swaggerBannerUpdate struct {
	Name          *string `json:"name,omitempty" example:"Летняя акция"`
	PhotoURL      *string `json:"photoUrl,omitempty" example:"https://cdn.example.com/banner.jpg"`
	Place         *string `json:"place,omitempty" example:"PRODUCT_FEED"`
	NavigateToURL *string `json:"navigateToUrl,omitempty" example:"https://torguisam.ru/product/6157119"`
}

// --- moderation extra ---

// ModerationSummary
// @Summary Сводка по модерации
// @Tags moderation-admin
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/summary [get]
func _swaggerModerationSummary() {}

// ModerationAuditLogs
// @Summary Журнал аудита модерации
// @Tags moderation-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/audit-logs [get]
func _swaggerModerationAuditLogs() {}

// ModerationAppeals
// @Summary Список апелляций модерации
// @Tags moderation-admin
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/appeals [get]
func _swaggerModerationAppeals() {}

// ModerationReviewAppeal
// @Summary Рассмотреть апелляцию модерации
// @Tags moderation-admin
// @Accept json
// @Produce json
// @Param id path int true "Appeal id"
// @Param body body swaggerModerationAppealReview true "Решение"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /admin/moderation/appeals/{id}/review [put]
func _swaggerModerationReviewAppeal() {}

type swaggerModerationAppealReview struct {
	Status        string `json:"status" example:"APPROVED"`
	ReviewComment string `json:"reviewComment" example:"Апелляция рассмотрена модератором"`
}

// --- deals extra ---

// DealCreate
// @Summary Создать безопасную сделку
// @Security SessionId
// @Tags deals
// @Accept json
// @Produce json
// @Param body body swaggerDealCreate true "Тело"
// @Success 201 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /deals/ [post]
func _swaggerDealCreate() {}

// DealMyPurchases
// @Summary Мои покупки
// @Security SessionId
// @Tags deals
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /deals/my-purchases [get]
func _swaggerDealMyPurchases() {}

// DealMySales
// @Summary Мои продажи
// @Security SessionId
// @Tags deals
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /deals/my-sales [get]
func _swaggerDealMySales() {}

// DealMyAll
// @Summary Все мои сделки
// @Security SessionId
// @Tags deals
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]interface{}
// @Router /deals/my [get]
func _swaggerDealMyAll() {}

// DealPay
// @Summary Оплатить сделку
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
// @Summary Синхронизировать статус оплаты сделки
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
// @Summary Указать способ передачи товара в CDEK
// @Security SessionId
// @Tags deals
// @Accept json
// @Produce json
// @Param id path int true "Deal id"
// @Param body body swaggerDealHandoff true "Тело"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /deals/{id}/cdek-handoff [post]
func _swaggerDealSetCDEKHandoff() {}

// DealConfirmDelivery
// @Summary Подтвердить получение товара
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
// @Summary Открыть спор по сделке
// @Security SessionId
// @Tags deals
// @Accept json
// @Produce json
// @Param id path int true "Deal id"
// @Param body body swaggerDealDispute true "Причина"
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
// @Summary Получить доступные тарифы CDEK
// @Tags cdek
// @Accept json
// @Produce json
// @Param body body swaggerCDEKCalculateRequest true "Body"
// @Success 200 {object} object
// @Failure 400 {object} map[string]interface{}
// @Router /cdek/tariffs [post]
func _swaggerCDEKTariffs() {}
