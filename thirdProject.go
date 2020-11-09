package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"io"
	"os/exec"
	"strconv"
	"time"
)

const (
	frameSize = 960 * 720 * 3
)

func setupFfmpeg(frameX int, frameY int) (*exec.Cmd, io.WriteCloser, io.ReadCloser) {

	ffmpeg := exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
		"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameX)+"x"+strconv.Itoa(frameY), "-f", "rawvideo", "pipe:1")

	ffmpegIn, err := ffmpeg.StdinPipe()

	if err != nil {
		fmt.Printf("Error creating input of ffmpeg: %+v\n", err)
	}

	ffmpegOut, err := ffmpeg.StdoutPipe()

	if err != nil {
		fmt.Printf("Error creating output of mplayer: %+v\n", err)
	}

	return ffmpeg, ffmpegIn, ffmpegOut
}

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

func work(drone Drone, window *gocv.Window, ffmpeg *exec.Cmd, ffmpegIn io.WriteCloser, ffmpegOut io.ReadCloser) {
	go func() {
		//drone.SetupCameraWithMplayer(4, 0)
		drone.SetupCameraWithFfmpeg(window, ffmpeg, ffmpegIn, ffmpegOut, 60, 0, frameSize, 960, 720)
	}()

	//go func() {
	//	println("Started move drone goroutine!")
	//	drone.TakeOff()
	//	SleepSeconds(1)
	//	drone.Hover()
	//	SleepSeconds(3)
	//	drone.Land()
	//}()
}

func main() {
	driver := tello.NewDriver("8888")
	window := gocv.NewWindow("Project 3")

	drone := Drone{driver: driver}

	ffmpeg, ffmpegIn, ffmpegOut := setupFfmpeg(960, 720)

	job := func() {
		//basic(drone)
		work(drone, window, ffmpeg, ffmpegIn, ffmpegOut)
	}

	robot := gobot.NewRobot("Project 3: Drone",
		[]gobot.Connection{},
		[]gobot.Device{driver},
		job,
	)

	err := robot.Start(false)
	if err != nil {
		fmt.Printf("Error starting the Drone: %+v\n", err)
	}

	for {
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}
		img, _ := gocv.NewMatFromBytes(720, 960, gocv.MatTypeCV8UC3, buf)
		if img.Empty() {
			continue
		}

		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
