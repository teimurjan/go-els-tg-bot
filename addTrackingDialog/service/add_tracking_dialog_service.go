package service

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-els-tg-bot/addTrackingDialog"
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type addTrackingDialogService struct {
	addTrackingDialogRepo addTrackingDialog.AddTrackingDialogRepository
	userRepo              user.UserRepository
	trackingRepo          tracking.TrackingRepository
	statusFetcher         tracking.TrackingStatusFetcher
	logger                *logrus.Logger
}

func NewAddTrackingDialogService(
	addTrackingDialogRepo addTrackingDialog.AddTrackingDialogRepository,
	userRepo user.UserRepository,
	trackingRepo tracking.TrackingRepository,
	statusFetcher tracking.TrackingStatusFetcher,
	logger *logrus.Logger,
) *addTrackingDialogService {
	return &addTrackingDialogService{
		addTrackingDialogRepo,
		userRepo,
		trackingRepo,
		statusFetcher,
		logger,
	}
}

func (s *addTrackingDialogService) GetDialogForChat(chatID int64) (*models.AddTrackingDialog, error) {
	user, err := s.userRepo.GetByChatID(chatID)
	if err != nil {
		s.logger.Errorf("Couldn't find a user with chatID=%d", chatID)
		return nil, err
	}
	return s.addTrackingDialogRepo.GetForUser(user.ID)
}

func (s *addTrackingDialogService) StartDialog(chatID int64) (*models.AddTrackingDialog, error) {
	user, err := s.userRepo.GetByChatID(chatID)
	if err != nil {
		s.logger.Errorf("Couldn't find a user with chatID=%d", chatID)
		return nil, err
	}
	dialog := models.AddTrackingDialog{
		UserID: user.ID,
		Step:   1,
	}
	id, err := s.addTrackingDialogRepo.Store(&dialog)
	if err != nil {
		dialogJSON, _ := json.Marshal(dialog)
		s.logger.Errorf("Couldn't store dialog %s", string(dialogJSON))
		return nil, err
	}
	dialog.ID = id
	return &dialog, nil
}

func (s *addTrackingDialogService) UpdateDialogName(dialog *models.AddTrackingDialog, name string) error {
	dialog.FutureTrackingName = name
	dialog.Step = 2
	return s.addTrackingDialogRepo.Update(dialog)
}

func (s *addTrackingDialogService) UpdateDialogTracking(dialog *models.AddTrackingDialog, trackingNumber string) (*models.Tracking, error) {
	status, err := s.statusFetcher.Fetch(trackingNumber)
	if err != nil {
		return nil, err
	}

	tracking := models.Tracking{
		UserID: dialog.UserID,
		Name:   dialog.FutureTrackingName,
		Value:  trackingNumber,
		Status: status,
	}
	ID, err := s.trackingRepo.Store(&tracking)
	tracking.ID = ID
	if err != nil {
		return nil, err
	}

	return &tracking, s.ResetDialog(dialog)
}

func (s *addTrackingDialogService) ResetDialog(dialog *models.AddTrackingDialog) error {
	dialog.Step = 0
	dialog.FutureTrackingName = ""
	return s.addTrackingDialogRepo.Update(dialog)
}
