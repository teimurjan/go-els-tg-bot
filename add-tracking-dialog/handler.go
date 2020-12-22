package addTrackingDialog

type AddTrackingDialogHandler interface {
	StartDialog(chatID int64)
	UpdateDialogIfActive(text string, chatID int64)
	ResetDialog(chatID int64)
}
