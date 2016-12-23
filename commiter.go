package main

import "os"

func do_absorb(args *Args) {
	os.Setenv("GIT_EDITOR", "true")

	cleanup, undo := commit_changes(args.target_commit)
	defer cleanup()

	if err := rebase_to_ref(args.target_commit); err != nil {
		rebase_abort()
		undo()
	}
}
