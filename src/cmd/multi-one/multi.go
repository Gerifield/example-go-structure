package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	ex := flag.Int("example", 1, "Example number")
	flag.Parse()

	switch *ex {
	case 1:
		example1()
	case 2:
		example2()
	case 3:
		example3()
	case 4:
		example4()
	case 5:
		example5()
	case 6:
		example6()
	case 7:
		example7()
	case 8:
		example8()
	default:
		log.Fatalln("Unknown example")
		// Some more examples:
		// https://www.opsdash.com/blog/job-queues-in-go.html
		// https://www.leolara.me/blog/closing_a_go_channel_written_by_several_goroutines/
		// https://callistaenterprise.se/blogg/teknik/2019/10/05/go-worker-cancellation/
	}
}

// example1 shows how to start simple goroutines
func example1() {
	go func() { fmt.Println("Hello 1") }()
	go func() { fmt.Println("Hello 2") }()
	go func() { fmt.Println("Hello 3") }()
	go func() { fmt.Println("Hello 4") }()
	go func() { fmt.Println("Hello 5") }()

	time.Sleep(1 * time.Second)
	fmt.Println("Finish")
}

// example2 shows how to make some basic processing with sleep and a common variable overwrite mistake
func example2() {
	for i := 1; i <= 5; i++ {
		go func() { fmt.Printf("Hello %d\n", i) }() // WRONG!
		//go func(i int) { fmt.Printf("Hello %d\n", i) }(i)
	}

	time.Sleep(1 * time.Second)
	fmt.Println("Finish")
}

// example3 shows how NOT to make job processing
func example3() {
	var finish chan struct{} // fatal error: all goroutines are asleep - deadlock!
	//finish := make(chan struct{}) // fatal error: all goroutines are asleep - deadlock!
	//finish := make(chan struct{}, 1) // still deadlock if not closed!

	for i := 1; i <= 5; i++ {
		go func(i int) {
			fmt.Printf("Hello %d\n", i)
			finish <- struct{}{}
			//if i == 5 { // Very ugly and hacky close which does not even work most of the times!
			//	close(finish)
			//}
		}(i)
	}

	// This loops until the channel is closed (draining the channel), if not closed and there are no writer -> deadlock
	for range finish {
		fmt.Println("OK")
	}

	// Alternate finish for
	//for i := 0; i < 5; i++ {
	//	fmt.Println("OK") // Works, but sometimes even this would not print the results properly
	//}

	fmt.Println("Finish")
}

// example4 shows how to wait for job results without sleep
func example4() {
	var wg sync.WaitGroup

	wg.Add(5)
	for i := 1; i <= 5; i++ {
		go func(i int) {
			defer wg.Done()

			fmt.Printf("Hello %d\n", i)
		}(i)
	}

	wg.Wait() // Without sleep!
	fmt.Println("DONE")
}

// example5 shows a simple multi threaded job queue processing
func example5() {
	var wg sync.WaitGroup
	workerNum := 2
	start := time.Now()

	// Define the queue processor
	worker := func(jobs <-chan int) {
		defer wg.Done()

		for j := range jobs {
			fmt.Printf("Hello %d\n", j)

			// Emulate some working time to make a difference between buffered and non buffered channels
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Setup the job queue (I'd prefer a buffered one here for this!)
	//var jobs chan int // Nil channel -> error, don't forget to make it!
	jobs := make(chan int)
	//jobs := make(chan int, 50) // The sender part takes much less time

	// Start some worker process (Add 1 to the work group counter for each)
	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		go worker(jobs)
	}

	// Add some jobs to the queue
	for i := 0; i < 100; i++ {
		jobs <- i
	}

	//Signal the end of the jobs, this will stop the workers after they finish the queue
	// Always try to close at the producer side!
	fmt.Printf("Sending finished after %s\n", time.Now().Sub(start))
	close(jobs)

	// Wait until all the workers finish their jobs
	wg.Wait()
	fmt.Printf("Finished in %s\n", time.Now().Sub(start))
}

// example6 shows how to signal a finish with a simple channel
func example6() {
	stop := make(chan struct{})

	go func(stop chan struct{}) {
		time.Sleep(1 * time.Second)
		close(stop)
	}(stop)

	fmt.Println("Wait for the worker...")

	<-stop
	fmt.Println("Finished!")
}

// example7 shows and non-blocked sending with limited channel length (it'll drop the overflowing values)
func example7() {
	var wg sync.WaitGroup
	jobs := make(chan int, 5)
	start := time.Now()
	workerNum := 2

	var atomicCounter int64 // This is also an example for the atomic integer changes (use this or a mutex)

	worker := func(jobs <-chan int) {
		defer wg.Done()

		for j := range jobs {
			fmt.Printf("Hello %d\n", j)
			time.Sleep(100 * time.Millisecond)

			atomic.AddInt64(&atomicCounter, 1)
		}
	}

	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		go worker(jobs)
	}

	// Spam some jobs in a non-blocking way
	// You could read a channel with select in a non-blocking way too!
	for i := 0; i < 100; i++ {
		select {
		case jobs <- i:
			fmt.Printf("Added %d\n", i)
		default:
			fmt.Printf("Failed %d\n", i)
		}
	}

	fmt.Printf("Sending finished after %s\n", time.Now().Sub(start))
	close(jobs)

	// Wait for the processes
	wg.Wait()
	fmt.Printf("Finished in %s\nAtomic integer counter at: %d\n", time.Now().Sub(start), atomicCounter)
}

//example8 shows how to add a timeout/cancel for the job processing
func example8() {
	jobs := make(chan int, 5)
	var wg sync.WaitGroup
	start := time.Now()

	worker := func(ctx context.Context, jobs <-chan int) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case j, ok := <-jobs: // Check if the channel closed
				if !ok || ctx.Err() != nil { // Make sure we exit if we run on this path after the cancel call
					return
				}
				fmt.Printf("Hello %d (I'll go to sleep)\n", j)
				time.Sleep(4 * time.Second)
				fmt.Printf("Hello %d finished!\n", j)
			}
		}
	}

	// Same as below, but in a different, more complex way
	//ctx, cancel := context.WithCancel(context.Background())
	//// Tricky, but will call cancel... Use context.WithTimeout() instead, this is only an example
	//time.AfterFunc(5*time.Second, func() {
	//	fmt.Println("Calling cancel...")
	//	cancel()
	//	fmt.Println("Called")
	//})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // You should never forget to call this at least at the finis
	// Imagine that you have a fast job and a long timeout, when the job finishes, it should stop the timeout counter in the background to not hold unused resources

	wg.Add(2)
	go worker(ctx, jobs)
	go worker(ctx, jobs)

	for i := 0; i < 100; i++ {
		// At some point we'll force stop the worker, so we'll have no listeners on the channels -> this could cause a deadlock
		select {
		case jobs <- i:
		default:
		}
	}
	close(jobs) // I don't like to keep this open

	wg.Wait() // Wait no matter what for the processes
	fmt.Printf("Finished in %s\n", time.Now().Sub(start))
}
