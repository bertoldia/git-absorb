package main

// Test if super contains all of sub
func contains(super []string, sub []string) bool {
	_super := set(super)

	for _, e := range sub {
		if _, ok := _super[e]; !ok {
			return false
		}
	}
	return true
}

func filter_commits_by_referenced_file(commits []string) []string {
	result := make([]string, 0)
	files, _ := list_modified_files()

	for _, commit := range commits {
		modified_by_commit := list_files_modified_by_commit(commit)
		if contains(modified_by_commit, files) {
			result = append(result, commit)
		}
	}
	return result
}

func filter_candidate_commits(ccs []string) []string {
	return filter_commits_by_referenced_file(ccs)
}

func get_candidate_absorb_commits() []string {
	candidate_commits := list_commits_in_branch()
	if len(candidate_commits) < 1 {
		return candidate_commits
	}
	// Is this a good idea?
	return filter_candidate_commits(candidate_commits)
}

func find_possible_merge_commits(args *Args) []string {
	candidate_commits := get_candidate_absorb_commits()
	if len(candidate_commits) == 0 {
		return candidate_commits
	}
	return approach1(candidate_commits, args)
}

func approach1(commits []string, args *Args) []string {
	var valid_commits = make([]string, 0)
	cleanup := commit_changes(commits[0])
	head_reflog := get_head_ref()
	//defer restore_to_reflog_point(head_reflog)
	defer cleanup()

	for _, commit := range commits {
		update_commit_msg(commit)
		if err := rebase_to_ref(commit); err != nil {
			rebase_abort()
		} else {
			valid_commits = append(valid_commits, commit)
			restore_to_reflog_point(head_reflog)
		}
	}
	return valid_commits
}
