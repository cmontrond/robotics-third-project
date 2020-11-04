package main

import (
	"fmt"
	"io"
	"os/exec"
	"time"
)

func SleepSeconds(seconds time.Duration) {
	duration := seconds * time.Second
	time.Sleep(duration)
}

func SleepMilliSeconds(milliSeconds time.Duration) {
	duration := milliSeconds * time.Millisecond
	time.Sleep(duration)
}

func SetupMplayer() io.WriteCloser {
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
