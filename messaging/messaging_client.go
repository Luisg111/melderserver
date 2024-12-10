package messaging

type MessagingClient interface {
	SendAlarm(isTest bool)
	SendAlarmConfirmed()
}
