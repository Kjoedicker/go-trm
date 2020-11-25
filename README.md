# GO-TRM [![Build Status](https://travis-ci.com/Kjoedicker/go-trm.svg?branch=master)](https://travis-ci.com/Kjoedicker/go-trm)
A simple CLI trash manager written in Golang.

## About
Allows you to delete, restore and view deleted files.

## Installation

In your `.bashrc` define the following enviromental variables

```
export XDG_DATA_HOME=$HOME'/.local/share'
export XDG_TRASH_HOME=$XDG_HOME_DATA'/Trash/Trash'
```

This is the Freedesktop standard for where trash files are stored. A lot of programs use this standard, so ```GO-TRM``` follows suite.

```
go get github.com/Kjoedicker/go-trm
```

## Use

#### Deleting

By default ./trm will delete a file

```
$ ./go-trm <file_name>
```

#### Restoring 

```
$ ./go-trm -r <file_name>
```

#### Viewing

```
$ ./go-trm -l
```
