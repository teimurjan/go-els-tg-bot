package user

// UserHandler is an interface for user handler
type UserHandler interface {
	Join(chatID int64)
}
