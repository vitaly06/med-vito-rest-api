package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSupportTicketNotFound = errors.New("support: ticket not found")

type SupportPG struct {
	pool *pgxpool.Pool
}

func NewSupportPG(pool *pgxpool.Pool) *SupportPG {
	return &SupportPG{pool: pool}
}

// SupportUser — вложенный user/moderator/author (как Prisma select).
type SupportUser struct {
	ID       int32  `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

type SupportRoleName struct {
	Name string `json:"name"`
}

// SupportTicketFull — тикет с user/moderator (без сообщений).
type SupportTicketFull struct {
	ID          int32        `json:"id"`
	Theme       string       `json:"theme"`
	Subject     string       `json:"subject"`
	Status      string       `json:"status"`
	Priority    string       `json:"priority"`
	UserID      int32        `json:"userId"`
	ModeratorID *int32       `json:"moderatorId,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	User        SupportUser  `json:"user"`
	Moderator   *SupportUser `json:"moderator,omitempty"`
}

// SupportLastMessage — последнее сообщение в списке тикетов.
type SupportLastMessage struct {
	ID        int32       `json:"id"`
	Text      string      `json:"text"`
	SentAt    time.Time   `json:"sentAt"`
	AuthorID  int32       `json:"authorId"`
	Author    SupportUser `json:"author"`
}

// SupportTicketListItem — элемент списка с lastMessage.
type SupportTicketListItem struct {
	SupportTicketFull
	LastMessage *SupportLastMessage `json:"lastMessage"`
}

// SupportMessageOut — сообщение с автором и role.name.
type SupportMessageOut struct {
	ID       int32              `json:"id"`
	TicketID int32              `json:"ticketId"`
	AuthorID int32              `json:"authorId"`
	Text     string             `json:"text"`
	SentAt   time.Time          `json:"sentAt"`
	Author   SupportAuthorOut `json:"author"`
}

type SupportAuthorOut struct {
	ID       int32            `json:"id"`
	FullName string           `json:"fullName"`
	Email    string           `json:"email"`
	Role     *SupportRoleName `json:"role,omitempty"`
}

