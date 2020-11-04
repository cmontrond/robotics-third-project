package main

import "time"

func SleepSeconds(seconds int) {
	duration := time.Duration(seconds) * time.Second
	time.Sleep(duration)
}

func SleepMilliSeconds(milliSeconds int) {
	duration := time.Duration(milliSeconds) * time.Millisecond
	time.Sleep(duration)
}
