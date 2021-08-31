// Muster is a tool to automate the creation of Must funcs. Must funcs panic
// if the last result of a method is not nil. Muster should only be used on
// funcs/methods that return an error as the last result.
//
// Muster works on normal funcs as well as receiver funcs (methods).
//
// For example, given this snippet,
//
//	package person
//
//	func DoGood(name str) error {
//		return errors.New("bad")
//	}
//
// running this command
//
//	muster -func=DoGood
//
// in the same directory will create the file dogood_must.go, in package person,
// containing a definition of
//
//	func MustDoGood(name string)
//
// That method will panic if the result of DoGood is not nil.
//
// Typically this process would be run using go generate, like this:
//
//	//go:generate muster -func=DoGood
//
// With no arguments, it processes the package in the current directory.
// Otherwise, the arguments must name a single directory holding a Go package
// or a set of Go source files that represent a single Go package.
//
// The -func flag accepts a comma-separated list of types so a single run can
// generate methods for multiple types. The default output file is f_must.go,
// where f is the lower-cased name of the first func listed. It can be overridden
// with the -output flag.
package main // import "github.com/vivalldi/muster"
