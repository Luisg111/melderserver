package main

import (
	"log"
	"luis/melderserver/constants"
	"luis/melderserver/garage_server"
	"luis/melderserver/messaging"
	"luis/melderserver/serial"
	"luis/melderserver/weather"
	"time"
)

func main() {
	var serialImpl serial.Serial
	var weatherClientImpl weather.WeatherClient
	var messagingClientImpl messaging.MessagingClient

	handleGarageStatusUpdate := func(status int) {
		log.Println("Garage status update received: ", status)
		serialImpl.SendGarageStatus(status)
	}

	handleLightCommand := func(command int) {
		log.Println("Light command received: ", command)
		serialImpl.SendLightCommand(command)
	}

	handleVentilatorCommand := func(command int) {
		log.Println("Ventilator command received: ", command)
		serialImpl.SendVentilatorCommand(command)
	}

	handleTestCommand := func() {
		log.Println("Test command received")
		serialImpl.SendTestCommand()
	}

	onAlarm := func() {
		log.Println("Alarm received")
		alarmTime := time.Now()
		isTest := alarmTime.Day() == int(time.Wednesday) && alarmTime.Sub(time.Date(alarmTime.Year(), alarmTime.Month(), alarmTime.Day(), 19, 0, 0, 0, time.Local)) < 10*time.Minute
		messagingClientImpl.SendAlarm(isTest)
		serialImpl.SendLightCommand(constants.LIGHT_ON)
		serialImpl.SendVentilatorCommand(constants.VENTILATOR_OFF)
		go func() {
			time.Sleep(5 * time.Minute)
			serialImpl.SendLightCommand(constants.LIGHT_OFF)
			log.Println("Alarm light off")
		}()
	}

	onAlarmConfirmed := func() {
		log.Println("Alarm confirmed")
		messagingClientImpl.SendAlarmConfirmed()
	}

	serialImpl = serial.NewSerialImpl(onAlarm, onAlarmConfirmed)
	garage_server.NewGarageServerImpl(serialImpl, handleGarageStatusUpdate)
	weatherClientImpl = weather.NewDWDWeatherClient()
	messagingClientImpl = messaging.NewTelegramMessagingClient(handleLightCommand, handleVentilatorCommand, handleTestCommand, weatherClientImpl)

	for {
		time.Sleep(1 * time.Minute)
	}
}
