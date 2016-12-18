# git-absorb

A git command insipred by the mercurial command of the same name described
[here](
https://groups.google.com/forum/#!topic/mozilla.dev.version-control/nh4fITFlEMk).

The gist of this command is to automagically fixup or squash uncommitted (though
possibly staged) modifications into the right ancestor commit (or a user
specified commit) in a working branch with no user interaction.

The common use case or workflow is for e.g. to modifyg commits in response to
issues raised during a code review, or when you change your mind about the
content of existing commits in your working branch.

An alternate workflow for the above use-cases is to do an interactive rebase,
mark the relevant commits with (m)odify, make changes, then do git add + git
rebase --continue.

## Phase 1
* user must specify target commit into which changes should be absorbed.
* no squash, fixup only.

## Phase 2
* find (if it exists) the single commit that can cleanly (i.e. without merge
  conflicts) absorb the outstanding changes. Fail if more than one such commit
  exists.
* no squash, fixup only.

## Phase 3
* Add support for --squash.
