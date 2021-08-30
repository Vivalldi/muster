package main

import (
	"fmt"
	"go/ast"
	"log"
	"strings"
)

// FormatMode is used by formatting helpers
type FormatMode int

const (
	// Formatting options
	NeedName  FormatMode = 1 << iota // Include name in result (some formatters will gen the name, others will error)
	NeedType                         // NeedType ensures that the type of a field in included in the return
	NeedError                        // NeedError adds `err` | `err error` to the end of the returned string
)

// FormatReceiver builds a string from the given ast receiver. NeedName
// should always be set.
func FormatReceiver(fieldList *ast.FieldList, mode FormatMode) string {
	if fieldList == nil {
		return ""
	}

	if len(fieldList.List) != 1 {
		log.Fatalf("FormatReceiver: field list should only have 1 field")
	}

	var name string
	if mode&NeedName != 0 {
		if fieldList.List[0].Names == nil || len(fieldList.List[0].Names) == 0 {
			log.Fatalf("FormatReceiver: NeedName but no name for receiver")
		} else {
			name = fieldList.List[0].Names[0].Name
		}
	}

	var typ string
	if mode&NeedType != 0 {
		if fieldList.List[0].Type == nil {
			log.Fatalf("FormatReceiver: NeedType but no type for receiver")
		} else {
			typ = formatType(fieldList.List[0].Type)
		}
	}

	// Return a period after the name if we don't need type
	// becuase this is probably being used for a method call (and not declaration)
	var ret string
	if mode&NeedType != 0 {
		ret = fmt.Sprintf("(%s %s)", name, typ)
	} else {
		ret = fmt.Sprintf("%s.", name)
	}

	return strings.TrimSpace(ret)
}

// FormatParams builds a string from the given ast params. NeedError does not
// apply to this function. NeedName behaves like NeedType.
func FormatParams(fieldList *ast.FieldList, mode FormatMode) string {
	if fieldList == nil {
		return ""
	}

	var params []string
	for _, field := range fieldList.List {
		// Create a name variable. If NeedName is set then we need to
		// error if we don't have a name.
		var name string
		if mode&NeedName != 0 {
			if field.Names == nil || len(field.Names) == 0 {
				log.Fatalf("FormatParams: NeedName but no name for param")
			} else {
				name = field.Names[0].Name
			}
		}

		// Create a type variable. If NeedType is set then we need to
		// error if we don't have a type.
		var typ string
		if mode&NeedType != 0 {
			if field.Type == nil {
				log.Fatalf("FormatParams: NeedType but no type for param")
			} else {
				typ = formatType(field.Type)
			}
		}

		// Join the name & type together. We then trim any whitespace from the
		// result to ensure that we don't have any extra spaces (like when only name is set)
		param := fmt.Sprintf("%s %s", name, typ)
		params = append(params, strings.TrimSpace(param))
	}

	return strings.Join(params, ", ")
}

// FormatResults builds a string from the given ast results. NeedName is
// handled differently than most other formatters. If NeedName is set
// a name can be generated, but if it is not set, the name still a part of
// the return if the field had it set. FormatResults also always skips the
// last field. It is assumed that the last field is the error.
func FormatResults(fieldList *ast.FieldList, mode FormatMode) string {
	if fieldList == nil {
		return ""
	}

	var results []string
	for i, field := range fieldList.List {
		if i == len(fieldList.List)-1 {
			continue
		}

		// Create a name variable. If NeedName mode is set then we are allowed
		// to generate a name, as this will most likely be used as part of an assignment
		// call or return statement. Even if NeedName isn't set, we still return the
		// name if the original field had one.
		var name string
		if mode&NeedName != 0 {
			if field.Names == nil || len(field.Names) == 0 {
				// Generate a name since we don't have one & we need one
				field.Names = []*ast.Ident{ast.NewIdent(fmt.Sprintf("r%d", i))}
			}
		}
		// We also want to always use the name of the field if it was set
		if field.Names != nil && len(field.Names) > 0 {
			name = field.Names[0].Name
		}

		// Create a type variable. If NeedType mode is set then we need to
		// error if we don't have a type.
		var typ string
		if mode&NeedType != 0 {
			if field.Type == nil {
				log.Fatalf("FormatResults: NeedType but no type for result %d", i)
			} else {
				typ = formatType(field.Type)
			}
		}

		// Join the name & type together. We then trim any whitespace from the
		// result to ensure that we don't have any extra spaces (like when only name is set)
		result := fmt.Sprintf("%s %s", name, typ)
		results = append(results, strings.TrimSpace(result))
	}

	// If NeedError is set then we need to add an error return
	if mode&NeedError != 0 {
		if mode&NeedType != 0 {
			results = append(results, "err error")
		} else {
			results = append(results, "err")
		}
	}

	return strings.Join(results, ", ")
}

// formatType is a helper to format a type from a field
func formatType(typ ast.Expr) string {
	switch t := typ.(type) {
	case *ast.StarExpr:
		return fmt.Sprintf("*%s", t.X)
	case *ast.Ident:
		return t.Name
	case *ast.Ellipsis:
		return fmt.Sprintf("...%s", t.Elt)
	default:
		log.Fatalf("formatFieldType: unknown type %T", t)
	}

	return "" // not reached but compiler complains
}
