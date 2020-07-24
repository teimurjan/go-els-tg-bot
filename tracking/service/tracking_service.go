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

// NewTrackingService creates new trackingService instance
func NewTrackingService(
	trackingRepo tracking.TrackingRepository,
	userRepo user.UserRepository,
	statusFetcher tracking.TrackingStatusFetcher,
	logger *logrus.Logger,
) tracking.TrackingService {
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
		s.logger.Error(err)
		return nil, err
	}

	user, err := s.userRepo.GetByChatID(chatID)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	tracking := models.Tracking{
		UserID: user.ID,
		Value:  value,
		Status: status.Status,
		Weight: status.Weight,
		Name:   name,
	}
	id, err := s.trackingRepo.Store(&tracking)
	if err != nil {
		s.logger.Error(err)
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
		s.logger.Error(err)
		return nil, err
	}

	trackings, err := s.trackingRepo.GetForUser(user.ID)
	if err != nil {
		s.logger.Error(err)
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
		s.logger.Error(err)
		return nil, err
	}

	var trackingUpdates []*tracking.TrackingUpdate

	for _, user := range users {
		trackings, err := s.trackingRepo.GetForUser(user.ID)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}

		var updatedTrackings []*models.Tracking
		for i := 0; i < len(trackings); i++ {
			updated, err := s.updateStatus(trackings[i])

			if err != nil {
				s.logger.Error(err)
				return nil, err
			}

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

func (s *trackingService) updateStatus(t *models.Tracking) (bool, error) {
	newStatus, err := s.statusFetcher.Fetch(t.Value)
	if err != nil {
		return false, err
	}

	if newStatus.Status != t.Status {
		trackingJSON, _ := json.Marshal(t)
		s.logger.Info(fmt.Sprintf("%s status changed to %s", trackingJSON, newStatus))

		t.Status = newStatus.Status

		s.trackingRepo.Update(t)
		return true, nil
	}

	if newStatus.Weight != t.Weight {
		t.Weight = newStatus.Weight
		s.trackingRepo.Update(t)
		return false, nil
	}

	return false, nil
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
