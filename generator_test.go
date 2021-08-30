package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
)

type GeneratorTest struct {
	name  string
	input string // input
	want  string // expected output
}

var tests = []GeneratorTest{
	{"Args", args_in, args_out},
	{"Results", results_in, results_out},
	{"NamedResults", named_results_in, named_results_out},
	{"Receiver", receiver_in, receiver_out},
	{"PtrReceiver", ptr_receiver_in, ptr_receiver_out},
}

func TestGenerator(t *testing.T) {
	is := is.New(t)
	dir, err := ioutil.TempDir("", "muster")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir)

	for _, test := range tests {
		is := is.NewRelaxed(t)

		// Some setup
		g := Generator{}
		input := "package test\n" + test.input
		file := test.name + ".go"
		absFile := filepath.Join(dir, file)
		is.NoErr(ioutil.WriteFile(absFile, []byte(input), 0644)) // create file without error please

		// Run the generator
		g.parsePackage([]string{absFile}, nil)
		g.generate(test.name)

		// Comparisons
		got := string(g.format())
		is.Equal(got, test.want)
	}
}

const args_in = `func Args(str string, n int) error {
	return nil
}`

const args_out = `// MustArgs is like Args except it panics if Args errors.
func MustArgs(str string, n int) {
	err := Args(str, n)
	if err != nil {
		panic(err)
	}
	return
}`

const results_in = `func Results() (string, int, error) {
	return "", 0, nil
}`

const results_out = `// MustResults is like Results except it panics if Results errors.
func MustResults() (string, int) {
	r0, r1, err := Results()
	if err != nil {
		panic(err)
	}
	return r0, r1
}`

const named_results_in = `func NamedResults() (str string, n int, err error) {
	return "", 0, nil
}`

const named_results_out = `// MustNamedResults is like NamedResults except it panics if NamedResults errors.
func MustNamedResults() (str string, n int) {
	str, n, err := NamedResults()
	if err != nil {
		panic(err)
	}
	return str, n
}`

const receiver_in = `func (t T) Receiver() error {
	return nil
}`

const receiver_out = `// MustReceiver is like Receiver except it panics if Receiver errors.
func (t T) MustReceiver() {
	err := t.Receiver()
	if err != nil {
		panic(err)
	}
	return
}`

const ptr_receiver_in = `func (t *T) PtrReceiver() error {
	return nil
}`

const ptr_receiver_out = `// MustPtrReceiver is like PtrReceiver except it panics if PtrReceiver errors.
func (t *T) MustPtrReceiver() {
	err := t.PtrReceiver()
	if err != nil {
		panic(err)
	}
	return
}`
