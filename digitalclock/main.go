//go:build !solution

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func getDigit(charStr string) string {
	digit := ""
	if charStr == "0" {
		digit = Zero
	} else if charStr == "1" {
		digit = One
	} else if charStr == "2" {
		digit = Two
	} else if charStr == "3" {
		digit = Three
	} else if charStr == "4" {
		digit = Four
	} else if charStr == "5" {
		digit = Five
	} else if charStr == "6" {
		digit = Six
	} else if charStr == "7" {
		digit = Seven
	} else if charStr == "8" {
		digit = Eight
	} else if charStr == "9" {
		digit = Nine
	} else {
		digit = Colon
	}
	return digit
}
func GenerateImg(size int, time string) *image.RGBA {
	width := 8
	height := 12
	clnWidth := 4

	imgWidth := (6*width + 2*clnWidth) * size
	imgHeight := height * size
	lup := image.Point{X: 0, Y: 0}
	rlow := image.Point{X: imgWidth, Y: imgHeight}

	img := image.NewRGBA(image.Rectangle{Min: lup, Max: rlow})

	offset := 0

	for _, char := range time {
		charStr := string(char)
		digit := getDigit(charStr)
		lines := strings.Split(digit, "\n")

		for i := 0; i < len(lines); i++ {
			for w, char := range lines[i] {
				j := w + offset
				charStr := string(char)
				draw(img, j, i, size, charStr)
			}
		}
		if digit != Colon {
			offset += width
		} else {
			offset += clnWidth
		}
	}

	return img
}
func draw(img *image.RGBA, x, y, size int, char string) {
	startX := x * size
	startY := y * size
	endX := (x + 1) * size
	endY := (y + 1) * size

	var targetColor color.RGBA
	if char == "." {
		targetColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	} else {
		targetColor = Cyan
	}

	for i := startX; i < endX; i++ {
		for j := startY; j < endY; j++ {
			img.Set(i, j, targetColor)
		}
	}
}
func getPng(w http.ResponseWriter, r *http.Request) {
	var timeP string
	var size = 1
	var err error
	queries := r.URL.Query()
	if k, ok := queries["k"]; ok && len(k) > 0 {
		size, err = strconv.Atoi(k[0])
		if err != nil || size < 1 || size > 30 {
			http.Error(w, "invalid k", http.StatusBadRequest)
			return
		}
	}
	if tm, ok := queries["time"]; ok {
		timeP = tm[0]
		if ok, err = regexp.MatchString("^(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$", timeP); err != nil || !ok {
			http.Error(w, "incorrect time", http.StatusBadRequest)
			return
		}
	} else {
		tm := time.Now()
		timeP = fmt.Sprintf("%02d:%02d:%02d", tm.Hour(), tm.Minute(), tm.Second())
	}
	img := GenerateImg(size, timeP)
	w.Header().Set("Content-Type", "image/png")
	if err = png.Encode(w, img); err != nil {
		http.Error(w, "image error", http.StatusBadRequest)
	}
}

func main() {
	port := flag.String("port", "80", "port http server")
	flag.Parse()
	http.HandleFunc("/", getPng)
	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}
