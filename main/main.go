package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Dowtai/pr-reviewer-service/internal/api"
	"github.com/Dowtai/pr-reviewer-service/internal/repo/memory_repo"
	"github.com/Dowtai/pr-reviewer-service/internal/service"
)

func NewServer() *http.Server {
	repo := memory_repo.NewMemoryRepo()

	svc := service.NewService(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/team/add", api.TeamAddHandler(svc))
	mux.HandleFunc("/team/get", api.TeamGetHandler(svc))
	mux.HandleFunc("/users/setIsActive", api.UsersSetIsActiveHandler(svc))
	mux.HandleFunc("/pullRequest/create", api.PullRequestCreateHandler(svc))
	mux.HandleFunc("/pullRequest/merge", api.PullRequestMergeHandler(svc))
	mux.HandleFunc("/pullRequest/reassign", api.PullRequestReassignHandler(svc))
	mux.HandleFunc("/users/getReview", api.UsersGetReviewHandler(svc))

	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}
	return server
}

func ShutdownServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server: %v", err)
	}
}

func main() {
	server := NewServer()
	log.Println("Server started on port", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
