package main

import (
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type Args struct {
	PrintCandidates bool `short:"p" long:"print-candidates" description:"Print the SHA1 of all candidate commits (i.e. all commits in the branch) in a human readable format."`
	Machine         bool `short:"m" long:"machine-parsable" description:"Print the candidate commits in a machine-parsable format (i.e. just the full SHA1)."`
	Force           bool `short:"f" long:"force" description:"Do not undo the absorb attempt if it failed (e.g because of a merge conflict)."`
	Target          struct {
		Commit string `positional-arg-name:"target-commit" description:"The SHA1 of the commit into which outstanding changes should be absorbed."`
	} `positional-args:"yes"`
}

var args Args
var parser = flags.NewParser(&args, flags.Default)
var help = `Absorb (i.e. merge) outstanding changes into the specified commit.`

func init() {
	parser.LongDescription = strings.Replace(help, "\n", " ", -1)
}

func print_candidates(machine bool) {
	if machine {
		exit(0, "%s", strings.Join(commits_in_branch(), " "))
	}
	exit(0, "%s", strings.Join(human_commits(commits_in_branch()), "\n"))
}

func parse_args() {
	if _, err := parser.Parse(); err != nil {
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
	} else if !args.PrintCandidates {
		exit(5, "One of --target-commit or --print-candidates is required.")
	}
}

func main() {
	parse_args()

	wd_is_git_repo()

	if args.PrintCandidates {
		print_candidates(args.Machine)
	}

	files, _ := files_to_absorb()
	if len(files) == 0 {
		exit(0, "Nothing to do...")
	}

	if err := do_absorb(&args); err != nil {
		exit(4, err.Error())
	}
	exit(0, "Successfully Absorbed changes into commit '%s'",
		human_commit(args.Target.Commit))
}
