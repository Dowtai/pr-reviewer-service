package main

import (
	"log"
	"net/http"

	"github.com/Dowtai/pr-reviewer-service/internal/api"
	"github.com/Dowtai/pr-reviewer-service/internal/repo/memory_repo"
	"github.com/Dowtai/pr-reviewer-service/internal/service"
)

func main() {
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
	log.Println("Server started on port", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
