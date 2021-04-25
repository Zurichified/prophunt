package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"strings"
)

// CreatePackages takes a valid path in the projects folder and creates pkgs.
// pkgs has the resulting packages from the folder: if tests are found, pkgs[1], else pkgs[0] is the valid package from the folder
// pkgs[1].GoFiles gives all the GoFiles as strings
// pkgs[1].Syntax gives all the *ast.File nodes of the files to further analyze.
// pkgs[1].TypesInfo gives the *types.Info of the package to type check later.
func CreatePackages(path string) []*packages.Package {
	var validpkgs []*packages.Package
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo, Dir:path, Tests: true}
	pkgs, err := packages.Load(cfg)
	if err != nil {
		fmt.Println(err)
	}
	//packages.PrintErrors(pkgs)

	//Printing resulting packages and their GoFiles --> To get all the files including test files, take pkgs[1] because this is the one with tests
	for _, pkg := range pkgs{
		//fmt.Println(pkg.ID)
		//for index, file := range pkg.GoFiles {
		//	fmt.Println("\t ", index, file)
		//}
		if !strings.HasSuffix(pkg.ID, ".test") {
			//fmt.Println(index, pkg.ID, pkg.IllTyped)
			validpkgs = append(validpkgs, pkg)

		}
	}

	// Mechanism to take only right packages out of a folder.
	// If there is only 1 package resulting from folder --> return this package. {do nothing}
	// If there are 2 packages resulting from folder, there can be 2 possibilities: {check if validpkg[0] and [1] have the same ID}
	// 		1) First package is standard package, second is standard package including test files.
	// 				In that case --> Return validpkgs[1]since this has all the files + test files.
	//      2) First package is standard package, second is _test of the standard package.
	//      		In that case --> Return validpkgs[0:2] since [1] has normal [2] has test package.
	// If there are 3 packages resulting from folder, then it's necessarily ordered like [0] = standard, [1] = standard_with_tests, [2] = _test.
	// 		In that case --> Return validpkgs[1:]

	if len(validpkgs) == 3 {
		validpkgs = validpkgs[1:]
	} else if len(validpkgs) == 2 {
		if strings.Split(validpkgs[0].ID, " ")[0] == strings.Split(validpkgs[1].ID, " ")[0]{
			validpkgs = validpkgs[1:2]
		}
	}

	return validpkgs
}

// Main iterator of the files inside a packages.Package
// Now only printing
func IterateFiles(files []string){
	for index, filename := range files{
		fmt.Println(index, filename)
	}
}

// Function to get all the FuncDecl Nodes out of a File Node
func GetFunctionNodes(mainNode *ast.File) []*ast.FuncDecl{
	funcs := []*ast.FuncDecl{}
	for _, node := range mainNode.Decls {
		if fn, isFunc := node.(*ast.FuncDecl); isFunc {
			funcs = append(funcs, fn)
		}
	}
	return funcs
}

