package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func exit_ok() {
	exit(0, "")
}

func exit(code int, format string, args ...interface{}) {
	if len(format) > 0 {
		if !strings.HasSuffix(format, "\n") {
			format = format + "\n"
		}
		fmt.Printf(format, args...)
	}
	os.Exit(code)
}

func git_cmd(args ...string) []string {
	return exec_cmd("git", args...)
}

func exec_cmd(cmd string, args ...string) []string {
	res, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Fatal("'"+cmd+" "+strings.Join(args, " ")+"' failed. ", err)
	}
	if len(res) == 0 {
		return make([]string, 0)
	}
	return strings.Split(strings.Trim(string(res), " \n"), "\n")
}

func list_files_modified_by_commit(sha1 string) []string {
	return git_cmd("diff-tree", "--no-commit-id", "--name-only", "-r", sha1)
}

// Return the sha1 of all commits in the current branch.
func list_commits_in_branch() []string {
	cb := git_cmd("rev-parse", "--abbrev-ref", "HEAD")[0]
	us := git_cmd("rev-parse", "--abbrev-ref", cb+"@{upstream}")[0]
	mb := git_cmd("merge-base", us, cb)[0]
	return git_cmd("rev-list", mb+".."+cb)
}

// returns (possibly partially) staged files
func list_staged_files() []string {
	return git_cmd("diff-index", "--cached", "--name-only", "HEAD", "--")
}

// returns all dirty files (staged or otherwise)
func list_dirty_files() []string {
	return git_cmd("diff-index", "--name-only", "HEAD", "--")
}

func list_modified_files() ([]string, bool) {
	staged := list_staged_files()
	if len(staged) > 0 {
		return staged, true
	}
	return list_dirty_files(), false
}

func changes_staged() bool {
	_, err := exec.Command("git", "diff-index", "--cached", "--quiet", "HEAD", "--").Output()
	return err != nil
}

func stash() {
	git_cmd("stash", "--quiet")
}

func stash_pop() {
	git_cmd("stash", "pop", "--quiet")
}

func ensure_valid_sha1(sha1 string, commits []string) (string, error) {
	for _, c := range commits {
		if strings.HasPrefix(c, sha1) {
			return c, nil
		}
	}
	return "", errors.New(sha1 + " not in [" + strings.Join(commits, " ") + "].")
}

// Turn a slice into a set(ish)
func set(a []string) map[string]bool {
	res := make(map[string]bool, len(a))
	for _, val := range a {
		res[val] = true
	}
	return res
}
