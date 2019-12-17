To build a stripped binary -- removal of debug statements

go build -ldflags="-s -w"


Running tests
There are two ways. The easy one is to use the -run flag and provide a pattern matching names of the tests you want to run.

ex: $ go test -run NameOfTest. See the docs for more info.

The other way is to name the specific file, containing the tests you want to run:

$ go test foo_test.go
But there's a catch. This works well if

foo.go is package foo
foo_test.go is package foo_test and imports 'foo'.
If 'foo_test.go' and 'foo.go' are the same package (a common case), then you must name all other files required to build 'foo_test'. In this example it would be:

$ go test foo_test.go foo.go
I'd recommend to use the name pattern. Or, where/when possible, always run all package tests.


Example: Run JUST kwiselists_test.go
	go test -v kwiselists_test.go kwiselist.go interface.go

