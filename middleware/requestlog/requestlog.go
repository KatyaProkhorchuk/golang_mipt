//go:build !solution

package requestlog

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	// "github.com/google/uuid"

	"go.uber.org/zap"
)

func generateRequestID() string {
	// Генерируем случайный байтовый массив длиной 16 байт
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		// В случае ошибки генерации используем текущее время как запасной вариант
		return time.Now().Format(time.RFC3339Nano)
	}
	// Преобразуем байтовый массив в строку в шестнадцатеричном формате
	return hex.EncodeToString(randomBytes)
}
func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := generateRequestID()

			startTime := time.Now()

			l.Info("request started",
				zap.String("request_id", requestID),
				zap.String("path", r.URL.Path),
				zap.String("method", r.Method),
			)

			rw := &responseWriterWithStatus{ResponseWriter: w}

			defer func() {
				duration := time.Since(startTime)

				if rw.status != 0 {
					l.Info("request finished",
						zap.String("request_id", requestID),
						zap.String("path", r.URL.Path),
						zap.String("method", r.Method),
						zap.Int("status_code", rw.status),
						zap.Duration("duration", duration),
					)
				} else {
					path := r.URL.Path
					method := r.Method
					if r := recover(); r != nil {
						l.Info("request panicked",
							zap.String("request_id", requestID),
							zap.String("path", path),
							zap.String("method", method),
						)
						panic(r) // Пробрасываем панику дальше
					} else {
						// Если не было вызова WriteHeader, считаем запрос успешно завершенным
						l.Info("request finished",
							zap.String("request_id", requestID),
							zap.String("path", path),
							zap.String("method", method),
							zap.Int("status_code", http.StatusOK), // Предполагаем успешный статус код
							zap.Duration("duration", duration),
						)
					}
				}
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

type responseWriterWithStatus struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriterWithStatus) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
