package main

import (
	"fmt"
	"gontainer/cmd/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
