package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Zuf4ik/Z_JobScout/internal/httpapi"
	"github.com/Zuf4ik/Z_JobScout/internal/service"
)

func main() {
	searchService := service.NewSearchService()

	handler := httpapi.NewHandler(searchService)

	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("JobScout API запущен: http://127.0.0.1:8080")

	if err := server.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatal("ошибка HTTP-сервера:", err)
	}
}
