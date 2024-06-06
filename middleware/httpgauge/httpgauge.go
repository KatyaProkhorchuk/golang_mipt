package httpgauge

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Gauge struct {
	metrics map[string]int
	mutex   sync.Mutex
}

func New() *Gauge {
	return &Gauge{
		metrics: make(map[string]int),
	}
}

func (g *Gauge) Snapshot() map[string]int {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.metrics
}

func (g *Gauge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	metrics := g.Snapshot()
	jsonData := ""
	var keys []string
	for key := range metrics {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		jsonData += fmt.Sprintf("%s %d\n", key, metrics[key])
	}

	jsonData = strings.TrimSuffix(jsonData, "/")

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(jsonData))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				// Паника произошла, обрабатываем ее
				// Например, можем залогировать информацию о панике
				g.mutex.Lock()
				defer g.mutex.Unlock()
				g.metrics["/panic"]++
				defer panic("")
			}
		}()
		if r != nil {
			next.ServeHTTP(w, r)
			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			g.mutex.Lock()
			defer g.mutex.Unlock()
			g.metrics[routePattern]++
		}
	})
}
