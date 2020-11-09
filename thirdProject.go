package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"time"
)

const (
	frameSize = 960 * 720 * 3
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

func work(drone Drone, window *gocv.Window) {
	go func() {
		//drone.SetupCameraWithMplayer(4, 0)
		drone.SetupCameraWithFfmpeg(window, 4, 0, frameSize, 960, 720)

	}()

	//go func() {
	//	drone.TakeOff()
	//	SleepSeconds(3)
	//	drone.Hover()
	//	SleepSeconds(3)
	//	drone.Land()
	//}()
}

func main() {
	driver := tello.NewDriver("8888")
	window := gocv.NewWindow("Face Detect")

	drone := Drone{driver: driver}

	job := func() {
		//basic(drone)
		work(drone, window)
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
