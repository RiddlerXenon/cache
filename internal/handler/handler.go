package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/RiddlerXenon/cache/internal/cache"
	"go.uber.org/zap"
)

func parseKey(u *url.URL) (string, error) {
	url, err := u.Parse("/")
	if err != nil {
		return "", err
	}

	params := url.Query()
	key := params.Get("key")

	return key, nil
}

func SetHandler(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			key, err := parseKey(r.URL)
			if err != nil {
				http.Error(w, "Unable to parse URL", http.StatusBadRequest)
				zap.S().Error(err)
				return
			}
			if key == "" {
				http.Error(w, "Key is empty", http.StatusBadRequest)
				zap.S().Errorf("Key is empty")
				return
			}

			decoder := json.NewDecoder(r.Body)
			var request Request
			err = decoder.Decode(&request)

			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				zap.S().Error(err)
				return
			}

			if request.Value != "" {
				if request.TTL > 0 {
					c.Set(key, request.Value, request.TTL)
				} else {
					c.Set(key, request.Value, 0)
				}
			}

			http.Error(w, "Value is empty", http.StatusBadRequest)
			zap.S().Errorf("Value is empty")
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			zap.S().Errorf("Method not allowed")
		}
	}
}

func GetHandler(c *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			key, err := parseKey(r.URL)
			if err != nil {
				http.Error(w, "Unable to parse URL", http.StatusBadRequest)
				zap.S().Error(err)
				return
			}
			if key == "" {
				http.Error(w, "Key is empty", http.StatusBadRequest)
				zap.S().Errorf("Key is empty")
				return
			}

			value, err := c.Get(key)

			if err != nil {
				http.Error(w, "Key not found", http.StatusBadRequest)
				zap.S().Errorf("Key not found")
				return
			}

			response := Response{
				Value: value,
			}
			json.NewEncoder(w).Encode(response)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			zap.S().Errorf("Method not allowed")
		}
	}
}

// func GetCacheHandler(*cache.Cache) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case htt.MethodGet:
// 			return
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			zap.S().Errorf("Method not allowed")
// 		}
// 	}
// }