// Main iterator of the *ast.File nodes inside a packages.Package
func IterateFileNodes(pkg *packages.Package, functions FunctionMap){
	// Iterating over every *ast.File
	for i, node := range pkg.Syntax{
		function_declarations := GetFunctionNodes(node)
		// Iterating over every *ast.FuncDecl within a node
		for _, function_declaration := range function_declarations{

			// Create a NewFunction as the key of a map, a NewProperties as the value of this key
			var NewFunction Function
			var NewProperties FunctionProperties

			// Fill the key properties
			NewFunction.Pkg = pkg.PkgPath
			//if i > len(pkg.GoFiles)-1 {
			//	NewFunction.File = pkg.PkgPath + "/" + "???"
			//} else { NewFunction.File = pkg.PkgPath + "/" + filepath.Base(pkg.GoFiles[i]) }
			NewFunction.File = pkg.CompiledGoFiles[i]
			NewFunction.Name = function_declaration.Name.Name
			NewFunction.Type = GetTypeOfFunc(function_declaration)
			NewFunction.Receiver, NewFunction.Ptr = GetReceiverType(function_declaration, pkg.TypesInfo)

			//fmt.Printf("%+v", NewFunction)
			//fmt.Println()

			// Fill the value properties
			NewProperties.PkgFiles = len(pkg.Syntax)
			NewProperties.FileLoc = pkg.Fset.Position(node.End()).Line - pkg.Fset.Position(node.Pos()).Line + 1
			NewProperties.NameLength = len(function_declaration.Name.Name)
			NewProperties.Parameters = len(function_declaration.Type.Params.List)
			if NewProperties.Parameters != 0 {
				//fmt.Printf("%+v\n", NewFunction)
				//fmt.Printf("%+v\n", NewProperties)
				NewProperties.ParameterTypes = GetParameterTypes(function_declaration, pkg.TypesInfo)
			}
			if ok := function_declaration.Type.Results; ok != nil {
				NewProperties.Returns = len(function_declaration.Type.Results.List)
				//fmt.Printf("%+v\n", NewFunction)
				//fmt.Printf("%+v\n", NewProperties)
				NewProperties.ReturnTypes = GetReturnTypes(function_declaration, pkg.TypesInfo)
			}
			NewProperties.LoC = pkg.Fset.Position(function_declaration.End()).Line - pkg.Fset.Position(function_declaration.Pos()).Line
			NewProperties.FuncCalls, NewProperties.FuncNames = GetFunctionCalls(function_declaration, pkg.TypesInfo)

			//fmt.Println(pkg.GoFiles[i], node.Name, function_declaration.Name)
			NewProperties.Loops, NewProperties.NestedLoops = GetLoops(function_declaration)
			NewProperties.Variables, NewProperties.Pointers = GetDeclarations(function_declaration)

			//fmt.Printf("%+v\n", NewFunction)
			//fmt.Printf("%+v\n", NewProperties)
			NewProperties.Channels, NewProperties.Sends, NewProperties.Receives, NewProperties.Closes, NewProperties.Gos,
			NewProperties.ConcrRange, NewProperties.Selects, NewProperties.SelectCases = GetConcrProps(function_declaration, pkg.TypesInfo)

			NewProperties.Slices, NewProperties.Maps = GetSliceMaps(function_declaration)

			NewProperties.IfElses, NewProperties.Switches, NewProperties.SwitchCases = GetControlFlow(function_declaration)

			NewProperties.Panics, NewProperties.Recovers, NewProperties.Defers = GetControlMechanisms(function_declaration)

			NewProperties.CyclomaticComplexity = GetCyclomaticComplexity(function_declaration)

			//if NewFunction.Type == "Normal" {
			//	maps["Normal"][NewFunction] = NewProperties
			//} else if NewFunction.Type =="Benchmark" {
			//	maps["Benchmark"][NewFunction] = NewProperties
			//}	else { maps["Test"][NewFunction] = NewProperties }

			functions[NewFunction] = NewProperties

		}
			//fmt.Printf("%+v", NewProperties)
			//fmt.Println()

		}
		//fmt.Println()
	}


// Function to get the type of a function (Normal, Test or Benchmark)
func GetTypeOfFunc(function *ast.FuncDecl) string {

	if strings.HasPrefix(function.Name.Name, "Benchmark"){
		return "Benchmark"
	} else if strings.HasPrefix(function.Name.Name, "Test"){
		return "Test"
	} else { return "Normal" }

	//// First check if there are parameters
	//if len(function.Type.Params.List) != 0 {
	//	switch f := function.Type.Params.List[0].Type.(type) {
	//	case *ast.StarExpr:
	//		if selector := f.X.(*ast.SelectorExpr).X.(*ast.Ident).Name; selector == "testing" {
	//			if x := f.X.(*ast.SelectorExpr).Sel.Name; x == "T" {
	//				functype = "Test"
	//			} else if x == "B" {
	//				functype = "Benchmark"
	//			}
	//		}
	//	}
	//}
	//return functype
}

// Function to get the function calls within a node and the count of them (handling *ast.CallExpr)
func GetFunctionCalls(function *ast.FuncDecl, pkginfo *types.Info) (int, []string) {
	var fcallnumber int
	var funcnames []string
	ast.Inspect(function, func(node ast.Node) bool {
		switch n := node.(type) {
			case *ast.CallExpr:
				//fmt.Printf("%T \n", n.Fun)
				switch nodetype := n.Fun.(type) {
					case *ast.Ident:
						if obj := pkginfo.ObjectOf(nodetype); obj != nil {
							if obj.Pkg() != nil {
								funcname := obj.Pkg().Path() + "." + obj.Name()
								funcnames = append(funcnames, funcname)
							} else { funcnames = append(funcnames, "builtin."+obj.Name()) }
						} else { funcnames = append(funcnames, "otherPackage." + nodetype.Name)}

					case *ast.SelectorExpr:
						if obj := pkginfo.ObjectOf(nodetype.Sel); obj != nil {
							if obj.Pkg() != nil {
								funcname := obj.Pkg().Path() + "."  + obj.Name()
								funcnames = append(funcnames, funcname)
							} else { funcnames = append(funcnames, "builtin."+obj.Name())}
						} else {funcnames = append(funcnames, "otherPackage." + nodetype.Sel.Name)}
				}
				fcallnumber += 1
		}
		return true
	})
	return fcallnumber, funcnames
}

