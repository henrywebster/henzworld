package github

import (
	"fmt"
	"henzworld/internal/model"
	"sort"
	"time"
)

func (r *Response) ToCommits() ([]model.Commit, error) {
	var commits []model.Commit

	if len(r.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL errors: %v", r.Errors)
	}

	for _, repo := range r.Data.Viewer.Repositories.Nodes {
		if repo.DefaultBranchRef.Target.History.Nodes == nil {
			// no commits
			continue
		}

		for _, node := range repo.DefaultBranchRef.Target.History.Nodes {
			commitDate, err := time.Parse(time.RFC3339, node.CommittedDate)
			if err != nil {
				return nil, fmt.Errorf("invalid commit date %s: %w", node.CommittedDate, err)
			}

			commit := model.Commit{
				Message:  node.MessageHeadline,
				URL:      node.CommitURL,
				RepoURL:  repo.URL,
				RepoName: repo.Name,
				Date:     commitDate,
			}

			commits = append(commits, commit)
		}
	}

	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})

	return commits, nil
}
