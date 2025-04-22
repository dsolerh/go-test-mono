package publisher

import (
	"fmt"
	"os/exec"
	"strings"
)

func CommitAndTagChanges(pmap PackagesMap, packagesName []string) error {
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
		packageVersion := pmap[packageName].Version // TODO: increase the version by semantic versioning
		// tag the commit
		tag := fmt.Sprintf("%s/%s", packageName, packageVersion)
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

func PushChanges() error {
	// push the changes
	output, err := exec.Command("git", "push", "--follow-tags").CombinedOutput()
	if err != nil {
		return fmt.Errorf("error pushing the changes: %w, output: %s\n", err, output)
	}
	return nil
}
