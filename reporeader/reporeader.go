package reporeader

import (
	"fmt"
	"os"
	"time"

	"github.com/go-enry/go-license-detector/v4/licensedb"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type RepoReader struct {
	repository *git.Repository
}

type RepoDetails struct {
	CreatedDate    time.Time
	AuthorsCommits map[Author][]object.Commit
	License        string
}

type Author struct {
	Name  string
	Email string
}

func NewRepoReader(dir string) (*RepoReader, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("NewRepoReader: error detected in attempting to open repository: %w", err)
	}

	return &RepoReader{repository: repo}, nil
}

func NewRepoReaderRepository(repo *git.Repository) (*RepoReader, error) {
	_, err := ValidateRepository(repo)
	if err != nil {
		return nil, fmt.Errorf("NewRepoReaderRepository: received an invalid repository: %w", err)
	}

	return &RepoReader{repository: repo}, nil
}

func (r *RepoReader) GetRepoDetails() (RepoDetails, error) {
	head, err := r.repository.Head()
	if err != nil {
		return RepoDetails{}, fmt.Errorf("GetRepoDetails: unable to get the repository head: %w", err)
	}

	commits := make([]*object.Commit, 0)
	cIter, _ := r.repository.Log(&git.LogOptions{From: head.Hash(), Order: git.LogOrderCommitterTime})
	_ = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	wt, err := r.repository.Worktree()
	if err != nil {
		return RepoDetails{}, fmt.Errorf("GetRepoDetails: unable to get the worktree from the repository: %w", err)
	}

	createdDate := r.getCreatedDate(commits)
	authorsCommits := r.getAuthorsByCommits(commits)

	license, err := r.getLicenseFromRoot(wt.Filesystem)
	if err != nil {
		return RepoDetails{}, fmt.Errorf("GetRepoDetails: unable to get the license for the repository: %w", err)
	}

	details := RepoDetails{
		CreatedDate:    createdDate,
		AuthorsCommits: authorsCommits,
		License:        license,
	}

	return details, nil
}

func (r *RepoReader) getCreatedDate(commits []*object.Commit) time.Time {
	if len(commits) == 0 {
		return time.Time{}
	}

	oldestTime := commits[0].Author.When

	for _, commit := range commits {
		commitTime := commit.Author.When
		if commitTime.Before(oldestTime) {
			oldestTime = commit.Author.When
		}
	}

	return oldestTime
}

func (r *RepoReader) getAuthorsByCommits(commits []*object.Commit) map[Author][]object.Commit {
	contributorCommits := make(map[Author][]object.Commit)

	for _, commit := range commits {
		author := Author{
			commit.Author.Name,
			commit.Author.Email,
		}

		contributorCommits[author] = append(contributorCommits[author], *commit)
	}

	return contributorCommits
}

func (r *RepoReader) getLicenseFromRoot(fs billy.Filesystem) (string, error) {
	files, err := fs.ReadDir(".")
	if err != nil {
		return "", fmt.Errorf("getLicenseFromRoot: could not read root directory: %w", err)
	}

	var licenseInfo os.FileInfo
	for _, file := range files {
		if file.Name() == "LICENSE" {
			licenseInfo = file
			break
		}
	}
	if licenseInfo == nil {
		return "NO LICENSE", nil
	}

	licensePath := fs.Join(fs.Root(), licenseInfo.Name())

	contents, err := os.ReadFile(licensePath)
	if err != nil {
		return "", fmt.Errorf("getLicenseFromRoot: could not read file from filepath %s: %w", licensePath, err)
	}

	results := licensedb.InvestigateLicenseText(contents)

	type licenseConfidence struct {
		license    string
		confidence float32
	}

	bestConfidence := licenseConfidence{}
	for license, confidence := range results {
		if bestConfidence.confidence < confidence {
			bestConfidence.license = license
			bestConfidence.confidence = confidence
		}
	}

	return bestConfidence.license, nil
}

// GetCreatedDate returns the time that the repository was first created.
func (r *RepoReader) GetCreatedDate() (time.Time, error) {
	head, err := r.repository.Head()
	if err != nil {
		return time.Time{}, fmt.Errorf("GetCreatedDate: unable to get the repository head: %w", err)
	}

	commits := make([]*object.Commit, 0)
	cIter, _ := r.repository.Log(&git.LogOptions{From: head.Hash(), Order: git.LogOrderCommitterTime})
	_ = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	return r.getCreatedDate(commits), nil
}

// GetAuthorsByCommits returns the authors and their commits they made.
func (r *RepoReader) GetAuthorsByCommits() (map[Author][]object.Commit, error) {
	head, err := r.repository.Head()
	if err != nil {
		defaultContributorCommits := make(map[Author][]object.Commit)
		return defaultContributorCommits, fmt.Errorf("GetAuthorsByCommits: unable to get the repository head: %w", err)
	}

	commits := make([]*object.Commit, 0)
	cIter, _ := r.repository.Log(&git.LogOptions{From: head.Hash(), Order: git.LogOrderCommitterTime})
	_ = cIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	authorCommits := r.getAuthorsByCommits(commits)

	return authorCommits, nil
}

// GetLicense attempts to determine the license type of the repository.
func (r *RepoReader) GetLicense() (string, error) {
	wt, err := r.repository.Worktree()
	if err != nil {
		return "", fmt.Errorf("GetLicense: unable to get the worktree from the repository: %w", err)
	}

	fs := wt.Filesystem

	license, err := r.getLicenseFromRoot(fs)
	if err != nil {
		return "", fmt.Errorf("GetLicense: error getting license from root: %w", err)
	}

	return license, err
}

// ValidateRepository validates the given repository and returns the head of the repository if valid and a nil error.
// If the repository is invalid a nil head reference and a non-nil error are returned.
func ValidateRepository(repo *git.Repository) (*plumbing.Reference, error) {
	if repo == nil {
		return nil, fmt.Errorf("ValidateRepository: received a nil repository")
	}
	if repo.Storer == nil {
		return nil, fmt.Errorf("ValidateRepository: invalid repository - Storer is nil")
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("ValidateRepository: received a repository without a head: %w", err)
	}
	if head == nil {
		return nil, fmt.Errorf("ValidateRepository: received a repository without a head")
	}

	return head, nil
}
