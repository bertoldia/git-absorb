package main

import "os"

func doAbsorb(args *Args) error {
	os.Setenv("GIT_EDITOR", "true")

	cleanupFn, recoverFn := CommitChanges(args.Target.Commit)

	if err := rebaseToRef(args.Target.Commit); err != nil {
		if args.NoRecover {
			return err
		}
		rebase_abort()
		recoverFn()
		cleanupFn()
		return err
	}

	cleanupFn()
	return nil
}
