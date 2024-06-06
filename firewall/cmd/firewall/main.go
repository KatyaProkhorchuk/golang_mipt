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

func GenerateImg(time string, size int) *image.RGBA {
	digitWidth := 8
	digitHeight := 12
	clnWidth := 4

	imgWidth := (6*digitWidth + 2*clnWidth) * size
	imgHeight := digitHeight * size
	upLeft := image.Point{X: 0, Y: 0}
	lowRight := image.Point{X: imgWidth, Y: imgHeight}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	offset := 0
	digit := ""

	for _, char := range time {
		switch string(char) {
		case "0":
			digit = Zero
		case "1":
			digit = One
		case "2":
			digit = Two
		case "3":
			digit = Three
		case "4":
			digit = Four
		case "5":
			digit = Five
		case "6":
			digit = Six
		case "7":
			digit = Seven
		case "8":
			digit = Eight
		case "9":
			digit = Nine
		default:
			digit = Colon
		}

		drawDigit(img, digit, size, offset)

		if digit == Colon {
			offset += clnWidth
		} else {
			offset += digitWidth
		}
	}

	return img
}

func drawDigit(img *image.RGBA, digit string, k int, offset int) {
	lines := strings.Split(digit, "\n")

	for h := 0; h < len(lines); h++ {
		for w, char := range lines[h] {
			fillCell(img, w+offset, h, k, string(char))
		}
	}
}

func fillCell(img *image.RGBA, x int, y, size int, char string) {
	for i := x * size; i < (x+1)*size; i++ {
		for j := y * size; j < (y+1)*size; j++ {
			if char != "." {
				img.Set(i, j, Cyan)
			} else {
				img.Set(i, j, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			}
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
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		if size < 1 || size > 30 {
			http.Error(w, "invalid k", http.StatusBadRequest)
			return
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
		fmt.Println(size, timeP)
		img := GenerateImg(timeP, size)
		w.Header().Set("Content-Type", "image/png")
		if err = png.Encode(w, img); err != nil {
			http.Error(w, "image error", http.StatusBadRequest)
		}

	}
}

func main() {
	port := flag.String("port", "80", "port http server")
	flag.Parse()
	http.HandleFunc("/", getPng)
	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}
