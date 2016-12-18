package main

import "os"

func do_absorb(args *Args) {
	os.Setenv("GIT_EDITOR", "true")

	if !changes_staged() {
		git_cmd("add", "-u")
	}

	var action string = "--fixup"
	git_cmd("commit", action, args.target_commit, "--no-edit")

	if still_dirty := list_dirty_files(); len(still_dirty) != 0 {
		stash()
		defer stash_pop()
	}
	git_cmd("rebase", "-i", "--autosquash", args.target_commit+"~1")
}
