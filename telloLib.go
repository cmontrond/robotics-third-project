package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"time"
)

// Type Approach

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

func (drone Drone) Up(speed int) {
	Up(drone.drone, speed)
}

func (drone Drone) Down(speed int) {
	Down(drone.drone, speed)
}

func (drone Drone) Forward(speed int) {
	Forward(drone.drone, speed)
}

func (drone Drone) Backward(speed int) {
	Backward(drone.drone, speed)
}

func (drone Drone) Clockwise(speed int) {
	Clockwise(drone.drone, speed)
}

func (drone Drone) CounterClockwise(speed int) {
	CounterClockwise(drone.drone, speed)
}

func (drone Drone) Hover() {
	Hover(drone.drone)
}

func (drone Drone) StartVideo() {
	StartVideo(drone.drone)
}

func (drone Drone) SetVideoEncoderRate(rate tello.VideoBitRate) {
	SetVideoEncoderRate(drone.drone, rate)
}

func (drone Drone) SetupVideo(rate tello.VideoBitRate) {
	SetupVideo(drone.drone, rate)
}

func (drone Drone) SetupCamera(rate tello.VideoBitRate) {
	SetupCamera(drone.drone, rate)
}

// Functional Approach

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
		fmt.Printf("Error moving right: %+v\n", err)
	}
}

func Up(drone *tello.Driver, speed int) {
	err := drone.Up(speed)
	if err != nil {
		fmt.Printf("Error moving up: %+v\n", err)
	}
}

func Down(drone *tello.Driver, speed int) {
	err := drone.Down(speed)
	if err != nil {
		fmt.Printf("Error moving down: %+v\n", err)
	}
}

func Forward(drone *tello.Driver, speed int) {
	err := drone.Forward(speed)
	if err != nil {
		fmt.Printf("Error moving forward: %+v\n", err)
	}
}

func Backward(drone *tello.Driver, speed int) {
	err := drone.Backward(speed)
	if err != nil {
		fmt.Printf("Error moving backward: %+v\n", err)
	}
}

func Clockwise(drone *tello.Driver, speed int) {
	err := drone.Clockwise(speed)
	if err != nil {
		fmt.Printf("Error moving clockwise: %+v\n", err)
	}
}

func CounterClockwise(drone *tello.Driver, speed int) {
	err := drone.CounterClockwise(speed)
	if err != nil {
		fmt.Printf("Error moving counter-clockwise: %+v\n", err)
	}
}

func Hover(drone *tello.Driver) {
	drone.Hover()
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

func SetupVideo(drone *tello.Driver, rate tello.VideoBitRate) {
	StartVideo(drone)
	SetVideoEncoderRate(drone, rate)
	gobot.Every(100*time.Millisecond, func() {
		StartVideo(drone)
	})
}

func SetupCamera(drone *tello.Driver, rate tello.VideoBitRate) {
	mplayerInput := SetupMplayer()

	err := drone.On(tello.ConnectedEvent, func(data interface{}) {
		println("Connected")
		SetupVideo(drone, rate)
	})
	if err != nil {
		fmt.Printf("Error setting ConnectedEvent event for drone: %+v\n", err)
	}

	err = drone.On(tello.VideoFrameEvent, func(data interface{}) {
		packet := data.([]byte)
		if _, err := mplayerInput.Write(packet); err != nil {
			fmt.Printf("Error writing to mplayerInput: %+v\n", err)
		}
	})
	if err != nil {
		fmt.Printf("Error setting VideoFrameEvent event for drone: %+v\n", err)
	}
}
