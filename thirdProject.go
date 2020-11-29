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

const (
	frameX    = 720
	frameY    = 960
	frameSize = frameX * frameY * 3
)

var (
	// ffmpeg command
	ffmpeg = exec.Command("ffmpeg", "-i", "pipe:0", "-pix_fmt", "bgr24", "-vcodec", "rawvideo",
		"-an", "-sn", "-s", strconv.Itoa(frameY)+"x"+strconv.Itoa(frameX), "-f", "rawvideo", "pipe:1")

	ffmpegIn, _  = ffmpeg.StdinPipe()
	ffmpegOut, _ = ffmpeg.StdoutPipe()

	// gocv
	window     = gocv.NewWindow("Tello")
	classifier *gocv.CascadeClassifier
	green      = color.RGBA{G: 255}

	// tracking
	trackingEnabled          = true
	faceDetected             = false
	shouldDetectSize         = true
	distanceTolerance        = 0.001 * dist(0, 0, 120, 90)
	referenceDistance        float64
	left, top, right, bottom float64

	// drone
	drone      = tello.NewDriver("8890")
	flightData *tello.FlightData
)

func dist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
}

func init() {
	// drone events
	go func() {

		if err := ffmpeg.Start(); err != nil {
			fmt.Println(err)
			return
		}

		if err := drone.On(tello.FlightDataEvent, func(data interface{}) {
			flightData = data.(*tello.FlightData)
			println("Battery: ", flightData.BatteryPercentage)
		}); err != nil {
			println("Error in FlightDataEvent: ", err)
		}

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

func handleTrackFace(frame *gocv.Mat) {
	frameWidth := float64(frame.Cols())
	frameHeight := float64(frame.Rows())

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

		left = math.Min(math.Max(0.0, left), frameWidth-1.0)
		right = math.Min(math.Max(0.0, right), frameWidth-1.0)
		bottom = math.Min(math.Max(0.0, bottom), frameHeight-1.0)
		top = math.Min(math.Max(0.0, top), frameHeight-1.0)

		w := right - left
		h := top - bottom

		area := w * h

		faceDetected = true

		faces[area] = rect
	}

	if !faceDetected {
		return
	}

	if len(faces) > 0 {
		max := 0.0

		for key := range faces {
			if key > max {
				max = key
			}
		}

		gocv.Rectangle(frame, faces[max], green, 3)

		left = float64(faces[max].Min.X)
		top = float64(faces[max].Max.Y)
		right = float64(faces[max].Max.X)
		bottom = float64(faces[max].Min.Y)

		left = math.Min(math.Max(0.0, left), frameWidth-1.0)
		right = math.Min(math.Max(0.0, right), frameWidth-1.0)
		bottom = math.Min(math.Max(0.0, bottom), frameHeight-1.0)
		top = math.Min(math.Max(0.0, top), frameHeight-1.0)

		if shouldDetectSize {
			shouldDetectSize = false
			referenceDistance = dist(left, top, right, bottom)
		}

		distance := dist(left, top, right, bottom)

		// right and left
		switch {
		case right < frameWidth/2:
			//drone.CounterClockwise(20)
			drone.Left(15)
			println("Drone should go left...")
		case left > frameWidth/2:
			//drone.Clockwise(20)
			drone.Right(15)
			println("Drone should go right...")
		default:
			drone.Left(0)
			drone.Right(0)
			//drone.Clockwise(0)
		}

		// up and down
		switch {
		case top < frameHeight/10:
			drone.Up(25)
			println("Drone should go up...")
		case bottom > frameHeight-frameHeight/10:
			drone.Down(25)
			println("Drone should go down...")
		default:
			drone.Up(0)
			drone.Down(0)
		}

		// forward
		switch {
		case distance < referenceDistance-distanceTolerance:
			drone.Forward(20)
			println("Drone should move forward...")
		default:
			drone.Forward(0)
			println("Drone should not move forward...")
		}
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

	gocv.CvtColor(*img, &imgGrey, gocv.ColorBGRToGray)
	gocv.GaussianBlur(imgGrey, &imgBlur, image.Pt(35, 35), 0, 0, gocv.BorderDefault)
	gocv.Threshold(imgBlur, &imgThresh, 0, 255, gocv.ThresholdBinaryInv+gocv.ThresholdOtsu)

	contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	c := getBiggestContour(contours)

	gocv.ConvexHull(c, &hull, true, false)
	gocv.ConvexityDefects(c, hull, &defects)

	var angle float64
	numberOfFingers := 0
	for i := 0; i < defects.Rows(); i++ {
		start := c[defects.GetIntAt(i, 0)]
		end := c[defects.GetIntAt(i, 1)]
		far := c[defects.GetIntAt(i, 2)]

		a := math.Sqrt(math.Pow(float64(end.X-start.X), 2) + math.Pow(float64(end.Y-start.Y), 2))
		b := math.Sqrt(math.Pow(float64(far.X-start.X), 2) + math.Pow(float64(far.Y-start.Y), 2))
		c := math.Sqrt(math.Pow(float64(end.X-far.X), 2) + math.Pow(float64(end.Y-far.Y), 2))

		angle = math.Acos((math.Pow(b, 2)+math.Pow(c, 2)-math.Pow(a, 2))/(2*b*c)) * 57

		if angle <= 90 {
			numberOfFingers++
			gocv.Circle(img, far, 1, green, 2)
		}
	}

	status := fmt.Sprintf("numberOfFingers: %d", numberOfFingers+1)

	rect := gocv.BoundingRect(c)
	gocv.Rectangle(img, rect, color.RGBA{R: 255, G: 255, B: 255}, 2)

	gocv.PutText(img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, green, 2)

	switch {
	case (numberOfFingers + 1) == 3:
		// Make drone do a back flip
		if err := drone.BackFlip(); err != nil {
			println("Error making the drone perform a back flip...")
		}
	case (numberOfFingers + 1) == 1:
		drone.Hover()
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

		// Resize the image to increase performance
		img = resizeFrame(img, image.Point{
			X: 120,
			Y: 90,
		})

		handleTrackFace(&img)

		//handleGestures(&img)

		window.IMShow(img)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}
