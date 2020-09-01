package mappers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/containers"
	"miller/dsl"
	"miller/mapping"
	"miller/parsing/lexer"
	"miller/parsing/parser"
)

// ----------------------------------------------------------------
var PutSetup = mapping.MapperSetup{
	Verb:         "put",
	ParseCLIFunc: mapperPutParseCLI,
	IgnoresInput: false,
}

func mapperPutParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)
	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperPutUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentioally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	// Get the DSL string from the command line, after the flags
	if (argi >= argc) {
		flagSet.Usage()
		os.Exit(1)
	}
	dslString := args[argi]
	argi += 1

	mapper, err := NewMapperPut(dslString)
	if err != nil {
		// xxx make sure better parse-fail info is printed by the DSL parser
		fmt.Fprintf(os.Stderr, "%s %s: cannot parse DSL expression.\n",
			args[0], verb)
		os.Exit(1)
	}

	*pargi = argi
	return mapper
}

func mapperPutUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {DSL expression}\n", argv0, verb)
	fmt.Fprintf(o, "TODO: put detailed on-line help here.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperPut struct {
	ast         *dsl.AST
	interpreter *dsl.Interpreter
}

func NewMapperPut(dslString string) (*MapperPut, error) {
	ast, err := NewASTFromString(dslString)
	if err != nil {
		return nil, err
	}
	return &MapperPut{
		ast:         ast,
		interpreter: dsl.NewInterpreter(),
	}, nil
}

// xxx note (package cycle) why not a dsl.AST constructor :(
// xxx maybe split out dsl into two package ... and/or put the ast.go into miller/parsing -- ?
//   depends on TBD split-out of AST and CST ...
func NewASTFromString(dslString string) (*dsl.AST, error) {
	theLexer := lexer.NewLexer([]byte(dslString))
	theParser := parser.NewParser()
	interfaceAST, err := theParser.Parse(theLexer)
	if err != nil {
		return nil, err
	}
	ast := interfaceAST.(*dsl.AST)
	return ast, nil
}

func (this *MapperPut) Map(
	inrecAndContext *containers.LrecAndContext,
	outrecsAndContexts chan<- *containers.LrecAndContext,
) {
	inrec := inrecAndContext.Lrec
	context := inrecAndContext.Context
	if inrec != nil {
		// xxx maybe ast -> interpreter ctor
		outrec, err := this.interpreter.InterpretOnInputRecord(inrec, &context, this.ast)
		if err != nil {
			// need echan or what?
			fmt.Println(err)
		} else {
			outrecsAndContexts <- containers.NewLrecAndContext(
				outrec,
				&context,
			)
		}
	} else {
		outrecsAndContexts <- containers.NewLrecAndContext(
			nil, // signals end of input record stream
			&context,
		)

	}
}
