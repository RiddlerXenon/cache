package main

import (
	"net/http"

	"github.com/RiddlerXenon/cache/internal/cache"
	"github.com/RiddlerXenon/cache/internal/handler"
	"go.uber.org/zap"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	// repo, err := repository.New()
	// if err != nil {
	// 	zap.S().Error(err)
	// }
	// zap.S().Info("The database connected successful")

	// defer func() {
	// 	if err := repo.Close(); err != nil {
	// 		zap.S().Error(err)
	// 	}
	// 	zap.S().Info("The database closed successful")
	// }()
	c, err := cache.New()
	if err != nil {
		zap.S().Fatal(err)
	}

	http.HandleFunc("/api/get", handler.GetHandler(c))
	http.HandleFunc("/api/set", handler.SetHandler(c))
	// http.HandleFunc("/api/getCache", handler.GetCacheHandler(c))

	zap.S().Info("Server starting at http://127.0.0.1:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		zap.S().Fatal(err)
	}
}
