package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type Job struct {
	n int
}

// var threads := runtime.NumCPU()

var threads = 8
var rands = make([]*rand.Rand, 0, threads)

func init() {
	fmt.Printf("cpus: %d\n", threads)
	// runtime.GOMAXPROCS(threads)
	runtime.GOMAXPROCS(2)

	for i := 0; i < threads; i++ {
		rands = append(rands, rand.New(rand.NewSource(time.Now().UnixNano())))
	}
}

func MultiPI2(samples int) float64 {
	t1 := time.Now()

	threadSamples := samples / threads

	jobs := make(chan Job, 100)
	// rs := make(chan int, 1000)
	results := make(chan int, 100)

	for w := 0; w < threads; w++ {
		go worker2(w, jobs, results, threadSamples)
	}

	// fmt.Printf("before sleep...\n")
	// time.Sleep(1 * time.Second)
	// fmt.Printf("after sleep...\n")
	go func() {
		for i := 0; i < threads; i++ {
			jobs <- Job{
				n: i,
			}
		}
		close(jobs)
	}()

	var total int
	for i := 0; i < threads; i++ {
		total += <-results
	}

	result := float64(total) / float64(samples) * 4
	fmt.Printf("MultiPI2: %d times, value: %f, cost: %s\n", samples, result, time.Since(t1))
	return result
}
func worker2(id int, jobs <-chan Job, results chan<- int, threadSamples int) {
	for job := range jobs {
		// fmt.Printf("worker id: %d, job: %v, remain jobs: %d\n", id, job, len(jobs))
		var inside int
		// r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r := rands[id]
		fmt.Printf("n: %d\n", job.n)
		for i := 0; i < threadSamples; i++ {
			x, y := r.Float64(), r.Float64()

			if x*x+y*y <= 1 {
				inside++
			}
		}
		// results <- float64(inside) / float64(threadSamples) * 4
		results <- inside
	}
}

func MultiPI(samples int) float64 {
	t1 := time.Now()

	threadSamples := samples / threads
	results := make(chan int, threads)

	for j := 0; j < threads; j++ {
		go func() {
			var inside int
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			for i := 0; i < threadSamples; i++ {
				x, y := r.Float64(), r.Float64()

				if x*x+y*y <= 1 {
					inside++
				}
			}
			// results <- float64(inside) / float64(threadSamples) * 4
			results <- inside
		}()
	}

	var total int
	for i := 0; i < threads; i++ {
		total += <-results
	}

	result := float64(total) / float64(samples) * 4
	fmt.Printf("MultiPI: %d times, value: %f, cost: %s\n", samples, result, time.Since(t1))
	return result
}

func PI(samples int) (result float64) {
	t1 := time.Now()
	var inside int = 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < samples; i++ {
		x := r.Float64()
		y := r.Float64()
		if (x*x + y*y) < 1 {
			inside++
		}
	}

	ratio := float64(inside) / float64(samples)

	result = ratio * 4

	fmt.Printf("PI: %d times, value: %fr cost: %s\n", samples, result, time.Since(t1))

	return
}

// func init() {
// 	runtime.GOMAXPROCS(runtime.NumCPU())
// 	rand.Seed(time.Now().UnixNano())
// }

func main() {
	samples := 100000000
	// fmt.Println("Our value of Pi after 100 runs:\t\t\t", PI(100))
	// fmt.Println("Our value of Pi after 1,000 runs:\t\t", PI(1000))
	// fmt.Println("Our value of Pi after 10,000 runs:\t\t", PI(10000))
	// fmt.Println("Our value of Pi after 100,000 runs:\t\t", PI(100000))
	// fmt.Println("Our value of Pi after 1,000,000 runs:\t\t", PI(1000000))
	// fmt.Println("Our value of Pi after 200,000,000 runs:\t\t", PI(2000000000))
	// fmt.Println("Our value of Pi after 10,000,000 runs MultiPI:\t\t", MultiPI(1000000000))

	PI(samples)
	MultiPI(samples)
	MultiPI2(samples)
	// fmt.Printf("Our value of Pi after %d runs is:\t\t %f\n", samples, PI(samples))
	// fmt.Printf("Our value of Pi after %d runs is:\t\t %f\n", samples, MultiPI(samples))
	// fmt.Printf("Our value of Pi after %d runs is:\t\t %f\n", samples, MultiPI2(samples))

	// fmt.Println("Our value of Pi after 100,000,000 runs:\t\t", PI(100000000))
}
