package models

import (
	"fmt"
	"time"
)

type PullRequestStatus string

const (
	OPEN   PullRequestStatus = "OPEN"
	MERGED PullRequestStatus = "MERGED"
)

type ErrorDetailCode string

const (
	TEAM_EXISTS  ErrorDetailCode = "TEAM_EXISTS"
	PR_EXISTS    ErrorDetailCode = "PR_EXISTS"
	PR_MERGED    ErrorDetailCode = "PR_MERGED"
	NOT_ASSIGNED ErrorDetailCode = "NOT_ASSIGNED"
	NO_CANDIDATE ErrorDetailCode = "NO_CANDIDATE"
	NOT_FOUND    ErrorDetailCode = "NOT_FOUND"
	FATAL_ERROR  ErrorDetailCode = "FATAL_ERROR"
)

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type TeamMember struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type PullRequest struct {
	PullRequestId     string            `json:"pull_request_id"`
	PullRequestName   string            `json:"pull_request_name"`
	AuthorId          string            `json:"author_id"`
	Status            PullRequestStatus `json:"status"`
	AssignedReviewers []string          `json:"assigned_reviewers"`
	CreatedAt         *time.Time        `json:"createdAt,omitempty"`
	MergedAt          *time.Time        `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestId   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	AuthorId        string            `json:"author_id"`
	Status          PullRequestStatus `json:"status"`
}

type ErrorDetail struct {
	Code    ErrorDetailCode `json:"code"`
	Message string          `json:"message"`
}

type ErrorResponse struct {
	Detail ErrorDetail `json:"error"`
}

func NewUser(userId string, username string, teamName string, isActive bool) User {
	return User{
		UserId:   userId,
		Username: username,
		TeamName: teamName,
		IsActive: isActive,
	}
}

func NewTeamMember(user *User) TeamMember {
	return TeamMember{
		UserId:   user.UserId,
		Username: user.Username,
		IsActive: user.IsActive,
	}
}

func NewPR(pullRequestId, pullRequestName, authorId string, status PullRequestStatus, assignedReviewers []string, createdAt *time.Time) PullRequest {
	return PullRequest{
		PullRequestId:     pullRequestId,
		PullRequestName:   pullRequestName,
		AuthorId:          authorId,
		Status:            status,
		AssignedReviewers: assignedReviewers,
		CreatedAt:         createdAt,
	}
}

func NewPRShort(pr *PullRequest) PullRequestShort {
	return PullRequestShort{
		PullRequestId:   pr.PullRequestId,
		PullRequestName: pr.PullRequestName,
		AuthorId:        pr.AuthorId,
		Status:          pr.Status,
	}
}

func NewErrorDetail(code ErrorDetailCode, message string) ErrorDetail {
	return ErrorDetail{
		Code:    code,
		Message: message,
	}
}

func NewErrorResponse(code ErrorDetailCode, message string) ErrorResponse {
	return ErrorResponse{
		Detail: NewErrorDetail(code, message),
	}
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Detail.Code, e.Detail.Message)
}
