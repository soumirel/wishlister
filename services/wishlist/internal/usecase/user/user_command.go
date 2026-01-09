package user

type GetUsersCommand struct {
}

type GetUserCommand struct {
	UserID string
}

type CreateUserCommand struct {
	Name string
}

type DeleteUserCommand struct {
	UserID string
}
