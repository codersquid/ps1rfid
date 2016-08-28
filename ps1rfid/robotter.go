package ps1rfid

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/loansindi/ps1rfid/cfg"
	serial "github.com/tarm/goserial"
)

type Robotter interface {
	OpenDoor() error
	ReadFull() string
}

func NewRobotter(cfg cfg.Config, test bool) (Robotter, error) {
	if !test {
		robotter, err := NewBeagleboneRobotter(cfg)
		return robotter, err
	}
	return TestRobotter{f: os.Stdin}, nil
}

type BeagleboneRobotter struct {
	b      *beaglebone.BeagleboneAdaptor
	c      serial.Config
	port   io.ReadWriteCloser
	splate *gpio.DirectPinDriver
}

func (b BeagleboneRobotter) ReadFull() string {
	var buf []byte
	n, err := io.ReadFull(b.port, buf)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	// We need to strip the stop and start bytes from the tag, so we only assign a certain range of the slice
	return string(buf[1 : n-3])
}

func NewBeagleboneRobotter(cfg cfg.Config) (BeagleboneRobotter, error) {
	b := BeagleboneRobotter{
		b: beaglebone.NewBeagleboneAdaptor("beaglebone"),
		c: serial.Config{Name: cfg.SerialName, Baud: cfg.SerialBaud},
	}

	u, err := serial.OpenPort(&b.c)
	if err != nil {
		return b, err
	}

	b.splate = gpio.NewDirectPinDriver(b.b, "splate", cfg.TogglePin)
	b.port = u
	return b, nil
}

func (b BeagleboneRobotter) OpenDoor() error {
	err := b.splate.DigitalWrite(1)
	if err != nil {
		return fmt.Errorf("openDoor error: %v", err)
	}
	gobot.After(5*time.Second, func() {
		b.splate.DigitalWrite(0)
	})
	return nil
}

type TestRobotter struct {
	f *os.File
}

func (t TestRobotter) OpenDoor() error {
	fmt.Println("openDoor")
	return nil
}

func (t TestRobotter) ReadFull() string {
	// do something to get stdin here
	return "foo"
}
