# The King's Hand

The King's Hand, or kh for short, is a tool for organizing and executing
shellish scripts written in Go.  As the name suggests, hand gets common tasks
done for you w/out fuss and rarely with errors.

In essence, King's Hand make is easy to write small scripts in Go rather
than x scripting language. Note that King's Hand does not dynamically
execute go code like [gorun](https://wiki.ubuntu.com/gorun). All scripts
must be compiled to go binaries using the `kh update` command. More on this
later.

King's Hand is partially inspired by [sub](https://signalvnoise.com/posts/3264-automating-with-convention-introducing-sub)
from the great folks at Basecamp. The abbreviated name of King's Hand, kh,
is also an oblique reference to the great computer scientist and author
[Brian KernigHan](https://en.wikipedia.org/wiki/Brian_Kernighan).

## Installation

King's Hand assumes that you have a working Go development environment. First
execute,

```
go get -u github.com/bryanwb/kh
```

To install the default fingers and initialize your ~/.kh directory

```
kh init
```


## Rationale

If you're like me, you have to write a fair number of shell scripts as part of
your daily work as a developer or sysadmin.  It is tempting to write those
scripts in Bash but you think better of it, as Bash is a fucking mess that is
completely unmanageable once the script is longer than one hundred
lines. Further, managing command-line flags to bash scripts is a nightmare. You
could write those scripts in a higher-level language like Python or Ruby but
something still isn't right.  Your scripts have zero type safety and debugging
them is a chore for finding even minor typos.

Just as important, I want to organize my scripts logically as subcommands. Once
I have written scripts, it can be very hard for me to find them again later and
even recall how they work. I need an organizing structure for these tasks.
For our first example, let's write a bunch of scripts for git. Our top-level
program, let's call it `kh` has a subcommand `git` off of which all git-related
scripts hang.

```
kh git gerrit-hook  # download the gerrit pre-commit hook into the current project
kh git add-ignores  # add commonly ignored file globs to .gitignore
```

GPG is another program that I cannot use w/out looking up a cheatsheet despite using for several years.

```
kh gpg decrypt foo.asc  # decrypted contents of foo.asc are output to foo.asc.plain
kh gpg encrypt foo.plain # encrypt contents writtent to foo.plain.gpg w/out overwriting original!
```

Note that the above scripts might be better accomplished through scripting vim,
emacs, or sublimeText. However, in my experience there is zero consistency in
editor usage across a development team.

Let's call these subcommands **fingers** rather than scripts so we don't
confuse them with Bash.

I have been interested in writing shell(ish) scripts and further writing tools
that are easy to up-to-date in Google's Go programming language. The primary
benefit here is that you can use Go's higher level tooling and libraries AND
take advantage of static typing to catch common errors.

However, using Go for this purpose presents a couple problems. Firstly, where
do all these fingers(plugins) live? Go is pretty inflexible in how it expects
your code to be organized and further it can't really be used to execute code
on the fly. In fact, we don't really want or need our code to be executed on
the fly.

This tool requires that you have a code organization for Go present on your
machine.

User-defined fingers live in ~/.kh/

```                          
~/.kh/
          git/
              main.go
          gpg/
              main.go
          ruby/
              main.go      # your awesome scripts for manipulating ruby-related stuff
          python/
                main.go    # your awesome scripts for manipulating python-related stuff 
```

These plugins are not dynamically loaded! In fact, to use them you must first
update the hand binary. To do this just execute `kh update`. After updating, the
new fingers (plugins) will be avaiable for use.


## What's Missing

It would be great if King's Hand could fetch and install additional fingers. Haven't figured out how
to do that yet. Perhaps using `go get`?

How to support shell completions?

I would love some feedback on these ideas! Please let me know of any gotchas i have not considered
