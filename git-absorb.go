package main

import (
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type Args struct {
	Dryrun bool `short:"d" long:"dry-run" description:"Show the commit(s) the outstanding changes could successfully be absorbed into."`
	Squash bool `short:"s" long:"squash" description:"squash instead of fixup the changes into the relevant commit. NOT IMPLEMENTED"`
	Target struct {
		Commit string `positional-arg-name:"target-commit"`
	} `positional-args:"yes"`
}

var args Args

func parse_args() *Args {
	if _, err := flags.Parse(&args); err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			exit_ok()
		}
		exit(3, "%v", err)
	}

	if args.Target.Commit != "" {
		sha1, err := expand_ref(args.Target.Commit)
		if err != nil {
			exit(2, err.Error())
		}
		args.Target.Commit = sha1
	}
	return &args
}

func main() {
	args := parse_args()

	files, _ := files_to_absorb()
	if len(files) == 0 {
		exit(0, "Nothing to do...")
	}

	var target_commit string = args.Target.Commit
	if target_commit == "" {
		candidate_commits := successful_absorb_commits(args)

		if len(candidate_commits) != 1 {
			exit(1, "%d candidate commits found on this branch:\n%s",
				len(candidate_commits), strings.Join(human_commits(candidate_commits), "\n"))
		}
		target_commit = candidate_commits[0]
	}

	if !args.Dryrun {
		do_absorb(target_commit, args)
	}
	exit(0, "Absorbed changes into commit '%s'", human_commit(target_commit))
}
