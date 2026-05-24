# Роли и права доступа (RBAC)

Дата: 20.05.2026

## Иерархия ролей

- `GUEST` (0)
- `USER` (10)
- `USER_VERIFIED` (15)
- `SENIOR_MODERATOR` (70)
- `ADMIN` (90)
- `SUPERADMIN` (100)

## Матрица прав

### Пользователь (`USER`, `USER_VERIFIED`)
- Размещение/редактирование своих объявлений
- Просмотр объявлений
- Платежи и безопасные сделки
- Продвижение объявлений
- Подача заявок на рекламу
- Подача апелляций

### Модератор (`SENIOR_MODERATOR`)
- Очередь модерации (объявления/отзывы)
- Одобрение/отклонение контента
- Просмотр логов ИИ модерации
- Блокировка/разблокировка пользователей
- Доступ к списку сделок, деталям сделок, логам сделок
- Доступ к данным о резервах/бронях
- Доступ к балансовой информации пользователей

### Админ (`ADMIN`)
- Все права модератора
- Назначение ролей (модераторы/админы)
- Настройка тарифов/продвижения
- Аналитика модерации и платформы
- Редактирование пользователей
- Изменение баланса
- Настройка правил/фильтров ИИ

### Супер-админ (`SUPERADMIN`)
- Полный доступ
- Управление бэкапами и обновлениями
- Расширенные интеграции (в т.ч. госсервисы)
- Контур обучения/тюнинга ИИ модерации

## Реализация в коде

- Базовые роли/уровни: `api-go/internal/rbac/rbac.go`
- Пермишены: `api-go/internal/rbac/rbac.go` (`HasPermission`)
- Middleware:
  - `RequireSession`
  - `RequireRoleLevel`
  - `RequireModerator`
  - `RequireAdmin`
  - `RequirePermission`
  - Файл: `api-go/internal/httpserver/middleware/session.go`

## Изменения доступа в API

- Модератору открыт доступ:
  - `GET /admin/moderation/products`
  - `GET /admin/moderation/products/:id`
  - `GET /admin/moderation/audit-logs`
  - `GET /admin/moderation/appeals`
  - `PUT /admin/moderation/appeals/:id/review`
  - `GET /admin/deals/list`
  - `GET /admin/deals/:id`
  - `GET /admin/deals/:id/logs`
  - `GET /user/find-all`
  - `PUT /user/set-balance/:userId`
  - `PUT /user/toggle-banned/:id`

- Только админ и выше:
  - `PATCH /admin/deals/:id/status`
  - `PUT /user/:id/role`
  - `PATCH /user/:id`
  - `DELETE /user/:id`
  - `GET /admin/moderation/summary`

## Compliance и безопасность

- Логи модерации ИИ/модераторов сохраняются в `ModerationAuditLog`.
- Рекомендуемый срок хранения: 12–24 месяца (настройка ретенции БД/архивацией).
- Для защиты от подмены:
  - писать только append-only события в аудит
  - ограничить UPDATE/DELETE на аудит-таблицу
  - вынести архив в отдельное защищенное хранилище
- Шифрование:
  - TLS для внешних интеграций
  - шифрование бэкапов
  - секреты только через env/secret-store
