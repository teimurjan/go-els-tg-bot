package service

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type trackingService struct {
	trackingRepo  tracking.TrackingRepository
	userRepo      user.UserRepository
	statusFetcher tracking.TrackingStatusFetcher
	logger        *logrus.Logger
}

func NewTrackingService(
	trackingRepo tracking.TrackingRepository,
	userRepo user.UserRepository,
	statusFetcher tracking.TrackingStatusFetcher,
	logger *logrus.Logger,
) *trackingService {
	return &trackingService{
		trackingRepo,
		userRepo,
		statusFetcher,
		logger,
	}
}

func (s *trackingService) Create(value string, name string, chatID int64) (*models.Tracking, error) {
	status, err := s.statusFetcher.Fetch(value)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	tracking := models.Tracking{
		UserID: user.ID,
		Value:  value,
		Status: status,
		Name:   name,
	}
	id, err := s.trackingRepo.Store(&tracking)
	if err != nil {
		return nil, err
	}

	tracking.ID = id

	trackingJSON, _ := json.Marshal(tracking)
	s.logger.Info("Tracking created: " + string(trackingJSON))

	return &tracking, nil
}

func (s *trackingService) GetAll(chatID int64) ([]*models.Tracking, error) {
	user, err := s.userRepo.GetByChatID(chatID)
	if err != nil {
		return nil, err
	}

	trackings, err := s.trackingRepo.GetForUser(user.ID)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(trackings); i++ {
		s.updateStatus(trackings[i])
	}

	return trackings, nil
}

func (s *trackingService) GetUpdates() ([]*tracking.TrackingUpdate, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var trackingUpdates []*tracking.TrackingUpdate

	for _, user := range users {
		trackings, err := s.trackingRepo.GetForUser(user.ID)
		if err != nil {
			return nil, err
		}

		var updatedTrackings []*models.Tracking
		for i := 0; i < len(trackings); i++ {
			updated := s.updateStatus(trackings[i])
			if updated {
				updatedTrackings = append(updatedTrackings, trackings[i])
			}
		}

		for _, t := range updatedTrackings {
			trackingUpdates = append(
				trackingUpdates,
				tracking.NewTrackingUpdate(user, t),
			)
		}
	}

	if len(trackingUpdates) == 0 {
		s.logger.Info("No tracking updates found.")
	}

	return trackingUpdates, nil
}

func (s *trackingService) updateStatus(t *models.Tracking) bool {
	newStatus, _ := s.statusFetcher.Fetch(t.Value)
	if newStatus != t.Status {
		trackingJSON, _ := json.Marshal(t)
		s.logger.Info(fmt.Sprintf("%s status changed to %s", trackingJSON, newStatus))

		t.Status = newStatus
		s.trackingRepo.UpdateOne(t)
		return true
	}
	return false
}

func (s *trackingService) Delete(trackingID int64) error {
	err := s.trackingRepo.Delete(trackingID)
	if err == nil {
		s.logger.Info(fmt.Sprintf("Tracking(ID=%d) has been deleted.", trackingID))
	} else {
		s.logger.Error(fmt.Sprintf("Tracking(ID=%d) couldn't be deleted because of %s.", trackingID, err.Error()))
	}
	return err
}
