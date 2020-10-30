package main

import (
	"fmt"
	// "fmt"
	"time"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8888")

	work := func() {
		err := drone.TakeOff()
		if err != nil {
			fmt.Println("Error making the drone take off %+v", err)
		}
		gobot.After(3*time.Second, func() {
			err = drone.Left(50)
			if err != nil {
				fmt.Println("Error making the drone go left %+v", err)
			}
			time.Sleep(time.Second*3)
			err = drone.Right(50)
			if err != nil {
				fmt.Println("Error making the drone go right %+v", err)
			}
			time.Sleep(time.Second*3)
			err = drone.Land()
			if err != nil {
				fmt.Println("Error making the drone land %+v", err)
			}
		})
	}

	robot := gobot.NewRobot("Project 3: Drone",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	err := robot.Start()

	if err != nil {
		fmt.Errorf("Error starting the Drone %+v", err)
	}
}