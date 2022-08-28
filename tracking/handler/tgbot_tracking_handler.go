package handler

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	helper "github.com/teimurjan/go-els-tg-bot/helper/i18n"
	"github.com/teimurjan/go-els-tg-bot/tgbot"
	"github.com/teimurjan/go-els-tg-bot/tracking"
	argumentsUtil "github.com/teimurjan/go-els-tg-bot/utils/arguments"
	errsUtil "github.com/teimurjan/go-els-tg-bot/utils/errs"
)

type tgbotTrackingHandler struct {
	service    tracking.TrackingService
	bot        tgbot.TgBot
	i18nHelper helper.I18nHelper
}

// NewTgbotTrackingHandler creates new tgbotTrackingHandler instance
func NewTgbotTrackingHandler(
	service tracking.TrackingService,
	bot tgbot.TgBot,
	i18nHelper helper.I18nHelper,
) tracking.TrackingHandler {
	return &tgbotTrackingHandler{
		service,
		bot,
		i18nHelper,
	}
}

func (h *tgbotTrackingHandler) AddTracking(arguments string, chatID int64) {
	parsedArguments := argumentsUtil.ParseArguments(arguments)
	trackingNumber, trackingOk := parsedArguments["v"]
	name, nameOk := parsedArguments["n"]
	localizer := h.i18nHelper.MustGetLocalizer(chatID)

	if !trackingOk || !nameOk {
		text := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "notEnoughArgumentsForTracking",
				Other: "Please specify arguments in format:\n `/add_tracking -v=\"YOUR_TRACKING_NUMBER\" -n=\"NAME_OF_ORDER\"`",
			},
		})
		h.bot.Send(tgbotapi.NewMessage(chatID, text))
		return
	}

	tracking, err := h.service.Create(trackingNumber, name, chatID)
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, errsUtil.GetErrorMessage(err, localizer)))
		return
	}
	addedText := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "trackingAdded",
			Other: "Tracking was successfully added. ✅\nI will notify you about its changes ASAP.",
		},
	})
	infoText := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "trackingInfo",
			Other: "Name: *{{.Name}}*\nStatus: *{{.Status}}*\nWeight: *{{.Weight}}*\nTracking: *{{.Value}}*",
		},
		TemplateData: tracking,
	})
	msgAdded := tgbotapi.NewMessage(chatID, addedText)
	msgInfo := tgbotapi.NewMessage(chatID, infoText)
	msgInfo.ParseMode = tgbotapi.ModeMarkdown

	h.bot.Send(msgAdded)
	h.bot.Send(msgInfo)
}

func (h *tgbotTrackingHandler) GetAll(chatID int64) {
	localizer := h.i18nHelper.MustGetLocalizer(chatID)

	trackings, err := h.service.GetForChat(chatID)
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, errsUtil.GetErrorMessage(err, localizer)))
		return
	}

	if len(trackings) == 0 {
		text := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "noTrackings",
				Other: "You have 0️⃣ trackings added.\nCreate an order using command:\n/add_tracking",
			},
		})
		h.bot.Send(tgbotapi.NewMessage(chatID, text))
	}

	trackingCh, errCh, doneCh := h.service.SyncAll(trackings)

	for {
		select {
		case tracking := <-trackingCh:
			deleteText := localizer.MustLocalize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "delete",
					Other: "Delete ❓",
				},
			})

			inlineBtns := tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardButtonData(
						deleteText,
						fmt.Sprintf("/delete_tracking %d", tracking.ID),
					),
				},
			)

			text := localizer.MustLocalize(&i18n.LocalizeConfig{
				MessageID:    "trackingInfo",
				TemplateData: tracking,
			})
			msg := tgbotapi.NewMessage(chatID, text)
			msg.ParseMode = tgbotapi.ModeMarkdown
			msg.ReplyMarkup = inlineBtns
			h.bot.Send(msg)
		case err := <-errCh:
			h.bot.Send(tgbotapi.NewMessage(chatID, errsUtil.GetErrorMessage(err, localizer)))
			return
		case _ = <-doneCh:
			return
		}
	}
}

func (h *tgbotTrackingHandler) DeleteTracking(trackingID int64, chatID int64, messageID int64) {
	localizer := h.i18nHelper.MustGetLocalizer(chatID)
	err := h.service.Delete(trackingID)
	var msg tgbotapi.Chattable
	if err != nil {
		h.bot.Send(tgbotapi.NewMessage(chatID, errsUtil.GetErrorMessage(err, localizer)))
	} else {
		msg = tgbotapi.NewDeleteMessage(chatID, int(messageID))
	}
	h.bot.Send(msg)
}

func (h *tgbotTrackingHandler) CheckUpdates() {
	groupedByUser, err := h.service.GetAllGroupedByUser()
	if err != nil {
		return
	}

	for user, trackings := range groupedByUser {
		localizer := h.i18nHelper.MustGetLocalizer(user.ChatID)
		trackingCh, errCh, doneCh := h.service.SyncOnlyUpdated(trackings)

		for {
			select {
			case tracking := <-trackingCh:
				infoText := localizer.MustLocalize(&i18n.LocalizeConfig{
					MessageID:    "trackingInfo",
					TemplateData: tracking,
				})
				text := localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "trackingUpdated",
						Other: "❗️❗️❗️ Your order status has been changed ❗️❗️❗️",
					},
				}) + "\n\n" + infoText

				msg := tgbotapi.NewMessage(
					user.ChatID,
					text,
				)
				msg.ParseMode = tgbotapi.ModeMarkdown
				h.bot.Send(msg)
			case err := <-errCh:
				h.bot.Send(tgbotapi.NewMessage(user.ChatID, errsUtil.GetErrorMessage(err, localizer)))
			case <-doneCh:
				return
			}
		}
	}
}
