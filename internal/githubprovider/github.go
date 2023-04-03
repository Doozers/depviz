package githubprovider

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"moul.io/multipmuri"

	"github.com/Doozers/depviz/internal/dvmodel"
	"github.com/google/go-github/v30/github"
	"golang.org/x/oauth2"
)

type Opts struct {
	Since  *time.Time  `json:"since"`
	Logger *zap.Logger `json:"-"`
}

func FetchRepo(ctx context.Context, entity multipmuri.Entity, token string, out chan<- dvmodel.Batch, opts Opts) { // nolint:interfacer
	if opts.Logger == nil {
		opts.Logger = zap.NewNop()
	}

	type multipmuriMinimalInterface interface {
		Repo() *multipmuri.GitHubRepo
	}
	target, ok := entity.(multipmuriMinimalInterface)
	if !ok {
		opts.Logger.Warn("invalid entity", zap.String("entity", fmt.Sprintf("%v", entity.String())))
		return
	}
	repo := target.Repo()

	// create client
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// queries
	totalIssues := 0
	callOpts := &github.IssueListByRepoOptions{State: "all"}
	if opts.Since != nil {
		callOpts.Since = *opts.Since
	}
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, repo.OwnerID(), repo.RepoID(), callOpts)
		if err != nil {
			opts.Logger.Error("fetch GitHub issues", zap.Error(err))
			return
		}
		totalIssues += len(issues)
		opts.Logger.Debug("paginate",
			zap.Any("opts", opts),
			zap.String("provider", "github"),
			zap.String("repo", repo.String()),
			zap.Int("new-issues", len(issues)),
			zap.Int("total-issues", totalIssues),
		)

		if len(issues) > 0 {
			batch := fromIssues(issues, opts.Logger)
			out <- batch
		}

		// handle pagination
		if resp.NextPage == 0 {
			break
		}
		callOpts.Page = resp.NextPage
	}

	if rateLimits, _, err := client.RateLimits(ctx); err == nil {
		opts.Logger.Debug("github API rate limiting", zap.Stringer("limit", rateLimits.GetCore()))
	}

	// FIXME: fetch incomplete/old users, orgs, teams & repos
}
