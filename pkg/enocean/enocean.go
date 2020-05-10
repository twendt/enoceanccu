package enocean

import (
	"context"

	"github.com/jonaz/goenocean"
	"github.com/sirupsen/logrus"
)

type EnOcean struct {
	tty        string
	enoceanIn  chan goenocean.Packet
	enoceanOut chan goenocean.Encoder
	adapterIn  chan goenocean.Packet
	senderID   [4]byte
	ctx        context.Context
}

func NewEnOcean(ctx context.Context, tty string, incomingPacket chan goenocean.Packet, sendChan chan goenocean.Encoder) (*EnOcean, error) {
	eo := &EnOcean{
		tty:        tty,
		enoceanIn:  incomingPacket,
		enoceanOut: sendChan,
		adapterIn:  make(chan goenocean.Packet, 100),
		ctx:        ctx,
	}
	err := eo.connect()
	return eo, err
}

func (e *EnOcean) connect() error {

	err := goenocean.Serial(e.tty, e.enoceanOut, e.adapterIn)
	if err != nil {
		return err
	}

	e.getIDBase()
	go e.reciever()
	return nil
}

func (e *EnOcean) getIDBase() {
	p := goenocean.NewPacket()
	p.SetPacketType(goenocean.PacketTypeCommonCommand)
	p.SetData([]byte{0x08})
	e.enoceanOut <- p
}

func (e *EnOcean) reciever() {
	for {
		select {
		case p := <-e.adapterIn:
			logrus.Infof("Received: %x\n", p.Encode())
			logrus.Infof("Sender ID: %x\n", p.SenderId())
			if p.PacketType() == goenocean.PacketTypeResponse && len(p.Data()) == 5 {
				copy(e.senderID[:], p.Data()[1:4])
				logrus.Infof("senderid: % x ( % x )", e.senderID, p.Data())
				continue
			}
			if p.SenderId() != [4]byte{0, 0, 0, 0} {
				e.incomingPacket(p)
			}
		case <-e.ctx.Done():
			return
		}
	}
}

func (e *EnOcean) incomingPacket(p goenocean.Packet) {
	logrus.Infof("Incoming packet")
	e.enoceanIn <- p
	//if t, ok := p.(goenocean.TelegramRps); ok {
	//	logrus.Infof(hex.EncodeToString([]byte{t.TelegramType()}))
	//	eep := goenocean.EepF60201{}
	//	eep.SetTelegram(t)
	//	db0Str := strconv.Itoa(int(eep.TelegramData()[0]))
	//	logrus.Infof("DB0: %s\n", db0Str)
	//e.mqttClient.Publish(fmt.Sprintf("hm/status/BidCos-RF/%x:1/DB0", t.SenderId()), false, fmt.Sprintf("{ \"val\": \"%s\" }", db0Str))
	//e.mqttClient.Publish(fmt.Sprintf("hm/status/BidCos-RF/%x:1/PRESS_SHORT", t.SenderId()), false, fmt.Sprintf("{ \"val\": %s }", "true"))
	//}
}
