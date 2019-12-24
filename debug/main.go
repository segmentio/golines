package main

import "log"

func main() {
	err := genNodeToGraphNode()
	if err != nil {
		log.Fatal(err)
	}
}
