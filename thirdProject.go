package main

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"io"
	"log"
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
	frameX    = 720
	frameY    = 960
	frameSize = frameX * frameY * 3
)

var (
	// ffmpeg command to decode video stream from drone
	//ffmpeg = exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
	//	"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameY)+"x"+strconv.Itoa(frameX), "-f", "rawvideo", "pipe:1")

	ffmpeg = exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-vcodec", "rawvideo",
		"-an", "-sn", "-s", strconv.Itoa(frameY)+"x"+strconv.Itoa(frameX), "-f", "rawvideo", "pipe:1")

	//ffmpeg = exec.Command("ffmpeg", "-hwaccel", "auto", "-hwaccel_device", "opencl", "-i", "pipe:0",
	//	"-pix_fmt", "bgr24", "-s", strconv.Itoa(frameX)+"x"+strconv.Itoa(frameY), "-f", "rawvideo", "pipe:1")

	ffmpegIn, _  = ffmpeg.StdinPipe()
	ffmpegOut, _ = ffmpeg.StdoutPipe()

	// gocv
	window     = gocv.NewWindow("Tello")
	classifier *gocv.CascadeClassifier
	green      = color.RGBA{G: 255}

	// tracking
	tracking   = true
	detected   = false
	detectSize = true
	//distTolerance            = 0.05 * dist(0, 0, frameX, frameY)
	distTolerance            = 0.05 * dist(0, 0, 120, 90) // TODO: Change things here, maybe to 0.10
	refDistance              float64
	left, top, right, bottom float64

	// drone
	drone      = tello.NewDriver("8890")
	flightData *tello.FlightData
)

func dist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}

func init() {
	// process drone events in separate goroutine for concurrency
	go func() {

		if err := ffmpeg.Start(); err != nil {
			fmt.Println(err)
			return
		}

		//if err := drone.On(tello.FlightDataEvent, func(data interface{}) {
		//	// TODO: protect flight data from race condition
		//	flightData = data.(*tello.FlightData)
		//	//println("Battery: ", flightData.BatteryPercentage)
		//}); err != nil {
		//	println("Error in FlightDataEvent: ", err)
		//}

		if err := drone.On(tello.ConnectedEvent, func(data interface{}) {
			fmt.Println("Connected")
			if err := drone.StartVideo(); err != nil {
				println("Error in StartVideo: ", err)
			}

			if err := drone.SetVideoEncoderRate(tello.VideoBitRateAuto); err != nil {
				println("Error in SetVideoEncoderRate: ", err)
			}

			if err := drone.SetExposure(0); err != nil {
				println("Error in SetExposure: ", err)
			}

			gobot.Every(100*time.Millisecond, func() {
				if err := drone.StartVideo(); err != nil {
					println("Error in StartVideo: ", err)
				}
			})
		}); err != nil {
			println("Error in ConnectedEvent: ", err)
		}

		if err := drone.On(tello.VideoFrameEvent, func(data interface{}) {
			pkt := data.([]byte)
			if _, err := ffmpegIn.Write(pkt); err != nil {
				fmt.Println(err)
			}
		}); err != nil {
			println("Error in VideFrameEvent: ", err)
		}

		robot := gobot.NewRobot("Project 3 - Drone",
			[]gobot.Connection{},
			[]gobot.Device{drone},
		)

		if err := robot.Start(false); err != nil {
			println("Error in robot.Start: : ", err)
		}
	}()
}

func trackFace(frame *gocv.Mat) {
	W := float64(frame.Cols())
	H := float64(frame.Rows())

	imageRectangles := classifier.DetectMultiScale(*frame)

	if len(imageRectangles) == 0 {
		return
	}

	faces := make(map[float64]image.Rectangle)

	for _, rect := range imageRectangles {

		left = float64(rect.Min.X)
		top = float64(rect.Max.Y)
		right = float64(rect.Max.X)
		bottom = float64(rect.Min.Y)

		left = math.Min(math.Max(0.0, left), W-1.0)
		right = math.Min(math.Max(0.0, right), W-1.0)
		bottom = math.Min(math.Max(0.0, bottom), H-1.0)
		top = math.Min(math.Max(0.0, top), H-1.0)

		w := right - left
		h := top - bottom

		area := w * h

		detected = true

		faces[area] = rect
	}

	if !detected {
		return
	}

	if len(faces) > 0 {
		max := 0.0

		for key := range faces {
			if key > max {
				max = key
			}
		}

		//fmt.Printf("Face Rectangle: ", faces[max])

		gocv.Rectangle(frame, faces[max], green, 3)

		//println("Found a face!")

		left = float64(faces[max].Min.X)
		top = float64(faces[max].Max.Y)
		right = float64(faces[max].Max.X)
		bottom = float64(faces[max].Min.Y)

		left = math.Min(math.Max(0.0, left), W-1.0)
		right = math.Min(math.Max(0.0, right), W-1.0)
		bottom = math.Min(math.Max(0.0, bottom), H-1.0)
		top = math.Min(math.Max(0.0, top), H-1.0)

		if detectSize {
			detectSize = false
			refDistance = dist(left, top, right, bottom)
		}

		distance := dist(left, top, right, bottom)

		// x axis
		//switch {
		//case right < W/2:
		//	//drone.CounterClockwise(50)
		//	println("Drone moving counter clockwise...")
		//case left > W/2:
		//	//drone.Clockwise(50)
		//	println("Drone moving clockwise")
		//default:
		//	//drone.Clockwise(0)
		//	println("Drone not moving clockwise")
		//}

		// y axis
		//switch {
		//case top < H/10:
		//	//drone.Up(25)
		//	println("Drone moving up...")
		//case bottom > H-H/10:
		//	//drone.Down(25)
		//	println("Drone moving Down...")
		//default:
		//	//drone.Up(0)
		//	println("Drone not moving up or down...")
		//}

		// z axis
		switch {
		case distance < refDistance-distTolerance:
			//drone.Forward(20)
			println("Drone should move forward...")
			SleepSeconds(2)
		//case distance > refDistance+distTolerance:
		//	//drone.Backward(20)
		//	println("Drone should move backward...")
		//	//SleepSeconds(2)
		default:
			//drone.Forward(0)
			//drone.Backward(0)
			println("Drone should not move forward...")
			// TODO: Maybe turn around when you can't find a face
			//SleepSeconds(2)
		}

		// TODO: Do this only if the drone is at a safe enough distance
		//handleGestures(frame)
	}
}

