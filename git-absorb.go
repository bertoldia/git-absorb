package main

import "os"

type Args struct {
	target_commit string
}

func parse_args() *Args {
	var args = &Args{}

	if len(os.Args) > 1 {
		sha1, err := expand_ref(os.Args[1], list_commits_in_branch())
		if err != nil {
			exit(2, err.Error())
		}
		args.target_commit = sha1
	}
	return args
}

func main() {
	files, _ := list_modified_files()
	if len(files) == 0 {
		exit(0, "Nothing to do...")
	}

	args := parse_args()

	if args.target_commit != "" {
		do_absorb(args)
	} else {
		candidate_commits := find_possible_merge_commits()
		if len(candidate_commits) != 1 {
			exit(1, "%d candidate commits found on this branch.", len(candidate_commits))
		}
	}
	exit_ok()
}
