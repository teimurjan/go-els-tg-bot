package service

import (
	"database/sql"
	"fmt"

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

func (s *usaAddressService) GetUpdates() ([]string, *models.UsaAddress, error) {
	updatedFields := make([]string, 0)

	address, err := s.GetFirst()
	if err != nil {
		return updatedFields, nil, err
	}
	newAddress, err := s.fetcher.Fetch()
	if err != nil {
		return updatedFields, nil, err
	}

	logMsg := "UsaAddress data has changed."

	if newAddress.Street != address.Street {
		logMsg += fmt.Sprintf("\nStreet: %s to %s.", address.Street, newAddress.Street)
		updatedFields = append(updatedFields, "Street")
	}
	if newAddress.City != address.City {
		logMsg += fmt.Sprintf("\nCity: %s to %s.", address.City, newAddress.City)
		updatedFields = append(updatedFields, "City")
	}
	if newAddress.State != address.State {
		logMsg += fmt.Sprintf("\nState: %s to %s.", address.State, newAddress.State)
		updatedFields = append(updatedFields, "State")
	}
	if newAddress.Zip != address.Zip {
		logMsg += fmt.Sprintf("\nZip: %s to %s.", address.Zip, newAddress.Zip)
		updatedFields = append(updatedFields, "Zip")
	}
	if newAddress.PhoneNumber != address.PhoneNumber {
		logMsg += fmt.Sprintf("\nPhoneNumber: %s to %s.", address.PhoneNumber, newAddress.PhoneNumber)
		updatedFields = append(updatedFields, "PhoneNumber")
	}

	if len(updatedFields) > 0 {
		s.logger.Info(logMsg)
		s.usaAddressRepo.Update(newAddress)
	}

	return updatedFields, newAddress, nil
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