// Function to get the receiver type of a function (value or pointer)
func GetReceiverType(function *ast.FuncDecl, pkginfo *types.Info) (string, bool)  {
	name := ""
	ptr := false
	if recv := function.Recv; recv != nil {
		switch fn := function.Recv.List[0].Type.(type) {
			case *ast.Ident:
				if obj := pkginfo.ObjectOf(fn); obj != nil {
					if obj.Pkg() != nil {
						name = obj.Pkg().Path() + "." + obj.Name()
					} else { name = "builtin." + obj.Name() }
				} else { name = "otherPackage." + fn.Name }
			case *ast.StarExpr:
				if obj := pkginfo.ObjectOf(fn.X.(*ast.Ident)); obj != nil {
					if obj.Pkg() != nil {
					name = obj.Pkg().Path() + "." + obj.Name()
					} else { name = "builtin." + obj.Name() }
				} else { name = "otherPackage." + fn.X.(*ast.Ident).Name }
				ptr = true
		}
	}
	return name, ptr
}

// Function to get the parameter types of a function (ident, selector, pointer, interface, function, array, map)
func GetParameterTypes(function *ast.FuncDecl, pkginfo *types.Info) []string {
	var parameters []string
	for _, parameter := range function.Type.Params.List {
		switch ptype := parameter.Type.(type) {
		case *ast.SelectorExpr:
			if objtype := pkginfo.ObjectOf(ptype.Sel); objtype != nil {
				if ok := objtype.Pkg(); ok != nil {
					name := objtype.Pkg().Name() + "." + objtype.Name()
					parameters = append(parameters, name)
				} else { name := "builtin." + objtype.Name()
					parameters = append(parameters, name)}
			} else { name := "otherPackage." + ptype.Sel.Name
				parameters = append(parameters, name)}

			//name := ptype.X.(*ast.Ident).Name + "." +  ptype.Sel.Name
		case *ast.Ident:
			if objtype := pkginfo.ObjectOf(ptype); objtype != nil {
			if ok := objtype.Pkg(); ok != nil {
				name := objtype.Pkg().Name() + "." + objtype.Name()
				parameters = append(parameters, name)
			} else {name := "builtin." + objtype.Name()
				parameters = append(parameters, name)}
			} else { name := "otherPackage." + ptype.Name
				parameters = append(parameters, name)}
		case *ast.StarExpr:
			switch further := ptype.X.(type) {
			case *ast.SelectorExpr:
				if objtype := pkginfo.ObjectOf(further.Sel); objtype != nil {
				if ok := objtype.Pkg(); ok != nil {
					name := objtype.Pkg().Name() + "." + objtype.Name()
					parameters = append(parameters, name)
				} else { name := "builtin." + objtype.Name()
					parameters = append(parameters, name)}
				} else { name := "otherPackage." + further.Sel.Name
					parameters = append(parameters, name)}
			case *ast.Ident:
				if objtype := pkginfo.ObjectOf(further); objtype != nil {
				if ok := objtype.Pkg(); ok != nil {
					name := objtype.Pkg().Name() + "." + objtype.Name()
					parameters = append(parameters, name)
				} else { name := "builtin." + objtype.Name()
					parameters = append(parameters, name)}
				} else { name := "otherPackage." + further.Name
				parameters = append(parameters, name)}
			}
		case *ast.InterfaceType:
			parameters = append(parameters, "interface")
		case *ast.FuncType:
			if len(parameter.Names) > 0 {
				if objtype := pkginfo.ObjectOf(parameter.Names[0]); objtype != nil {
					if ok := objtype.Pkg(); ok != nil {
						name := objtype.Pkg().Name() + "." + objtype.Name()
						parameters = append(parameters, name)
					} else {
						name := "builtin." + objtype.Name()
						parameters = append(parameters, name)
					}
				} else {
					name := "otherPackage." + parameter.Names[0].Name
					parameters = append(parameters, name)
				}
			} else { parameters = append(parameters, "interface")}
		case *ast.ArrayType:
			switch further := ptype.Elt.(type) {
			case *ast.Ident:
				if objtype := pkginfo.ObjectOf(further); objtype != nil {
					if ok := objtype.Pkg(); ok != nil {
						name := "[]" + objtype.Pkg().Name() + "." + objtype.Name()
						parameters = append(parameters, name)
					} else {
						name := "[]builtin." + objtype.Name()
						parameters = append(parameters, name)
					}
				}else {name := "otherPackage." + further.Name
					parameters = append(parameters, name)}
			case *ast.InterfaceType:
				parameters = append(parameters, "[]interface")
			}
		case *ast.MapType:
			switch further := ptype.Value.(type) {
			case *ast.Ident:
				if objtype := pkginfo.ObjectOf(further); objtype != nil {
				if ok := objtype.Pkg(); ok != nil {
					name := "map[]"+objtype.Pkg().Name() + "." + objtype.Name()
					parameters = append(parameters, name)
				} else { name := "map[]builtin." + objtype.Name()
					parameters = append(parameters, name)}
				} else { name := "otherPackage." + further.Name
					parameters = append(parameters, name)}
			case *ast.InterfaceType:
				parameters = append(parameters, "map[]interface")
			}
		case *ast.ChanType:
			switch further := ptype.Value.(type) {
			case *ast.Ident:
			if objtype := pkginfo.ObjectOf(further); objtype != nil {
			if ok := objtype.Pkg(); ok != nil {
				name := "chan_" + objtype.Pkg().Name() + "." + objtype.Name()
				parameters = append(parameters, name)
			} else { name := "chan_builtin." + objtype.Name()
				parameters = append(parameters, name)}
			} else { name := "otherPackage." + further.Name
				parameters = append(parameters, name)}
			}
		}
	}

	return parameters
}

