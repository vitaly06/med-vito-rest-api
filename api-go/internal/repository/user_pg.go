package repository

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"med-vito/api-go/internal/domain"
)

type UserPG struct {
	pool *pgxpool.Pool
}

func NewUserPG(pool *pgxpool.Pool) *UserPG {
	return &UserPG{pool: pool}
}

func (r *UserPG) FindUserByLogin(ctx context.Context, login string) (*domain.UserEntity, error) {
	const q = `
		SELECT u."id", u."fullName", u."email", u."phoneNumber", u."password", u."profileType"::text,
		       u."photo", u."rating", u."isAnswersCall", u."roleId", r."name",
		       u."isBanned", u."isResetVerified"
		FROM "User" u
		LEFT JOIN "Role" r ON u."roleId" = r."id"
		WHERE u."email" = $1 OR u."phoneNumber" = $1
		LIMIT 1`
	return r.scanUser(ctx, r.pool.QueryRow(ctx, q, login))
}

func (r *UserPG) FindUserByEmailOrPhone(ctx context.Context, email, phone string) (bool, error) {
	const q = `SELECT 1 FROM "User" WHERE "email" = $1 OR "phoneNumber" = $2 LIMIT 1`
	var one int
	err := r.pool.QueryRow(ctx, q, email, phone).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserPG) FindUserByEmail(ctx context.Context, email string) (*domain.UserEntity, error) {
	const q = `
		SELECT u."id", u."fullName", u."email", u."phoneNumber", u."password", u."profileType"::text,
		       u."photo", u."rating", u."isAnswersCall", u."roleId", r."name",
		       u."isBanned", u."isResetVerified"
		FROM "User" u
		LEFT JOIN "Role" r ON u."roleId" = r."id"
		WHERE u."email" = $1`
	return r.scanUser(ctx, r.pool.QueryRow(ctx, q, email))
}

func (r *UserPG) FindUserByID(ctx context.Context, id int32) (*domain.UserEntity, error) {
	const q = `
		SELECT u."id", u."fullName", u."email", u."phoneNumber", u."password", u."profileType"::text,
		       u."photo", u."rating", u."isAnswersCall", u."roleId", r."name",
		       u."isBanned", u."isResetVerified"
		FROM "User" u
		LEFT JOIN "Role" r ON u."roleId" = r."id"
		WHERE u."id" = $1`
	return r.scanUser(ctx, r.pool.QueryRow(ctx, q, id))
}

func (r *UserPG) scanUser(ctx context.Context, row pgx.Row) (*domain.UserEntity, error) {
	var u domain.UserEntity
	var roleName *string
	err := row.Scan(
		&u.ID, &u.FullName, &u.Email, &u.PhoneNumber, &u.PasswordHash, &u.ProfileType,
		&u.Photo, &u.Rating, &u.IsAnswersCall, &u.RoleID, &roleName,
		&u.IsBanned, &u.IsResetVerified,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	u.RoleName = roleName
	return &u, nil
}

func (r *UserPG) RoleIDByName(ctx context.Context, name string) (int32, error) {
	const q = `SELECT "id" FROM "Role" WHERE "name" = $1`
	var id int32
	if err := r.pool.QueryRow(ctx, q, name).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, err
	}
	return id, nil
}

