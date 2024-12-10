package serial

type Serial interface {
	SendGarageStatus(garageStatus int) error
	SendVentilatorCommand(ventilatorCommand int) error
	SendLightCommand(lightCommand int) error
	SendTestCommand() error
}
