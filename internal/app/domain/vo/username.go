package vo

type Username struct {
	value string
}

func (u Username) Value() string {
	return u.value
}
