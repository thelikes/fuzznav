# fuzznav

A utility for parsing and visualizing the web server endpoints discovered with ffuf and the wordlists used against each target.

## Features
* Parse many, many ffuf files at once
* Show target FUZZ endpoints
* Show discovered endpoints
* Pretty print - full URIs, color, status, size, words, & lines

## Install

```
$ GO111MODULE=on go get -u -v github.com/thelikes/fuzznav
```

## Run

Use shell commands (`find`, `ls`, etc) to list ffuf JSON files and pipe the
output to `fuzznav`. Throw `e` or `t` for endpoints and targets, respectively. 

**TIP** Use the firefox addon [Open Multiple
URLs](https://addons.mozilla.org/en-US/firefox/addon/open-multiple-urls/), or
similar, to pop a list of results into the browser quickly.

### Basics

Print found endpoints:

```
$ ls
ffuf-example.com.json
ffuf-examplecorporate.com.json
ffuf-admin.example.com.json

# print all endpoints found for "example.com"
$ ls | grep example.com | fuzznav
http://example.com/.hta                   [Status: 403, Size: 274, Words: 20, Lines: 10]
http://example.com/admin                  [Status: 301, Size: 306, Words: 20, Lines: 10]
http://example.com/doc                    [Status: 301, Size: 304, Words: 20, Lines: 10]
http://example.com/index.html             [Status: 200, Size: 94, Words: 2, Lines: 9]
```

Print fuzzed targets:

```
$ ls
ffuf-victim.com.json
ffuf-victim.com-admin.json

# print all FUZZ targets
$ ls | grep victim.com | fuzznav targs
http://example.com/admin/FUZZ             common.txt
http://example.com/FUZZ                   common.txt,raft-small-files.txt
```

## Background
This tool's aim is to aid in the mapping of fuzzing efforts. Instead of keeping a mental representation of what endpoints were fuzzed with what wordlists, this tool will make it easy to visualize where a server has been fuzzed and with what. Additionally, as ffuf (helpfully) stores a lot of data about discovered results and scanning in general, it can be cumbersome to get just what you need from the resulting json - fuzznav makes it simple to extract and parse. 

## To Do

### General
- [x] integreate [cobra](https://github.com/spf13/cobra)

### Endpoints
- [x] color
- [ ] filters
- [ ] show file found in
- [ ] tree view

        example.com/
            /login
            /user
            /admin
                /manage
                /upload
        
### Targets
- [ ] handle extensions
- [ ] handle multi custerbomb

### Stats
- [ ] no. of requests (ffuf provide?)
- [ ] no. of sessions (ffuf provide?)
- [ ] scanning time (ffuf provide?)
- [ ] endpoints discovered


## Thanks
Special thanks to [joohoi](https://twitter.com/joohoi) for such an amazing tool.