// Function to get the return types of a function (ident, selector, pointer, array, map, interface)
func GetReturnTypes(function *ast.FuncDecl, pkginfo *types.Info) []string {
	var returns []string
	for _, returnvalue := range function.Type.Results.List {
		switch rtype := returnvalue.Type.(type) {
		case *ast.SelectorExpr:
			if objtype := pkginfo.ObjectOf(rtype.Sel); objtype != nil {
				if ok := objtype.Pkg(); ok != nil {
					name := objtype.Pkg().Path() + "." + objtype.Name()
					returns = append(returns, name)
				} else {
					name := "builtin." + objtype.Name()
					returns = append(returns, name)
				}
			} else { name := "otherPackage." + rtype.Sel.Name
				returns = append(returns, name)}
		case *ast.Ident:
			if objtype := pkginfo.ObjectOf(rtype); objtype != nil {
				if ok := objtype.Pkg(); ok != nil {
					name := objtype.Pkg().Path() + "." + objtype.Name()
					returns = append(returns, name)
				} else {
					name := "builtin." + objtype.Name()
					returns = append(returns, name)
				}
			} else { name := "otherPackage." + rtype.Name
				returns = append(returns, name)}
		case *ast.StarExpr:
			switch further := rtype.X.(type) {
			case *ast.SelectorExpr:
				if objtype := pkginfo.ObjectOf(further.Sel); objtype != nil {
					if ok := objtype.Pkg(); ok != nil {
						name := objtype.Pkg().Path() + "." + objtype.Name()
						returns = append(returns, name)
					} else {
						name := "builtin." + objtype.Name()
						returns = append(returns, name)
					}
				} else { name := "otherPackage." + further.Sel.Name
				returns = append(returns, name) }
			case *ast.Ident:
				if objtype := pkginfo.ObjectOf(further); objtype != nil {
					if ok := objtype.Pkg(); ok != nil {
						name := objtype.Pkg().Path() + "." + objtype.Name()
						returns = append(returns, name)
					} else {
						name := "builtin." + objtype.Name()
						returns = append(returns, name)
					}
				} else { name := "otherPackage." + further.Name
				returns = append(returns, name)}
			}
		case *ast.ArrayType:
			switch further := rtype.Elt.(type) {
			case *ast.Ident:
				if objtype := pkginfo.ObjectOf(further); objtype != nil {
					if ok := objtype.Pkg(); ok != nil {
						name := "[]" + objtype.Pkg().Path() + "." + objtype.Name()
						returns = append(returns, name)
					} else {
						name := "[]builtin." + objtype.Name()
						returns = append(returns, name)
					}
				} else { name := "otherPackage." + further.Name
				returns = append(returns, name)}
			case *ast.InterfaceType:
				returns = append(returns, "[]interface")
			}
		case *ast.MapType:
			switch further := rtype.Value.(type) {
			case *ast.Ident:
				if objtype := pkginfo.ObjectOf(further); objtype != nil {
					if ok := objtype.Pkg(); ok != nil {
						name := "map[]" + objtype.Pkg().Path() + "." + objtype.Name()
						returns = append(returns, name)
					} else {
						name := "map[]builtin." + objtype.Name()
						returns = append(returns, name)
					}
				} else { name := "otherPackage." + further.Name
				returns = append(returns, name)}
			case *ast.InterfaceType:
				returns = append(returns, "map[]interface")
			}
		case *ast.InterfaceType:
			returns = append(returns, "interface")
		case *ast.ChanType:
			switch further := rtype.Value.(type) {
			case *ast.Ident:
				if objtype := pkginfo.ObjectOf(further); objtype != nil {
					if ok := objtype.Pkg(); ok != nil {
						name := "chan_" + objtype.Pkg().Path() + "." + objtype.Name()
						returns = append(returns, name)
					} else {
						name := "chan_builtin." + objtype.Name()
						returns = append(returns, name)
					}
				} else { name := "otherPackage." + further.Name
				returns = append(returns, name)}
			}
		}
	}
	return returns
}

