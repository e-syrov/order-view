package main

import (
	"fmt"
	"wb/internal/app"
)

func main() {
	err := app.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
