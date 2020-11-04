package main

import (
	"fmt"
	"gobot.io/x/gobot/platforms/dji/tello"
)

type Drone struct {
	drone *tello.Driver
}

func (drone Drone) TakeOff() {
	TakeOff(drone.drone)
}

func (drone Drone) Land() {
	Land(drone.drone)
}

func (drone Drone) Left(speed int) {
	Left(drone.drone, speed)
}

func (drone Drone) Right(speed int) {
	Right(drone.drone, speed)
}

func (drone Drone) StartVideo() {
	StartVideo(drone.drone)
}

func (drone Drone) SetVideoEncoderRate(rate tello.VideoBitRate) {
	SetVideoEncoderRate(drone.drone, rate)
}

func TakeOff(drone *tello.Driver) {
	err := drone.TakeOff()
	if err != nil {
		fmt.Printf("Error taking off: %+v\n", err)
	}
}

func Land(drone *tello.Driver) {
	err := drone.Land()
	if err != nil {
		fmt.Printf("Error landing: %+v\n", err)
	}
}

func Left(drone *tello.Driver, speed int) {
	err := drone.Left(speed)
	if err != nil {
		fmt.Printf("Error moving to the left: %+v\n", err)
	}
}

func Right(drone *tello.Driver, speed int) {
	err := drone.Right(speed)
	if err != nil {
		fmt.Printf("Error moving right something: %+v\n", err)
	}
}

func StartVideo(drone *tello.Driver) {
	err := drone.StartVideo()
	if err != nil {
		fmt.Printf("Error starting video: %+v\n", err)
	}
}

func SetVideoEncoderRate(drone *tello.Driver, rate tello.VideoBitRate) {
	err := drone.SetVideoEncoderRate(rate)
	if err != nil {
		fmt.Printf("Error setting video encoder rate: %+v\n", err)
	}
}
