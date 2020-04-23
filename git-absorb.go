package main

import (
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type Args struct {
	PrintCandidates bool `short:"p" long:"print-candidates" description:"Print the SHA1 of all candidate commits (i.e. all commits in the branch) in a human readable format."`
	Machine         bool `short:"m" long:"machine-parsable" description:"Print the candidate commits in a machine-parsable format (i.e. just the full SHA1)."`
	NoRecover       bool `short:"n" long:"no-recover" description:"Do not undo the absorb attempt if it failed (e.g because of a merge conflict)."`
	Force           bool `short:"f" long:"force" description:"Attempt the absorb even if the target commit is not from the current working-set (branch)"`
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
		Exit(0, "%s", strings.Join(CommitsInBranch(), " "))
	}
	Exit(0, "%s", strings.Join(HumanCommits(CommitsInBranch()), "\n"))
}

func parse_args() {
	if _, err := parser.Parse(); err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			ExitOK()
		}
		Exit(3, "%v", err)
	}

	if args.Target.Commit != "" {
		sha1, err := expandRef(args.Target.Commit)
		if err != nil {
			Exit(2, err.Error())
		}
		args.Target.Commit = sha1
	} else if !args.PrintCandidates {
		Exit(5, "One of --target-commit or --print-candidates is required.")
	}
}

func main() {
	parse_args()

	EnsureCwdIsGitRepo()

	if args.PrintCandidates {
		print_candidates(args.Machine)
	}

	if !CommitIsInWorkingSet(args.Target.Commit) && !args.Force {
		Exit(1, "Commit %s is not in working-set.", args.Target.Commit)
	}

	files, _ := FilesToAbsorb()
	if len(files) == 0 {
		Exit(0, "Nothing to do...")
	}

	if err := doAbsorb(&args); err != nil {
		Exit(4, err.Error())
	}

	Exit(0, "Successfully Absorbed changes into commit '%s'",
		human_commit(args.Target.Commit))
}
