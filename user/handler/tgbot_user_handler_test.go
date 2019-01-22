package handler

import (
	"errors"
	"testing"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/golang/mock/gomock"
	"github.com/teimurjan/go-els-tg-bot/mocks"
	"github.com/teimurjan/go-els-tg-bot/texts"
)

func TestJoinError(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	var chatID int64 = 1

	mockService := mocks.NewMockUserService(ctrl)
	err := errors.New("Test error")
	mockService.
		EXPECT().
		Create(chatID).
		Return(nil, err).
		Times(1)

	mockBot := mocks.NewMockBotAPI(ctrl)
	mockBot.
		EXPECT().
		Send(tgbotapi.NewMessage(chatID, texts.GetErrorMessage(err))).
		Times(1)

	handler := NewTgbotUserHandler(mockService, mockBot)

	handler.Join(chatID)
}

func TestJoinSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	var chatID int64 = 1

	mockService := mocks.NewMockUserService(ctrl)
	mockService.
		EXPECT().
		Create(gomock.Any()).
		Return(nil, nil).
		Times(1)

	mockBot := mocks.NewMockBotAPI(ctrl)
	mockBot.
		EXPECT().
		Send(tgbotapi.NewMessage(chatID, texts.GetWelcomeMessage())).
		Times(1)

	handler := NewTgbotUserHandler(mockService, mockBot)

	handler.Join(chatID)
}
