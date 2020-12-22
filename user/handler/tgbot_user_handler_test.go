package handler

import (
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/golang/mock/gomock"
	"github.com/teimurjan/go-els-tg-bot/mocks"
)

func TestJoinError(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	var chatID int64 = 1

	mockService := mocks.NewMockUserService(ctrl)
	mockI18nHelper := mocks.NewMockI18nHelper(ctrl)
	err := errors.New("Test error")
	mockService.
		EXPECT().
		Create(chatID).
		Return(nil, err).
		Times(1)

	mockBot := mocks.NewMockBotAPI(ctrl)
	mockBot.
		EXPECT().
		Send(tgbotapi.NewMessage(chatID, "Something went wrong. 🆘")).
		Times(1)

	handler := NewTgbotUserHandler(mockService, mockBot, mockI18nHelper)

	handler.Join(chatID)
}

func TestJoinSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	var chatID int64 = 1

	mockService := mocks.NewMockUserService(ctrl)
	mockI18nHelper := mocks.NewMockI18nHelper(ctrl)
	mockService.
		EXPECT().
		Create(gomock.Any()).
		Return(nil, nil).
		Times(1)

	mockBot := mocks.NewMockBotAPI(ctrl)
	mockBot.
		EXPECT().
		Send(tgbotapi.NewMessage(chatID, "Hi there! 👋\nStart monitoring your orders by typing:\n/add_tracking")).
		Times(1)

	handler := NewTgbotUserHandler(mockService, mockBot, mockI18nHelper)

	handler.Join(chatID)
}
