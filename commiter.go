package main

import "os"

func do_absorb(args *Args) error {
	os.Setenv("GIT_EDITOR", "true")

	cleanup, undo := commit_changes(args.Target.Commit)

	if err := rebase_to_ref(args.Target.Commit); err != nil {
		if !args.Force {
			rebase_abort()
			undo()
			cleanup()
		}
		return err
	}

	cleanup()
	return nil
}
