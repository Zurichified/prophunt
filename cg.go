package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"gonum.org/v1/gonum/graph/simple"
	"os"
	"strconv"
	"strings"
)

func main() {
	csvfile, _ := os.Open("csv1/index.csv")
	r := csv.NewReader(bufio.NewReader(csvfile))
	r.Comma = ','
	// Get rid of header
	r.Read()
	// Get all the projects
	lines, _ := r.ReadAll()

	if _, err := os.Stat("csv2"); os.IsNotExist(err) {
		os.Mkdir("csv2", os.ModePerm)
	}

	indexName := "csv2" + string(os.PathSeparator) + "index.csv"
	var indexFile *os.File
	var indexWriter *csv.Writer
	if _, err := os.Stat(indexName); err == nil {
		indexFile, _ = os.OpenFile(indexName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		indexWriter = csv.NewWriter(indexFile)
	} else {
		indexFile, _ = os.Create(indexName)
		indexWriter = csv.NewWriter(indexFile)
		indexWriter.Write([]string{"name", "csvpath"})
		indexWriter.Flush()
	}


	// Start main loop
	for count, line := range lines {
		// Reading all functions
		project_csvpath := line[1]
		project_csvfile, _ := os.Open(project_csvpath)
		f := csv.NewReader(bufio.NewReader(project_csvfile))
		f.Comma = ','
		// Get rid of header
		f.Read()
		functions, _ := f.ReadAll()

		// Setup for the output folder
		project_name := line[0]
		nn_project_name := strings.Replace(project_name, "/", "&", -1)
		output_name := "csv2" + string(os.PathSeparator) + strconv.Itoa(count+1) + "&" + nn_project_name + ".csv"
		if _, err := os.Stat(output_name); err == nil {
			continue
		}

		// Getting benchmarks from csv file
		var benchmarks [][]string
		for _, function := range functions {
			if strings.Contains(function[1], "Benchmark") {
				benchmarks = append(benchmarks, function)
			}
		}

		project_path := line[2]
		callgraph := simple.NewDirectedGraph()
		nodemap := make(NodeMap)

		WalkPathCG(project_path, callgraph, nodemap)


		var records [][]string

		// Iterating each Benchmark
		for _, benchmark := range benchmarks {
			reachables := make(ReachableNodes)
			benchmark_node := nodemap[benchmark[1]]

			// If this benchmark wasn't a part of callgraph result, skip this benchmark
			if benchmark_node == nil{
				continue
			}


			// Collecting reachable nodes for the benchmark, in form of CSV1 Entries
			FromNext(callgraph, benchmark_node, reachables)

			var reachables_csv [][]string

			for node, _ := range reachables {
				func_id := NodeLookUp(node, nodemap)
				if func_id != "" {
					// Search for the function in csv1
					for _, function := range functions {
						if func_id == function[1] {
							reachables_csv = append(reachables_csv, function)
							break
						}
					}
				}
			}

			newvaluesmap := make(map[int]int)
			for _, function := range reachables_csv{
				for i := 2; i<len(function); i++ {
					numval, _ := strconv.Atoi(function[i])
					newvaluesmap[i] += numval
				}
			}
			records = append(records, SumToRecord(project_name,newvaluesmap,benchmark))
		}
		WriteCSV2(output_name, records)
		indexWriter.Write([]string{project_name,output_name})
		indexWriter.Flush()

		fmt.Println("Completed projects:" + strconv.Itoa(count+1))
	}
}