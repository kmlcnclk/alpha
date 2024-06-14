package user

type Command struct {
	Id        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Age       int32
}

type CommandSignIn struct {
	Email    string
	Password string
}
