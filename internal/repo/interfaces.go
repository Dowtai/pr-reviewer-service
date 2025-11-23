package repo

import "github.com/Dowtai/pr-reviewer-service/internal/models"

type Repo interface {
	TeamExists(teamName string) bool
	GetTeamByName(teamName string) *models.Team
	GetUserById(userId string) *models.User
	GetPullRequestById(prId string) *models.PullRequest
	GetPullRequestsByUserId(userId string) []*models.PullRequest
	UpdateUser(user *models.User) error
	UpdateTeamMember(team *models.Team, user *models.User) error
	UpdatePR(pr *models.PullRequest) error
	CreateUser(user models.User) error
	CreateTeam(team models.Team) error
	CreatePR(pr models.PullRequest) error
	AddPRToUser(userId string, prId string) error
	RemovePRFromUser(userId string, prId string) error
}
