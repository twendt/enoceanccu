package handler

import (
	"encoding/hex"
	"enoceanccu/pkg/entities"
	"fmt"
	"strconv"

	"github.com/jonaz/goenocean"
	"github.com/sirupsen/logrus"
)

type EEPF60201Handler struct {
}

func (h EEPF60201Handler) Process(out chan entities.HmState, d entities.Device, p goenocean.Packet) error {
	if t, ok := p.(goenocean.TelegramRps); ok {
		logrus.Infof(hex.EncodeToString([]byte{t.TelegramType()}))
		eep := goenocean.EepF60201{}
		eep.SetTelegram(t)
		db0Str := strconv.Itoa(int(eep.TelegramData()[0]))
		logrus.Infof("DB0: %s\n", db0Str)
		logrus.Infof("Sending %s to %s", db0Str, d.HMAddress)
		state := entities.HmState{
			Address:   d.HMAddress,
			DataPoint: "DB0",
			Channel:   1,
			Value:     db0Str,
		}
		out <- state
		//err := client.Publish(fmt.Sprintf("hm/status/BidCos-RF/%s:1/DB0", d.HMAddress), false, fmt.Sprintf("{ \"val\": \"%s\" }", db0Str))
		//if err != nil {
		//	return err
		//}
		state = entities.HmState{
			Address:   d.HMAddress,
			DataPoint: "PRESS_SHORT",
			Channel:   1,
			Value:     "true",
		}
		out <- state

		//err = client.Publish(fmt.Sprintf("hm/status/BidCos-RF/%s:1/PRESS_SHORT", d.HMAddress), false, fmt.Sprintf("{ \"val\": %s }", "true"))
		//if err != nil {
		//	return err
		//}
		return nil
	}

	return fmt.Errorf("Not a RPS telegram")
}

func init() {
	Handlers["eepf60201"] = &EEPF60201Handler{}
}
