package job

type Command struct {
	Id                string
	BusinessAccountID string
	Name              string
	Description       string
	Price             float32
	Category          string
}