// Function to get count of total loops and loops with at least one other loop inside (nested loop)
// TODO: Get level specific count of nested loops
func GetLoops(function *ast.FuncDecl) (int,int) {
	loops, nestedloops := 0,0
	if function.Body != nil {
		ast.Inspect(function.Body, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.ForStmt:
				ast.Inspect(n.Body, func(node ast.Node) bool {
					switch node.(type) {
					case *ast.ForStmt:
						nestedloops += 1
						return false
					case *ast.RangeStmt:
						nestedloops += 1
						return false
					}
					return true
				})
				loops += 1

			case *ast.RangeStmt:
				ast.Inspect(n.Body, func(node ast.Node) bool {
					switch node.(type) {
					case *ast.ForStmt:
						nestedloops += 1
						return false
					case *ast.RangeStmt:
						nestedloops += 1
						return false
					}
					return true
				})
				loops += 1
			}
			return true
		})
	}
	return loops, nestedloops
}

// Tried a function with recursion that looks for the bodies of for statements as soon as it finds one
func LoopLevels(fornode *ast.ForStmt) int {

	found := containsForNode(fornode)

	if !found {
		return 1
	} else {
		for _, node := range fornode.Body.List {
			switch n := node.(type) {
			case *ast.ForStmt:
				return LoopLevels(n) + 1
			}
		}
	}
	panic("should never happen")
}

// Check if there is another for statement in the body of the *ast.ForStmt
func containsForNode(stmt *ast.ForStmt) bool {
	for _, node := range stmt.Body.List {
		if _, ok := node.(*ast.ForStmt); ok {
			return true
		}
	}
	return false
}

// Function to get the variable and pointer declarations inside the body
func GetDeclarations(function *ast.FuncDecl) (int,int) {
	variables := 0
	pointers := 0
	if function.Body != nil {
		ast.Inspect(function.Body, func(node ast.Node) bool {
			switch n := node.(type) {
			// Declarations with var keyword, pointers not possible --> Extract slices and maps
			case *ast.DeclStmt:
				switch decl := n.Decl.(type) {
				case *ast.GenDecl:
					switch specs := decl.Specs[0].(type) {
					case *ast.ValueSpec:
						switch specs.Type.(type) {
						case *ast.ArrayType:
							break
						case *ast.MapType:
							break
						default:
							variables += len(specs.Names)

						}
					}
				}
			// Assigns can be short var. declaration or assingments, we take ":=" and then check if variable or pointer.
			// Check is made for every item on the Rhs of the assignment, if something is not basitlit or pointer, it won't be count.
			case *ast.AssignStmt:
				if n.Tok.String() == ":=" {
					for _, expr := range n.Rhs {
						switch exprtype := expr.(type) {
						case *ast.BasicLit:
							variables += 1
						case *ast.UnaryExpr:
							if exprtype.Op.String() == "&" {
								pointers += 1
							}
							if exprtype.Op.String() == "<-" {
								variables += 1
							}
						}
					}
				}
			}
			return true
		})
	}
	return variables, pointers
}

