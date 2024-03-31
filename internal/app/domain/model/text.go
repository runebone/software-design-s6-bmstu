package model

import (
	"acltpcipciksntl/internal/app/domain/vo"
	"time"
)

type Text struct {
	ID         vo.ID
	Content    string // XXX
	Title      string
	Language   string    // FIXME
	DateLoaded time.Time // XXX
}
