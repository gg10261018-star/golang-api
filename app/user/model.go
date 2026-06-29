package user

import "time"

type User struct {
	Email    string    `xorm:"pk varchar(50) 'email'"`
	Password string    `xorm:"varchar(255) 'password'"`
	Created  time.Time `xorm:"datetime 'created'"`
	Updated  time.Time `xorm:"datetime 'updated'"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
