package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

type CachedResponse struct {
	StatusCode int
	Body       []byte
	Completed  bool
}

type Store struct {
	mu   sync.Mutex
	data map[string]*CachedResponse
}

func NewStore() *Store {
	return &Store{data: make(map[string]*CachedResponse)}
}

func (s *Store) Get(key string) (*CachedResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *Store) Start(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; exists {
		return false
	}

	s.data[key] = &CachedResponse{Completed: false}
	return true
}

func (s *Store) Finish(key string, code int, body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = &CachedResponse{
		StatusCode: code,
		Body:       body,
		Completed:  true,
	}
}

func IdempotencyMiddleware(store *Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			http.Error(w, "Missing Idempotency-Key", 400)
			return
		}

		if val, ok := store.Get(key); ok {
			if val.Completed {
				w.WriteHeader(val.StatusCode)
				w.Write(val.Body)
			} else {
				http.Error(w, "Processing", 409)
			}
			return
		}

		if !store.Start(key) {
			http.Error(w, "Conflict", 409)
			return
		}

		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		store.Finish(key, rec.Code, rec.Body.Bytes())

		w.WriteHeader(rec.Code)
		w.Write(rec.Body.Bytes())
	})
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing started...")
	time.Sleep(2 * time.Second)

	resp := map[string]interface{}{
		"status":         "paid",
		"amount":         1000,
		"transaction_id": "uuid-12345",
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	store := NewStore()

	handler := IdempotencyMiddleware(store, http.HandlerFunc(paymentHandler))

	server := httptest.NewServer(handler)
	defer server.Close()

	client := &http.Client{}
	key := "abc-123"

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			req, _ := http.NewRequest("GET", server.URL, nil)
			req.Header.Set("Idempotency-Key", key)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("error:", err)
				return
			}

			fmt.Println("Request", i, "Status:", resp.StatusCode)
		}(i)
	}

	wg.Wait()
}