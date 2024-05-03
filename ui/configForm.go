package ui

import (
	"JotunBack/model"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func ConfigForm(bot *tgbotapi.BotAPI, acState *model.ACState) {
	msg := tgbotapi.NewMessage(acState.ChatID, ConfigMsg(acState))

	keyboard := ConfigKeyboard(acState)
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func ConfigKeyboard(acState *model.ACState) tgbotapi.InlineKeyboardMarkup {
	temp := strconv.Itoa(acState.Config.Degrees) + "°C"
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("-", "decrease"),
		tgbotapi.NewInlineKeyboardButtonData(temp, "temperature"),
		tgbotapi.NewInlineKeyboardButtonData("+", "increase"),
	)

	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("❄️", "cool"),
		tgbotapi.NewInlineKeyboardButtonData("☀️", "heat"),
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(row1, row2)
	return keyboard
}

func ConfigMsg(acState *model.ACState) string {
	modeEmoji := "❄️"
	if acState.Config.Mode == 2 {
		modeEmoji = "☀️"
	}

	return fmt.Sprintf("Режим: %s\nТемпература: %d°C", modeEmoji, acState.Config.Degrees)
}
