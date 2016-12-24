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
	return _parse_cmd_exec_result(res)
}

func _parse_cmd_exec_result(res []byte) []string {
	if len(res) == 0 {
		return make([]string, 0)
	}
	return strings.Split(strings.Trim(string(res), " \n"), "\n")
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

func files_modified_by_commit(sha1 string) []string {
	return git_cmd("diff-tree", "--no-commit-id", "--name-only", "-r", sha1)
}

// Return the sha1 of all commits in the current branch.
func commits_in_branch() []string {
	cb := current_branch()
	us := branch_upstream(cb)
	mb := merge_base(us, cb)
	return git_cmd("rev-list", mb+".."+cb)
}

// returns (possibly partially) staged files
func _staged_files() []string {
	return git_cmd("diff-index", "--cached", "--name-only", "HEAD", "--")
}

// returns all dirty files (staged or otherwise)
func _dirty_files() []string {
	return git_cmd("diff-index", "--name-only", "HEAD", "--")
}

// List either file with staged changes or, if no changes have been staged, all
// dirty files.
func files_to_absorb() ([]string, bool) {
	staged := _staged_files()
	if len(staged) > 0 {
		return staged, true
	}
	return _dirty_files(), false
}

func reset_hard_to(sha1 string) {
	git_cmd("reset", "--hard", sha1)
}

func reset_head_soft() {
	git_cmd("reset", "--soft", "HEAD~1")
}

func reset_head() {
	git_cmd("reset", "HEAD~1")
}

func reflog_head() string {
	return git_cmd("log", "-g", "--format=%H", "-n", "1")[0]
}

func are_changes_staged() bool {
	_, err := exec.Command("git", "diff-index", "--cached", "--quiet", "HEAD", "--").Output()
	return err != nil
}

func rebase_abort() {
	git_cmd("rebase", "--abort")
}

func rebase_to_ref(sha1 string) error {
	args := []string{"rebase", "-i", "--autosquash", sha1 + "~1"}
	_, err := exec.Command("git", args...).Output()
	return err
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
	return _parse_cmd_exec_result(res)[0], nil
}

func human_commit(sha1 string) string {
	return git_cmd("log", "--pretty=oneline", "--abbrev-commit", "-1", sha1)[0]
}

func human_commits(commits []string) []string {
	var result = make([]string, 0)
	for _, sha1 := range commits {
		result = append(result, human_commit(sha1))
	}
	return result
}

type cleanup_func func()
type undo_func func()

func noop() {}

func update_commit_msg(sha1 string) {
	git_cmd("commit", "--amend", "--fixup", sha1, "--no-edit")
}

// Commit uncommitted changes and return cleanup and undo functions. If any
// changes have been staged, only commit those and stash the remaining changes.
// In this case the cleanup operation is to stash pop and the undo operation is
// a soft reset on head. If none of the outstanding changed have been staged,
// commit them all. In this case the cleanup function is a noop and the undo
// operation is a (regular) reset of head.
func commit_changes(sha1 string) (cleanup_func, undo_func) {
	if !are_changes_staged() {
		git_cmd("add", "-u")
	}

	var action string = "--fixup"
	git_cmd("commit", action, sha1, "--no-edit")

	if still_dirty := _dirty_files(); len(still_dirty) != 0 {
		stash()
		return stash_pop, reset_head_soft
	}
	return noop, reset_head
}

// Turn a slice into a set(ish)
func set(a []string) map[string]bool {
	res := make(map[string]bool, len(a))
	for _, val := range a {
		res[val] = true
	}
	return res
}
