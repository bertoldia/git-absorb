package main

import (
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type Args struct {
	Force           bool `short:"f" long:"force" description:"Perform the absorb operation instead of just showing the commit(s) the outstanding changes could successfully be absorbed into."`
	Squash          bool `short:"s" long:"squash" description:"Squash instead of fixup the changes into the relevant commit. NOT IMPLEMENTED"`
	PrintCandidates bool `short:"p" long:"print-candidates" description:"Print the sha1 of all candidate commits (i.e. all commits in the branch) in a human readable format."`
	Machine         bool `short:"m" long:"machine-readable" description:"Print the candidate commits in a machine-parsable format."`
	Target          struct {
		Commit string `positional-arg-name:"target-commit"`
	} `positional-args:"yes"`
}

var args Args

func print_candidates(machine bool) {
	if machine {
		exit(0, "%s", strings.Join(commits_in_branch(), " "))
	}
	exit(0, "%s", strings.Join(human_commits(commits_in_branch()), "\n"))
}

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

	wd_is_git_repo()

	if args.PrintCandidates {
		print_candidates(args.Machine)
	}

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

	if args.Force {
		do_absorb(target_commit, args)
	}
	exit(0, "Absorbed changes into commit '%s'", human_commit(target_commit))
}
