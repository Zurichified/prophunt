package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/sealuzh/GoABS/deps"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {

	csvpath := os.Args[1]
	csvfile, _ := os.Open(csvpath)
	r := csv.NewReader(bufio.NewReader(csvfile))
	r.Comma = ';'
	lines, _ := r.ReadAll()
	errs := 0
	get_cmd := exec.Command("go", "get", "-t", "./...")

	for count, line := range lines {
		github_name := line[0]
		fmt.Println(github_name)
		project_path := line[2]

		gopath := strings.Split(line[2], string(os.PathSeparator)+"src")[0]
		os.Setenv("GOPATH", gopath)

		get_cmd.Dir = project_path
		get_cmd.Env = os.Environ()
		get_cmd.Output()

		err := deps.Fetch(project_path)
		if err != nil {
			fmt.Println(err)
			errs += 1
		}

		fmt.Println(strconv.Itoa(count+1) + " projects' dependencies fetched.")

	}
	fmt.Println(strconv.Itoa(errs) + " projects gave an error while fetching dependencies.")
}
