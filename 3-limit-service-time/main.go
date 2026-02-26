//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

var timelimit int64 = 10

var mu sync.Mutex

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {

	if u.IsPremium {
		process()
		return true
	}

	mu.Lock()
	timeRemaining := timelimit - u.TimeUsed
	mu.Unlock()

	fmt.Printf("Uwer %d has %d seconds \n", u.ID, timeRemaining)

	if timeRemaining <= 0 {
		return false
	}

	timer := time.NewTimer(time.Duration(timeRemaining) * time.Second)
	defer timer.Stop()

	start := time.Now()
	done := make(chan struct{})

	go func() {
		defer close(done)
		process()

	}()

	select {
	case <-timer.C:
		mu.Lock()
		u.TimeUsed = timelimit
		mu.Unlock()
		return false
	case <-done:
		mu.Lock()

		remaining := int64((time.Since(start) + time.Second - 1) / time.Second)

		if remaining+u.TimeUsed > timelimit {
			u.TimeUsed = timelimit
		} else {
			u.TimeUsed += remaining
		}
		mu.Unlock()
		return true
	}
}

func main() {
	RunMockServer()
}
