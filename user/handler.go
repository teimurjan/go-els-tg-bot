package user

// UserHandler is an interface for user handler
type UserHandler interface {
	Join(chatID int64)
	RequestLanguageChange(chatID int64)
	ChangeLanguage(language string, chatID int64, messageID int64)
}
