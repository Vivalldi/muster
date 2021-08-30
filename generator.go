package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	PackagesMode = packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo
)

type Generator struct {
	buf bytes.Buffer      // Accumulated output.
	pkg *packages.Package // Package we are scanning
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode:       PackagesMode,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatalf("parsePackage: load %s: %s", patterns, err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	g.addPackage(pkgs[0])
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = pkg
}

// generate produces the Must method for the named function.
func (g *Generator) generate(funcName string) {
	var decl *ast.FuncDecl

	// Find our function.
	for _, file := range g.pkg.Syntax {
		ast.Inspect(file, func(node ast.Node) bool {
			if n, ok := node.(*ast.FuncDecl); ok {
				if n.Name.Name == funcName {
					decl = n
					return false
				}
			}
			return true
		})
	}

	if decl == nil {
		log.Fatalf("couldn't find declaration of %s", funcName)
	}

	// Reusable variables
	mustName := "Must" + funcName

	// Add function comment
	g.Printf("// %[1]s is like %[2]s except it panics if %[2]s errors.\n", mustName, funcName)

	// Create the function declaration.
	g.Printf(
		"func %s %s(%s) (%s) { \n", // Receiver, Name, Params, Result
		FormatReceiver(decl.Recv, NeedName|NeedType),
		mustName,
		FormatParams(decl.Type.Params, NeedName|NeedType),
		FormatResults(decl.Type.Results, NeedType),
	)

	// Create the function body
	g.Printf(
		"\t%s := %s%s(%s)\n", // Result, Receiver, Name, Params
		FormatResults(decl.Type.Results, NeedName|NeedError),
		FormatReceiver(decl.Recv, NeedName),
		funcName,
		FormatParams(decl.Type.Params, NeedName),
	)

	// Check for err
	g.Printf("\tif err != nil {\n")
	g.Printf("\t\tpanic(err)\n")
	g.Printf("\t}\n")

	// Return results
	g.Printf("\treturn %s\n", FormatResults(decl.Type.Results, NeedName))
	g.Printf("}")
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() []byte {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		return g.buf.Bytes()
	}
	return src
}
