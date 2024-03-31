package vo

type Email struct {
	value string
}

func (e Email) Value() string {
	return e.value
}
