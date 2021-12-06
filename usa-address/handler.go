package usaAddress

type UsaAddressHandler interface {
	GetAddress(chatID int64)
	CheckDiff()
}
