package shared

import (
	"encoding/json"
	"testing"
)

func TestGithubToCommits_NoRepos(t *testing.T) {
	jsonData := `{
        "data": {
            "viewer": {
                "repositories": {
                    "nodes": []
                }
            }
        }
    }`

	var resp GitHubAPIResponse
	json.Unmarshal([]byte(jsonData), &resp)

	commits, err := resp.ToCommits()
	if err != nil {
		t.Error("Expected success", err)
	}
	if len(commits) > 0 {
		t.Error("Expected empty commits")
	}
}

func TestGithubToCommits_Success(t *testing.T) {
	jsonData := `{
        "data": {
            "viewer": {
                "repositories": {
                    "nodes": [
                        {
                            "name": "test-repo",
                            "url": "https://github.com/user/test-repo",
                            "updatedAt": "2023-01-15T10:30:00Z",
                            "defaultBranchRef": {
                                "target": {
                                    "history": {
                                        "nodes": [
                                            {
                                                "messageHeadline": "Fix bug in parser",
                                                "committedDate": "2023-01-15T10:30:00Z",
                                                "commitUrl": "https://github.com/user/test-repo/commit/abc123"
                                            },
                                            {
                                                "messageHeadline": "Add new feature",
                                                "committedDate": "2023-01-14T09:15:00Z",
                                                "commitUrl": "https://github.com/user/test-repo/commit/def456"
                                            }
                                        ]
                                    }
                                }
                            }
                        }
                    ]
                }
            }
        }
    }`

	var resp GitHubAPIResponse
	json.Unmarshal([]byte(jsonData), &resp)

	commits, err := resp.ToCommits()
	if err != nil {
		t.Error("Expected success", err)
	}
	if len(commits) != 2 {
		t.Errorf("Expected 2 commits, got %d", len(commits))
	}
}
