package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrForbiddenChat — участник не buyer/seller этого чата.
var ErrForbiddenChat = errors.New("forbidden chat access")

type ChatPG struct {
	pool *pgxpool.Pool
}

func NewChatPG(pool *pgxpool.Pool) *ChatPG {
	return &ChatPG{pool: pool}
}

func (r *ChatPG) ProductSellerID(ctx context.Context, productID int32) (sellerID int32, ok bool, err error) {
	err = r.pool.QueryRow(ctx, `SELECT "userId" FROM "Product" WHERE id = $1`, productID).Scan(&sellerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return sellerID, true, nil
}

// ChatFullRow — чат с товаром и участниками (для start / ответов).
type ChatFullRow struct {
	ID               int32
	ProductID        *int32
	BuyerID          int32
	SellerID         int32
	IsModerationChat bool
	LastMessageAt    time.Time
	CreatedAt        time.Time

	JoinProductID *int32 // p.id из JOIN (для ответа API)
	ProductName   *string
	ProductImages []string
	ProductPrice  *int32
	ProductDesc   *string

	BuyerName, BuyerPhone   string
	SellerName, SellerPhone string
}

func (r *ChatPG) scanChatFull(row pgx.Row) (*ChatFullRow, error) {
	var c ChatFullRow
	err := row.Scan(
		&c.ID, &c.ProductID, &c.BuyerID, &c.SellerID, &c.IsModerationChat, &c.LastMessageAt, &c.CreatedAt,
		&c.JoinProductID, &c.ProductName, &c.ProductImages, &c.ProductPrice, &c.ProductDesc,
		&c.BuyerName, &c.BuyerPhone, &c.SellerName, &c.SellerPhone,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

const chatFullSelect = `
	c.id, c."productId", c."buyerId", c."sellerId", c."isModerationChat", c."lastMessageAt", c."createdAt",
	p.id, p.name, p.images, p.price, p.description,
	b."fullName", b."phoneNumber", s."fullName", s."phoneNumber"`

const chatFullJoins = `
	FROM "Chat" c
	LEFT JOIN "Product" p ON p.id = c."productId"
	JOIN "User" b ON b.id = c."buyerId"
	JOIN "User" s ON s.id = c."sellerId"`

func (r *ChatPG) FindChatByProductParticipants(ctx context.Context, productID, buyerID, sellerID int32) (*ChatFullRow, error) {
	q := `SELECT ` + chatFullSelect + chatFullJoins + `
		WHERE c."productId" = $1 AND c."buyerId" = $2 AND c."sellerId" = $3`
	row := r.pool.QueryRow(ctx, q, productID, buyerID, sellerID)
	c, err := r.scanChatFull(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ChatPG) InsertProductChat(ctx context.Context, productID, buyerID, sellerID int32) (int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `
		INSERT INTO "Chat" ("buyerId", "sellerId", "productId", "isModerationChat", "lastMessageAt", "createdAt", "updatedAt")
		VALUES ($1,$2,$3,false,NOW(),NOW(),NOW())
		RETURNING id`, buyerID, sellerID, productID).Scan(&id)
	return id, err
}

func (r *ChatPG) GetChatFull(ctx context.Context, chatID int32) (*ChatFullRow, error) {
	q := `SELECT ` + chatFullSelect + chatFullJoins + ` WHERE c.id = $1`
	c, err := r.scanChatFull(r.pool.QueryRow(ctx, q, chatID))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return c, err
}

type ChatListRow struct {
	ChatFullRow
	UnreadBuyer  int32
	UnreadSeller int32
	LMID         *int32
	LMContent    *string
	LMCreated    *time.Time
	LMSender     *int32
	LMRead       *bool
}

func (r *ChatPG) ListChatsForUser(ctx context.Context, userID int32) ([]ChatListRow, error) {
	q := `
		SELECT c.id, c."productId", c."buyerId", c."sellerId", c."isModerationChat", c."lastMessageAt", c."createdAt",
			p.id, p.name, p.images, p.price, p.description,
			b."fullName", b."phoneNumber", s."fullName", s."phoneNumber",
			c."unreadCountBuyer", c."unreadCountSeller",
			m.id, m.content, m."createdAt", m."senderId", m."isRead"
		FROM "Chat" c
		LEFT JOIN "Product" p ON p.id = c."productId"
		JOIN "User" b ON b.id = c."buyerId"
		JOIN "User" s ON s.id = c."sellerId"
		LEFT JOIN "Message" m ON m.id = c."lastMessageId"
		WHERE c."buyerId" = $1 OR c."sellerId" = $1
		ORDER BY c."lastMessageAt" DESC`
	rows, err := r.pool.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ChatListRow
	for rows.Next() {
		var row ChatListRow
		err := rows.Scan(
			&row.ID, &row.ProductID, &row.BuyerID, &row.SellerID, &row.IsModerationChat, &row.LastMessageAt, &row.CreatedAt,
			&row.JoinProductID, &row.ProductName, &row.ProductImages, &row.ProductPrice, &row.ProductDesc,
			&row.BuyerName, &row.BuyerPhone, &row.SellerName, &row.SellerPhone,
			&row.UnreadBuyer, &row.UnreadSeller,
			&row.LMID, &row.LMContent, &row.LMCreated, &row.LMSender, &row.LMRead,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *ChatPG) ChatParticipants(ctx context.Context, chatID int32) (buyerID, sellerID int32, err error) {
	err = r.pool.QueryRow(ctx, `SELECT "buyerId", "sellerId" FROM "Chat" WHERE id = $1`, chatID).Scan(&buyerID, &sellerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, 0, ErrNotFound
	}
	return buyerID, sellerID, err
}

// ChatUnreadCounts — счётчики непрочитанного для getChatInfo (как в Nest).
func (r *ChatPG) ChatUnreadCounts(ctx context.Context, chatID int32) (unreadBuyer, unreadSeller int32, err error) {
	err = r.pool.QueryRow(ctx, `SELECT "unreadCountBuyer", "unreadCountSeller" FROM "Chat" WHERE id = $1`, chatID).Scan(&unreadBuyer, &unreadSeller)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, 0, ErrNotFound
	}
	return unreadBuyer, unreadSeller, err
}

type MessageRow struct {
	ID           int32
	Content      string
	SenderID     int32
	SenderName   string
	IsRead       bool
	ReadAt       *time.Time
	CreatedAt    time.Time
	RelProductID *int32
	RelName      *string
	RelImages    []string
	RelPrice     *int32
	RelModerate  *string
}

func (r *ChatPG) ListMessagesPage(ctx context.Context, chatID int32, limit, offset int) ([]MessageRow, int32, error) {
	var total int32
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*)::int FROM "Message" WHERE "chatId" = $1`, chatID).Scan(&total); err != nil {
		return nil, 0, err
	}
	q := `
		SELECT m.id, m.content, m."senderId", u."fullName", m."isRead", m."readAt", m."createdAt",
			rp.id, rp.name, rp.images, rp.price, rp."moderateState"::text
		FROM "Message" m
		JOIN "User" u ON u.id = m."senderId"
		LEFT JOIN "Product" rp ON rp.id = m."relatedProductId"
		WHERE m."chatId" = $1
		ORDER BY m."createdAt" DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, q, chatID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []MessageRow
	for rows.Next() {
		var m MessageRow
		if err := rows.Scan(&m.ID, &m.Content, &m.SenderID, &m.SenderName, &m.IsRead, &m.ReadAt, &m.CreatedAt,
			&m.RelProductID, &m.RelName, &m.RelImages, &m.RelPrice, &m.RelModerate); err != nil {
			return nil, 0, err
		}
		list = append(list, m)
	}
	return list, total, rows.Err()
}

// InsertChatMessage — сообщение + обновление Chat (как Nest sendMessage).
func (r *ChatPG) InsertChatMessage(ctx context.Context, chatID, senderID int32, content string) (
	msgID int32,
	senderName string,
	createdAt time.Time,
	buyerID, sellerID int32,
	err error,
) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, "", time.Time{}, 0, 0, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `SELECT "buyerId", "sellerId" FROM "Chat" WHERE id = $1`, chatID).Scan(&buyerID, &sellerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, "", time.Time{}, 0, 0, ErrNotFound
	}
	if err != nil {
		return 0, "", time.Time{}, 0, 0, err
	}
	if buyerID != senderID && sellerID != senderID {
		return 0, "", time.Time{}, 0, 0, ErrForbiddenChat
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO "Message" (content, "senderId", "chatId", "isRead", "createdAt", "updatedAt")
		VALUES ($1, $2, $3, false, NOW(), NOW())
		RETURNING id, "createdAt"`, content, senderID, chatID).Scan(&msgID, &createdAt)
	if err != nil {
		return 0, "", time.Time{}, 0, 0, err
	}

	err = tx.QueryRow(ctx, `SELECT "fullName" FROM "User" WHERE id = $1`, senderID).Scan(&senderName)
	if err != nil {
		return 0, "", time.Time{}, 0, 0, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE "Chat" SET
			"lastMessageId" = $1,
			"lastMessageAt" = NOW(),
			"unreadCountSeller" = "unreadCountSeller" + CASE WHEN $2 = "buyerId" THEN 1 ELSE 0 END,
			"unreadCountBuyer" = "unreadCountBuyer" + CASE WHEN $2 = "sellerId" THEN 1 ELSE 0 END,
			"updatedAt" = NOW()
		WHERE id = $3`, msgID, senderID, chatID)
	if err != nil {
		return 0, "", time.Time{}, 0, 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, "", time.Time{}, 0, 0, err
	}
	return msgID, senderName, createdAt, buyerID, sellerID, nil
}

// MarkChatMessagesRead — как Nest markMessagesAsRead.
func (r *ChatPG) MarkChatMessagesRead(ctx context.Context, chatID, userID int32) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var buyerID, sellerID int32
	err = tx.QueryRow(ctx, `SELECT "buyerId", "sellerId" FROM "Chat" WHERE id = $1`, chatID).Scan(&buyerID, &sellerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if buyerID != userID && sellerID != userID {
		return ErrForbiddenChat
	}

	_, err = tx.Exec(ctx, `
		UPDATE "Message" SET "isRead" = true, "readAt" = NOW(), "updatedAt" = NOW()
		WHERE "chatId" = $1 AND "senderId" <> $2 AND "isRead" = false`, chatID, userID)
	if err != nil {
		return err
	}

	if buyerID == userID {
		_, err = tx.Exec(ctx, `UPDATE "Chat" SET "unreadCountBuyer" = 0, "updatedAt" = NOW() WHERE id = $1`, chatID)
	} else {
		_, err = tx.Exec(ctx, `UPDATE "Chat" SET "unreadCountSeller" = 0, "updatedAt" = NOW() WHERE id = $1`, chatID)
	}
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