func (r *UserPG) UserIDExists(ctx context.Context, id int32) (bool, error) {
	const q = `SELECT 1 FROM "User" WHERE "id" = $1`
	var one int
	err := r.pool.QueryRow(ctx, q, id).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// GenerateUniqueUserID — 7-значный id как в Nest id-generator.
func (r *UserPG) GenerateUniqueUserID(ctx context.Context) (int32, error) {
	const minID, maxID = 1_000_000, 9_999_999
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for attempt := 0; attempt < 100; attempt++ {
		id := int32(minID + rng.Intn(maxID-minID+1))
		ok, err := r.UserIDExists(ctx, id)
		if err != nil {
			return 0, err
		}
		if !ok {
			return id, nil
		}
	}
	return 0, fmt.Errorf("не удалось сгенерировать уникальный user id")
}

func (r *UserPG) InsertUser(ctx context.Context, id int32, fullName, email, phone, passwordHash string, roleID int32) error {
	const q = `
		INSERT INTO "User" ("id","fullName","email","phoneNumber","password","roleId","updatedAt")
		VALUES ($1,$2,$3,$4,$5,$6,NOW())`
	_, err := r.pool.Exec(ctx, q, id, fullName, email, phone, passwordHash, roleID)
	return err
}

func (r *UserPG) SetResetVerified(ctx context.Context, userID int32, v bool) error {
	tag, err := r.pool.Exec(ctx, `UPDATE "User" SET "isResetVerified" = $2, "updatedAt" = NOW() WHERE "id" = $1`, userID, v)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *UserPG) UpdatePassword(ctx context.Context, userID int32, passwordHash string) error {
	const q = `UPDATE "User" SET "password" = $2, "isResetVerified" = false, "updatedAt" = NOW() WHERE "id" = $1`
	tag, err := r.pool.Exec(ctx, q, userID, passwordHash)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// UserAdminRow — список пользователей для админки (как findAll в Nest).
type UserAdminRow struct {
	ID            int32
	IsBanned      bool
	FullName      string
	Email         string
	PhoneNumber   string
	Rating        *int32
	ProfileType   string
	Photo         *string
	Balance       float64
	BonusBalance  float64
	ProductsCount int32
}

func (r *UserPG) ListUsersAdmin(ctx context.Context) ([]UserAdminRow, error) {
	const q = `
		SELECT u."id", u."isBanned", u."fullName", u."email", u."phoneNumber", u."rating",
		       u."profileType"::text, u."photo", u."balance", u."bonusBalance",
		       COUNT(p."id")::int
		FROM "User" u
		LEFT JOIN "Product" p ON p."userId" = u."id"
		GROUP BY u."id"`
	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []UserAdminRow
	for rows.Next() {
		var row UserAdminRow
		if err := rows.Scan(&row.ID, &row.IsBanned, &row.FullName, &row.Email, &row.PhoneNumber,
			&row.Rating, &row.ProfileType, &row.Photo, &row.Balance, &row.BonusBalance, &row.ProductsCount); err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, rows.Err()
}

func (r *UserPG) GetBannedState(ctx context.Context, id int32) (bool, error) {
	var b bool
	err := r.pool.QueryRow(ctx, `SELECT "isBanned" FROM "User" WHERE "id" = $1`, id).Scan(&b)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, ErrNotFound
	}
	return b, err
}

func (r *UserPG) SetBanned(ctx context.Context, id int32, banned bool) error {
	tag, err := r.pool.Exec(ctx, `UPDATE "User" SET "isBanned" = $2, "updatedAt" = NOW() WHERE "id" = $1`, id, banned)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *UserPG) DeleteUserByID(ctx context.Context, id int32) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM "User" WHERE "id" = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// UserInfoRow — поля для публичной карточки / сброса лимита.
type UserInfoRow struct {
	ID               int32
	FullName         string
	ProfileType      string
	PhoneNumber      string
	Balance          float64
	BonusBalance     float64
	Photo            *string
	FreeAdsLimit     int32
	UsedFreeAds      int32
	LastAdLimitReset time.Time
}

func (r *UserPG) GetUserInfoRow(ctx context.Context, id int32) (*UserInfoRow, error) {
	const q = `
		SELECT u."id", u."fullName", u."profileType"::text, u."phoneNumber", u."balance", u."bonusBalance",
		       u."photo", u."freeAdsLimit", u."usedFreeAds", u."lastAdLimitReset"
		FROM "User" u WHERE u."id" = $1`
	var row UserInfoRow
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&row.ID, &row.FullName, &row.ProfileType, &row.PhoneNumber, &row.Balance, &row.BonusBalance,
		&row.Photo, &row.FreeAdsLimit, &row.UsedFreeAds, &row.LastAdLimitReset,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *UserPG) ResetMonthlyFreeAds(ctx context.Context, userID int32, now time.Time) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE "User" SET "usedFreeAds" = 0, "lastAdLimitReset" = $2, "updatedAt" = NOW() WHERE "id" = $1`,
		userID, now)
	return err
}

func (r *UserPG) ApprovedReviewsStats(ctx context.Context, reviewedUserID int32) (avg float64, count int32, err error) {
	const q = `
		SELECT COALESCE(AVG("rating"), 0)::float8, COUNT(*)::int
		FROM "Review"
		WHERE "reviewedUserId" = $1 AND "moderateState"::text = 'APPROVED'`
	err = r.pool.QueryRow(ctx, q, reviewedUserID).Scan(&avg, &count)
	return avg, count, err
}

// FreeAdsState — для remaining-free-ads.
type FreeAdsState struct {
	FreeAdsLimit     int32
	UsedFreeAds      int32
	LastAdLimitReset time.Time
}

func (r *UserPG) GetFreeAdsState(ctx context.Context, userID int32) (*FreeAdsState, error) {
	const q = `SELECT "freeAdsLimit", "usedFreeAds", "lastAdLimitReset" FROM "User" WHERE "id" = $1`
	var s FreeAdsState
	err := r.pool.QueryRow(ctx, q, userID).Scan(&s.FreeAdsLimit, &s.UsedFreeAds, &s.LastAdLimitReset)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *UserPG) GetPhoneForUser(ctx context.Context, userID int32) (string, error) {
	var p string
	err := r.pool.QueryRow(ctx, `SELECT "phoneNumber" FROM "User" WHERE "id" = $1`, userID).Scan(&p)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrNotFound
	}
	return p, err
}

func (r *UserPG) PhoneViewExists(ctx context.Context, viewedByID, viewedUserID int32) (bool, error) {
	const q = `SELECT 1 FROM "PhoneNumberView" WHERE "viewedById" = $1 AND "viewedUserId" = $2 LIMIT 1`
	var one int
	err := r.pool.QueryRow(ctx, q, viewedByID, viewedUserID).Scan(&one)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

func (r *UserPG) InsertPhoneView(ctx context.Context, viewedByID, viewedUserID int32) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO "PhoneNumberView" ("viewedById", "viewedUserId") VALUES ($1, $2)
		 ON CONFLICT ("viewedById", "viewedUserId") DO NOTHING`,
		viewedByID, viewedUserID)
	return err
}

