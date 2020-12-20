package errs

import (
	"errors"

	"github.com/teimurjan/go-els-tg-bot/texts"
)

type ErrNoSuchTracking struct {
	msg string
}

func (e *ErrNoSuchTracking) Error() string {
	return e.msg
}

func NewBaseError(msg string) error {
	return errors.New(msg)
}

func NewNoSuchTrackingError(tracking string) error {
	return &ErrNoSuchTracking{msg: texts.GetTrackingNotExistsMessage(tracking)}
}

func ErrToHumanReadableMessage(e error) string {
	switch e.(type) {
	case *ErrNoSuchTracking:
		return e.Error()
	default:
		return texts.GetCommonErrorMessage()
	}
}
