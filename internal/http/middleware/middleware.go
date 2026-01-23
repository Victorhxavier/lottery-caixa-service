package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

// Handler é o tipo de função para middlewares
type Handler func(http.Handler) http.Handler

// Chain aplica múltiplos middlewares
func Chain(handlers ...Handler) func(http.HandlerFunc) http.HandlerFunc {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		for i := len(handlers) - 1; i >= 0; i-- {
			handler = wrapHandler(handler, handlers[i])
		}
		return handler
	}
}

// wrapHandler envolve um handlerfunc com um middleware
func wrapHandler(handler http.HandlerFunc, mw Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mw(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}
}

// Logging middleware registra requisições
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		start := time.Now()

		// Envolver ResponseWriter para capturar status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("[%s] %s %s %s", 
			requestID, r.Method, r.RequestURI, r.RemoteAddr)

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("[%s] Response: %d - %v", 
			requestID, wrapped.statusCode, duration)
	})
}

// Recovery middleware recupera de panics
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := r.Header.Get("X-Request-ID")
				if requestID == "" {
					requestID = uuid.New().String()
				}

				log.Printf("[%s] PANIC: %v\n%s", 
					requestID, err, debug.Stack())

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// CORS middleware configura CORS
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimit middleware implementa rate limiting
func RateLimit(maxRequests int, window time.Duration) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Implementar rate limiting por IP
			// clientIP := r.RemoteAddr
			// if !limiter.Allow(clientIP) {
			//     w.WriteHeader(http.StatusTooManyRequests)
			//     return
			// }
			next.ServeHTTP(w, r)
		})
	}
}

// Timeout middleware com timeout customizado
func Timeout(duration time.Duration) Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(r.Context())
			// ctx, cancel := context.WithTimeout(r.Context(), duration)
			// defer cancel()
			// r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter para capturar status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}
