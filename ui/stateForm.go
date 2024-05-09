package ui

import (
	"JotunBack/model"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func StateForm(bot *tgbotapi.BotAPI, acState *model.ACState, isOnline bool) {
	msg := tgbotapi.NewMessage(acState.ChatID, StateMsg(acState, isOnline))

	keyboard := StateKeyboard(acState)
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
}

func StateKeyboard(acState *model.ACState) tgbotapi.InlineKeyboardMarkup {
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Викл.", "stop"),
	)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row1)
	return keyboard
}

func StateMsg(acState *model.ACState, isOnline bool) string {
	modeEmoji := "Cool❄️"
	switch acState.Config.Mode {
	case 1:
		modeEmoji = "Cool❄️"
	case 2:
		modeEmoji = "Heat☀️"
	case 3:
		modeEmoji = "Dry♨️"
	}

	onlineEmoji := "❌"
	if isOnline {
		onlineEmoji = "✅"
	}

	return fmt.Sprintf("Підключен: %s\nРежим: %s\nТемпература: %s -> %s °C \nПрацювати до: %s",
		onlineEmoji,
		modeEmoji,
		acState.GetTargetTemp(),
		acState.GetTemp(),
		acState.Until.Format("15:04"))
}
