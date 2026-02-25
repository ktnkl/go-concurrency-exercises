// //////////////////////////////////////////////////////////////////////
// //
// // Given is a producer-consumer scenario, where a producer reads in
// // tweets from a mockstream and a consumer is processing the
// // data. Your task is to change the code so that the producer as well
// // as the consumer can run concurrently
// //

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func producer(stream Stream, ch chan<- *Tweet) {
// 	for {
// 		tweet, err := stream.Next()
// 		if err == ErrEOF {
// 			return
// 		}
// 		ch <- tweet
// 	}
// }

// func consumer(in <-chan *Tweet, out chan<- string) {
// 	for t := range in {
// 		if t.IsTalkingAboutGo() {
// 			out <- fmt.Sprintln(t.Username, "\ttweets about golang")
// 		} else {
// 			out <- fmt.Sprintln(t.Username, "\tdoes not tweet about golang")
// 		}
// 	}
// }

// func main() {
// 	start := time.Now()
// 	var in chan *Tweet
// 	var out chan string
// 	var wg sync.WaitGroup

// 	stream := GetMockStream()

// 	// Consumer
// 	go func() {
// 		wg.Go(func() {
// 			consumer(in, out)
// 		})
// 	}()

// 	// Producer
// 	go func() {
// 		defer wg.Done()
// 		producer(stream, in)
// 	}()

// 	go func() {
// 		close(out)
// 		wg.Wait()

// 	}()

// 	fmt.Printf("Process took %s\n", time.Since(start))

// }

//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

var ch = make(chan *Tweet)

func producer(stream Stream, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return
		}
		ch <- tweet
	}
}

func consumer(wg *sync.WaitGroup) {
	defer wg.Done()

	for t := range ch {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	wg := &sync.WaitGroup{}

	// Producer
	wg.Add(1)
	go producer(stream, wg)

	// Consumer
	wg.Add(1)
	go consumer(wg)

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
