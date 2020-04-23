# pingo
A pi-calculus interpreter written in Go

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

For lambda-calculus+ref:
+ `fun x -> l` : the function that associates l to the variable x
+ `m n` : apply the function m to argument n
+ `ref v` : a reference that is initialized containing v
+ `x <- a` (read): assign the value contained in reference a to x
+ `a := x` (write): assign the value x to ref a
+ `x <- a := v` (swap): evaluate v, assign the value contained in a to x and replace it by v

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
+ `-outcode` to show the code that's going to be executed (different from showsrc only of a translation happened)
+ `-translation` to parse lambda-calculus code and translate it into pi-calculus before execution

The executable reads from a file if provided, or stdin if none was provided or reading from the file failed, until EOF (or Ctrl+D in command line) is met