func (r *UserPG) GetUserPhoto(ctx context.Context, userID int32) (*string, error) {
	var ph *string
	err := r.pool.QueryRow(ctx, `SELECT "photo" FROM "User" WHERE "id" = $1`, userID).Scan(&ph)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return ph, err
}

// UserSettingsPatch — PATCH настроек (только заданные поля).
type UserSettingsPatch struct {
	FullName      *string
	PhoneNumber   *string
	IsAnswersCall *bool
	ProfileType   *string
	Photo         *string
}

func (r *UserPG) UpdateUserSettings(ctx context.Context, userID int32, p UserSettingsPatch) error {
	// Собираем SET динамически — как spread dto в Prisma.
	var sets []string
	var args []any
	n := 1
	if p.FullName != nil {
		sets = append(sets, fmt.Sprintf(`"fullName" = $%d`, n))
		args = append(args, *p.FullName)
		n++
	}
	if p.PhoneNumber != nil {
		sets = append(sets, fmt.Sprintf(`"phoneNumber" = $%d`, n))
		args = append(args, *p.PhoneNumber)
		n++
	}
	if p.IsAnswersCall != nil {
		sets = append(sets, fmt.Sprintf(`"isAnswersCall" = $%d`, n))
		args = append(args, *p.IsAnswersCall)
		n++
	}
	if p.ProfileType != nil {
		sets = append(sets, fmt.Sprintf(`"profileType" = $%d::"ProfileType"`, n))
		args = append(args, *p.ProfileType)
		n++
	}
	if p.Photo != nil {
		sets = append(sets, fmt.Sprintf(`"photo" = $%d`, n))
		args = append(args, *p.Photo)
		n++
	}
	if len(sets) == 0 {
		return nil
	}
	sets = append(sets, `"updatedAt" = NOW()`)
	args = append(args, userID)
	q := fmt.Sprintf(`UPDATE "User" SET %s WHERE "id" = $%d`, strings.Join(sets, ", "), n)
	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *UserPG) SetEmailVerified(ctx context.Context, userID int32, v bool) error {
	tag, err := r.pool.Exec(ctx, `UPDATE "User" SET "isEmailVerified" = $2, "updatedAt" = NOW() WHERE "id" = $1`, userID, v)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *UserPG) SetBonusBalance(ctx context.Context, userID int32, v float64) error {
	tag, err := r.pool.Exec(ctx, `UPDATE "User" SET "bonusBalance" = $2, "updatedAt" = NOW() WHERE "id" = $1`, userID, v)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// AdminUserPatch — админское обновление пользователя.
type AdminUserPatch struct {
	FullName     *string
	Email        *string
	PhoneNumber  *string
	ProfileType  *string
	BonusBalance *float64
}

func (r *UserPG) FindUserIDByEmail(ctx context.Context, email string) (*int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `SELECT "id" FROM "User" WHERE "email" = $1`, email).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *UserPG) FindUserIDByPhone(ctx context.Context, phone string) (*int32, error) {
	var id int32
	err := r.pool.QueryRow(ctx, `SELECT "id" FROM "User" WHERE "phoneNumber" = $1`, phone).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *UserPG) GetUserEmailPhone(ctx context.Context, userID int32) (email, phone string, err error) {
	err = r.pool.QueryRow(ctx, `SELECT "email", "phoneNumber" FROM "User" WHERE "id" = $1`, userID).Scan(&email, &phone)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", ErrNotFound
	}
	return email, phone, err
}

func (r *UserPG) UpdateUserAdmin(ctx context.Context, userID int32, p AdminUserPatch) error {
	var sets []string
	var args []any
	n := 1
	if p.FullName != nil {
		sets = append(sets, fmt.Sprintf(`"fullName" = $%d`, n))
		args = append(args, *p.FullName)
		n++
	}
	if p.Email != nil {
		sets = append(sets, fmt.Sprintf(`"email" = $%d`, n))
		args = append(args, *p.Email)
		n++
	}
	if p.PhoneNumber != nil {
		sets = append(sets, fmt.Sprintf(`"phoneNumber" = $%d`, n))
		args = append(args, *p.PhoneNumber)
		n++
	}
	if p.ProfileType != nil {
		sets = append(sets, fmt.Sprintf(`"profileType" = $%d::"ProfileType"`, n))
		args = append(args, *p.ProfileType)
		n++
	}
	if p.BonusBalance != nil {
		sets = append(sets, fmt.Sprintf(`"bonusBalance" = $%d`, n))
		args = append(args, *p.BonusBalance)
		n++
	}
	if len(sets) == 0 {
		return nil
	}
	sets = append(sets, `"updatedAt" = NOW()`)
	args = append(args, userID)
	q := fmt.Sprintf(`UPDATE "User" SET %s WHERE "id" = $%d`, strings.Join(sets, ", "), n)
	tag, err := r.pool.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
