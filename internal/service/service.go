package service

import (
	"time"

	"github.com/Dowtai/pr-reviewer-service/internal/models"
	"github.com/Dowtai/pr-reviewer-service/internal/repo"
)

const (
	OBJECT_EXISTS    int = 400
	OBJECT_NOT_FOUND int = 404
	DOMAIN_ERROR     int = 409
	INTERNAL_ERROR   int = 500
)

type ErrorService struct {
	Code    int
	Message string
}

func (e ErrorService) Error() string {
	return e.Message
}

func NewErrorService(code int, message string) ErrorService {
	return ErrorService{
		Code:    code,
		Message: message,
	}
}

type PrReviewerService struct {
	repo repo.Repo
}

func (s *PrReviewerService) TeamAdd(team models.Team) (models.Team, error) {
	if s.repo.TeamExists(team.TeamName) {
		return team, NewErrorService(OBJECT_EXISTS, "Team already exists")
	}

	var err error
	for _, member := range team.Members {
		if user := s.repo.GetUserById(member.UserId); user != nil {
			user.Username = member.Username
			user.TeamName = team.TeamName
			user.IsActive = member.IsActive

			err = s.repo.UpdateUser(*user)
		} else {
			newUser := models.NewUser(member.UserId, member.Username, team.TeamName, member.IsActive)
			err = s.repo.CreateUser(newUser)

		}

		if err != nil {
			return team, NewErrorService(INTERNAL_ERROR, err.Error())
		}
	}

	err = s.repo.CreateTeam(team)
	if err != nil {
		return team, NewErrorService(INTERNAL_ERROR, err.Error())
	}

	return team, nil
}

func (s *PrReviewerService) TeamGet(team models.Team) (models.Team, error) {
	if t := s.repo.GetTeamByName(team.TeamName); t != nil {
		return *t, nil
	}

	return team, NewErrorService(OBJECT_NOT_FOUND, "Team not found")
}

func (s *PrReviewerService) UsersSetIsActive(userId string, isActive bool) (models.User, error) {
	if user := s.repo.SetIsActiveUser(userId, isActive); user != nil {
		return *user, nil
	}
	return models.User{}, NewErrorService(OBJECT_NOT_FOUND, "User not found")
}

func (s *PrReviewerService) PullRequestCreate(pullRequestId, pullRequestName, authorId string) (models.PullRequest, error) {
	if pr := s.repo.GetPullRequestById(pullRequestId); pr != nil {
		return *pr, NewErrorService(DOMAIN_ERROR, "Pull request already exists")
	}

	team := s.repo.GetTeamByUserId(authorId)
	if team == nil {
		return models.PullRequest{}, NewErrorService(OBJECT_NOT_FOUND, "Team not found")
	}

	reviewers := make([]string, 0, 2)

	for _, member := range team.Members {
		if member.UserId == authorId {
			continue
		}
		if member.IsActive {
			reviewers = append(reviewers, member.Username)
		}
		if len(reviewers) == 2 {
			break
		}
	}

	now := time.Now()
	pr := models.NewPR(pullRequestId, pullRequestName, authorId, models.OPEN, reviewers, &now)

	if err := s.repo.CreatePR(pr); err != nil {
		return pr, NewErrorService(INTERNAL_ERROR, err.Error())
	}

	return pr, nil
}

func (s *PrReviewerService) PullRequestMerge(pullRequestId string) (models.PullRequest, error) {
	if pr := s.repo.GetPullRequestById(pullRequestId); pr != nil {
		if pr.Status != models.MERGED {
			now := time.Now()
			pr.Status = models.MERGED
			pr.MergedAt = &now
		}
		return *pr, nil
	}

	return models.PullRequest{}, NewErrorService(OBJECT_NOT_FOUND, "Pull request not found")
}

func (s *PrReviewerService) PullRequestReassign(pullRequestId, oldUserId string) (models.PullRequest, string, error) {
	pr, user := s.repo.GetPullRequestById(pullRequestId), s.repo.GetUserById(oldUserId)
	if pr == nil || user == nil {
		return models.PullRequest{}, "", NewErrorService(OBJECT_NOT_FOUND, "Pull request or user not found")
	}
	if pr.Status == models.MERGED {
		return *pr, "", NewErrorService(DOMAIN_ERROR, "Cannot reassign on merged PR")
	}

	for i, reviewer := range pr.AssignedReviewers {
		if reviewer == oldUserId {
			team := s.repo.GetTeamByName(user.TeamName)
			for _, candidate := range team.Members {
				if candidate.UserId != reviewer && candidate.IsActive {
					// I see 2 solutions here: hash map or just for iterating
					// number of reviewers of PR - only 2
					// I can hardcode with ifs, but I think code with simple extension will be better
					// there are no real reasons to use hash map instead of iterating on such small numbers
					collision := false
					for _, r := range pr.AssignedReviewers {
						if candidate.UserId == r {
							collision = true
						}
					}
					if !collision {
						pr.AssignedReviewers[i] = candidate.UserId
						return *pr, candidate.UserId, nil
					}
				}
			}
			break
		}
	}

	return *pr, "", NewErrorService(DOMAIN_ERROR, "reviewer is not assigned to this PR")
}

func (s *PrReviewerService) UsersGetReview(userId string) ([]models.PullRequestShort, error) {
	prs := s.repo.GetPullRequestsByUserId(userId)
	if prs == nil {
		return nil, NewErrorService(OBJECT_NOT_FOUND, "User not found")
	}

	prsShort := make([]models.PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		prsShort = append(prsShort, models.NewPRShort(pr))
	}

	return prsShort, nil
}
