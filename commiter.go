package main

import "os"

func do_absorb(commit string, args *Args) {
	os.Setenv("GIT_EDITOR", "true")

	cleanup, undo := commit_changes(commit)
	defer cleanup()

	if err := rebase_to_ref(commit); err != nil {
		rebase_abort()
		undo()
	}
}
