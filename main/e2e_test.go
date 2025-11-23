package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/Dowtai/pr-reviewer-service/internal/models"
)

func doRequest(t *testing.T, method, url string, body interface{}) *http.Response {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal body: %v", err)
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	return resp
}

func assertJSONEqual[T any](t *testing.T, resp *http.Response, expected T) T {
	defer resp.Body.Close()
	var actual T
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	if err := json.Unmarshal(body, &actual); err != nil {
		t.Fatalf("Failed to unmarshal: %v\nBody: %s", err, string(body))
	}

	if !reflect.DeepEqual(actual, expected) {
		expBytes, _ := json.MarshalIndent(expected, "", "  ")
		actBytes, _ := json.MarshalIndent(actual, "", "  ")
		t.Fatalf("JSON not equal\nExpected:\n%s\nActual:\n%s", expBytes, actBytes)
	}
	return actual
}

const baseURL = "http://localhost:8080"

func createTeam(t *testing.T, teamName string, members []models.TeamMember) models.Team {
	teamReq := map[string]interface{}{
		"team_name": teamName,
		"members":   members,
	}
	expected := models.Team{
		TeamName: teamName,
		Members:  members,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/team/add", teamReq)
	if resp.StatusCode != 201 {
		t.Fatalf("Expected 201, got %d", resp.StatusCode)
	}
	assertJSONEqual(t, resp, expected)
	return expected
}

func createTeamExpectError(t *testing.T, teamName string, members []models.TeamMember, expectedStatus int, code models.ErrorDetailCode, message string) {
	teamReq := map[string]interface{}{
		"team_name": teamName,
		"members":   members,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/team/add", teamReq)
	if resp.StatusCode != expectedStatus {
		t.Fatalf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}
	expectedErr := models.ErrorResponse{
		Detail: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	assertJSONEqual(t, resp, expectedErr)
}

func getTeam(t *testing.T, teamName string, expected models.Team) models.Team {
	url := baseURL + "/team/get?team_name=" + teamName
	resp := doRequest(t, http.MethodGet, url, nil)

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	return assertJSONEqual(t, resp, expected)
}

func getTeamExpectError(t *testing.T, teamName string, expectedStatus int, code models.ErrorDetailCode, message string) {
	url := baseURL + "/team/get?team_name=" + teamName
	resp := doRequest(t, http.MethodGet, url, nil)

	if resp.StatusCode != expectedStatus {
		t.Fatalf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	expectedErr := models.ErrorResponse{
		Detail: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}

	assertJSONEqual(t, resp, expectedErr)
}

func setUserIsActive(t *testing.T, userId string, isActive bool, expected models.User) models.User {
	req := map[string]interface{}{
		"user_id":   userId,
		"is_active": isActive,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/users/setIsActive", req)
	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	return assertJSONEqual(t, resp, expected)
}

func setUserIsActiveExpectError(t *testing.T, userId string, isActive bool, expectedStatus int, code models.ErrorDetailCode, message string) {
	req := map[string]interface{}{
		"user_id":   userId,
		"is_active": isActive,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/users/setIsActive", req)
	if resp.StatusCode != expectedStatus {
		t.Fatalf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	expectedErr := models.ErrorResponse{
		Detail: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}

	assertJSONEqual(t, resp, expectedErr)
}

func createPullRequest(t *testing.T, pullRequestId, pullRequestName, authorId string, expected *models.PullRequest) models.PullRequest {
	prReq := map[string]interface{}{
		"pull_request_id":   pullRequestId,
		"pull_request_name": pullRequestName,
		"author_id":         authorId,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/pullRequest/create", prReq)
	if resp.StatusCode != 201 {
		t.Fatalf("Expected 201, got %d", resp.StatusCode)
	}

	var actual models.PullRequest
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &actual); err != nil {
		t.Fatalf("Failed to unmarshal: %v\nBody: %s", err, string(body))
	}

	expected.CreatedAt = actual.CreatedAt
	if !reflect.DeepEqual(actual, *expected) {
		expBytes, _ := json.MarshalIndent(expected, "", "  ")
		actBytes, _ := json.MarshalIndent(actual, "", "  ")
		t.Fatalf("JSON not equal\nExpected:\n%s\nActual:\n%s", expBytes, actBytes)
	}
	return actual
}

func createPullRequestExpectError(t *testing.T, pullRequestId, pullRequestName, authorId string, expectedStatus int, code models.ErrorDetailCode, message string) {
	prReq := map[string]interface{}{
		"pull_request_id":   pullRequestId,
		"pull_request_name": pullRequestName,
		"author_id":         authorId,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/pullRequest/create", prReq)
	if resp.StatusCode != expectedStatus {
		t.Fatalf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}
	expectedErr := models.ErrorResponse{
		Detail: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	assertJSONEqual(t, resp, expectedErr)
}

func mergePullRequest(t *testing.T, pullRequestId string, expected *models.PullRequest) {
	prReq := map[string]interface{}{
		"pull_request_id": pullRequestId,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/pullRequest/merge", prReq)
	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var actual models.PullRequest
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &actual); err != nil {
		t.Fatalf("Failed to unmarshal: %v\nBody: %s", err, string(body))
	}

	if expected.MergedAt == nil {
		expected.MergedAt = actual.MergedAt
	}
	if !reflect.DeepEqual(actual, *expected) {
		expBytes, _ := json.MarshalIndent(expected, "", "  ")
		actBytes, _ := json.MarshalIndent(actual, "", "  ")
		t.Fatalf("JSON not equal\nExpected:\n%s\nActual:\n%s", expBytes, actBytes)
	}
}

func mergePullRequestExpectError(t *testing.T, pullRequestId string, expectedStatus int, code models.ErrorDetailCode, message string) {
	prReq := map[string]interface{}{
		"pull_request_id": pullRequestId,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/pullRequest/merge", prReq)
	if resp.StatusCode != expectedStatus {
		t.Fatalf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}
	expectedErr := models.ErrorResponse{
		Detail: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	assertJSONEqual(t, resp, expectedErr)
}

func reassignPullRequest(t *testing.T, pullRequestId, oldUserId string, expectedPR models.PullRequest, expectedUserId string) {
	prReq := map[string]interface{}{
		"pull_request_id": pullRequestId,
		"old_user_id":     oldUserId,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/pullRequest/reassign", prReq)
	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var actual struct {
		PR         models.PullRequest `json:"pr"`
		ReplacedBy string             `json:"replaced_by"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	defer resp.Body.Close()

	if err := json.Unmarshal(body, &actual); err != nil {
		t.Fatalf("Failed to unmarshal: %v\nBody: %s", err, string(body))
	}

	if !reflect.DeepEqual(actual.PR, expectedPR) {
		expBytes, _ := json.MarshalIndent(expectedPR, "", "  ")
		actBytes, _ := json.MarshalIndent(actual, "", "  ")
		t.Fatalf("JSON not equal\nExpected:\n%s\nActual:\n%s", expBytes, actBytes)
	}
	if !reflect.DeepEqual(actual.ReplacedBy, expectedUserId) {
		expBytes, _ := json.MarshalIndent(expectedPR, "", "  ")
		actBytes, _ := json.MarshalIndent(actual, "", "  ")
		t.Fatalf("JSON not equal\nExpected:\n%s\nActual:\n%s", expBytes, actBytes)
	}
}

func reassignPullRequestExpectError(t *testing.T, pullRequestId, oldUserId string, expectedStatus int, code models.ErrorDetailCode, message string) {
	prReq := map[string]interface{}{
		"pull_request_id": pullRequestId,
		"old_user_id":     oldUserId,
	}

	resp := doRequest(t, http.MethodPost, baseURL+"/pullRequest/reassign", prReq)
	if resp.StatusCode != expectedStatus {
		t.Fatalf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}
	expectedErr := models.ErrorResponse{
		Detail: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
	assertJSONEqual(t, resp, expectedErr)
}

func getReview(t *testing.T, userId string, expected []models.PullRequestShort) {
	url := baseURL + "/users/getReview?user_id=" + userId
	resp := doRequest(t, http.MethodGet, url, nil)

	if resp.StatusCode != 200 {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	assertJSONEqual(t, resp, expected)
}

func getReviewExpectError(t *testing.T, userId string, expectedStatus int, code models.ErrorDetailCode, message string) {
	url := baseURL + "/users/getReview?user_id=" + userId
	resp := doRequest(t, http.MethodGet, url, nil)

	if resp.StatusCode != expectedStatus {
		t.Fatalf("Expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	expectedErr := models.ErrorResponse{
		Detail: models.ErrorDetail{
			Code:    code,
			Message: message,
		},
	}

	assertJSONEqual(t, resp, expectedErr)
}

func TestTeam(t *testing.T) {
	server := NewServer()
	go func() { _ = server.ListenAndServe() }()
	defer ShutdownServer(server)
	var team1, team2 models.Team

	t.Run("AddUniqueTeams", func(t *testing.T) {
		members := []models.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
			{UserId: "u2", Username: "Bob", IsActive: true},
		}
		members2 := []models.TeamMember{
			{UserId: "u3", Username: "Alice", IsActive: true},
			{UserId: "u4", Username: "Bob", IsActive: true},
		}
		team1 = createTeam(t, "backend", members)
		team2 = createTeam(t, "backend2", members2)
	})

	t.Run("AddDuplicate", func(t *testing.T) {
		members := []models.TeamMember{
			{UserId: "u3", Username: "Alice", IsActive: true},
			{UserId: "u4", Username: "Bob", IsActive: true},
		}
		createTeamExpectError(t, "backend", members, 400, models.TEAM_EXISTS, "team_name already exists")
		createTeamExpectError(t, "backend3", members, 500, models.FATAL_ERROR, "user already exists")
	})

	t.Run("GetExistingTeam", func(t *testing.T) {
		getTeam(t, "backend", team1)
		getTeam(t, "backend2", team2)
	})

	t.Run("GetNonExistingTeam", func(t *testing.T) {
		getTeamExpectError(t, "backend3", 404, models.NOT_FOUND, "team_name not found")
	})

}

func TestUsers(t *testing.T) {
	server := NewServer()
	go func() { _ = server.ListenAndServe() }()
	defer ShutdownServer(server)
	var team1 models.Team

	t.Run("SetIsActive", func(t *testing.T) {
		members := []models.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
			{UserId: "u2", Username: "Bob", IsActive: true},
		}
		members2 := []models.TeamMember{
			{UserId: "u3", Username: "Alice", IsActive: true},
			{UserId: "u4", Username: "Bob", IsActive: true},
		}
		team1 = createTeam(t, "backend", members)
		createTeam(t, "backend2", members2)

		expectedUser1 := models.User{
			UserId:   "u1",
			Username: "Alice",
			TeamName: "backend",
			IsActive: true,
		}

		setUserIsActive(t, "u1", true, expectedUser1)
		getTeam(t, "backend", team1)
		expectedUser1.IsActive = false
		setUserIsActive(t, "u1", false, expectedUser1)
		team1.Members[0].IsActive = false
		getTeam(t, "backend", team1)
	})

	t.Run("SetIsActiveNonExisting", func(t *testing.T) {
		setUserIsActiveExpectError(t, "u5", true, 404, models.NOT_FOUND, "user_id not found")
	})
}

func TestPullRequest(t *testing.T) {
	server := NewServer()
	go func() { _ = server.ListenAndServe() }()
	defer ShutdownServer(server)
	var pr, pr1, pr2 models.PullRequest

	t.Run("Create", func(t *testing.T) {
		members := []models.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
			{UserId: "u2", Username: "Bob", IsActive: true},
		}
		members2 := []models.TeamMember{
			{UserId: "u3", Username: "Alice", IsActive: true},
			{UserId: "u4", Username: "Bob", IsActive: true},
		}
		createTeam(t, "backend", members)
		createTeam(t, "backend2", members2)

		expectedPR := models.PullRequest{
			PullRequestId:     "r1",
			PullRequestName:   "req1",
			AuthorId:          "u1",
			Status:            models.OPEN,
			AssignedReviewers: []string{"u2"},
		}
		pr1 = createPullRequest(t, "r1", "req1", "u1", &expectedPR)

		expectedUser3 := models.User{
			UserId:   "u3",
			Username: "Alice",
			TeamName: "backend2",
			IsActive: false,
		}
		setUserIsActive(t, "u3", false, expectedUser3)

		expectedPR2 := models.PullRequest{
			PullRequestId:     "r2",
			PullRequestName:   "req2",
			AuthorId:          "u4",
			Status:            models.OPEN,
			AssignedReviewers: []string{},
		}
		pr2 = createPullRequest(t, "r2", "req2", "u4", &expectedPR2)
	})

	t.Run("CreateNotFound", func(t *testing.T) {
		createPullRequestExpectError(t, "r3", "req3", "u5", 404, models.NOT_FOUND, "author or team not found")
	})

	t.Run("CreateExistingPullRequest", func(t *testing.T) {
		_ = pr2
		createPullRequestExpectError(t, "r1", "req1", "u1", 409, models.PR_EXISTS, "pull_request already exists")
		createPullRequestExpectError(t, "r1", "req3", "u4", 409, models.PR_EXISTS, "pull_request already exists")
	})

	t.Run("MergeIdempotency", func(t *testing.T) {
		pr1.Status = models.MERGED
		mergePullRequest(t, "r1", &pr1)
		timestamp := pr1.MergedAt
		for i := 0; i < 200; i++ {
			mergePullRequest(t, "r1", &pr1)
			if pr1.MergedAt != timestamp {
				t.Fatal("Something went wrong with tests")
			}
		}
	})

	t.Run("MergeNonExisting", func(t *testing.T) {
		mergePullRequestExpectError(t, "r3", 404, models.NOT_FOUND, "pull_request not found")
	})

	t.Run("Reassign", func(t *testing.T) {
		members := []models.TeamMember{
			{UserId: "u5", Username: "user1", IsActive: true},
			{UserId: "u6", Username: "user2", IsActive: true},
			{UserId: "u7", Username: "user3", IsActive: true},
			{UserId: "u8", Username: "user4", IsActive: true},
			{UserId: "u9", Username: "user5", IsActive: true},
			{UserId: "u10", Username: "user6", IsActive: true},
			{UserId: "u11", Username: "user7", IsActive: true},
			{UserId: "u12", Username: "user8", IsActive: true},
			{UserId: "u13", Username: "user9", IsActive: true},
			{UserId: "u14", Username: "user10", IsActive: true},
		}
		createTeam(t, "backend3", members)

		expectedPR := models.PullRequest{
			PullRequestId:     "r3",
			PullRequestName:   "req3",
			AuthorId:          "u5",
			Status:            models.OPEN,
			AssignedReviewers: []string{"u6", "u7"},
		}
		pr = createPullRequest(t, "r3", "req3", "u5", &expectedPR)

		for i := 0; i < 2; i++ {
			for j := 0; j < 10; j++ {
				old := pr.AssignedReviewers[i]
				for k := 6; k <= 14; k++ {
					newId := fmt.Sprintf("u%d", k)
					if pr.AssignedReviewers[0] != newId && pr.AssignedReviewers[1] != newId {
						pr.AssignedReviewers[i] = newId
						break
					}
				}
				reassignPullRequest(t, "r3", old, pr, pr.AssignedReviewers[i])
			}
		}
	})

	t.Run("ReassignNonExisting", func(t *testing.T) {
		reassignPullRequestExpectError(t, "r4", "u1", 404, models.NOT_FOUND, "pull_request or user not found")
	})

	t.Run("ReassignMerged", func(t *testing.T) {
		reassignPullRequestExpectError(t, "r1", "u2", 409, models.PR_MERGED, "cannot reassign on merged PR")
	})

	t.Run("ReassignNotAssigned", func(t *testing.T) {
		for k := 6; k <= 14; k++ {
			newId := fmt.Sprintf("u%d", k)
			if pr.AssignedReviewers[0] != newId && pr.AssignedReviewers[1] != newId {
				reassignPullRequestExpectError(t, "r3", newId, 409, models.NOT_ASSIGNED, "reviewer is not assigned to this PR")
				break
			}
		}
	})

	t.Run("NoActive", func(t *testing.T) {
		var reviewer string
		for k := 6; k <= 14; k++ {
			id := fmt.Sprintf("u%d", k)
			name := fmt.Sprintf("user%d", k-4)
			expectedUser := models.User{
				UserId:   id,
				Username: name,
				TeamName: "backend3",
				IsActive: false,
			}
			setUserIsActive(t, id, false, expectedUser)
			if id == pr.AssignedReviewers[0] || id == pr.AssignedReviewers[1] {
				reviewer = id
			}
		}
		reassignPullRequestExpectError(t, "r3", reviewer, 409, models.NO_CANDIDATE, "no active replacement candidate in team")
	})

	t.Run("GetReview", func(t *testing.T) {
		for k := 6; k <= 14; k++ {
			id := fmt.Sprintf("u%d", k)
			name := fmt.Sprintf("user%d", k-4)
			expectedUser := models.User{
				UserId:   id,
				Username: name,
				TeamName: "backend3",
				IsActive: true,
			}
			setUserIsActive(t, id, true, expectedUser)
		}
		expectedPR := models.PullRequest{
			PullRequestId:     "r4",
			PullRequestName:   "req4",
			AuthorId:          "u7",
			Status:            models.OPEN,
			AssignedReviewers: []string{"u5", "u6"},
		}
		createPullRequest(t, "r4", "req4", "u7", &expectedPR)
		fmt.Println(pr, pr1, pr2, expectedPR)

		getReview(t, "u1", []models.PullRequestShort{})
		getReview(t, "u2", []models.PullRequestShort{models.NewPRShort(&pr1)})
		getReview(t, "u5", []models.PullRequestShort{models.NewPRShort(&expectedPR)})
		getReview(t, "u6", []models.PullRequestShort{models.NewPRShort(&pr), models.NewPRShort(&expectedPR)})
		getReview(t, "u7", []models.PullRequestShort{models.NewPRShort(&pr)})
	})

	t.Run("GetReviewNotFound", func(t *testing.T) {
		getReviewExpectError(t, "u23", 404, models.NOT_FOUND, "user not found")
	})
}
