package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// save a repo and try to load it using stdout
func TestSaveAndLoadRepoStdout(t *testing.T) {
	runCmd := exec.Command(dockerBinary, "run", "-d", "busybox", "true")
	out, _, err := runCommandWithOutput(runCmd)
	errorOut(err, t, fmt.Sprintf("failed to create a container: %v %v", out, err))

	cleanedContainerID := stripTrailingCharacters(out)

	repoName := "foobar-save-load-test"

	inspectCmd := exec.Command(dockerBinary, "inspect", cleanedContainerID)
	out, _, err = runCommandWithOutput(inspectCmd)
	errorOut(err, t, fmt.Sprintf("output should've been a container id: %v %v", cleanedContainerID, err))

	commitCmd := exec.Command(dockerBinary, "commit", cleanedContainerID, repoName)
	out, _, err = runCommandWithOutput(commitCmd)
	errorOut(err, t, fmt.Sprintf("failed to commit container: %v %v", out, err))

	inspectCmd = exec.Command(dockerBinary, "inspect", repoName)
	before, _, err := runCommandWithOutput(inspectCmd)
	errorOut(err, t, fmt.Sprintf("the repo should exist before saving it: %v %v", before, err))

	saveCmdTemplate := `%v save %v > /tmp/foobar-save-load-test.tar`
	saveCmdFinal := fmt.Sprintf(saveCmdTemplate, dockerBinary, repoName)
	saveCmd := exec.Command("bash", "-c", saveCmdFinal)
	out, _, err = runCommandWithOutput(saveCmd)
	errorOut(err, t, fmt.Sprintf("failed to save repo: %v %v", out, err))

	deleteImages(repoName)

	loadCmdFinal := `cat /tmp/foobar-save-load-test.tar | docker load`
	loadCmd := exec.Command("bash", "-c", loadCmdFinal)
	out, _, err = runCommandWithOutput(loadCmd)
	errorOut(err, t, fmt.Sprintf("failed to load repo: %v %v", out, err))

	inspectCmd = exec.Command(dockerBinary, "inspect", repoName)
	after, _, err := runCommandWithOutput(inspectCmd)
	errorOut(err, t, fmt.Sprintf("the repo should exist after loading it: %v %v", after, err))

	if before != after {
		t.Fatalf("inspect is not the same after a save / load")
	}

	deleteContainer(cleanedContainerID)
	deleteImages(repoName)

	os.Remove("/tmp/foobar-save-load-test.tar")

	logDone("save - save a repo using stdout")
	logDone("load - load a repo using stdout")
}

// save a repo and try to load it using flags
func TestSaveAndLoadRepoFlags(t *testing.T) {
	runCmd := exec.Command(dockerBinary, "run", "-d", "busybox", "true")
	out, _, err := runCommandWithOutput(runCmd)
	errorOut(err, t, fmt.Sprintf("failed to create a container: %v %v", out, err))

	cleanedContainerID := stripTrailingCharacters(out)

	repoName := "foobar-save-load-test"

	inspectCmd := exec.Command(dockerBinary, "inspect", cleanedContainerID)
	out, _, err = runCommandWithOutput(inspectCmd)
	errorOut(err, t, fmt.Sprintf("output should've been a container id: %v %v", cleanedContainerID, err))

	commitCmd := exec.Command(dockerBinary, "commit", cleanedContainerID, repoName)
	out, _, err = runCommandWithOutput(commitCmd)
	errorOut(err, t, fmt.Sprintf("failed to commit container: %v %v", out, err))

	inspectCmd = exec.Command(dockerBinary, "inspect", repoName)
	before, _, err := runCommandWithOutput(inspectCmd)
	errorOut(err, t, fmt.Sprintf("the repo should exist before saving it: %v %v", before, err))

	saveCmdTemplate := `%v save -o /tmp/foobar-save-load-test.tar %v`
	saveCmdFinal := fmt.Sprintf(saveCmdTemplate, dockerBinary, repoName)
	saveCmd := exec.Command("bash", "-c", saveCmdFinal)
	out, _, err = runCommandWithOutput(saveCmd)
	errorOut(err, t, fmt.Sprintf("failed to save repo: %v %v", out, err))

	deleteImages(repoName)

	loadCmdFinal := `docker load -i /tmp/foobar-save-load-test.tar`
	loadCmd := exec.Command("bash", "-c", loadCmdFinal)
	out, _, err = runCommandWithOutput(loadCmd)
	errorOut(err, t, fmt.Sprintf("failed to load repo: %v %v", out, err))

	inspectCmd = exec.Command(dockerBinary, "inspect", repoName)
	after, _, err := runCommandWithOutput(inspectCmd)
	errorOut(err, t, fmt.Sprintf("the repo should exist after loading it: %v %v", after, err))

	if before != after {
		t.Fatalf("inspect is not the same after a save / load")
	}

	deleteContainer(cleanedContainerID)
	deleteImages(repoName)

	os.Remove("/tmp/foobar-save-load-test.tar")

	logDone("save - save a repo using -o")
	logDone("load - load a repo using -i")
}
