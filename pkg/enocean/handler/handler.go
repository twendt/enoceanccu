package handler

import (
	"enoceanccu/pkg/entities"
	"fmt"

	"github.com/jonaz/goenocean"
)

type EEPHandler interface {
	Process(mqttOut chan entities.HmState, d entities.Device, p goenocean.Packet) error
}

var Handlers = make(map[string]EEPHandler)

func GetHandler(eep string) (EEPHandler, error) {
	if h, ok := Handlers[eep]; ok {
		return h, nil
	}

	return nil, fmt.Errorf("No handler found for EEP: %s", eep)
}
