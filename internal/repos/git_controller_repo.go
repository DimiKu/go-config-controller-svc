package repos

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"go.uber.org/zap"
	"os"
	"time"
)

type GitControllerRepo struct {
	localPath string
	userName  string
	password  string
	repoURL   string
	repoRef   *plumbing.Reference
	log       *zap.Logger
	repoObj   *git.Repository
}

func NewGitControllerRepo(
	localPath string,
	userName string,
	password string,
	repoURL string,
	log *zap.Logger,
) *GitControllerRepo {
	return &GitControllerRepo{
		localPath: localPath,
		userName:  userName,
		password:  password,
		repoURL:   repoURL,
		log:       log,
	}
}

func (g *GitControllerRepo) InitRepo(branch string) error {
	if err := g.cloneRepo(g.localPath); err != nil {
		return err
	}

	repo, err := git.PlainOpen(g.localPath)
	if err != nil {
		g.log.Error("Failed to open repo: ", zap.Error(err))
	}

	worktree, err := g.getWorktreeAndPull(repo)
	if err != nil {
		return err
	}

	if err := g.getNeededBranch(branch, worktree, repo); err != nil {
		return err
	}

	g.repoObj = repo
	return nil
}

func (g *GitControllerRepo) GetLastCommitTime() (time.Time, error) {
	commit, err := g.repoObj.CommitObject(g.repoRef.Hash())
	if err != nil {
		g.log.Error("Failed to get commit time", zap.Error(err))
		return time.Time{}, err
	}

	lastCommitTime := commit.Author.When
	return lastCommitTime, nil
}

func (g *GitControllerRepo) getWorktreeAndPull(repo *git.Repository) (*git.Worktree, error) {
	worktree, err := repo.Worktree()
	if err != nil {
		g.log.Error("Failed to get Worktree: ", zap.Error(err))
		return nil, err

	}

	if err := worktree.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
		if !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return nil, err
		}
	}

	return worktree, nil
}

func (g *GitControllerRepo) getNeededBranch(branch string, worktree *git.Worktree, repo *git.Repository) error {
	var err error
	branchRef := plumbing.ReferenceName("refs/heads/" + branch)
	g.repoRef, err = repo.Reference(branchRef, true)
	if err != nil {
		g.log.Error("Branch doesn't exist", zap.String("branch", branch), zap.Error(err))
		return err
	}

	if branch != "main" {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: branchRef,
		})
		if err != nil {
			g.log.Error("Failed checkout branch", zap.String("branch", branch), zap.Error(err))
			return err
		}
	}

	return nil
}

func (g *GitControllerRepo) cloneRepo(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			g.log.Error("Failed to stat repo", zap.Error(err), zap.String("path", g.localPath))
			return err
		} else {
			_, errA := git.PlainClone(path, false, &git.CloneOptions{
				URL: g.repoURL,
			})
			if errA != nil {
				g.log.Error("Failed to open repo: %v", zap.Error(err))
				return errA
			}
		}
	}

	return nil
}
