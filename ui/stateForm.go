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
		tgbotapi.NewInlineKeyboardButtonData("Стоп", "stop"),
	)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(row1)
	return keyboard
}

func StateMsg(acState *model.ACState, isOnline bool) string {
	modeEmoji := "Cool❄️"
	switch acState.Config.Mode {
	case 1:
		modeEmoji = "Cool"
		for i := 0; i < acState.EmojiNum; i++ {
			modeEmoji += "❄️"
		}
	case 2:
		modeEmoji = "Heat"
		for i := 0; i < acState.EmojiNum; i++ {
			modeEmoji += "☀️"
		}
	case 3:
		modeEmoji = "Dry"
		for i := 0; i < acState.EmojiNum; i++ {
			modeEmoji += "♨️"
		}
	}

	onlineEmoji := "❌"
	if isOnline {
		onlineEmoji = "✅"
	}

	return fmt.Sprintf("Підключен: %s\nРежим: %s\nТемпература: %s -> %s °C \nПрацювати до: %s",
		onlineEmoji,
		modeEmoji,
		acState.GetTemp(),
		acState.GetTargetTemp(),
		acState.Until.Format("15:04"))
}
