package service

import (
	"github.com/teimurjan/go-els-tg-bot/models"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	"github.com/teimurjan/go-els-tg-bot/user"
)

type trackingService struct {
	trackingRepo  tracking.TrackingRepository
	userRepo      user.UserRepository
	statusFetcher tracking.TrackingStatusFetcher
}

func NewTrackingService(
	trackingRepo tracking.TrackingRepository,
	userRepo user.UserRepository,
	statusFetcher tracking.TrackingStatusFetcher,
) *trackingService {
	return &trackingService{
		trackingRepo,
		userRepo,
		statusFetcher,
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

func (s *trackingService) CheckUpdates() ([]*tracking.TrackingUpdate, error) {
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

	return trackingUpdates, nil
}

func (s *trackingService) updateStatus(t *models.Tracking) bool {
	newStatus, _ := s.statusFetcher.Fetch(t.Value)
	if newStatus != t.Status {
		t.Status = newStatus
		s.trackingRepo.UpdateOne(t)
		return true
	}
	return false
}

func (s *trackingService) Delete(trackingID int64) error {
	return s.trackingRepo.Delete(trackingID)
}
