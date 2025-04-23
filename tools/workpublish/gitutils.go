package workpublish

import (
	"fmt"
	"os/exec"
	"strings"
)

// HasUncommittedChanges checks if the git repository at the given path
// has any uncommitted changes.
func HasUncommittedChanges(repoPath string) (bool, error) {
	// Run git status --porcelain
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("error running git status: %w", err)
	}

	// If the output is not empty, there are uncommitted changes
	return len(strings.TrimSpace(string(output))) > 0, nil
}

func CommitAndTagChanges(pmap PackagesMap, packagesName []string, versionUpdater func(string) string) error {
	// add the changes
	baseArgs := []string{"add", "--"}
	fullArgs := append(baseArgs, packagesName...)
	fullArgs = append(fullArgs, "go.work")
	output, err := exec.Command("git", fullArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding package to git: %w, output: %s\n", err, output)
	}

	// commit the changes
	var plural = func() string {
		if len(packagesName) > 1 {
			return "s"
		}
		return ""
	}
	message := fmt.Sprintf(`"ci: publish package%s %s"`, plural(), strings.Join(packagesName, ","))
	output, err = exec.Command("git", "commit", "-m", message).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error commiting the changes: %w, output: %s\n", err, output)
	}

	for _, packageName := range packagesName {
		pmap[packageName].Version = versionUpdater(pmap[packageName].Version)
		// tag the commit
		tag := fmt.Sprintf("%s/%s", packageName, pmap[packageName].Version)
		output, err = exec.Command("git", "tag", "--", tag, "HEAD").CombinedOutput()
		if err != nil {
			return fmt.Errorf("error tagging the changes: %w, output: %s\n", err, output)
		}
	}

	return nil
}

func CleanUpCommit(packagesName []string) error {
	// add the changes
	baseArgs := []string{"add", "--"}
	fullArgs := append(baseArgs, packagesName...)
	fullArgs = append(fullArgs, "go.work")
	output, err := exec.Command("git", fullArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error adding package to git: %w, output: %s\n", err, output)
	}

	// commit the changes
	output, err = exec.Command("git", "commit", "-m", `"ci: cleanup publish"`).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error commiting the changes: %w, output: %s\n", err, output)
	}

	return nil
}

func PushChanges(pmap PackagesMap, packagesName []string) error {
	// push the changes
	output, err := exec.Command("git", "push").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error pushing the changes: %w, output: %s\n", err, output)
	}

	tags := make([]string, 0, len(packagesName))
	for _, packageName := range packagesName {
		tags = append(tags, fmt.Sprintf("%s/%s", packageName, pmap[packageName].Version))
	}
	baseArgs := []string{"push", "origin"}
	fullArgs := append(baseArgs, tags...)
	output, err = exec.Command("git", fullArgs...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error tagging the changes: %w, output: %s\n", err, output)
	}
	return nil
}
