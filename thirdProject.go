package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"time"
)

func basic(drone *tello.Driver) {
	// Functional Approach
	TakeOff(drone)
	gobot.After(3*time.Second, func() {
		Left(drone, 30)
		time.Sleep(2 * time.Second)
		Right(drone, 30)
		time.Sleep(2 * time.Second)
		Land(drone)
	})

	// Type Approach
	//myDrone := Drone{drone: drone}
	//
	//myDrone.TakeOff()
	//gobot.After(3*time.Second, func() {
	//	myDrone.Left(30)
	//	SleepSeconds(2)
	//	myDrone.Right(30)
	//	SleepSeconds(2)
	//	myDrone.Land()
	//})
}

func work(drone *tello.Driver) {
	// Functional approach
	SetupCamera(drone, 4)

	// Type approach
	//myDrone := Drone{drone: drone}
	//myDrone.SetupCamera(4)
}

func main() {
	drone := tello.NewDriver("8888")

	mainFunc := func() {
		basic(drone)
	}

	robot := gobot.NewRobot("Project 3: Drone",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		mainFunc,
	)

	err := robot.Start()
	if err != nil {
		fmt.Printf("Error starting the Drone: %+v\n", err)
	}
}
