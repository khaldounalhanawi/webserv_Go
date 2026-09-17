package main

import (
	"os"
)

func main() {
	i,err := os.Stat("./training")
	if err != nil {
		println(err.Error())
		return
	}

	println(i.Name(), "//Is Directory:", i.IsDir())
}