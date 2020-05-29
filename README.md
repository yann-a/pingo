# pingo
A pi-calculus interpreter written in Go
It also accepts input in `lambda + ref` that's translated to pi before execution

<img src="media/pingo.png" align="right" alt="pingo" width="250"/>

## Requirements
To compile, **your Go version must be greater than 1.14**.
To verify that, type `go version`

The parsing is done using goyacc.
To install it, run one of the following :
~~~
export GOPATH=$HOME/gocode
export PATH=$PATH:$GOPATH/bin
go get -u golang.org/x/tools/cmd/goyacc
~~~
or
~~~
sudo apt install golang-golang-x-tools
~~~

## Run
+ compile with `make`
+ run with `./pingo [options] [file]`

## Syntax
For pi-calculus:
+ `0` : the null instruction
+ `P | Q` : run P and Q in parallel
+ `^n m` : send m on channel n
+ `n(m).P` : read m on channel n then run P
+ `(n)P` : create a new (private) channel called n then run P
+ `print x` : print the integer x
+ `print x; P` : print the integer x then run P
+ `!n(m).P` : read m on channel n then in parrallel:
	+ run P
	+ run `!n(m).P` (the same instruction)
+ `[x=y]P` or `[x!=y]P` : run P if the condition between brackets is met
+ `a(x).P + b(x).Q` : non-deterministic choice : we go on with P if we read from `a` first, or Q if we read from `b`

For lambda-calculus+ref+"ocaml":
+ `fun x -> l` : the function that associates l to the variable x
+ `m n` : apply the function m to argument n
+ `new r = v in l` : initialize r to a reference containing v, then execute l
+ `x <- a` (read): assign the value contained in reference a to x. Note that a can also be an expression that returns a reference.
+ `a := x` (write): assign the value x to ref a. Note that a can also be an expression that returns a reference.
+ `x <- a := v` (swap): evaluate v, assign the value contained in a to x and replace it by v. Note that a can also be an expression that returns a reference.
+ `!r` : get the value contained in ref r
+ `let x = v in l` : binds the name x to v in l

## Reserved names ##
Some names are reserved in lambda. They are listed here, and forbidden by the parser anyway :

## Running the tests

Pyyaml is required:
```
pip3 install pyyaml
```

Then run:
```
make tests
```

## Features
The available options are :
+ `-showsrc` to show the parsed input
+ `-outcode` to show the code that's going to be executed (different from showsrc only if -lambda is enabled)
+ `-lambda` to parse lambda-calculus code and translate it into pi-calculus before execution

The executable reads from a file if provided, or stdin if none was provided or reading from the file failed, until EOF (or Ctrl+D in command line) is met
