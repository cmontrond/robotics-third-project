package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"os/exec"
	"strconv"
	"time"
)

// Type Approach

type Drone struct {
	driver *tello.Driver
}

func (drone Drone) TakeOff() {
	TakeOff(drone.driver)
}

func (drone Drone) Land() {
	Land(drone.driver)
}

func (drone Drone) Left(speed int) {
	Left(drone.driver, speed)
}

func (drone Drone) Right(speed int) {
	Right(drone.driver, speed)
}

func (drone Drone) Up(speed int) {
	Up(drone.driver, speed)
}

func (drone Drone) Down(speed int) {
	Down(drone.driver, speed)
}

func (drone Drone) Forward(speed int) {
	Forward(drone.driver, speed)
}

func (drone Drone) Backward(speed int) {
	Backward(drone.driver, speed)
}

func (drone Drone) Clockwise(speed int) {
	Clockwise(drone.driver, speed)
}

func (drone Drone) CounterClockwise(speed int) {
	CounterClockwise(drone.driver, speed)
}

func (drone Drone) Hover() {
	Hover(drone.driver)
}

func (drone Drone) StartVideo() {
	StartVideo(drone.driver)
}

func (drone Drone) SetVideoEncoderRate(rate tello.VideoBitRate) {
	SetVideoEncoderRate(drone.driver, rate)
}

func (drone Drone) SetExposureLevel(level int) {
	SetExposureLevel(drone.driver, level)
}

func (drone Drone) SetupVideo(rate tello.VideoBitRate, level int) {
	SetupVideo(drone.driver, rate, level)
}

func (drone Drone) SetupCamera(rate tello.VideoBitRate, level int) {
	SetupCameraWithMplayer(drone.driver, rate, level)
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

func SetExposureLevel(drone *tello.Driver, level int) {
	err := drone.SetExposure(level)
	if err != nil {
		fmt.Printf("Error setting exposure level: %+v\n", err)
	}
}

func SetupVideo(drone *tello.Driver, rate tello.VideoBitRate, level int) {
	StartVideo(drone)
	SetVideoEncoderRate(drone, rate)
	SetExposureLevel(drone, level)
	gobot.Every(100*time.Millisecond, func() {
		StartVideo(drone)
	})
}

func SetupCameraWithMplayer(drone *tello.Driver, rate tello.VideoBitRate, level int) {

	mplayer := exec.Command("mplayer", "-fps", "25", "-")

	mplayerInput, err := mplayer.StdinPipe()
	if err != nil {
		fmt.Printf("Error creating input of mplayer: %+v\n", err)
	}

	err = mplayer.Start()
	if err != nil {
		fmt.Printf("Error starting mplayer: %+v\n", err)
	}

	err = drone.On(tello.ConnectedEvent, func(data interface{}) {
		println("Connected")
		SetupVideo(drone, rate, level)
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

func SetupCameraWithFfmpeg(drone *tello.Driver, rate tello.VideoBitRate, level int, frameX int, frameY int) {
	ffmpeg := exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
		"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameX)+"x"+strconv.Itoa(frameY), "-f", "rawvideo", "pipe:1")

	ffmpegIn, err := ffmpeg.StdinPipe()

	if err != nil {
		fmt.Printf("Error creating input of ffmpeg: %+v\n", err)
	}

	//ffmpegOut,err := ffmpeg.StdoutPipe()

	if err != nil {
		fmt.Printf("Error creating output of mplayer: %+v\n", err)
	}

	err = ffmpeg.Start()
	if err != nil {
		fmt.Printf("Error starting ffmpeg: %+v\n", err)
	}

	err = drone.On(tello.ConnectedEvent, func(data interface{}) {
		println("Connected")
		SetupVideo(drone, rate, level)
	})
	if err != nil {
		fmt.Printf("Error setting ConnectedEvent event for drone: %+v\n", err)
	}

	err = drone.On(tello.VideoFrameEvent, func(data interface{}) {
		packet := data.([]byte)
		if _, err := ffmpegIn.Write(packet); err != nil {
			fmt.Printf("Error writing to ffmpeg input: %+v\n", err)
		}
	})
	if err != nil {
		fmt.Printf("Error setting VideoFrameEvent event for drone: %+v\n", err)
	}
}
