package main

import (
	"fmt"
	"gonum.org/v1/gonum/graph/simple"
	"os"
	"path/filepath"
	"strings"
)


const (
	vendorFolder = "vendor"
	usVendorFolder = "_vendor"
	workspaceFolder = "_workspace"
)

func WalkPath(path string, functions FunctionMap) {
	walkerr := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if isValidDir(path) {
				pkgs := CreatePackages(path)
				for _ , pkg := range pkgs {
				fmt.Println("--> \t" + pkg.ID)
				IterateFiles(pkg.CompiledGoFiles)
				fmt.Println()
				IterateFileNodes(pkg, functions)
				}
			} else {
				return filepath.SkipDir
			}
		}
	return nil})

	if walkerr != nil {
		fmt.Println("Error occured: " + walkerr.Error())
	}
}


func isValidDir(path string) bool {
	pathElems := strings.Split(path, string(filepath.Separator))

	for _, el := range pathElems {
		// remove everything from dependencies folders
		if el == vendorFolder || el == workspaceFolder || el == usVendorFolder {
			return false
		}
		// remove all hidden folders
		if strings.HasPrefix(el, ".") {
			return false
		}
	}
	return true
}

func WalkPathCG(path string, directedGraph *simple.DirectedGraph, nodemap NodeMap) {
	walkerr := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if isValidDir(path) {
				fmt.Println(path)
				output, err := CallgraphCommand(path)
				if err == nil {
					GetCallgraph(output, directedGraph, nodemap)
				} else { fmt.Println(err.Error())}
			} else {
				return filepath.SkipDir
			}
		}

		return nil})

	if walkerr != nil {
		fmt.Println("Error occured: " + walkerr.Error())
	}
}
