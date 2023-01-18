package reporeader

import (
	"fmt"
	"os"

	"github.com/go-enry/go-license-detector/v4/licensedb"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
)

// GetLicense attempts to determine the license type of the repository.
func GetLicense(repo *git.Repository) (string, error) {
	head, err := ValidateRepository(repo)
	if err != nil || head == nil {
		return "", fmt.Errorf("GetLicense: received an invalid repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("GetLicense: received an invalid repository worktree: %w", err)
	}
	fs := worktree.Filesystem

	license, err := getLicenseFromRoot(fs)
	if err != nil {
		return "", fmt.Errorf("GetLicense: error getting license from root: %w", err)
	}

	return license, err
}

func getLicenseFromRoot(fs billy.Filesystem) (string, error) {
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
	best := licenseConfidence{}
	for license, confidence := range results {
		if best.confidence < confidence {
			best.license = license
			best.confidence = confidence
		}
	}

	return best.license, nil
}
