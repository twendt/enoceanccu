package controller

import (
	"enoceanccu/pkg/enocean"
	"enoceanccu/pkg/enocean/handler"
	"enoceanccu/pkg/entities"
	"enoceanccu/pkg/repository"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jonaz/goenocean"
	"golang.org/x/net/context"
)

type Controller struct {
	mqttIn        chan entities.HmState
	mqttOut       chan entities.HmState
	enoceanClient *enocean.EnOcean
	repo          repository.Repository
	devices       map[string]entities.Device
	enoceanIn     chan goenocean.Packet
	enoceanOut    chan goenocean.Packet
	ctx           context.Context
}

func New(ctx context.Context, repo repository.Repository, enoceanIn chan goenocean.Packet, enoceanOut chan goenocean.Encoder, mqttIn chan entities.HmState, mqttOut chan entities.HmState) (*Controller, error) {
	devices, err := repo.GetDevices()
	if err != nil {
		return nil, err
	}

	return &Controller{
		devices:   devices,
		enoceanIn: enoceanIn,
		repo:      repo,
		ctx:       ctx,
		mqttIn:    mqttIn,
		mqttOut:   mqttOut,
	}, nil
}

func (c *Controller) Run() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case p := <-c.enoceanIn:
			logrus.Infof("Received Enocean packet")
			c.processPacket(p)
		}
	}
}

func (c *Controller) processPacket(p goenocean.Packet) {
	senderID := fmt.Sprintf("%x", p.SenderId())
	device, ok := c.devices[senderID]
	if ok {
		logrus.Infof("%#v", device)
		eeps := device.RcvEEPs
		if len(eeps) > 0 {
			handler, err := handler.GetHandler(eeps[0])
			if err != nil {
				logrus.Infof("No handler found for EEP: %s", eeps[0])
				return
			}
			err = handler.Process(c.mqttOut, device, p)
			if err != nil {
				logrus.Infof("EEP Handler failed: %s", err)
			}
			return
		}
		logrus.Infof("Device has no receive EEPs configured: %s", senderID)
		return
	}
	logrus.Infof("No device found for sender ID: %s", senderID)
}

func (c *Controller) GetDevices() (entities.DevicesResponse, error) {
	result := entities.DevicesResponse{
		InterfaceName: "BidCos-RF",
		Devices:       []entities.DeviceDef{},
	}
	devs, err := c.repo.GetDevices()
	if err != nil {
		return result, err
	}

	for _, d := range devs {
		result.Devices = append(result.Devices, entities.DeviceDef{
			Type:    d.HMType,
			Address: d.HMAddress,
		})
	}

	return result, nil
}
