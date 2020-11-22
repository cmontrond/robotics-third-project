package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
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
	// gocv
	window     = gocv.NewWindow("Tello")
	classifier *gocv.CascadeClassifier
	green      = color.RGBA{G: 255}

	// tracking
	tracking                 = true
	detected                 = false
	detectSize               = true
	distTolerance            = 0.05 * dist(0, 0, frameX, frameY)
	refDistance              float64
	left, top, right, bottom float64
)

func dist(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
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
		//gocv.Rectangle(frame, rect, green, 3)
		//fmt.Printf("Found a face: %v\n", rect)

		//left = float64(rect.Min.X) * W
		//top = float64(rect.Max.Y) * H
		//right = float64(rect.Max.X) * W
		//bottom = float64(rect.Min.Y) * H

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

		// TODO: Maybe display rectangle here and then overrride left/top/right using the W and H values
		// TODO: Should disregard other faces
		//fmt.Printf("face w: %v\n", w)
		//fmt.Printf("face h: %v\n", h)
		//fmt.Printf("face area: %v\n", area)
		//fmt.Printf("W: %v\n", W)
		//fmt.Printf("H: %v\n", H)
		//fmt.Printf("Left: %v\n", left)
		//fmt.Printf("Right: %v\n", right)
		//fmt.Printf("Bottom: %v\n", bottom)
		//fmt.Printf("Top: %v\n\n", top)
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

		gocv.Rectangle(frame, faces[max], green, 3)
	}

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
	case distance > refDistance+distTolerance:
		//drone.Backward(20)
		println("Drone should move backward...")
	default:
		//drone.Forward(0)
		println("Drone should not move forward...")
	}

	// TODO: Do this only if the drone is at a safe enough distance
	handleGestures(frame)
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

func main() {

	cascadeClassifier := gocv.NewCascadeClassifier()
	cascadeClassifier.Load("haarcascade_frontalface_default.xml")

	classifier = &cascadeClassifier
	defer classifier.Close()

	// open webcam
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	//doTakeOff := true

	for {

		//if doTakeOff {
		//	drone.TakeOff()
		//	drone.Up(30)
		//	SleepSeconds(2)
		//	doTakeOff = false
		//} else {
		//	drone.Hover()
		//}

		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", 0)
			return
		}
		if img.Empty() {
			continue
		}

		trackFace(&img)

		window.IMShow(img)
		if window.WaitKey(10) >= 0 {
			println("Land the robot...")
			break
		}
	}
}
