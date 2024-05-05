package ui

import (
	"JotunBack/model"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
)

func ConfigForm(bot *tgbotapi.BotAPI, acState *model.ACState, isOnline bool) {
	msg := tgbotapi.NewMessage(acState.ChatID, ConfigMsg(acState, isOnline))

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
		tgbotapi.NewInlineKeyboardButtonData("❄️Cool", "cool"),
		tgbotapi.NewInlineKeyboardButtonData("☀️Heat", "heat"),
		tgbotapi.NewInlineKeyboardButtonData("♨️Dry", "dry"),
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(row1, row2)
	return keyboard
}

func ConfigMsg(acState *model.ACState, isOnline bool) string {
	modeEmoji := "❄️"
	switch acState.Config.Mode {
	case 1:
		modeEmoji = "❄️Cool"
	case 2:
		modeEmoji = "☀️Heat"
	case 3:
		modeEmoji = "♨️Dry"
	}

	onlineEmoji := "❌"
	if isOnline {
		onlineEmoji = "✅"
	}

	return fmt.Sprintf("Підключен: %s\nРежим: %s\nТемпература: %d°C", onlineEmoji, modeEmoji, acState.Config.Degrees)
}
