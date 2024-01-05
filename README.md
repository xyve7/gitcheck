# gitcheck
Checks for uncommited git repositories on your local machine.  

## Basic Usage
```bash
$ gitcheck -user <NAME> -root <ROOT>
```

## How it works
Probably the worst way I could do this, but oh well.  
All it does is parse the remote links (assumes https, will support ssh later) and finds the username to match the one you passed.  
I plan to have a more robust way to handle it later.  
Then, it runs the `git status` command and checks if there are any uncommited changes.  

## Plans
Seriously speed this up

## Compilation

Requirements:
    - go
    - git

Clone this repository:  
```bash
$ git clone https://github.com/xyve7/gitcheck.git
$ cd gitcheck
$ go build
```