func handleGestures(img *gocv.Mat) {

	imgGrey := gocv.NewMat()
	defer imgGrey.Close()

	imgBlur := gocv.NewMat()
	defer imgBlur.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	hull := gocv.NewMat()
	defer hull.Close()

	defects := gocv.NewMat()
	defer defects.Close()

	// cleaning up image
	gocv.CvtColor(*img, &imgGrey, gocv.ColorBGRToGray)
	gocv.GaussianBlur(imgGrey, &imgBlur, image.Pt(35, 35), 0, 0, gocv.BorderDefault)
	gocv.Threshold(imgBlur, &imgThresh, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)

	// now find biggest contour
	contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	c := getBiggestContour(contours)

	gocv.ConvexHull(c, &hull, true, false)
	gocv.ConvexityDefects(c, hull, &defects)

	var angle float64
	defectCount := 0
	for i := 0; i < defects.Rows(); i++ {
		start := c[defects.GetIntAt(i, 0)]
		end := c[defects.GetIntAt(i, 1)]
		far := c[defects.GetIntAt(i, 2)]

		a := math.Sqrt(math.Pow(float64(end.X-start.X), 2) + math.Pow(float64(end.Y-start.Y), 2))
		b := math.Sqrt(math.Pow(float64(far.X-start.X), 2) + math.Pow(float64(far.Y-start.Y), 2))
		c := math.Sqrt(math.Pow(float64(end.X-far.X), 2) + math.Pow(float64(end.Y-far.Y), 2))

		// apply cosine rule here
		angle = math.Acos((math.Pow(b, 2)+math.Pow(c, 2)-math.Pow(a, 2))/(2*b*c)) * 57

		// ignore angles > 90 and highlight rest with dots
		if angle <= 90 {
			defectCount++
			gocv.Circle(img, far, 1, green, 2)
		}
	}

	status := fmt.Sprintf("defectCount: %d", defectCount+1)

	//rect := gocv.BoundingRect(c)
	//gocv.Rectangle(img, rect, color.RGBA{R: 255, G: 255, B: 255}, 2)

	gocv.PutText(img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, green, 2)

	switch {
	case (defectCount + 1) == 3:
		// TODO: Make drone do something, like go back or forward, or do a back flip
	case (defectCount + 1) == 1:
		// TODO: Make drone do something else
	}
}

func getBiggestContour(contours [][]image.Point) []image.Point {
	var area float64
	index := 0
	for i, c := range contours {
		newArea := gocv.ContourArea(c)
		if newArea > area {
			area = newArea
			index = i
		}
	}
	return contours[index]
}

func resizeFrame(frame gocv.Mat, targetSize image.Point) gocv.Mat {
	result := gocv.NewMatWithSize(targetSize.X, targetSize.Y, gocv.MatTypeCV8UC3)
	gocv.Resize(frame, &result, image.Pt(targetSize.X, targetSize.Y), 0, 0, gocv.InterpolationNearestNeighbor)
	return result
}

func main() {

	cascadeClassifier := gocv.NewCascadeClassifier()
	cascadeClassifier.Load("haarcascade_frontalface_default.xml")

	classifier = &cascadeClassifier
	defer classifier.Close()

	doTakeOff := false

	for {

		if doTakeOff {
			if err := drone.TakeOff(); err != nil {
				println("Error in TakeOff: ", err)
			}
			doTakeOff = false
		}

		// get next frame from stream
		buf := make([]byte, frameSize)
		if _, err := io.ReadFull(ffmpegOut, buf); err != nil {
			fmt.Println(err)
			continue
		}

		img, err1 := gocv.NewMatFromBytes(frameX, frameY, gocv.MatTypeCV8UC3, buf)
		if err1 != nil {
			log.Print(err1)
			continue
		}
		if img.Empty() {
			continue
		}

		img = resizeFrame(img, image.Point{
			X: 120,
			Y: 90,
		})

		trackFace(&img)
		//handleGestures(&img)

		window.IMShow(img)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}
