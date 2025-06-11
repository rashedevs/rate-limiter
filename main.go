package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type RateLimitRequest struct {
	Key    string `json:"key"`
	Limit  int    `json:"limit"`
	Window int    `json:"window"`
}

type RateLimitResponse struct {
	Allowed        bool `json:"allowed"`
	Remaining      int  `json:"remaining"`
	ResetInSeconds int  `json:"resetInSeconds"`
}

var limiter *RateLimiter

func main() {
	limiter = NewRateLimiter()

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RateLimitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		allowed, remaining, reset := limiter.Check(req.Key, req.Limit, req.Window)

		resp := RateLimitResponse{
			Allowed:        allowed,
			Remaining:      remaining,
			ResetInSeconds: reset,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Server started on :8058")
	log.Fatal(http.ListenAndServe(":8058", nil))
}
