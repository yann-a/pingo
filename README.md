# pingo
A pi-calculus interpreter written in Go

<img src="media/pingo.png" align="right" alt="pingo" width="250"/>

## Requirements
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

## Features
The available options are :
+ `-showsrc` to show the parsed input before running it

The executable reads from a file if provided, or stdin if none was provided or reading from the file failed, until EOF (or Ctrl+D in command line) is met