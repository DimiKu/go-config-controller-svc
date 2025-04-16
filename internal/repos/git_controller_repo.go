package repos

import (
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"go.uber.org/zap"
	"log"
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
	_, err := os.Stat(g.localPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			g.log.Error("Failed to stat repo: %v, with path %v", zap.Error(err), zap.String("path", g.localPath))
			return err
		} else {
			_, errA := git.PlainClone(g.localPath, false, &git.CloneOptions{
				URL: g.repoURL,
			})
			if errA != nil {
				g.log.Error("Failed to open repo: %v", zap.Error(err))
				return errA
			}
		}
	}

	repo, err := git.PlainOpen(g.localPath)
	if err != nil {
		log.Fatalf("Failed to open repo: %v", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Failed to get Worktree: %v", err)
	}

	if err := worktree.Pull(&git.PullOptions{RemoteName: "origin"}); err != nil {
		if !errors.Is(err, git.NoErrAlreadyUpToDate) {
			return err
		}
	}

	branchRef := plumbing.ReferenceName("refs/heads/" + branch)
	g.repoRef, err = repo.Reference(branchRef, true)
	if err != nil {
		log.Fatalf("Ветка '%s' не найдена: %v", branch, err)
	}

	if branch != "main" {
		err = worktree.Checkout(&git.CheckoutOptions{
			Branch: branchRef,
		})
		if err != nil {
			log.Fatalf("Ошибка при переключении на ветку '%s': %v", branch, err)
		}
	}

	g.repoObj = repo
	return nil
}

func (g *GitControllerRepo) GetLastCommitTime() (time.Time, error) {
	commit, err := g.repoObj.CommitObject(g.repoRef.Hash())
	if err != nil {
		log.Fatalf("Ошибка при получении коммита: %v", err)
	}

	lastCommitTime := commit.Author.When
	return lastCommitTime, nil
}

//func (g *GitControllerRepo) GitPull(repo *git.Repository) error {
//	remote, err := repo.Remote("origin")
//	if err != nil {
//		log.Fatalf("Ошибка при получении удалённого репозитория: %v", err)
//	}
//
//	fetchOptions := &git.FetchOptions{
//		RemoteName: "origin",
//		Auth: &http.BasicAuth{
//			Username: g.userName,
//			Password: g.password,
//		},
//	}
//	if err := remote.Fetch(fetchOptions); err != nil && err != git.NoErrAlreadyUpToDate {
//		log.Fatalf("Ошибка при выполнении fetch: %v", err)
//	}
//
//	head, err := repo.Head()
//	if err != nil {
//		log.Fatalf("Ошибка при получении HEAD: %v", err)
//	}
//
//	worktree, err := repo.Worktree()
//	if err != nil {
//		log.Fatalf("Ошибка при получении рабочей директории: %v", err)
//	}
//
//	mergeOptions := &git.MergeOptions{
//		//Strategy: 1,
//	}
//
//	worktree.
//
//	if err := worktree.Merge(&plumbing.ReferenceName(head.Name()), mergeOptions); err != nil {
//		log.Fatalf("Ошибка при выполнении merge: %v", err)
//	}
