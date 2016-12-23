package main

import "os"

// Test if super contains all of sub
func _contains(super []string, sub []string) bool {
	_super := set(super)

	for _, e := range sub {
		if _, ok := _super[e]; !ok {
			return false
		}
	}
	return true
}

func _filter_by_modified_files(commits []string) []string {
	result := make([]string, 0)
	files_to_absorb, _ := files_to_absorb()

	for _, commit := range commits {
		modified_by_commit := files_modified_by_commit(commit)
		if _contains(modified_by_commit, files_to_absorb) {
			result = append(result, commit)
		}
	}
	return result
}

func _filter_commits(commits []string) []string {
	// Is this a good idea?
	return _filter_by_modified_files(commits)
}

// The commits in the working branch that the outstanding changes could
// potentially be absorbed into. This is a subset of all commits in the working
// branch.
func _candidate_absorb_commits() []string {
	candidate_commits := commits_in_branch()
	if len(candidate_commits) < 1 {
		return candidate_commits
	}
	return _filter_commits(candidate_commits)
}

func _successful_absorb_commits_by_brute_force(commits []string, args *Args) []string {
	os.Setenv("GIT_EDITOR", "true")

	var valid_commits = make([]string, 0)
	cleanup, undo := commit_changes(commits[0])
	defer undo()
	defer cleanup()
	head_reflog := reflog_head()

	for _, commit := range commits {
		update_commit_msg(commit)
		if err := rebase_to_ref(commit); err != nil {
			rebase_abort()
		} else {
			valid_commits = append(valid_commits, commit)
			reset_hard_to(head_reflog)
		}
	}
	return valid_commits
}

// The commits in the working branch that the outstanding changes could
// successfully be absorbed (merged) into (i.e. without causing merge conflicts
// at any point in the rebase).
func successful_absorb_commits(args *Args) []string {
	candidate_commits := _candidate_absorb_commits()
	if len(candidate_commits) == 0 {
		return candidate_commits
	}
	return _successful_absorb_commits_by_brute_force(candidate_commits, args)
}
