package domain

// Log — строка из таблицы "Log" (Prisma), те же поля в JSON что у Nest.
type Log struct {
	ID     int32  `json:"id"`
	UserID int32  `json:"userId"`
	Action string `json:"action"`
}
