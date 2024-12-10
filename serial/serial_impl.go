package serial

import (
	"errors"
	"log"
	"luis/melderserver/constants"
	"strings"

	"go.bug.st/serial"
)

const (
	alarmKeyword    = "ALARM"
	alarmEndKeyword = "ENDE"

	garageOpenKeyword   = "GarageOffen"
	garageClosedKeyword = "GarageZu"

	ventilatorOffKeyword = "OFF"
	ventilatorLowKeyword = "LOW"
	ventilatorMedKeyword = "MED"
	ventilatorHiKeyword  = "HI"

	lightOnKeyword  = "ON"
	lightOffKeyword = "LIGHTOFF"

	testKeyword = "TEST"
)

type SerialImpl struct {
	serialPort      serial.Port
	onAlarmReceived func()
	onAlarmEnd      func()
}

func NewSerialImpl(onAlarmReceived func(), onAlarmEnd func()) *SerialImpl {
	impl := &SerialImpl{
		onAlarmReceived: onAlarmReceived,
		onAlarmEnd:      onAlarmEnd,
	}
	impl.openConnection()
	return impl
}

func (impl *SerialImpl) openConnection() {
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.EvenParity,
	}
	port, err := serial.Open("/dev/serial/by-id/usb-1a86_USB2.0-Serial-if00-port0", mode)
	if err != nil {
		log.Fatal("could not open serial connection ", err)
	}
	impl.serialPort = port

	go func() {
		for {
			buf := make([]byte, 128)
			n, err := impl.serialPort.Read(buf)
			if err != nil {
				log.Fatal("error reading from serial connection!")
			}
			received := string(buf[:n])
			received = strings.TrimSuffix(received, "\n")
			log.Println("received: " + received)

			if received == alarmKeyword {
				if impl.onAlarmReceived != nil {
					impl.onAlarmReceived()
				}
			} else if received == alarmEndKeyword {
				if impl.onAlarmEnd != nil {
					impl.onAlarmEnd()
				}
			}
		}
	}()
}

func (impl *SerialImpl) sendCommand(command string) error {
	commandString := command + "\n"
	_, err := impl.serialPort.Write([]byte(commandString))
	if err != nil {
		log.Println("error sending command " + commandString + " : " + err.Error())
	}
	return err

}

func (impl *SerialImpl) SendGarageStatus(garageStatus int) error {
	switch garageStatus {
	case constants.GARAGE_OPEN:
		return impl.sendCommand(garageOpenKeyword)
	case constants.GARAGE_CLOSED:
		return impl.sendCommand(garageClosedKeyword)
	default:
		return errors.New("unknown garage status")
	}

}

func (impl *SerialImpl) SendVentilatorCommand(ventilatorCommand int) error {
	switch ventilatorCommand {
	case constants.VENTILATOR_OFF:
		return impl.sendCommand(ventilatorOffKeyword)
	case constants.VENTILATOR_LOW:
		return impl.sendCommand(ventilatorLowKeyword)
	case constants.VENTILATOR_MED:
		return impl.sendCommand(ventilatorMedKeyword)
	case constants.VENTILATOR_HI:
		return impl.sendCommand(ventilatorHiKeyword)
	default:
		return errors.New("unknown ventilator command")
	}
}

func (impl *SerialImpl) SendLightCommand(lightCommand int) error {
	switch lightCommand {
	case constants.LIGHT_ON:
		return impl.sendCommand(lightOnKeyword)
	case constants.LIGHT_OFF:
		return impl.sendCommand(lightOffKeyword)
	default:
		return errors.New("unknown light command")
	}
}

func (impl *SerialImpl) SendTestCommand() error {
	return impl.sendCommand(testKeyword)
}
