# Robotics Project 3

Third Project for COMP 570: Robotics at Bridgewater State University.<br>

This project aims to create a drone, and a program that finds a face and then will follow that person, without 
running into the person. <br>
By doing this project, the student will learn more advanced work in robotics with drones and video / computer vision.

## Details

The robot / drone needs to be able to:

* Take off
* Look for the target person
* Follow the target person at a safe distance
* Move left, right, forward and back as needed
* Be directed by hand gesture recognition

## Getting Started

These instructions will get you a copy of the project up and running on your local drone for development 
and testing purposes.

### Prerequisites

* DJI Tello Drone
* Have ffmpeg Installed
* Have Go Installed
* Have Gobot Installed
* Have OpenCV / GoCV Installed

### Running the Project

Git Clone this repository:

```
git clone https://github.com/cmontrond/robotics-third-project
```

CD into the project folder:

```
cd robotics-third-project
```

Compile the source code:

```
go build .
```

Run the project:

```
./robotics-third-project
```

## Built With

* [Go](https://golang.org//) - The Programming Language
* [Gobot](https://gobot.io/) - The Robotics/IoT Framework
* [DJI Tello](https://www.ryzerobotics.com/tello) - The Drone
* [GoCV](https://gocv.io/) - Computer Vision Library

## Author

**Christopher Montrond da Veiga Fernandes** - [Contact](mailto:cmontronddaveigafern@student.bridgew.edu)

## Instructor

**Dr. John Santore** - [Contact](mailto:jsantore@bridgew.edu)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
