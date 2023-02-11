package main

import "fmt"

func main() {
	now := Position{
		Pan:  45.02,
		Tilt: 0.34,
		Zoom: 0.33,
	}

	expect := Position{
		Pan:  45.01,
		Tilt: 0.33,
		Zoom: 0.33,
	}

	if now.Pan >= expect.Pan-2 {
		fmt.Println("now.Pan >= expect.Pan-2 ok")
	}

	if now.Pan <= expect.Pan+2 {
		fmt.Println("now.Pan <= expect.Pan+2 ok")
	}

	if now.Tilt >= expect.Tilt-2 {
		fmt.Println("now.Tilt >= expect.Tilt-2 ok")
	}

	if now.Tilt <= expect.Tilt+2 {
		fmt.Println("now.Tilt <= expect.Tilt+2 ok")
	}

	if now.Zoom >= expect.Zoom-2 {
		fmt.Println("now.Zoom >= expect.Zoom-2 ok")
	}

	if now.Zoom <= expect.Zoom+2 {
		fmt.Println("now.Zoom <= expect.Zoom+2 ok")
	}

	if now.Pan >= expect.Pan-2 && now.Pan <= expect.Pan+2 &&
		now.Tilt >= expect.Tilt-2 && now.Tilt <= expect.Tilt+2 &&
		now.Zoom >= expect.Zoom-2 && now.Zoom <= expect.Zoom+2 {

		fmt.Println("ok")
		return
	}

	fmt.Println("failed")
}

type Position struct {
	Pan  float64
	Tilt float64
	Zoom float64
}
