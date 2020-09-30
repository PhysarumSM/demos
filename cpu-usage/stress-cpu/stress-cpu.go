// Adapted from:
// https://www.reddit.com/r/golang/comments/crt5ou/using_algorithms_to_simulate_cpu_intensive_work/
// https://gist.github.com/danielcasler/a1c451b3b0422d219ff598c72621550c

// Do some useless, but CPU intensive work.
// Some short delays are added throughout so work gets spread out a little bit.
// Otherwise, the only pattern is 0% utilization one second, 100% the next (few) second(s),
// then back to 0%.

package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// Max fib number
	// Keeps it to a 32 bit int
	max := 40

	fmt.Println("Usage: %s [<max fib number>]")
	fmt.Println("Optionally provide the max fib number to calculate, between [1, 40]")
	fmt.Println("This will determine how stressful this will be for the CPU")
	if len(os.Args) > 1 {
		maxI64, err := strconv.ParseInt(os.Args[1], 10, 32)
		if err != nil {
			panic(err)
		}
		max = int(maxI64)
		if max < 1 || max > 40 {
			panic("Error: choose number between [1, 40]")
		}
	} else {
		fmt.Println("Using default max fib value of 40")
	}

	rand.Seed(time.Now().UnixNano())
    for {
		timesInc := rand.Intn(3000) + 500
		fmt.Println("Trying", timesInc, "times increasing fib")
		for num := 0; num <= max; num++ {
			CPU(timesInc, num)
			time.Sleep(time.Millisecond)
		}
		timesDec := rand.Intn(3000) + 500
		fmt.Println("Trying", timesDec, "times decreasing fib")
		for num := max; num >= 0; num-- {
			CPU(timesDec, num)
			time.Sleep(time.Millisecond)
		}
		pauseMs := rand.Intn(2000)
		fmt.Println("Pause for", pauseMs, "ms")
		time.Sleep(time.Millisecond * time.Duration(pauseMs))
    }
}

// CPU does stuff
func CPU(times, num int) []int {
	pauseNs := rand.Intn(1500) + 1500
	var r []int
	for i := 0; i < times; i++ {
		d := bubble(expand(reversePrime(fib(num))))
		if len(d) > 0 {
			r = append(r, d[len(d)-1])
		}
		time.Sleep(time.Nanosecond * time.Duration(pauseNs))
	}
	return r
}

func fib(n int) []int {
	s := []int{1}
	c := 1
	p := 0
	i := n - 1
	for i >= 1 {
		c += p
		p = c - p
		s = append(s, c)
		i--
	}
	return s
}

func prime(n int) bool {
	if n%1 != 0 {
		return false
	} else if n <= 1 {
		return false
	} else if n <= 3 {
		return true
	} else if n%2 == 0 {
		return false
	}
	dl := int(math.Sqrt(float64(n)))
	for d := 3; d <= dl; d += 2 {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func reversePrime(slice []int) []int {
	l := len(slice) - 1
	var r []int
	for l >= 0 {
		if prime(slice[l]) {
			r = append(r, slice[l])
		}
		l--
	}
	return r
}

func expand(slice []int) []int {
	ol := len(slice)
	oc := 0
	l := ol * 10
	var r []int
	for i := 0; i < l; i++ {
		r = append(r, (slice[oc] + (i * 100)))
		if oc < ol-1 {
			oc++
		} else {
			oc = 0
		}
	}
	return r
}

func bubble(slice []int) []int {
	for i := 0; i < len(slice); i++ {
		for y := 0; y < len(slice)-1; y++ {
			if slice[y+1] < slice[y] {
				t := slice[y]
				slice[y] = slice[y+1]
				slice[y+1] = t
			}
		}
	}
	return slice
}
