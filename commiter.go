package main

import "os"

func do_absorb(args *Args) {
	os.Setenv("GIT_EDITOR", "true")

	cleanup := commit_changes(args.target_commit)
	defer cleanup()

	git_cmd("rebase", "-i", "--autosquash", args.target_commit+"~1")
}
