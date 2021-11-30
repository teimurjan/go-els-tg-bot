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
	trackingRepo tracking.TrackingRepository
	userRepo     user.UserRepository
	fetcher      tracking.TrackingNumberFetcher
	logger       *logrus.Logger
}

// NewTrackingService creates new trackingService instance
func NewTrackingService(
	trackingRepo tracking.TrackingRepository,
	userRepo user.UserRepository,
	fetcher tracking.TrackingNumberFetcher,
	logger *logrus.Logger,
) tracking.TrackingService {
	return &trackingService{
		trackingRepo,
		userRepo,
		fetcher,
		logger,
	}
}

func (s *trackingService) Create(value string, name string, chatID int64) (*models.Tracking, error) {
	status, err := s.fetcher.Fetch(value)
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
		s.fetchAndUpdateTrackingData(trackings[i])
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
			statusUpdated, _, err := s.fetchAndUpdateTrackingData(trackings[i])

			if err != nil {
				s.logger.Error(err)
				return nil, err
			}

			if statusUpdated {
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

func (s *trackingService) fetchAndUpdateTrackingData(t *models.Tracking) (bool, bool, error) {
	fetchedTracking, err := s.fetcher.Fetch(t.Value)
	if err != nil {
		return false, false, err
	}

	isStatusChanged := fetchedTracking.Status != t.Status
	isWeightChanged := fetchedTracking.Weight != t.Weight

	logMsg := "Tracking number has changed."

	if isWeightChanged {
		logMsg += fmt.Sprintf(" Weight: %s to %s.", t.Weight, fetchedTracking.Weight)
		t.Weight = fetchedTracking.Weight
	}
	if isStatusChanged {
		logMsg += fmt.Sprintf(" Status: %s to %s.", t.Status, fetchedTracking.Status)
		t.Status = fetchedTracking.Status
	}

	if isStatusChanged || isWeightChanged {
		s.logger.Info(logMsg)
		s.trackingRepo.Update(t)
	}

	return isStatusChanged, isWeightChanged, nil
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
