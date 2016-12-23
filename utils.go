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
	return parse_cmd_exec_result(res)
}

func current_branch() string {
	//return git_cmd("rev-parse", "--abbrev-ref", "HEAD")[0]
	return git_cmd("name-rev", "--name-only", "HEAD")[0]
}

func branch_upstream(branch string) string {
	return git_cmd("rev-parse", "--abbrev-ref", branch+"@{upstream}")[0]
}

func merge_base(branch, upstream string) string {
	return git_cmd("merge-base", upstream, branch)[0]
}

func parse_cmd_exec_result(res []byte) []string {
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
	cb := current_branch()
	us := branch_upstream(cb)
	mb := merge_base(us, cb)
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

// List either file with staged changes or, if no changes have been staged, all
// dirty files.
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

// Expand a ref (in the form of a shortened sha1 or other valid ref spec like
// HEAD) to the full sha1. Also useful for verifying the specified ref is valid.
func expand_ref(sha1 string) (string, error) {
	res, err := exec.Command("git", "rev-parse", "--verify", sha1).Output()
	if err != nil {
		return "", errors.New("'" + sha1 + "' is not a valid sha1.")
	}
	return parse_cmd_exec_result(res)[0], nil
}

type cleanup_func func()

func noop() {}

// Commit uncommitted changes and return a cleanup function. If any changes have
// been stages, only commit those and stash the remaining changes. In this case
// the cleanup operation is to stash pop the remaining changes. If none of the
// outstanding changed have been stages, commit them all. In this case the
// cleanup function is a noop.
func commit_changes(sha1 string) cleanup_func {
	if !changes_staged() {
		git_cmd("add", "-u")
	}

	var action string = "--fixup"
	git_cmd("commit", action, sha1, "--no-edit")

	if still_dirty := list_dirty_files(); len(still_dirty) != 0 {
		stash()
		return stash_pop
	}
	return noop
}

// Turn a slice into a set(ish)
func set(a []string) map[string]bool {
	res := make(map[string]bool, len(a))
	for _, val := range a {
		res[val] = true
	}
	return res
}
