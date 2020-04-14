package main

type expr interface{}

type parallel []expr
type send struct { // send a value on a given channel
  channel string
  value expr
}
type receiveThen struct { // receive a value then execute a process
  channel string
  received string
  then expr
}
type lambda struct { // used to define anonymous channels
  channel string
  then expr
}
type print expr

type skip int
type constant int
type channel string
