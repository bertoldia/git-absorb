package main

import "os"

type Args struct {
	target_commit string
}

func parse_args() *Args {
	var args = &Args{}

	if len(os.Args) > 1 {
		sha1, err := ensure_valid_sha1(os.Args[1], list_commits_in_branch())
		if err != nil {
			exit(1, err.Error())
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

	do_absorb(args)
	exit_ok()
}
