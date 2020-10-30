package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"io"
	"os/exec"
	"time"
)

func setupMplayer() io.WriteCloser {
	mplayer := exec.Command("mplyaer", "-fps", "25", "-")

	mplayerInput, err := mplayer.StdinPipe()
	if err != nil {
		fmt.Printf("Error creating input of mplayer: %+v\n", err)
	}

	err = mplayer.Start()
	if err != nil {
		fmt.Printf("Error starting mplayer: %+v\n", err)
	}

	return mplayerInput
}

func setupDroneVideo(drone *tello.Driver) {
	err := drone.StartVideo()
	if err != nil {
		fmt.Printf("Error starting video: %+v\n", err)
	}
	err = drone.SetVideoEncoderRate(4)
	if err != nil {
		fmt.Printf("Error setting video encoder rate: %+v\n", err)
	}
	gobot.Every(100*time.Millisecond, func() {
		err = drone.StartVideo()
		if err != nil {
			fmt.Printf("Error starting video: %+v\n", err)
		}
	})
}

func work(drone *tello.Driver) {
	mplayerInput := setupMplayer()

	err := drone.On(tello.ConnectedEvent, func(data interface{}) {
		println("Connected")
		setupDroneVideo(drone)
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

func main() {
	drone := tello.NewDriver("8890")

	robot := gobot.NewRobot("Project 3: Drone",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Printf("Error starting the Drone: %+v\n", err)
	}
}
