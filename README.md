# fuzznav

## What

A go script to aid in navigating various directory fuzzing output files. 

## Why

When fuzzing multiple domains, multiple targets, and using multiple wordlists, it
becomes cumbersome to eyeball output filenames in an attempt to identify which
domains and endpoints have been fuzzed with which wordlists. This tool aims to
make the endeavour of navigating fuzzing efforts easier.

## When (to use this)

Mainly, this tool should be used when multiple domains and endpoints have been
fuzzed with multiple wordlists and tools.

This tool also comes in handy as a catffuf utility, in that currently, ffuf
only outputs json, html, markdown, and csv. This tool can be used to print
the simple, grep-able stuff from a ffuf file (ie endpoint, status, length).

## Where (to use this)

I made this tool to help with my personal workflow. You'll have a hard time
using it if you do not match this workflow. Generally, I create a directory and
within it, store all fuzz output files (gobuster, wfuzz, ffuf). Sometimes,
I create sub-directories that are specific to one task of fuzzing.

When I want to see all discovered endpoints of a URL, I'll typically run
something along the lines of:

`find -name 'gobuster*' | grep admin.targ.com | xargs cat | sort -u`

This will find all output files made from fuzzing 'admin.targ.com' using
gobuster, and output all discovered endpoints, uniquely. As I have started to
use ffuf more and more, the above method to view discovered endpoints will not
work. Cat'ing a ffuf json blob does not make for readable output. You might
ask, why not use gron or jq? Because without a pretty long bash command to
parse json with gron/jq, I cannot output discovered endpoints AND grep on
size/status. 

Even more particular, the output file of both gobuster and ffuf is lacking.
Gobuster outputs only discovered endpoints and status/size, if you chose those
flags. This information does not disclose the wordlists used, nor the target
(URL) that was fuzzed. Ffuf does a little better, including the entire command
executed when running ffuf.

Therefore, extracting the Target and the Wordlist for a directory fuzz is
currently relying on the output file's naming scheme. For me, I name the output
file like this:

`tool-wordlist-target.txt` 

Fuzznav will parse out the Target and the Wordlist based on the output file's
name. To make use of the tool, your output files need to match this format. The
following bash will semi-automate this:

```
targ=https://admin.targ.com
list=/opt/SecLists/Discover/Web-Content/big.txt
ffuf -u $targ/FUZZ -w $list -o ffuf-$(basename $list)-$(echo $targ | sed 's/\//_/g' | sed 's/\:/_/g').txt
```

Resulting in:

`ffuf-big.txt-https___admin.targ.com.txt`

Its not pretty, but it gets the job done. Here is an example:

```
$ find | grep admin.targ.com
ffuf-raft-large-words.txt-https___admin.targ.com_management.txt
gobuster-raft-large-directories.txt-https___admin.targ.com_management.txt

$ find | grep admin.targ.com | fuzznav targs
https://admin.targ.com/                 big.txt, raft-large-words.txt
https://admin.targ.com/management/      raft-large-words.txt

$ find | grep admin.targ.com | fuzznav eps
/management/WEB-INF             (Status: 400)   [Size: 869]
/management/test.txt            (Status: 200)   [Size: 521]
/management/backup              (Status: 301)   [Size: 233]
/management/backup/sql          (Status: 301)   [Size: 233]
```

---

## Usage

Print a list of files to stdin and pipe it to fuzznav, selecting a mode.

Currently, two modes are supported: 
1. `targs` - Show targets and wordlists
2. `eps`   - Show endpoints

Example:

`ls fuzz-output/ | fuzznav targs`

### ***Highly dependent on the filename***

```
gobuster-directory-list-2.3-small.txt-https.jss-dev.getpwnd.xyz.txt
[ type ] [        wordlist          ] [          targ        ]
```

---

### TODO
- [ ] add colors to eps mode output
- [x] parse the Target from the results
- [ ] add wfuzz support
- [ ] add ffuf cat
- [ ] add ability to process multiple targets on input
