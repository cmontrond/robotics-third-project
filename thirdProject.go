package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"image/color"
	"io"
	"log"
	"math"
	"os/exec"
	"strconv"
	"time"
)

const (
	frameX    = 360
	frameY    = 240
	frameSize = frameX * frameY * 3
)

var (
	// ffmpeg
	ffmpeg = exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
		"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameX)+"x"+strconv.Itoa(frameY), "-f", "rawvideo", "pipe:1")

	ffmpegIn, _  = ffmpeg.StdinPipe()
	ffmpegOut, _ = ffmpeg.StdoutPipe()

	// GO CV
	window     = gocv.NewWindow("Project 3 - Drone")
	net        *gocv.Net
	green      = color.RGBA{G: 255}
	classifier *gocv.CascadeClassifier

	// tracking
	tracking                 = false
	detected                 = false
	detectSize               = false
	distTolerance            = 0.05 * dist(0, 0, frameX, frameY)
	refDistance              float64
	left, top, right, bottom float64

	// Drone
	droneDriver = tello.NewDriver("8890")
	drone       = Drone{droneDriver}
	flightData  *tello.FlightData
)

func init() {
	// Drone events
	go func() {

		if err := ffmpeg.Start(); err != nil {
			fmt.Println(err)
			return
		}

		droneDriver.On(tello.FlightDataEvent, func(data interface{}) {
			flightData = data.(*tello.FlightData)
			fmt.Printf("Battery Percentage: %v %", flightData.BatteryPercentage)
		})

		droneDriver.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			droneDriver.StartVideo()
			droneDriver.SetVideoEncoderRate(tello.VideoBitRateAuto)
			droneDriver.SetExposure(0)
			gobot.Every(100*time.Millisecond, func() {
				droneDriver.StartVideo()
			})
		})

		droneDriver.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		})

		robot := gobot.NewRobot("Project 3 - Drone",
			[]gobot.Connection{},
			[]gobot.Device{droneDriver},
		)

		robot.Start()
	}()
}

func dist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}

func trackFace(frame *gocv.Mat) {

	W := float64(frame.Cols())
	H := float64(frame.Rows())

	//blob := gocv.BlobFromImage(*frame, 1.0, image.Pt(300, 300), gocv.NewScalar(104, 177, 123, 0), false, false)
	//defer blob.Close()
	//
	//net.SetInput(blob, "data")
	//
	//detBlob := net.Forward("detection_out")
	//defer detBlob.Close()
	//
	//detections := gocv.GetBlobChannel(detBlob, 0, 0)
	//defer detections.Close()

	imageRectangles := classifier.DetectMultiScale(*frame)

	for _, rect := range imageRectangles {
		log.Println("found a face,", rect)
		gocv.Rectangle(frame, rect, green, 3)
	}

	//for r := 0; r < detections.Rows(); r++ {
	//	confidence := detections.GetFloatAt(r, 2)
	//	if confidence < 0.5 {
	//		continue
	//	}
	//
	//	left = float64(detections.GetFloatAt(r, 3)) * W
	//	top = float64(detections.GetFloatAt(r, 4)) * H
	//	right = float64(detections.GetFloatAt(r, 5)) * W
	//	bottom = float64(detections.GetFloatAt(r, 6)) * H
	//
	//	left = math.Min(math.Max(0.0, left), W-1.0)
	//	right = math.Min(math.Max(0.0, right), W-1.0)
	//	bottom = math.Min(math.Max(0.0, bottom), H-1.0)
	//	top = math.Min(math.Max(0.0, top), H-1.0)
	//
	//	detected = true
	//	rect := image.Rect(int(left), int(top), int(right), int(bottom))
	//	gocv.Rectangle(frame, rect, green, 3)
	//}

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
	case right < W/2:
		//droneDriver.CounterClockwise(50)
		println("Drone should turn counter clockwise 50 ...")
	case left > W/2:
		//droneDriver.Clockwise(50)
		println("Drone should turn clockwise 50 ...")
	default:
		//droneDriver.Clockwise(0)
		println("Drone should turn counter clockwise 0 ...")
	}

	// y axis
	switch {
	case top < H/10:
		//droneDriver.Up(25)
		println("Drone should move up...")
	case bottom > H-H/10:
		//droneDriver.Down(25)
		println("Drone should move down...")
	default:
		//droneDriver.Up(0)
		println("Drone should move up...")
	}

	// z axis
	switch {
	case distance < refDistance-distTolerance:
		//droneDriver.Forward(20)
		println("Drone should move forward...")
	case distance > refDistance+distTolerance:
		//droneDriver.Backward(20)
		println("Drone should move backward...")
	default:
		//droneDriver.Forward(0)
		println("Drone should move forward")
	}
}

func main() {

	cascadeClassifier := gocv.NewCascadeClassifier()
	cascadeClassifier.Load("haarcascade_frontalface_default.xml")
	defer cascadeClassifier.Close()

	classifier = &cascadeClassifier

	//model := "model.caffemodel"
	//proto := "proto.txt"
	//
	//// open DNN classifier
	//n := gocv.ReadNetFromCaffe(proto, model)
	//if n.Empty() {
	//	fmt.Printf("Error reading network model from : %v %v\n", proto, model)
	//	return
	//}
	//net = &n
	//defer net.Close()
	//net.SetPreferableBackend(gocv.NetBackendDefault)
	//net.SetPreferableTarget(gocv.NetTargetCPU)

	for {
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
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