// CreateTicket — транзакция: тикет + первое сообщение; возвращает тикет с user/moderator (без messages).
func (r *SupportPG) CreateTicket(ctx context.Context, userID int32, theme, subject, priority, message string) (*SupportTicketFull, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var id int32
	err = tx.QueryRow(ctx, `
		INSERT INTO "SupportTicket" (theme, subject, status, priority, "userId", "moderatorId", "createdAt", "updatedAt")
		VALUES ($1::"TicketTheme", $2, 'OPEN'::"TicketStatus", $3::"TicketPriority", $4, NULL, NOW(), NOW())
		RETURNING id`,
		theme, subject, priority, userID,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO "SupportMessage" ("ticketId", "authorId", text, "sentAt") VALUES ($1, $2, $3, NOW())`, id, userID, message); err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return r.GetTicketFull(ctx, id)
}

func (r *SupportPG) GetTicketFull(ctx context.Context, ticketID int32) (*SupportTicketFull, error) {
	q := `
		SELECT t.id, t.theme::text, t.subject, t.status::text, t.priority::text, t."userId", t."moderatorId", t."createdAt", t."updatedAt",
			u.id, u."fullName", u.email,
			m.id, m."fullName", m.email
		FROM "SupportTicket" t
		JOIN "User" u ON u.id = t."userId"
		LEFT JOIN "User" m ON m.id = t."moderatorId"
		WHERE t.id = $1`
	row := r.pool.QueryRow(ctx, q, ticketID)
	var t SupportTicketFull
	var mID *int32
	var mFn, mEm *string
	err := row.Scan(
		&t.ID, &t.Theme, &t.Subject, &t.Status, &t.Priority, &t.UserID, &t.ModeratorID, &t.CreatedAt, &t.UpdatedAt,
		&t.User.ID, &t.User.FullName, &t.User.Email,
		&mID, &mFn, &mEm,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrSupportTicketNotFound
	}
	if err != nil {
		return nil, err
	}
	if mID != nil && mFn != nil && mEm != nil {
		t.Moderator = &SupportUser{ID: *mID, FullName: *mFn, Email: *mEm}
	}
	return &t, nil
}

func (r *SupportPG) GetTicketOwnerID(ctx context.Context, ticketID int32) (userID int32, status string, err error) {
	err = r.pool.QueryRow(ctx, `SELECT "userId", status::text FROM "SupportTicket" WHERE id = $1`, ticketID).Scan(&userID, &status)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, "", ErrSupportTicketNotFound
	}
	return userID, status, err
}

// ListTickets — список с lastMessage; onlyUserID != nil — только тикеты пользователя.
func (r *SupportPG) ListTickets(ctx context.Context, onlyUserID *int32, status, priority, theme *string, search *string, searchUserName bool, page, limit int) ([]SupportTicketListItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	var where strings.Builder
	var args []any
	where.WriteString(` WHERE 1=1 `)
	n := 1
	if onlyUserID != nil {
		fmt.Fprintf(&where, ` AND t."userId" = $%d `, n)
		args = append(args, *onlyUserID)
		n++
	}
	if status != nil && *status != "" {
		fmt.Fprintf(&where, ` AND t.status = $%d::"TicketStatus" `, n)
		args = append(args, *status)
		n++
	}
	if priority != nil && *priority != "" {
		fmt.Fprintf(&where, ` AND t.priority = $%d::"TicketPriority" `, n)
		args = append(args, *priority)
		n++
	}
	if theme != nil && *theme != "" {
		fmt.Fprintf(&where, ` AND t.theme = $%d::"TicketTheme" `, n)
		args = append(args, *theme)
		n++
	}
	if search != nil && strings.TrimSpace(*search) != "" {
		pat := "%" + strings.TrimSpace(*search) + "%"
		if searchUserName {
			fmt.Fprintf(&where, ` AND (t.subject ILIKE $%d OR u."fullName" ILIKE $%d OR EXISTS (SELECT 1 FROM "SupportMessage" sm WHERE sm."ticketId" = t.id AND sm.text ILIKE $%d)) `, n, n, n)
		} else {
			fmt.Fprintf(&where, ` AND (t.subject ILIKE $%d OR EXISTS (SELECT 1 FROM "SupportMessage" sm WHERE sm."ticketId" = t.id AND sm.text ILIKE $%d)) `, n, n)
		}
		args = append(args, pat)
		n++
	}

	countQ := `SELECT COUNT(*)::int FROM "SupportTicket" t JOIN "User" u ON u.id = t."userId" ` + where.String()
	var total int
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	dataQ := `
		SELECT t.id, t.theme::text, t.subject, t.status::text, t.priority::text, t."userId", t."moderatorId", t."createdAt", t."updatedAt",
			u.id, u."fullName", u.email,
			modu.id, modu."fullName", modu.email,
			lm.id, lm.text, lm."sentAt", lm."authorId",
			lma.id, lma."fullName", lma.email
		FROM "SupportTicket" t
		JOIN "User" u ON u.id = t."userId"
		LEFT JOIN "User" modu ON modu.id = t."moderatorId"
		LEFT JOIN LATERAL (
			SELECT sm.id, sm.text, sm."sentAt", sm."authorId"
			FROM "SupportMessage" sm
			WHERE sm."ticketId" = t.id
			ORDER BY sm."sentAt" DESC
			LIMIT 1
		) lm ON true
		LEFT JOIN "User" lma ON lma.id = lm."authorId"
	` + where.String() + fmt.Sprintf(` ORDER BY t."updatedAt" DESC LIMIT $%d OFFSET $%d`, n, n+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, dataQ, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []SupportTicketListItem
	for rows.Next() {
		var item SupportTicketListItem
		var mID *int32
		var mFn, mEm *string
		var lmID *int32
		var lmText *string
		var lmAt *time.Time
		var lmAuth *int32
		var lmaID *int32
		var lmaFN, lmaEm *string

		err := rows.Scan(
			&item.ID, &item.Theme, &item.Subject, &item.Status, &item.Priority, &item.UserID, &item.ModeratorID, &item.CreatedAt, &item.UpdatedAt,
			&item.User.ID, &item.User.FullName, &item.User.Email,
			&mID, &mFn, &mEm,
			&lmID, &lmText, &lmAt, &lmAuth,
			&lmaID, &lmaFN, &lmaEm,
		)
		if err != nil {
			return nil, 0, err
		}
		if mID != nil && mFn != nil && mEm != nil {
			item.Moderator = &SupportUser{ID: *mID, FullName: *mFn, Email: *mEm}
		}
		if lmID != nil && lmText != nil && lmAt != nil && lmAuth != nil && lmaID != nil && lmaFN != nil && lmaEm != nil {
			item.LastMessage = &SupportLastMessage{
				ID:       *lmID,
				Text:     *lmText,
				SentAt:   *lmAt,
				AuthorID: *lmAuth,
				Author:   SupportUser{ID: *lmaID, FullName: *lmaFN, Email: *lmaEm},
			}
		}
		out = append(out, item)
	}
	return out, total, rows.Err()
}

func (r *SupportPG) ListMessagesForTicket(ctx context.Context, ticketID int32) ([]SupportMessageOut, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT m.id, m."ticketId", m."authorId", m.text, m."sentAt",
			u.id, u."fullName", u.email, r.name
		FROM "SupportMessage" m
		JOIN "User" u ON u.id = m."authorId"
		LEFT JOIN "Role" r ON r.id = u."roleId"
		WHERE m."ticketId" = $1
		ORDER BY m."sentAt" ASC`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []SupportMessageOut
	for rows.Next() {
		var m SupportMessageOut
		var roleName *string
		if err := rows.Scan(&m.ID, &m.TicketID, &m.AuthorID, &m.Text, &m.SentAt,
			&m.Author.ID, &m.Author.FullName, &m.Author.Email, &roleName); err != nil {
			return nil, err
		}
		if roleName != nil {
			m.Author.Role = &SupportRoleName{Name: *roleName}
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// SendMessageTx — вставка сообщения + обновление тикета (как Nest).
func (r *SupportPG) SendMessageTx(ctx context.Context, ticketID, authorID int32, text string, isModerator bool, currentStatus string) (*SupportMessageOut, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var mid int32
	var sentAt time.Time
	err = tx.QueryRow(ctx, `
		INSERT INTO "SupportMessage" ("ticketId", "authorId", text, "sentAt") VALUES ($1, $2, $3, NOW())
		RETURNING id, "sentAt"`, ticketID, authorID, text,
	).Scan(&mid, &sentAt)
	if err != nil {
		return nil, err
	}

	if isModerator {
		_, err = tx.Exec(ctx, `
			UPDATE "SupportTicket" SET status = 'IN_PROGRESS'::"TicketStatus", "moderatorId" = $2, "updatedAt" = NOW() WHERE id = $1`,
			ticketID, authorID)
	} else if currentStatus == "RESOLVED" {
		_, err = tx.Exec(ctx, `UPDATE "SupportTicket" SET status = 'OPEN'::"TicketStatus", "updatedAt" = NOW() WHERE id = $1`, ticketID)
	} else {
		_, err = tx.Exec(ctx, `UPDATE "SupportTicket" SET "updatedAt" = NOW() WHERE id = $1`, ticketID)
	}
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, `
		SELECT m.id, m."ticketId", m."authorId", m.text, m."sentAt",
			u.id, u."fullName", u.email, r.name
		FROM "SupportMessage" m
		JOIN "User" u ON u.id = m."authorId"
		LEFT JOIN "Role" r ON r.id = u."roleId"
		WHERE m.id = $1`, mid)
	var m SupportMessageOut
	var roleName *string
	if err := row.Scan(&m.ID, &m.TicketID, &m.AuthorID, &m.Text, &m.SentAt,
		&m.Author.ID, &m.Author.FullName, &m.Author.Email, &roleName); err != nil {
		return nil, err
	}
	if roleName != nil {
		m.Author.Role = &SupportRoleName{Name: *roleName}
	}
	return &m, nil
}

func (r *SupportPG) UpdateTicket(ctx context.Context, ticketID, moderatorID int32, status, priority *string) (*SupportTicketFull, error) {
	sets := []string{`"moderatorId" = $2`, `"updatedAt" = NOW()`}
	args := []any{ticketID, moderatorID}
	n := 3
	if status != nil && *status != "" {
		sets = append(sets, fmt.Sprintf(`status = $%d::"TicketStatus"`, n))
		args = append(args, *status)
		n++
	}
	if priority != nil && *priority != "" {
		sets = append(sets, fmt.Sprintf(`priority = $%d::"TicketPriority"`, n))
		args = append(args, *priority)
		n++
	}
	q := fmt.Sprintf(`UPDATE "SupportTicket" SET %s WHERE id = $1`, strings.Join(sets, ", "))
	cmd, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, ErrSupportTicketNotFound
	}
	return r.GetTicketFull(ctx, ticketID)
}

func (r *SupportPG) AssignTicket(ctx context.Context, ticketID, moderatorID int32) (*SupportTicketFull, error) {
	cmd, err := r.pool.Exec(ctx, `
		UPDATE "SupportTicket" SET "moderatorId" = $2, status = 'IN_PROGRESS'::"TicketStatus", "updatedAt" = NOW() WHERE id = $1`,
		ticketID, moderatorID)
	if err != nil {
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, ErrSupportTicketNotFound
	}
	return r.GetTicketFull(ctx, ticketID)
}

type SupportStats struct {
	Total      int `json:"total"`
	ByStatus   map[string]int `json:"byStatus"`
	ByTheme    []struct {
		Theme string `json:"theme"`
		Count int    `json:"count"`
	} `json:"byTheme"`
	ByPriority []struct {
		Priority string `json:"priority"`
		Count    int    `json:"count"`
	} `json:"byPriority"`
}

func (r *SupportPG) TicketStats(ctx context.Context) (*SupportStats, error) {
	var s SupportStats
	s.ByStatus = map[string]int{"open": 0, "inProgress": 0, "resolved": 0, "closed": 0}
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM "SupportTicket"`).Scan(&s.Total); err != nil {
		return nil, err
	}
	var o, ip, rs, cl int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM "SupportTicket" WHERE status = 'OPEN'`).Scan(&o); err != nil {
		return nil, err
	}
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM "SupportTicket" WHERE status = 'IN_PROGRESS'`).Scan(&ip); err != nil {
		return nil, err
	}
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM "SupportTicket" WHERE status = 'RESOLVED'`).Scan(&rs); err != nil {
		return nil, err
	}
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM "SupportTicket" WHERE status = 'CLOSED'`).Scan(&cl); err != nil {
		return nil, err
	}
	s.ByStatus["open"] = o
	s.ByStatus["inProgress"] = ip
	s.ByStatus["resolved"] = rs
	s.ByStatus["closed"] = cl

	trows, err := r.pool.Query(ctx, `SELECT theme::text, COUNT(*)::int FROM "SupportTicket" GROUP BY theme`)
	if err != nil {
		return nil, err
	}
	defer trows.Close()
	for trows.Next() {
		var th string
		var c int
		if err := trows.Scan(&th, &c); err != nil {
			return nil, err
		}
		s.ByTheme = append(s.ByTheme, struct {
			Theme string `json:"theme"`
			Count int    `json:"count"`
		}{Theme: th, Count: c})
	}
	if err := trows.Err(); err != nil {
		return nil, err
	}

	prows, err := r.pool.Query(ctx, `SELECT priority::text, COUNT(*)::int FROM "SupportTicket" GROUP BY priority`)
	if err != nil {
		return nil, err
	}
	defer prows.Close()
	for prows.Next() {
		var pr string
		var c int
		if err := prows.Scan(&pr, &c); err != nil {
			return nil, err
		}
		s.ByPriority = append(s.ByPriority, struct {
			Priority string `json:"priority"`
			Count    int    `json:"count"`
		}{Priority: pr, Count: c})
	}
	return &s, prows.Err()
}
