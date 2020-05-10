// Package main provides ...
package main

import (
	"context"
	"enoceanccu/pkg/controller"
	"enoceanccu/pkg/discovery"
	"enoceanccu/pkg/enocean"
	"enoceanccu/pkg/entities"
	"enoceanccu/pkg/mqtt"
	"enoceanccu/pkg/repository"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonaz/goenocean"
	"github.com/sirupsen/logrus"
)

func main() {
	tty := flag.String("device", "/dev/ttyUSB0", "Path of the device of the EnOcean Stick")
	serial := flag.String("serial", "abcdefg", "Serial of the CCU")
	devicesFile := flag.String("device-config", "devices.json", "Path of the file with the defined devices")

	logrus.SetLevel(logrus.InfoLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	discoveryServer := discovery.NewServer(*serial, ":43439")
	go discoveryServer.Listen(context.Background())

	mqttIn := make(chan entities.HmState)
	mqttOut := make(chan entities.HmState)
	enoceanIn := make(chan goenocean.Packet)
	enoceanOut := make(chan goenocean.Encoder, 100)

	repo := repository.New(*devicesFile)

	ctrl, err := controller.New(ctx, repo, enoceanIn, enoceanOut, mqttIn, mqttOut)
	if err != nil {
		panic(err)
	}

	mqttClient, err := mqtt.NewClient(ctx, "localhost", 1883, "enoceanccu", mqttIn, mqttOut, ctrl)
	if err != nil {
		panic(err)
	}

	_, err = enocean.NewEnOcean(ctx, *tty, enoceanIn, enoceanOut)
	if err != nil {
		logrus.Error(err)
		return
	}

	go ctrl.Run()

	mqttClient.Subscribe("#")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	cancel()
}
