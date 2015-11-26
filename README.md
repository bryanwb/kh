# The Hand of the King

The Hand of the King, or hand for short, is a tool for organizing and executing
shellish scripts written in Go.  As the name suggests, hand gets common tasks
done for you w/out fuss and rarely with errors.

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
program, let's call it `h` has a subcommand `git` off of which all git-related
scripts hang.

```
h git gerrit-hook  # download the gerrit pre-commit hook into the current project
h git add-ignores  # add commonly ignored file globs to .gitignore
```

GPG is another program that I cannot use w/out looking up a cheatsheet despite using for several years. Here are 
a few examples of how I might script GPG with `h`.

```
h gpg decrypt foo.asc  # decrypted contents of foo.asc are output to foo.asc.plain
h gpg encrypt foo.plain # encrypt contents writtent to foo.plain.gpg w/out overwriting original!
```

Note that the above scripts might be better accomplished through scripting vim,
emacs, or sublimeText. However, in my experience there is zero consistency in
editor usage across a development team.

Let's call these subcommands **fingers** rather than scripts so we don't confuse them with Bash.

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

User-defined fingers live in ~/.hand/

```                          
~/.hand/
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
update the hand binary. To do this just execute `h update`. After updating, the
new fingers (plugins) will be avaiable for use.

I would love some feedback on these ideas! Please let me know of any gotchas i have not considered
