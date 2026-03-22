package domain

// UserEntity — строка User + имя роли (для сессии / guards). Пароль не отдаём в JSON.
type UserEntity struct {
	ID            int32   `json:"id"`
	FullName      string  `json:"fullName"`
	Email         string  `json:"email"`
	PhoneNumber   string  `json:"phoneNumber"`
	PasswordHash  string  `json:"-"`
	ProfileType   string  `json:"profileType"`
	Photo         *string `json:"photo,omitempty"`
	Rating        *int32  `json:"rating,omitempty"`
	IsAnswersCall *bool   `json:"isAnswersCall,omitempty"`
	RoleID        *int32  `json:"roleId,omitempty"`
	RoleName      *string `json:"role,omitempty"`
	IsBanned      bool    `json:"isBanned"`
	IsResetVerified bool  `json:"isResetVerified"`
}

// MeResponse — как getCurrentUser в Nest.
type MeResponse struct {
	ID            int32   `json:"id"`
	Email         string  `json:"email"`
	FullName      string  `json:"fullName"`
	PhoneNumber   string  `json:"phoneNumber"`
	ProfileType   string  `json:"profileType"`
	Photo         *string `json:"photo"`
	Rating        *int32  `json:"rating,omitempty"`
	IsAnswersCall *bool   `json:"isAnswersCall,omitempty"`
	Role          *string `json:"role"`
}
