package model

import "acltpcipciksntl/internal/app/domain/vo"

// TODO
type User struct {
	ID       vo.ID
	Name     string
	Username vo.Username
	Email    vo.Email
}
