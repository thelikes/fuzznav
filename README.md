# fuzznav

Print a list of all known endpoints and the wordlist run against them.

---

### Usage

***Highly dependent on the filename***

1. Take a list of files from stdin
2. Parse the name of the file to collect the wordlist and the target
3. Print a unique list of endpoints along with the each wordlist the endpoint was found in

```
gobuster-directory-list-2.3-small.txt-https.jss-dev.getpwnd.xyz.txt
[ type ] [        wordlist          ] [          targ        ]
```

---

### TODO
- [ ] add colors to eps mode output
- [ ] parse the Target from the results
