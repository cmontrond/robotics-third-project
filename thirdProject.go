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
	//myDrone := Drone{driver: driver}
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
	//myDrone := Drone{driver: driver}
	//myDrone.SetupCamera(4)
}

func main() {
	driver := tello.NewDriver("8888")

	//driver := Drone{driver: driver}

	job := func() {
		basic(driver)
	}

	robot := gobot.NewRobot("Project 3: Drone",
		[]gobot.Connection{},
		[]gobot.Device{driver},
		job,
	)

	err := robot.Start()
	if err != nil {
		fmt.Printf("Error starting the Drone: %+v\n", err)
	}
}
