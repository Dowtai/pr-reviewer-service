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
	prsByUser map[string]map[string]struct{}
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		teams:     make(map[string]models.Team),
		users:     make(map[string]models.User),
		prs:       make(map[string]models.PullRequest),
		prsByUser: make(map[string]map[string]struct{}),
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
		i := 0
		for pr, _ := range prs {
			prModels[i] = r.GetPullRequestById(pr)
			if prModels[i] == nil {
				return nil
			}
			i++
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

func (r *MemoryRepo) UpdateTeamMember(team *models.Team, user *models.User) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.users[user.UserId]; !ok {
		return errors.New("updating non-existing user")
	}

	if _, ok := r.teams[team.TeamName]; !ok {
		return errors.New("updating non-existing team")
	}

	for i, member := range r.teams[team.TeamName].Members {
		if member.UserId == user.UserId {
			r.users[user.UserId] = *user
			r.teams[team.TeamName].Members[i] = models.NewTeamMember(user)
			return nil
		}
	}
	return errors.New("user is not member of team")
}

func (r *MemoryRepo) UpdatePR(pr *models.PullRequest) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.prs[pr.PullRequestId]; !ok {
		return errors.New("updating non-existing pull_request")
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

	if _, ok := r.teams[team.TeamName]; ok {
		return errors.New("creating already existing team")
	}

	r.teams[team.TeamName] = team
	return nil
}

func (r *MemoryRepo) CreatePR(pr models.PullRequest) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.prs[pr.PullRequestId]; ok {
		return errors.New("creating already existing pr")
	}

	r.prs[pr.PullRequestId] = pr
	return nil
}

func (r *MemoryRepo) AddPRToUser(userId string, prId string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.users[userId]; !ok {
		return errors.New("adding pr to non-existing user")
	}
	if _, ok := r.prs[prId]; !ok {
		return errors.New("adding non-existing pr to user")
	}

	if r.prsByUser[userId] == nil {
		r.prsByUser[userId] = make(map[string]struct{})
	}

	r.prsByUser[userId][prId] = struct{}{}
	return nil
}

func (r *MemoryRepo) RemovePRFromUser(userId string, prId string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.users[userId]; !ok {
		return errors.New("removing pr from non-existing user")
	}
	if _, ok := r.prs[prId]; !ok {
		return errors.New("removing non-existing pr from user")
	}

	if r.prsByUser[userId] == nil {
		return errors.New("user is not reviewer of this pr")
	}

	delete(r.prsByUser[userId], prId)
	return nil
}
