package command

import (
	"log"
	"os/exec"
)

func GetGitFilesAndDirs(revision, repository string) string {
	cmd := exec.Command("git", "ls-tree", revision)
	cmd.Dir = repository
	files, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("git ls-tree fatal %s", err)
	}
	return string(files)
}

func GetGitFiles(revision, repository string) string {
	cmd := exec.Command("git", "ls-tree", "--name-only", revision)
	cmd.Dir = repository
	files, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("git ls-tree --name-only fatal", err)
	}
	return string(files)
}

func GetGitEmptyFile(revision, file, repository string) string {
	cmd := exec.Command("git", "log", `--pretty=format:"%cn %H"`, revision, "--", file)
	cmd.Dir = repository
	files, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("git log fatal ", err)
	}
	return string(files)
}

func GetCommit(revision, file, repository string) string {
	cmd := exec.Command("git", "blame", "--porcelain", revision, file)
	cmd.Dir = repository
	files, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("git blame fatal", err)
	}
	return string(files)
}

func GetGitCommitter(hash, repository string) string {
	cmd := exec.Command("git", "show", "--format=%cn", hash[0:8])
	cmd.Dir = repository
	comitter, err := cmd.Output()
	if err != nil {
		log.Fatal("git blame fatal", err)
	}
	return string(comitter)
}
