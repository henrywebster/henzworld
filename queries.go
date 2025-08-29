package main

const PublicReposCommitsQuery = `
{
    viewer {
        repositories(first: 5, privacy: PUBLIC, orderBy: {field: PUSHED_AT, direction: DESC}) {
            nodes {
                name
                url
                updatedAt
                defaultBranchRef {
                    target {
                        ... on Commit {
                            history(first: 5) {
                                nodes {
                                    messageHeadline
                                    committedDate
                                    commitUrl
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}`
