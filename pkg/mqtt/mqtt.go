package mqtt

import (
	"context"
	"encoding/json"
	"enoceanccu/pkg/controller"
	"enoceanccu/pkg/entities"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttValue struct {
	Val interface{} `json:"val"`
}

type Client struct {
	uri    *url.URL
	client mqtt.Client
	ctrl   *controller.Controller
	id     string
	in     chan entities.HmState
	out    chan entities.HmState
	ctx    context.Context
}

func NewClient(ctx context.Context, server string, port int, id string, in chan entities.HmState, out chan entities.HmState, ctrl *controller.Controller) (*Client, error) {
	uri, err := url.Parse(fmt.Sprintf("mqtt://%s:%d", server, port))
	if err != nil {
		return nil, err
	}

	c := &Client{
		uri:  uri,
		id:   id,
		in:   in,
		out:  out,
		ctrl: ctrl,
		ctx:  ctx,
	}

	c.connect()

	return c, nil
}

func (c *Client) connect() {
	opts := c.createClientOptions()
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	c.client = client

	go c.publisher()

	return
}

func (c *Client) createClientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", c.uri.Host))
	opts.SetUsername(c.uri.User.Username())
	password, _ := c.uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(c.id)
	return opts
}

func (c *Client) Subscribe(topic string) {
	c.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		if msg.Topic() == "hm/devices" {
			logrus.Infof("Requesting device list")
			resp, err := c.ctrl.GetDevices()
			if err != nil {
				logrus.Errorf("Failed to get list of devices")
				return
			}

			b, err := json.Marshal(resp)
			if err != nil {
				logrus.Errorf("Failed to marshal device list: %s", err)
			}
			t := client.Publish("hm/devices/response", 1, false, string(b))
			if t.Error() != nil {
				logrus.Errorf("Failed to publish device list: %s", t.Error())
				return
			}
		}
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}

func (c *Client) publisher() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case state := <-c.out:
			topic := fmt.Sprintf("hm/status/BidCos-RF/%s:%d/%s", state.Address, state.Channel, state.DataPoint)

			retain := true
			switch state.DataPoint {
			case "PRESS_SHORT", "PRESS_LONG", "DB0":
				retain = false
			}

			val := mqttValue{Val: state.Value}
			valJSON, err := json.Marshal(val)
			if err != nil {
				logrus.Errorf("Failed to create val JSON: %v", state)
				continue
			}

			t := c.client.Publish(topic, 0, retain, string(valJSON))
			if t.Error() != nil {
				logrus.Errorf("Failed to publish value %#v to topic %s: %s", state.Value, topic, err)
				continue
			}
		}
	}
}
