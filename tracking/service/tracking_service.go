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

func (s *trackingService) GetForChat(chatID int64) ([]*models.Tracking, error) {
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

	return trackings, nil
}

func (s *trackingService) Update(tracking *models.Tracking) (bool, error) {
	fetchedTracking, err := s.fetcher.Fetch(tracking.Value)
	if err != nil {
		return false, err
	}

	if fetchedTracking.Status != tracking.Status {
		s.logger.Info(fmt.Sprintf("Tracking number has changed.\nStatus: %s to %s.", tracking.Status, fetchedTracking.Status))
		tracking.Status = fetchedTracking.Status
		s.trackingRepo.Update(tracking)

		return true, nil
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

func (s *trackingService) SyncAll(trackings []*models.Tracking) (chan *models.Tracking, chan error, chan bool) {
	trackingCh := make(chan *models.Tracking)
	errCh := make(chan error)
	doneCh := make(chan bool)

	workerFinishedCount := 0

	worker := func(tracking *models.Tracking) {
		s.logger.Info(fmt.Sprintf("Synchronizing tracking information for %s started.", tracking.Value))

		_, err := s.Update(tracking)
		if err != nil {
			s.logger.Info(fmt.Sprintf("Synchronizing tracking information for %s ended with an error %v.", tracking.Value, err))
			errCh <- err
		} else {
			s.logger.Info(fmt.Sprintf("Synchronizing tracking information for %s got a result.", tracking.Value))
			trackingCh <- tracking
		}

		if workerFinishedCount++; workerFinishedCount == len(trackings) {
			s.logger.Info(fmt.Sprintf("Synchronizing tracking information for %s ended.", tracking.Value))
			doneCh <- true
		}
	}

	for _, tracking := range trackings {
		go worker(tracking)
	}

	return trackingCh, errCh, doneCh
}

func (s *trackingService) GetOnlyUpdated(trackings []*models.Tracking) []*models.Tracking {
	updatedTrackings := make([]*models.Tracking, 0)

	for _, tracking := range trackings {
		updated, err := s.Update(tracking)
		if err == nil && updated {
			updatedTrackings = append(updatedTrackings, tracking)
		}
	}

	return updatedTrackings
}

func (s *trackingService) SyncOnlyUpdated(trackings []*models.Tracking) (chan *models.Tracking, chan error, chan bool) {
	trackingCh := make(chan *models.Tracking)
	errCh := make(chan error)
	doneCh := make(chan bool)

	workerFinishedCount := 0

	worker := func(tracking *models.Tracking) {
		s.logger.Info(fmt.Sprintf("Synchronizing tracking information for %s started.", tracking.Value))

		updated, err := s.Update(tracking)
		if err != nil {
			s.logger.Info(fmt.Sprintf("Synchronizing tracking information for %s ended with an error %v.", tracking.Value, err))
			errCh <- err
		} else if updated {
			s.logger.Info(fmt.Sprintf("Synchronizing tracking information for %s got a result.", tracking.Value))
			trackingCh <- tracking
		}

		if workerFinishedCount++; workerFinishedCount == len(trackings) {
			s.logger.Info(fmt.Sprintf("Synchronizing tracking information for %s ended.", tracking.Value))
			doneCh <- true
		}
	}

	for _, tracking := range trackings {
		s.logger.Info(fmt.Sprintf("Passing tracking %s to the worker.", tracking.Value))
		go worker(tracking)
	}

	return trackingCh, errCh, doneCh
}

func (s *trackingService) GetAllGroupedByUser() (map[*models.User][]*models.Tracking, error) {
	s.logger.Info("Getting trackings grouped by user.")

	groupedByUser := make(map[*models.User][]*models.Tracking)
	users, err := s.userRepo.GetAll()
	if err != nil {
		s.logger.Error(err)
		return groupedByUser, err
	}

	s.logger.Info("All the users are retreived.")

	for _, user := range users {
		s.logger.Info(fmt.Sprintf("Getting trackings for user %d.", user.ID))

		trackings, err := s.trackingRepo.GetForUser(user.ID)
		if err != nil {
			s.logger.Error(err)
			return groupedByUser, err
		}
		groupedByUser[user] = trackings

		s.logger.Info(fmt.Sprintf("User %d has %d trackings.", user.ID, len(trackings)))
	}

	return groupedByUser, nil
}