// Function to get concurrency related properties
func GetConcrProps(function *ast.FuncDecl, pkginfo *types.Info) (int, int, int, int, int, int, int, map[int]int) {
	var channels, sends, receives, selects, closes, gos, concrrange int
	selectcases := make(map[int]int)
	if function.Body != nil {
		ast.Inspect(function.Body, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.ChanType:
				channels += 1
			case *ast.SendStmt:
				sends += 1
			case *ast.UnaryExpr:
				if n.Op.String() == "<-" {
					receives += 1
				}
			case *ast.SelectStmt:
				selects += 1
				selectcases[selects] = len(n.Body.List)
			case *ast.CallExpr:
				switch call := n.Fun.(type) {
				case *ast.Ident:
					if call.Name == "close" {
						closes += 1
					}
				}
			case *ast.GoStmt:
				gos += 1
			case *ast.RangeStmt:
				switch key := n.Key.(type) {
				case *ast.Ident:
					if key.Obj != nil {
						switch obj := key.Obj.Decl.(type) {
						case *ast.AssignStmt:
							switch rhs := obj.Rhs[0].(type) {
							case *ast.UnaryExpr:
								switch xtype := rhs.X.(type) {
								case *ast.Ident:
									objtype := pkginfo.ObjectOf(xtype)
									if objtype != nil {
										if strings.HasPrefix(objtype.Type().String(), "chan") {
											concrrange += 1
										}
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return channels, sends, receives, closes, gos, concrrange, selects, selectcases
}

// Function to get number of slices and maps
func GetSliceMaps(function *ast.FuncDecl) (int,int) {
	var slices, maps int
	if function.Body != nil {
	ast.Inspect(function.Body, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.ArrayType:
			slices += 1
		case *ast.MapType:
			maps += 1

		}
		return true
	})
	}
	return slices, maps
}

// Function to get number of ifElses, Switches and case numbers per switch
func GetControlFlow(function *ast.FuncDecl) (int,int,map[int]int){
	var ifelses, switches int
	switchcases := make(map[int]int)
	ast.Inspect(function, func(node ast.Node) bool {
		switch n := node.(type){
		case *ast.IfStmt:
			ifelses += 1
		case *ast.SwitchStmt:
			switches += 1
			switchcases[switches]=len(n.Body.List)
		}
	return true})
	return ifelses, switches, switchcases
}

// Function to get panics, recovers, defers
func GetControlMechanisms(function *ast.FuncDecl) (int, int, int){
	var panics, recovers, defers int
	ast.Inspect(function, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.DeferStmt:
			defers += 1
		case *ast.CallExpr:
			switch f := n.Fun.(type){
			case *ast.Ident:
				if f.Name == "panic" { panics += 1}
				if f.Name == "recover" { recovers += 1}
		}
		}
	return true})
	return panics, recovers, defers
}

func GetCyclomaticComplexity(function *ast.FuncDecl) int {
	complexity := 0
	ast.Inspect(function, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.FuncDecl, *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
			complexity++
		case *ast.BinaryExpr:
			if n.Op == token.LAND || n.Op == token.LOR {
				complexity++
			}
		}
		return true
	})
	return complexity
}

func GetStdLibCounts(funcnames []string) map[string]int {
	// Create libmap
	temp := make(map[string]int)
	for _, name := range SortedPrelibraries{
		temp[name] = 0
	}

	// Updating element counts
	for _, fname := range funcnames {
		for libname, _ := range temp {
			if strings.HasPrefix(fname, libname){
				temp[libname] += 1
			}
		}
	}
	return temp
}

// Helper function to print maps pretty, good to have for debugging
func PrettyPrint(functionMap FunctionMap) {
	fmt.Printf("%+v \n", functionMap)
	fmt.Println()
}