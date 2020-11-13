package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"io"
	"math"
	"os/exec"
	"strconv"
	"time"
)

type pair struct {
	x float64
	y float64
}

const (
	frameX    = 400
	frameY    = 300
	frameSize = frameX * frameY * 3
	offset    = 32767.0
)

var (
	// ffmpeg command to decode video stream from drone
	ffmpeg = exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
		"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameX)+"x"+strconv.Itoa(frameY), "-f", "rawvideo", "pipe:1")
	ffmpegIn, _  = ffmpeg.StdinPipe()
	ffmpegOut, _ = ffmpeg.StdoutPipe()

	// gocv
	window     = gocv.NewWindow("Tello")
	classifier *gocv.CascadeClassifier
	green      = color.RGBA{G: 255}

	// tracking
	tracking                 = false
	detected                 = false
	detectSize               = false
	distTolerance            = 0.05 * dist(0, 0, frameX, frameY)
	refDistance              float64
	left, top, right, bottom int

	// drone
	drone      = tello.NewDriver("8890")
	flightData *tello.FlightData
)

func dist(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)))
}

func init() {
	// process drone events in separate goroutine for concurrency
	go func() {

		if err := ffmpeg.Start(); err != nil {
			fmt.Println(err)
			return
		}

		drone.On(tello.FlightDataEvent, func(data interface{}) {
			// TODO: protect flight data from race condition
			flightData = data.(*tello.FlightData)
		})

		drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			drone.StartVideo()
			drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
			drone.SetExposure(0)
			gobot.Every(100*time.Millisecond, func() {
				drone.StartVideo()
			})
		})

		drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})

		robot := gobot.NewRobot("Project 3 - Drone",
			[]gobot.Connection{},
			[]gobot.Device{drone},
		)

		robot.Start()
	}()
}

func trackFace(frame *gocv.Mat) {
	W := float64(frame.Cols())
	H := float64(frame.Rows())

	imageRectangles := classifier.DetectMultiScale(*frame)

	for _, rect := range imageRectangles {
		fmt.Printf("Found a face: %v\n", rect)

		left = rect.Min.X
		top = rect.Max.Y
		right = rect.Max.X
		bottom = rect.Min.Y

		detected = true

		face := image.Rect(left, top, right, bottom)
		gocv.Rectangle(frame, face, green, 3)
	}

	if !tracking || !detected {
		return
	}

	if detectSize {
		detectSize = false
		refDistance = dist(left, top, right, bottom)
	}

	distance := dist(left, top, right, bottom)

	// x axis
	switch {
	case float64(right) < W/2:
		//drone.CounterClockwise(50)
		println("Drone moving counter clockwise...")
	case float64(left) > W/2:
		//drone.Clockwise(50)
		println("Drone moving clockwise")
	default:
		//drone.Clockwise(0)
		println("Drone not moving clockwise")
	}

	// y axis
	switch {
	case float64(top) < H/10:
		//drone.Up(25)
		println("Drone moving up...")
	case float64(bottom) > H-H/10:
		//drone.Down(25)
		println("Drone moving Down...")
	default:
		//drone.Up(0)
		println("Drone not moving up or down...")
	}

	// z axis
	switch {
	case distance < refDistance-distTolerance:
		//drone.Forward(20)
		println("Drone should move forward...")
	case distance > refDistance+distTolerance:
		//drone.Backward(20)
		println("Drone should move backward...")
	default:
		//drone.Forward(0)
		println("Drone should not move forward...")
	}
}

func main() {

	cascadeClassifier := gocv.NewCascadeClassifier()
	cascadeClassifier.Load("haarcascade_frontalface_default.xml")

	classifier = &cascadeClassifier
	defer classifier.Close()

	for {
		// get next frame from stream
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}
		img, _ := gocv.NewMatFromBytes(frameY, frameX, gocv.MatTypeCV8UC3, buf)
		if img.Empty() {
			continue
		}

		trackFace(&img)

		window.IMShow(img)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}
