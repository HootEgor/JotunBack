package main

import (
	"JotunBack/handlers/botH"
	"JotunBack/model"
	"JotunBack/repository"
	"JotunBack/server"
	"JotunBack/ui"
	"context"
	firebase "firebase.google.com/go/v4"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

func main() {

	Ctx := context.Background()
	serviceAcc := option.WithCredentialsFile("serviceAccount/jotunn-8f418-firebase-adminsdk-ddgl2-cd17bb27c3.json")
	app, err := firebase.NewApp(Ctx, nil, serviceAcc)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(Ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	userRepo := repository.NewUserRepository(client, Ctx)

	hub := server.NewHub()

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			server.HandleWebSocket(hub, w, r, userRepo)
		})
		log.Println("WebSocket server started on :8081")
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	bot, err := tgbotapi.NewBotAPI("7000758343:AAHEh8KWo-hBPQVL0XvJ1i76_7yzWJUNnTQ")
	if err != nil {
		log.Panic(err)
	}

	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: 1,
	})

	var acStates = make(map[string]*model.ACState)

	for update := range updates {
		if update.Message != nil {
			go handleMessage(update, bot, acStates, hub, userRepo)
		} else if update.CallbackQuery != nil {
			go handleCallback(bot, update, acStates, userRepo)
		}
	}

}

func handleMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, acStates map[string]*model.ACState, hub *server.Hub,
	userRepo *repository.UserRepository) {

	if update.Message == nil {
		return
	}

	acState := acStates[update.Message.From.UserName]
	if acState == nil {
		acState = &model.ACState{
			Username: update.Message.From.UserName,
			ChatID:   update.Message.Chat.ID,
			Config:   model.AirConditionerConfig{},
		}
		acStates[update.Message.From.UserName] = acState
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Text {
	case "/test":
		err := botH.GetAcProtocol(update.Message.From.UserName, userRepo, hub)
		if err != nil {
			return
		}
		msg.Text = "Відскануйте будьяку кнонку на пульті"
	case "/start":
		newUser := model.User{
			Username: update.Message.From.UserName,
			ChatID:   update.Message.Chat.ID,
		}
		err := userRepo.CreateUser(newUser)
		if err != nil {
			log.Println(err)
			return
		}
		msg.Text = "Hello! I am Jotun. How can I help you?"
	case "/stop":
		msg.Text = "Goodbye!"
	case "/on":
		err := botH.TurnOnAc(update.Message.From.UserName, userRepo, hub)
		if err != nil {
			return
		}
		msg.Text = "Turning on the air conditioner."
	case "/off":
		err := botH.TurnOffAc(update.Message.From.UserName, userRepo, hub)
		if err != nil {
			return
		}
		msg.Text = "Turning off the air conditioner."
	case "/config":
		ui.ConfigForm(bot, acState)
	default:
		msg.Text = "I don't understand that command."
	}

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}

}

func handleCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update, acStates map[string]*model.ACState,
	userRepo *repository.UserRepository) {

	if update.CallbackQuery == nil {
		return
	}

	acState := acStates[update.CallbackQuery.From.UserName]
	if acState == nil {
		acState = &model.ACState{
			Username: update.CallbackQuery.From.UserName,
			ChatID:   update.CallbackQuery.Message.Chat.ID,
			Config:   model.AirConditionerConfig{},
		}
		acStates[update.CallbackQuery.From.UserName] = acState
	}

	ui.InLineKeyboardHandler(bot, acState, userRepo, update)
}
