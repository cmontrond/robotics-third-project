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
