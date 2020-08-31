我是光年实验室高级招聘经理。
我在github上访问了你的开源项目，你的代码超赞。你最近有没有在看工作机会，我们在招软件开发工程师，拉钩和BOSS等招聘网站也发布了相关岗位，有公司和职位的详细信息。
我们公司在杭州，业务主要做流量增长，是很多大型互联网公司的流量顾问。公司弹性工作制，福利齐全，发展潜力大，良好的办公环境和学习氛围。
公司官网是http://www.gnlab.com,公司地址是杭州市西湖区古墩路紫金广场B座，若你感兴趣，欢迎与我联系，
电话是0571-88839161，手机号：18668131388，微信号：echo 'bGhsaGxoMTEyNAo='|base64 -D ,静待佳音。如有打扰，还请见谅，祝生活愉快工作顺利。

# git-absorb

A git command inspired by the mercurial command of the same name described
[here](https://groups.google.com/forum/#!topic/mozilla.dev.version-control).

The gist of this command is to automagically fixup or squash uncommitted (though
possibly staged) modifications into the right ancestor commit (or a user
specified commit) in a working branch with no user interaction.

The common use case or workflow is for e.g. to modify commits in response to
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
