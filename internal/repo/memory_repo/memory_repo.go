package memory_repo

import (
	"errors"
	"sync"

	"github.com/Dowtai/pr-reviewer-service/internal/models"
)

type MemoryRepo struct {
	mx        sync.RWMutex
	teams     map[string]models.Team
	users     map[string]models.User
	prs       map[string]models.PullRequest
	prsByUser map[string][]string
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		teams:     make(map[string]models.Team),
		users:     make(map[string]models.User),
		prs:       make(map[string]models.PullRequest),
		prsByUser: make(map[string][]string),
	}
}

func (r *MemoryRepo) TeamExists(teamName string) bool {
	r.mx.RLock()
	defer r.mx.RUnlock()
	
	_, ok := r.teams[teamName]
	return ok
}

func (r *MemoryRepo) GetTeamByName(teamName string) *models.Team {
	r.mx.RLock()
	defer r.mx.RUnlock()

	if team, ok := r.teams[teamName]; ok {
		return &team
	}
	return nil
}

func (r *MemoryRepo) GetUserById(userId string) *models.User {
	r.mx.RLock()
	defer r.mx.RUnlock()

	if user, ok := r.users[userId]; ok {
		return &user
	}
	return nil
}

func (r *MemoryRepo) GetPullRequestById(prId string) *models.PullRequest {
	r.mx.RLock()
	defer r.mx.RUnlock()

	if pr, ok := r.prs[prId]; ok {
		return &pr
	}
	return nil
}

func (r *MemoryRepo) GetPullRequestsByUserId(userId string) []*models.PullRequest {
	r.mx.RLock()
	defer r.mx.RUnlock()

	if prs, ok := r.prsByUser[userId]; ok {
		prModels := make([]*models.PullRequest, len(prs))
		for i, pr := range prs {
			prModels[i] = r.GetPullRequestById(pr)
			if prModels[i] == nil {
				return nil
			}
		}
		return prModels
	}
	return nil
}

func (r *MemoryRepo) UpdateUser(user *models.User) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.users[user.UserId]; !ok {
		return errors.New("updating non-existing user")
	}

	r.users[user.UserId] = *user
	return nil
}

func (r *MemoryRepo) UpdatePR(pr *models.PullRequest) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.users[pr.PullRequestId]; !ok {
		return errors.New("updating non-existing user")
	}

	r.prs[pr.PullRequestId] = *pr
	return nil
}

func (r *MemoryRepo) CreateUser(user models.User) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.users[user.UserId]; ok {
		return errors.New("creating already existing user")
	}

	r.users[user.UserId] = user
	return nil
}

func (r *MemoryRepo) CreateTeam(team models.Team) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.users[team.TeamName]; ok {
		return errors.New("creating already existing team")
	}

	r.teams[team.TeamName] = team
	return nil
}

func (r *MemoryRepo) CreatePR(pr models.PullRequest) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.users[pr.PullRequestId]; ok {
		return errors.New("creating already existing pr")
	}

	r.prs[pr.PullRequestId] = pr
	return nil
}

func (r *MemoryRepo) AddPRToUser(userId string, prId string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.users[userId]; ok {
		return errors.New("adding pr to non-existing user")
	}
	if _, ok := r.prs[prId]; ok {
		return errors.New("adding non-existing pr to user")
	}

	r.prsByUser[prId] = append(r.prsByUser[prId], prId)
	return nil
}
