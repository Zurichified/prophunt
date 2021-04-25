package main

// Struct representing a distinct function found inside a .go file
type Function struct {
	Pkg string // *ast.File.Name.Name
	File string // GoFiles[]
	Name string // *ast.FuncDecl.Name.Name
	Type string // Normal, Bechmark, Test
	Receiver string // *ast.FuncDecl.Recv.List[0].Type.Name
	Ptr bool // is receiver a value --> false, is receiver a pointer --> true
}

// Struct representing properties of a function
type FunctionProperties struct {
	// Signature/lines properties
	PkgFiles int // len(nodes)
	FileLoc int // fset.Position(Node.Pos/End).Line difference
	NameLength int // len(Name)
	Parameters int // len(*ast.FuncDecl.Type.Params.List)
	ParameterTypes []string // *ast.FuncDecl.Type.Params.List LOOP over Types, if Type *ast.SelectorExpr --> X.Name + Sel.Name, else if Type *ast.Ident --> Name
	Returns int // len(*ast.FuncDecl.Type.Results.List)
	ReturnTypes []string //*ast.FuncDecl.Type.Results.List LOOP over Types, if Type *ast.StarExpr --> *X.Name, else if other Types? --> *ast.Ident
	LoC int // fset.Position(Node.Pos/End).Line difference

	// Body related features
	FuncNames []string  // get names of called functions from *ast.CallExpr: if *ast.Ident --> Name, if *ast.SelectorExpr --> Sel.Name
	FuncCalls int // number of *ast.CallExpr inside the node

	Loops int // *ast.ForStmt	|  *ast.RangeStmt
	NestedLoops int // if node == *ast.ForStmt or *ast.RangeStmt, inspect and look if there is another *ast.ForStmt or *ast.RangeStmt, if yes, add and continue

	Channels int // Inspect for *ast.ChanType. If this doesn't work, *ast.CallExpr.Fun.Name == make ? *ast.CallExpr.Args[0]
	Sends int // search for *ast.SendStmt, if found, add 1
	Receives int // *ast.UnaryExpr.Op == '<-'																										// for all the creations with make, a make check is made.
	Closes int // *ast.CallExpr.Fun.Name == close , then add 1
	Gos int // Inspect for *ast.GoStmt
	ConcrRange int // Ranging over concurrency related stuff?
	Selects int // *ast.SelectStmt, if found add 1
	SelectCases map[int]int // *ast.SelectStmt.Body.List length

	Variables int // If declared with ':=' : *ast.AssignStmt.Tok == ':=' and Rhs[0].Op != '&', add 1
	Pointers int //  If declared with ':=' : *ast.AssignStmt.Tok == ':=' and Rhs[0].Op == '&', add 1

	Slices int // Inspect for *ast.ArrayType
	Maps int // Inspect for *ast.MapType

	IfElses int // Inspect for *ast.IfStmt
	Switches int // Inspect for *ast.SwitchStmt
	SwitchCases map[int]int // *ast.SwitchStmt.Body.List length

	Panics int // *ast.CallExpr.Fun.Name == "panic", then add 1
	Recovers int // *ast.CallExpr.Fun.Name == "recover", then add 1
	Defers int // *ast.DeferStmt, then add 1

	// Other metrics to be calculated
	CyclomaticComplexity int
	//CallgraphDepth int
}

// Function --> FunctionProperties
type FunctionMap map[Function]FunctionProperties

// Key can be Normal, Benchmark or Test
//type Maps map[string]FunctionMap

// This type is for counting the usage of the libraries that I'm tracking
// For example: key->io value->Write int->3
type Libraries map[string]map[string]int

// Predefined libraries to look for
var Prelibraries = map[string]int{
	"database/sql": 0,
	"encoding":0,
	"encoding/binary": 0,
	"encoding/csv": 0,
	"encoding/json": 0,
	"encoding/xml":0,
	"io": 0,
	"io/ioutil": 0,
	"math": 0,
	"math/rand": 0,
	"net": 0,
	"net/http": 0,
	"net/http/httputil": 0,
	"net/http/httptest": 0,
	"net/http/httptrace": 0,
	"net/smtp": 0,
	"net/textproto": 0,
	"net/rpc": 0,
	"net/rpc/jsonrpc": 0,
	"os": 0,
	"os/exec": 0,
	"os/signal": 0,
	"sort": 0,
	"strconv": 0,
	"sync": 0,
	"sync/atomic": 0,
	"syscall": 0,
	"bufio": 0,
	"bytes": 0,
	"crypto": 0,
	"mime": 0,
}

var SortedPrelibraries = []string{
	"bufio",
	"bytes",
	"crypto",
	"database/sql",
	"encoding",
	"encoding/binary",
	"encoding/csv",
	"encoding/json",
	"encoding/xml",
	"io",
	"io/ioutil",
	"math",
	"math/rand",
	"mime",
	"net",
	"net/http",
	"net/http/httptest",
	"net/http/httptrace",
	"net/http/httputil",
	"net/rpc",
	"net/rpc/jsonrpc",
	"net/smtp",
	"net/textproto",
	"os",
	"os/exec",
	"os/signal",
	"sort",
	"strconv",
	"sync",
	"sync/atomic",
	"syscall",
}
