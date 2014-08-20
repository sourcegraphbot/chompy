package semantic

import (
	"log"

	"github.com/samertm/chompy/parse"
)

var _ = log.Fatal // debugging

func Gen(n parse.Node) []byte {
	t := check(n)
	return genCode(t)
}

// check is the "main" method for the semanic package. It runs all
// of the semanic checks and generates the IR for the backend.
func check(n parse.Node) *parse.Tree {
	t, ok := n.(*parse.Tree)
	if !ok {
		log.Fatal("Needed a tree.")
	}
	errs := treeWalks(t)
	if len(errs) != 0 {
		log.Fatal(errs)
	}
	return t
}

// Only deals with main.main right now.
func genCode(t *parse.Tree) []byte {
	code := emitStart()
	for _, n := range t.Kids {
		var f *parse.Funcdecl
		var ok bool
		if f, ok = n.(*parse.Funcdecl); !ok {
			continue
		}
		
		name := f.Name.Name
		code = append(code, emitFuncHeader(name)...)
		if name == "main" {
			code = append(code, emitFuncBody()...)
		}
	}
	return code
}

func emitStart() []byte {
	code := emitFuncHeader("_start")
	code = append(code, "\tbl\tmain\n" +
		"\tmov\tr0, #0\n" +
		"\tmov\tr7, #1\n" +
		"\tswi\t#0\n"...)
	return code
}

func emitFuncBody() []byte {
	return []byte("\tbx\tlr\n")
}

func emitFuncHeader(name string) []byte {
	return []byte("\t.align\t2\n"+
		"\t.global\t" + name + "\n" +
		name + ":\n")
}

//func emitMov(dest int, src int)
