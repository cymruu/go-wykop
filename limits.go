package wykop

import (
	"fmt"
	"time"
)

const (
	wykopDefaultSessionLength        = time.Hour
	wykopDefaultAllowedReqPerSession = 250
)

type limitOptional func(*wykopLimits)
type wykopLimits struct {
	sessionStart         time.Time
	requestsMade         uint32
	sessionLength        time.Duration
	allowedReqPerSession uint32
}

func optionalSesssionLength(v time.Duration) limitOptional {
	return func(target *wykopLimits) {
		target.sessionLength = v
	}
}

func optionalAllowedReqPerSession(v uint32) limitOptional {
	return func(target *wykopLimits) {
		target.allowedReqPerSession = v
	}
}
func initializeLimits(options ...limitOptional) *wykopLimits {
	limits := &wykopLimits{sessionLength: wykopDefaultSessionLength, allowedReqPerSession: wykopDefaultAllowedReqPerSession}
	for _, op := range options {
		op(limits)
	}
	return limits
}
func (limits *wykopLimits) register() {
	if time.Now().Sub(limits.sessionStart) > limits.sessionLength {
		limits.sessionStart = time.Now()
		limits.requestsMade = 0
	}
	limits.requestsMade++
}
func (limits *wykopLimits) GetResetTime() time.Time {
	return limits.sessionStart.Add(limits.sessionLength)
}
func (limits *wykopLimits) GetTimeToReset() time.Duration {
	return limits.sessionLength - time.Now().Sub(limits.sessionStart)
}
func (limits *wykopLimits) TimeSinceReset() time.Duration {
	return time.Now().Sub(limits.sessionStart)
}
func (limits *wykopLimits) GetUsage() float32 {
	return float32(limits.requestsMade) / float32(limits.allowedReqPerSession)
}
func (limits *wykopLimits) String() string {
	return fmt.Sprintf("%.2f%%", limits.GetUsage())
}
