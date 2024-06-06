//go:build !solution

package main

import (
	"fmt"
	"time"
	"os"
	"io"
	"net/http"
)

func get_url(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	defer resp.Body.Close()
	nbytes, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	ch <- fmt.Sprintf("%.2fs %7d %s", time.Since(start).Seconds(), nbytes, url)

}
func main() {
	start := time.Now()
	args := os.Args[1:]
	ch := make(chan string)
	for _, url := range args {
		go get_url(url, ch)
	}
	for range args {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
