package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"time"
)

func basic(drone Drone) {
	drone.TakeOff()
	gobot.After(3*time.Second, func() {
		drone.Left(30)
		SleepSeconds(2)
		drone.Right(30)
		SleepSeconds(2)
		drone.Land()
	})
}

func work(drone Drone) {
	go func() {
		drone.SetupCamera(4)
	}()
	go func() {
		drone.TakeOff()
		SleepSeconds(3)
		drone.Hover()
		SleepSeconds(3)
		drone.Land()
	}()
}

func main() {
	driver := tello.NewDriver("8888")

	drone := Drone{driver: driver}

	job := func() {
		//basic(drone)
		work(drone)
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
