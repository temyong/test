package main

import (
 //"encoding/json"
 //"errors"
 "fmt"
 "log"
 //"net/http"
 "reddit"
)

func main() {
	items, err := reddit.Get("golang")
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range items {
		fmt.Println(item)
	}
}
