//go:build !solution

package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type requestBody struct {
	URL string
}

type responseBody struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

var (
	store map[string]string
	mu    sync.Mutex
)

func genSha(key string) (string, error) {
	algo := sha1.New()
	_, err := algo.Write([]byte(key))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(algo.Sum(nil)), err
}

func createShortHandler(w http.ResponseWriter, r *http.Request) {
	var req requestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err == nil {
		shortLink, err := genSha(req.URL)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		mu.Lock()
		store[shortLink] = req.URL
		mu.Unlock()
		data := responseBody{URL: req.URL, Key: shortLink}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "invalid", http.StatusBadRequest)
	}
}

func shortRedirect(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	size := len(path)
	key := path[size-1]
	mu.Lock()
	url, err := store[key]
	mu.Unlock()
	if err {
		w.Header().Add("Location", url)
		w.WriteHeader(http.StatusFound)
	} else {
		http.Error(w, "not key", http.StatusNotFound)
	}
}

func main() {
	port := flag.String("port", "80", "port http server")
	flag.Parse()
	store = make(map[string]string)
	http.HandleFunc("/shorten", createShortHandler)
	http.HandleFunc("/go/", shortRedirect)
	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}
