package repo

import "github.com/Dowtai/pr-reviewer-service/internal/models"

type Repo interface {
	TeamExists(teamName string) bool
	GetTeamByName(teamName string) *models.Team
	GetTeamByUserId(userId string) *models.Team
	GetUserById(userId string) *models.User
	GetPullRequestById(prId string) *models.PullRequest
	GetPullRequestsByUserId(userId string) []*models.PullRequest
	UpdateUser(user models.User) error
	UpdatePR(user models.User) error
	CreateUser(user models.User) error
	CreateTeam(team models.Team) error
	CreatePR(pr models.PullRequest) error
	SetIsActiveUser(userId string, isActive bool) *models.User
}
