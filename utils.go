package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func ExitOK() {
	Exit(0, "")
}

func Exit(code int, format string, args ...interface{}) {
	if len(format) > 0 {
		if !strings.HasSuffix(format, "\n") {
			format = format + "\n"
		}
		fmt.Printf(format, args...)
	}
	os.Exit(code)
}

func gitCmd(args ...string) []string {
	return exec_cmd("git", args...)
}

func exec_cmd(cmd string, args ...string) []string {
	res, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Fatal("'"+cmd+" "+strings.Join(args, " ")+"' failed. ", err)
	}
	return parse_cmd_exec_result(res)
}

func parse_cmd_exec_result(res []byte) []string {
	if len(res) == 0 {
		return make([]string, 0)
	}
	return strings.Split(strings.Trim(string(res), " \n"), "\n")
}

func EnsureCwdIsGitRepo() {
	res, err := exec.Command("git", "rev-parse", "--show-toplevel").CombinedOutput()
	if err != nil {
		Exit(128, string(res))
	}
}

func CurrentBranch() string {
	//return gitCmd("rev-parse", "--abbrev-ref", "HEAD")[0]
	return gitCmd("name-rev", "--name-only", "HEAD")[0]
}

func Upstream(branch string) string {
	return gitCmd("rev-parse", "--abbrev-ref", branch+"@{upstream}")[0]
}

func MergeBase(branch, upstream string) string {
	return gitCmd("merge-base", upstream, branch)[0]
}

// Return the sha1 of all commits in the current branch.
func CommitsInBranch() []string {
	cb := CurrentBranch()
	us := Upstream(cb)
	mb := MergeBase(us, cb)
	return gitCmd("rev-list", mb+".."+cb)
}

// returns (possibly partially) staged files
func staged_files() []string {
	return gitCmd("diff-index", "--cached", "--name-only", "HEAD", "--")
}

// returns all dirty files (staged or otherwise)
func dirty_files() []string {
	return gitCmd("diff-index", "--name-only", "HEAD", "--")
}

// List either file with staged changes or, if no changes have been staged, all
// dirty files.
func FilesToAbsorb() ([]string, bool) {
	staged := staged_files()
	if len(staged) > 0 {
		return staged, true
	}
	return dirty_files(), false
}

func resetHeadSoft() {
	gitCmd("reset", "--soft", "HEAD~1")
}

func resetHead() {
	gitCmd("reset", "HEAD~1")
}

func areChangesStaged() bool {
	_, err := exec.Command("git", "diff-index", "--cached", "--quiet", "HEAD", "--").Output()
	return err != nil
}

func rebase_abort() {
	gitCmd("rebase", "--abort")
}

func rebaseToRef(sha1 string) error {
	args := []string{"rebase", "-i", "--autosquash", sha1 + "~1"}
	if msg, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		return errors.New(string(msg))
	}
	return nil
}

func stash() {
	gitCmd("stash", "--quiet")
}

func StashPop() {
	gitCmd("stash", "pop", "--quiet")
}

// Expand a ref (in the form of a shortened sha1 or other valid ref spec like
// HEAD) to the full sha1. Also useful for verifying the specified ref is valid.
func expandRef(sha1 string) (string, error) {
	res, err := exec.Command("git", "rev-parse", "--verify",
		strings.TrimRight(sha1, " ")).Output()
	if err != nil {
		return "", errors.New("'" + sha1 + "' is not a valid sha1.")
	}
	return parse_cmd_exec_result(res)[0], nil
}

func human_commit(sha1 string) string {
	return gitCmd("log", "--pretty=oneline", "--abbrev-commit", "-1", sha1)[0]
}

func HumanCommits(commits []string) []string {
	var result = make([]string, 0)
	for _, sha1 := range commits {
		result = append(result, human_commit(sha1))
	}
	return result
}

type cleanup_func func()
type recover_func func()

func noop() {}

func CommitIsInWorkingSet(sha1 string) bool {
	for _, s := range CommitsInBranch() {
		if sha1 == s {
			return true
		}
	}
	return false
}

// Commit uncommitted changes and return cleanup and undo functions. If any
// changes have been staged, only commit those and stash the remaining changes.
// In this case the cleanup operation is to stash pop and the undo operation is
// a soft reset on head. If none of the outstanding changed have been staged,
// commit them all. In this case the cleanup function is a noop and the undo
// operation is a (regular) reset of head.
func CommitChanges(sha1 string) (cleanup_func, recover_func) {
	if !areChangesStaged() {
		gitCmd("add", "-u")
	}

	var action string = "--fixup"
	gitCmd("commit", action, sha1, "--no-edit")

	if still_dirty := dirty_files(); len(still_dirty) != 0 {
		stash()
		return StashPop, resetHeadSoft
	}
	return noop, resetHead
}
