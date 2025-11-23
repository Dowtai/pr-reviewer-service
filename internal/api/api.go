package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Dowtai/pr-reviewer-service/internal/models"
	"github.com/Dowtai/pr-reviewer-service/internal/service"
)

func TeamAddHandler(svc *service.PrReviewerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var team models.Team
		if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdTeam, err := svc.TeamAdd(team)
		if err != nil {
			var svcErr service.ErrorService
			if errors.As(err, &svcErr) {
				switch svcErr.Code {
				case service.INTERNAL_ERROR:
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, svcErr.Error()))
				case service.OBJECT_EXISTS:
					w.WriteHeader(http.StatusBadRequest)
					json.NewEncoder(w).Encode(models.NewErrorResponse(svcErr.ApiCode, "team_name already exists"))
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, err.Error()))
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdTeam)
	}
}

func TeamGetHandler(svc *service.PrReviewerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		teamName := r.URL.Query().Get("team_name")
		if teamName == "" {
			http.Error(w, "wrong TeamNameQuery", http.StatusBadRequest)
			return
		}

		team, err := svc.TeamGet(teamName)
		if err != nil {
			var svcErr service.ErrorService
			if errors.As(err, &svcErr) {
				switch svcErr.Code {
				case service.INTERNAL_ERROR:
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, svcErr.Error()))
				case service.OBJECT_NOT_FOUND:
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(models.NewErrorResponse(svcErr.ApiCode, "team_name not found"))
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, err.Error()))
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(team)
	}
}

func UsersSetIsActiveHandler(svc *service.PrReviewerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			UserId   string `json:"user_id"`
			IsActive bool   `json:"is_active"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		updatedUser, err := svc.UsersSetIsActive(request.UserId, request.IsActive)
		if err != nil {
			var svcErr service.ErrorService
			if errors.As(err, &svcErr) {
				switch svcErr.Code {
				case service.INTERNAL_ERROR:
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, svcErr.Error()))
				case service.OBJECT_NOT_FOUND:
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(models.NewErrorResponse(svcErr.ApiCode, "user_id not found"))
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, err.Error()))
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedUser)
	}
}

func PullRequestCreateHandler(svc *service.PrReviewerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			PullRequestId   string `json:"pull_request_id"`
			PullRequestName string `json:"pull_request_name"`
			AuthorId        string `json:"author_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdPullRequest, err := svc.PullRequestCreate(request.PullRequestId, request.PullRequestName, request.AuthorId)
		if err != nil {
			var svcErr service.ErrorService
			if errors.As(err, &svcErr) {
				switch svcErr.Code {
				case service.INTERNAL_ERROR:
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, svcErr.Error()))
				case service.OBJECT_NOT_FOUND:
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(models.NewErrorResponse(svcErr.ApiCode, "author or team not found"))
				case service.DOMAIN_ERROR:
					w.WriteHeader(http.StatusConflict)
					json.NewEncoder(w).Encode(models.NewErrorResponse(svcErr.ApiCode, "pull_request already exists"))
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, err.Error()))
			}
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdPullRequest)
	}
}

func PullRequestMergeHandler(svc *service.PrReviewerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			PullRequestId string `json:"pull_request_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mergedPullRequest, err := svc.PullRequestMerge(request.PullRequestId)
		if err != nil {
			var svcErr service.ErrorService
			if errors.As(err, &svcErr) {
				switch svcErr.Code {
				case service.INTERNAL_ERROR:
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, svcErr.Error()))
				case service.OBJECT_NOT_FOUND:
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(models.NewErrorResponse(svcErr.ApiCode, "pull_request not found"))
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, err.Error()))
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mergedPullRequest)
	}
}

func PullRequestReassignHandler(svc *service.PrReviewerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request struct {
			PullRequestId string `json:"pull_request_id"`
			OldUserId     string `json:"old_user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		reassignedPullRequest, replacedBy, err := svc.PullRequestReassign(request.PullRequestId, request.OldUserId)
		if err != nil {
			var svcErr service.ErrorService
			if errors.As(err, &svcErr) {
				switch svcErr.Code {
				case service.INTERNAL_ERROR:
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, svcErr.Error()))
				case service.OBJECT_NOT_FOUND:
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(models.NewErrorResponse(svcErr.ApiCode, "pull_request or user not found"))
				case service.DOMAIN_ERROR:
					w.WriteHeader(http.StatusConflict)
					json.NewEncoder(w).Encode(models.NewErrorResponse(svcErr.ApiCode, svcErr.Error()))
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, err.Error()))
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			PR         models.PullRequest `json:"pr"`
			ReplacedBy string             `json:"replaced_by"`
		}{
			PR:         reassignedPullRequest,
			ReplacedBy: replacedBy,
		})
	}
}

func UsersGetReviewHandler(svc *service.PrReviewerService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		prs, err := svc.UsersGetReview(r.URL.Query().Get("user_id"))
		if err != nil {
			var svcErr service.ErrorService
			if errors.As(err, &svcErr) {
				switch svcErr.Code {
				case service.INTERNAL_ERROR:
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, svcErr.Error()))
				case service.OBJECT_NOT_FOUND:
					w.WriteHeader(http.StatusNotFound)
					json.NewEncoder(w).Encode(models.NewErrorResponse(svcErr.ApiCode, "user not found"))
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.NewErrorResponse(models.FATAL_ERROR, err.Error()))
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(prs)
	}
}
