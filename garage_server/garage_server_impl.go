package garage_server

import (
	"encoding/json"
	"io"
	"log"
	"luis/melderserver/constants"
	"luis/melderserver/serial"
	"net/http"
)

type GarageServerImpl struct {
	garageStatus   int
	serial         serial.Serial
	onStatusUpdate func(int)
}

func NewGarageServerImpl(serial serial.Serial, onStatusUpdate func(int)) *GarageServerImpl {
	impl := GarageServerImpl{
		serial:         serial,
		onStatusUpdate: onStatusUpdate,
	}
	go impl.startListening()
	return &impl
}

func (impl *GarageServerImpl) startListening() {
	http.HandleFunc("/garage_status", impl.handleGarageStatus)
	http.ListenAndServe(":8080", nil)
}

func (impl *GarageServerImpl) handleGarageStatus(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	switch r.Method {
	case http.MethodPost:
		if err != nil || len(body) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Body is empty or invalid")
			return
		}
		var jsonResponse GarageStatus
		if !json.Valid(body) {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Body is not a valid json")
			return
		}
		if err := json.Unmarshal(body, &jsonResponse); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("json could not be unmarshalled")
			return
		}
		if jsonResponse.IsOpen {
			impl.garageStatus = constants.GARAGE_OPEN
		} else {
			impl.garageStatus = constants.GARAGE_CLOSED
		}
		impl.serial.SendGarageStatus(impl.garageStatus)
		impl.onStatusUpdate(impl.garageStatus)
	case http.MethodGet:
		status := GarageStatus{
			IsOpen: impl.garageStatus == constants.GARAGE_OPEN,
		}
		jsonResponse, err := json.Marshal(status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Could not marshal json")
			return
		}
		w.Write(jsonResponse)
	}
}

func (impl *GarageServerImpl) GetStatus() int {
	return impl.garageStatus
}
