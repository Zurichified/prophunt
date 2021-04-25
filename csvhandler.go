package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var Csvheader = []string{
	"function",
	"id",
	"pkgfiles",
	"fileloc",
	"namelength",
	"parameters",
	"returns",
	"loc",
	"funccalls",
	"loops",
	"nestedloops",
	"channels",
	"sends",
	"receives",
	"closes",
	"gos",
	"concrranges",
	"selects",
	"selectcases",
	"variables",
	"pointers",
	"slices",
	"maps",
	"ifelses",
	"switches",
	"switchcases",
	"panics",
	"recovers",
	"defers",
	"cyclomaticcomplexity",
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
var BenchCsvheader = []string{
	"project_name",
	"function",
	"id",
	"pkgfiles",
	"fileloc",
	"namelength",
	"parameters",
	"returns",
	"loc",
	"funccalls",
	"loops",
	"nestedloops",
	"channels",
	"sends",
	"receives",
	"closes",
	"gos",
	"concrranges",
	"selects",
	"selectcases",
	"variables",
	"pointers",
	"slices",
	"maps",
	"ifelses",
	"switches",
	"switchcases",
	"panics",
	"recovers",
	"defers",
	"cyclomaticcomplexity",
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

// Function to create an entry to the csv file from the given function & function properties
func FuncToRecord(projectname string, function Function, properties FunctionProperties) []string {
	var records []string

	// Take relative part by splitting and then remove "/" if the benchmark is in root directory (for matching later)
	sep := strings.Split(function.File, projectname)
	rel_path := sep[len(sep)-1]
	if strings.Count(rel_path, "/") == 1 {
		rel_path = rel_path[1:]
	}
	fname := rel_path + "/" + function.Name

	var id string
	if function.Receiver != "" {
			if function.Ptr {
				id = fmt.Sprintf("\"(*%s).%s\"", function.Receiver, function.Name)
			} else {
				id = fmt.Sprintf("\"(%s).%s\"", function.Receiver, function.Name)
			}
		} else { id = fmt.Sprintf("\"%s.%s\"", function.Pkg, function.Name)
	}
	//fmt.Println(id)

	selectcases := 0
	for _, count := range properties.SelectCases {
		selectcases += count
	}

	switchcases := 0
	for _, count := range properties.SwitchCases {
		switchcases += count
	}

	propsarray := []string{
		strconv.Itoa(properties.PkgFiles),
		strconv.Itoa(properties.FileLoc),
		strconv.Itoa(properties.NameLength),
		strconv.Itoa(properties.Parameters),
		strconv.Itoa(properties.Returns),
		strconv.Itoa(properties.LoC),
		strconv.Itoa(properties.FuncCalls),
		strconv.Itoa(properties.Loops),
		strconv.Itoa(properties.NestedLoops),
		strconv.Itoa(properties.Channels),
		strconv.Itoa(properties.Sends),
		strconv.Itoa(properties.Receives),
		strconv.Itoa(properties.Closes),
		strconv.Itoa(properties.Gos),
		strconv.Itoa(properties.ConcrRange),
		strconv.Itoa(properties.Selects),
		strconv.Itoa(selectcases),
		strconv.Itoa(properties.Variables),
		strconv.Itoa(properties.Pointers),
		strconv.Itoa(properties.Slices),
		strconv.Itoa(properties.Maps),
		strconv.Itoa(properties.IfElses),
		strconv.Itoa(properties.Switches),
		strconv.Itoa(switchcases),
		strconv.Itoa(properties.Panics),
		strconv.Itoa(properties.Recovers),
		strconv.Itoa(properties.Defers),
		strconv.Itoa(properties.CyclomaticComplexity),
	}

	libsmap := GetStdLibCounts(properties.FuncNames)
	var libsarray []string
	for _, name := range SortedPrelibraries {
		libsarray = append(libsarray, strconv.Itoa(libsmap[name]))
	}

	records = append(records, fname)
	records = append(records, id)
	records = append(records, propsarray...)
	records = append(records, libsarray...)
	return records
}

func WriteCSV(output_name string, records [][]string){
	if _, err := os.Stat("csv1"); os.IsNotExist(err) {
		os.Mkdir("csv1", os.ModePerm)
	}

	// Open the csv file
	file, err := os.Create(output_name)
	if err != nil {
		log.Fatal("Cannot open file", err)
	}
	defer file.Close()

	// Create a writer object
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header and the records
	err = writer.Write(Csvheader)
	if err != nil {
		log.Fatal("Cannot write header to file.")
	}

	// Write records
	for _, value := range records {
		err = writer.Write(value)
		if err != nil {
			log.Fatal("Cannot write value to file", err)
		}
	}
}

func WriteCSV2(output_name string, records [][]string){
	if _, err := os.Stat("csv2"); os.IsNotExist(err) {
		os.Mkdir("csv2", os.ModePerm)
	}

	// Open the csv file
	file, err := os.Create(output_name)
	if err != nil {
		log.Fatal("Cannot open file", err)
	}
	defer file.Close()

	// Create a writer object
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header and the records
	err = writer.Write(BenchCsvheader)
	if err != nil {
		log.Fatal("Cannot write header to file.")
	}

	// Write records
	for _, value := range records {
		err = writer.Write(value)
		if err != nil {
			log.Fatal("Cannot write value to file", err)
		}
	}
}

func WriteIndexCsv(writer *csv.Writer, projectname, output_name, projectpath string){
	writer.Write([]string{projectname,output_name,projectpath})
	defer writer.Flush()
}

func SumToRecord(projectname string, valuesmap map[int]int, benchmark []string) []string {
	var record []string

	record = append(record, projectname)
	record = append(record ,benchmark[:2]...)

	startind := 2
	if len(valuesmap) == 0 {
		record = append(record, benchmark[2:]...)
	} else {
		for i := 2; i < startind+len(valuesmap); i++ {
			//fmt.Println(valuesmap[i])
			record = append(record, strconv.Itoa(valuesmap[i]))
		}
		// Set namelength and parameter back to original value (remind that record has initial new project name, that's why shift one
		record[5] = benchmark[4]
		record[6] = benchmark[5]
	}
	return record
}