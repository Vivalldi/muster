
# Muster

A golang tool to automate the creation of Must funcs


[![Go Reference](https://pkg.go.dev/badge/github.com/vivalldi/muster.svg)](https://pkg.go.dev/github.com/vivalldi/muster)


## Usage/Examples

Given this snippet,

```go
package person

func DoGood(name str) error {
	return errors.New("bad")
}
```
running this command

`muster -func=DoGood`

in the same directory will create the file dogood_must.go, in package person,
containing a definition of

```go
func MustDoGood(name string)
```

That method will panic if the result of DoGood is not nil.

### Go Generate

Add this snipet in a file to use with go generate
```go
//go:generate muster -func=DoGood
```