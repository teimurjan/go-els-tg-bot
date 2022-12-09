package service

import (
	"database/sql"
	"fmt"

	"github.com/r3labs/diff"
	"github.com/sirupsen/logrus"
	"github.com/teimurjan/go-els-tg-bot/models"
	usaAddress "github.com/teimurjan/go-els-tg-bot/usa-address"
)

type usaAddressService struct {
	usaAddressRepo usaAddress.UsaAddressRepository
	fetcher        usaAddress.UsaAddressFetcher
	logger         *logrus.Logger
}

// NewUsaAddressService creates new usaAddressService instance
func NewUsaAddressService(
	usaAddressRepo usaAddress.UsaAddressRepository,
	fetcher usaAddress.UsaAddressFetcher,
	logger *logrus.Logger,
) usaAddress.UsaAddressService {
	return &usaAddressService{
		usaAddressRepo,
		fetcher,
		logger,
	}
}

func (s *usaAddressService) GetFirst() (*models.UsaAddress, error) {
	address, err := s.usaAddressRepo.GetFirst()
	if err != nil {
		if err == sql.ErrNoRows {
			newAddress, err := s.fetcher.Fetch()
			if err != nil {
				return nil, err
			}

			ID, err := s.usaAddressRepo.Store(newAddress)
			if err != nil {
				return nil, err
			}

			newAddress.ID = ID

			return newAddress, nil
		} else {
			return nil, err
		}
	}

	return address, nil
}

func (s *usaAddressService) CheckAddressUpdates() (*models.UsaAddress, *models.UsaAddress, diff.Changelog, error) {
	address, err := s.GetFirst()
	if err != nil {
		return nil, nil, nil, err
	}
	newAddress, err := s.fetcher.Fetch()
	if err != nil {
		return nil, nil, nil, err
	}

	changelog, err := diff.Diff(address, newAddress)

	if len(changelog) > 0 {
		s.logger.Info("UsaAddress data has changed.")
		s.usaAddressRepo.Update(newAddress)
	}

	return address, newAddress, changelog, nil
}

func (s *usaAddressService) Delete() error {
	err := s.usaAddressRepo.Delete()
	if err == nil {
		s.logger.Info(fmt.Sprintf("UsaAddress has been deleted."))
	} else {
		s.logger.Error(fmt.Sprintf("UsaAddress couldn't be deleted because of %s.", err.Error()))
	}
	return err
}
