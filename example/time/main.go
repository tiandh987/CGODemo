package main

import (
	"fmt"
	"time"
)

func main() {
	timeStr := "15:08:10"
	timeStr2 := "15:08:09"

	start, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-%s %s",
		time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"), timeStr))
	if err != nil {
		panic(err)
	}

	end, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s-%s-%s %s",
		time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"), timeStr2))
	if err != nil {
		panic(err)
	}

	now := time.Now()

	if start.Before(end) {
		fmt.Printf("start before end\n\n")
	}

	if end.After(start) {
		fmt.Printf("end after before\n\n")
	}

	if end.Equal(start) {
		fmt.Printf("end equal before\n\n")
	}

	fmt.Printf("now: %d, \nstart: %d, \nend: %d\n\n",
		now.UnixMilli(), start.UnixMilli(), end.UnixMilli())

	fmt.Printf("now: %s, \nstart: %s, \nend: %s\n\n",
		now.Format("2006-01-02 15:04:05"), start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))

	parse, err := time.Parse("2006-01-02 15:04:05", "2023-02-24 10:30:55")
	if err != nil {
		panic(err)
	}

	parse2, err2 := time.Parse("2006-01-02 15:04:05", "2023-02-24 10:31:03")
	if err != nil {
		panic(err2)
	}

	sub := parse2.Sub(parse)

	fmt.Printf("sub: %d\n", sub)

	fmt.Printf("now: %s\n", time.Now().String())
	time.Sleep(sub)
	fmt.Printf("now: %s\n", time.Now().String())

}
