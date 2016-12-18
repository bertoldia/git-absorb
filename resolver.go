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

func find_possible_merge_commits() []string {
	candidate_commits := list_commits_in_branch()
	if len(candidate_commits) < 1 {
		return candidate_commits
	}
	return filter_candidate_commits(candidate_commits)
}
