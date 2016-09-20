package listeners

import (
	"github.com/AaronO/go-git-http"
	"gopkg.in/inconshreveable/log15.v2"

	"context"

	"sourcegraph.com/sourcegraph/sourcegraph/api/sourcegraph"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/conf/feature"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/conf/universe"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/langp"
	"sourcegraph.com/sourcegraph/sourcegraph/services/events"
)

func init() {
	events.RegisterListener(&gitHookListener{})
}

type gitHookListener struct{}

func (g *gitHookListener) Scopes() []string {
	return []string{"internal:githooks"}
}

func (g *gitHookListener) Start(ctx context.Context) {
	buildCallback := func(id events.EventID, p events.GitPayload) {
		buildHook(ctx, id, p)
	}
	events.Subscribe(events.GitPushEvent, buildCallback)
	events.Subscribe(events.GitCreateBranchEvent, buildCallback)

	inventoryCallback := func(id events.EventID, p events.GitPayload) {
		inventoryHook(ctx, id, p)
	}
	events.Subscribe(events.GitPushEvent, inventoryCallback)
	events.Subscribe(events.GitCreateBranchEvent, inventoryCallback)
}

func buildHook(ctx context.Context, id events.EventID, payload events.GitPayload) {
	cl, err := sourcegraph.NewClientFromContext(ctx)
	if err != nil {
		log15.Error("postPushHook: failed to create build", "err", err)
		return
	}
	repo := payload.Repo
	event := payload.Event

	if event.Type != githttp.PUSH && event.Type != githttp.PUSH_FORCE && event.Type != githttp.TAG {
		return
	}

	if feature.Features.Universe {
		repoFull, err := cl.Repos.Get(ctx, &sourcegraph.RepoSpec{ID: repo})
		if err != nil {
			log15.Error("postPushHook: failed to prepare workspace", "err", err)
			return
		}
		if universe.EnabledRepo(repoFull) || universe.Shadow(repoFull.URI) {
			// Ask the Language Processor to prepare the workspace.
			if err := langp.DefaultClient.Prepare(ctx, &langp.RepoRev{
				// TODO(slimsag): URI is correct only where the repo URI and clone
				// URI are directly equal.. but CloneURI is only correct (for Go)
				// when it directly matches the package import path.
				Repo:   repoFull.URI,
				Commit: event.Commit,
			}); err != nil {
				log15.Error("postPushHook: failed to prepare workspace", "err", err)
				return
			}
		}
	}
}

// inventoryHook triggers a Repos.GetInventory call that computes the
// repo's inventory and caches it in a commit status (and saves it to
// the repo's Language field for default branch pushes). Then it is
// available immediately for future callers (which generally expect
// that operation to be fast).
func inventoryHook(ctx context.Context, id events.EventID, payload events.GitPayload) {
	cl, err := sourcegraph.NewClientFromContext(ctx)
	if err != nil {
		log15.Error("inventoryHook error", "err", err)
	}
	event := payload.Event
	if event.Type == githttp.PUSH || event.Type == githttp.PUSH_FORCE {
		repoRev := &sourcegraph.RepoRevSpec{Repo: payload.Repo, CommitID: event.Commit}
		// Trigger a call to Repos.GetInventory so the inventory is
		// cached for subsequent calls.
		inv, err := cl.Repos.GetInventory(ctx, repoRev)
		if err != nil {
			log15.Warn("inventoryHook: call to Repos.GetInventory failed", "err", err, "repoRev", repoRev)
			return
		}

		// If this push is to the default branch, update the repo's
		// Language field with the primary language.
		repo, err := cl.Repos.Get(ctx, &sourcegraph.RepoSpec{ID: repoRev.Repo})
		if err != nil {
			log15.Warn("inventoryHook: call to Repos.Get failed", "err", err, "repoRev", repoRev)
			return
		}
		if event.Branch == repo.DefaultBranch {
			lang := inv.PrimaryProgrammingLanguage()
			if _, err := cl.Repos.Update(ctx, &sourcegraph.ReposUpdateOp{Repo: repo.ID, Language: lang}); err != nil {
				log15.Warn("inventoryHook: call to Repos.Update to set language failed", "err", err, "repoRev", repoRev, "language", lang)
			}
		}
	}
}
