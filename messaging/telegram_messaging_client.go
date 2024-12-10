package messaging

import (
	"log"
	"luis/melderserver/constants"
	"luis/melderserver/weather"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramMessagingClient struct {
	onLightCommandReceived      func(int)
	onVentilatorCommandReceived func(int)
	onTestCommandReceived       func()
	weatherClient               weather.WeatherClient
	alarmBot                    *tgbotapi.BotAPI
	weatherBot                  *tgbotapi.BotAPI
}

func NewTelegramMessagingClient(onLightCommandReceived func(int), onVentilatorCommandReceived func(int), onTestCommandReceived func(), weatherClient weather.WeatherClient) *TelegramMessagingClient {
	client := &TelegramMessagingClient{
		weatherClient:               weatherClient,
		onLightCommandReceived:      onLightCommandReceived,
		onVentilatorCommandReceived: onVentilatorCommandReceived,
		onTestCommandReceived:       onTestCommandReceived,
	}
	client.initBots()
	return client
}

func (client *TelegramMessagingClient) initBots() {
	alarmBotToken := os.Getenv("ALARM_BOT_TOKEN")
	weatherBotToken := os.Getenv("WEATHER_BOT_TOKEN")
	if alarmBotToken == "" || weatherBotToken == "" {
		log.Fatal("Could not find bot tokens")
	}
	alarmBot, err := tgbotapi.NewBotAPI(alarmBotToken)
	if err != nil {
		log.Fatal("Could not initialize alarm bot")
	}
	weatherBot, err := tgbotapi.NewBotAPI(weatherBotToken)
	if err != nil {
		log.Fatal("Could not initialize weather bot")
	}

	client.alarmBot = alarmBot
	client.weatherBot = weatherBot

	go func() {
		alarmUpdateConfig := tgbotapi.NewUpdate(0)
		alarmUpdateConfig.Timeout = 30
		updates := alarmBot.GetUpdatesChan(alarmUpdateConfig)
		for update := range updates {
			if update.Message == nil {
				continue
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			switch update.Message.Text {
			case "LOW":
				client.onVentilatorCommandReceived(constants.VENTILATOR_LOW)
			case "MED":
				client.onVentilatorCommandReceived(constants.VENTILATOR_MED)
			case "HI":
				client.onVentilatorCommandReceived(constants.VENTILATOR_HI)
			case "TEST":
				client.onTestCommandReceived()
			case "OFF":
				client.onVentilatorCommandReceived(constants.VENTILATOR_OFF)
			case "Licht An":
				client.onLightCommandReceived(constants.LIGHT_ON)
			case "Licht Aus":
				client.onLightCommandReceived(constants.LIGHT_OFF)
			default:
				log.Println("Unknown command")
			}
		}
	}()

	go func() {
		weatherUpdateConfig := tgbotapi.NewUpdate(0)
		weatherUpdateConfig.Timeout = 30
		updates := weatherBot.GetUpdatesChan(weatherUpdateConfig)
		for update := range updates {
			if update.Message == nil {
				continue
			}
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			var warningsText string = "Warnungen:\n"
			warnings := client.weatherClient.GetWarnings()
			for _, warning := range warnings {
				if !warning.IsActive {
					continue
				}
				warningsText += warning.Text + "\n"
			}
			if len(warnings) == 0 {
				warningsText += "Keine Warnungen\n"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, warningsText)
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := weatherBot.Send(msg); err != nil {
				log.Println("Could not answer weather bot message")
			}
		}
	}()
}

func (client *TelegramMessagingClient) SendAlarm(isTest bool) {
	chatId, err := strconv.ParseInt(os.Getenv("ALARM_CHAT_ID"), 10, 64)
	if err != nil || chatId == 0 {
		log.Fatal("Could not find alarm bot chat id")
	}
	msg := tgbotapi.NewMessage(chatId, "***ALARM***\n")
	if isTest {
		msg.Text += "Probealarm?\n"
	}
	msg.Text += time.Now().Format("02.01.2006 15:04:05")
	msg.Text += "\n"
	for _, warning := range client.weatherClient.GetWarnings() {
		msg.Text += warning.Title + "\n"
	}
	if _, err := client.alarmBot.Send(msg); err != nil {
		log.Println("Could not send alarm message")
	}
}

func (client *TelegramMessagingClient) SendAlarmConfirmed() {
	msg := tgbotapi.NewMessage(770527082, "***ALARM ENDE***\n")
	msg.Text += time.Now().Format("02.01.2006 15:04:05")
	if _, err := client.alarmBot.Send(msg); err != nil {
		log.Println("Could not send alarm end message")
	}
}
