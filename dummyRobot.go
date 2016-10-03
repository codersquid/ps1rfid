package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type DummyRobot struct {
}

func (d DummyRobot) configure() {

}

func (d DummyRobot) runRobot(shutdown chan bool) {

	for {
		select {
		case <-shutdown:
			log.Println("Caught the shutdown signal. Bailing out.")
			goto quit
		default:

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Code: ")
			text, _ := reader.ReadString('\n')
			urlString := "https://members.pumpingstationone.org/rfid/check/FrontDoor/" + text //let's remove a dependency on members.ps1.org
			fmt.Printf("I would send a request to %v but let's assume it just works.\n", urlString)
			d.openDoor()
		}
	}
quit:
	log.Println("Exited the runRobot loop successfully. Later, taters.")
}

func (d DummyRobot) openDoor() {
	log.Println("Opening door!")
	time.Sleep(5 * time.Second)
	log.Println("Locking door!")
}
