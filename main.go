package main

import (
	"fmt"
	"github.com/stachujone5/bencode"
)

func main() {
	fmt.Println("Hello, world!")
	data, err := bencode.Encode("Hello, world!")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(data))
}
