package tracking

type TrackingHandler interface {
	AddTracking(arguments string, chatID int64)
	GetAll(chatID int64)
	DeleteTracking(trackingID int64, chatID int64, messageID int64)
}
