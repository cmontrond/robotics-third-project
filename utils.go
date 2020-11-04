package main

import "time"

func SleepSeconds(seconds time.Duration) {
	duration := seconds * time.Second
	time.Sleep(duration)
}

func SleepMilliSeconds(milliSeconds time.Duration) {
	duration := milliSeconds * time.Millisecond
	time.Sleep(duration)
}
