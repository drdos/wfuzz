package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sync"
)

func main() {
	var site string
	if len(os.Args) > 1 {
		site = os.Args[1]
	} else {
		fmt.Println("Need the host. Example > http://www.google.com/")
		fmt.Println("wfuzz http://www.google.com/ wordlist.txt")
		site = "http://www.google.com"
		os.Exit(1)
	}
	f, err := os.Open(os.Args[2])

	fmt.Println(".+.")
	fmt.Println("Starting fuzz in ", site)
	fmt.Println(".'.'")

	wg := new(sync.WaitGroup)
	resultChan := make(chan string)

	if err != nil {
		fmt.Println("error opening file!", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		wg.Add(1)
		a := site + scanner.Text()
		go get_page(a, resultChan, wg)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("error", err)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for result := range resultChan {
		//deal with the result in some way
		fmt.Println(result)
	}
}

func get_page(site string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Get(site)
	if err != nil {
		ch <- fmt.Sprintf("error: %s", err)
		return
	}
	if res.StatusCode != 404 {
		ch <- fmt.Sprintf("Page Found!! +:: %s ::+ %d", site, res.StatusCode)
	}
}
