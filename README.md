# prophunt
Dependency Fetcher, Source Code Feature Extractor, Static Callgraph Analyzer for Projects written in Go

This Go project has 3 mains:
1) Dependency fetcher for projects
2) Prophunt (extracts source code features from all functions of a project and saves it in a .csv file)
3) Callgraph Analyzer (runs a callgraph analysis on the projects, which are in the output of Prophunt, to cumulate feature values for benchmarks. It then saves the benchmarks of a project with their values to a .csv file)


Usage of each main:

- **Dependency fetcher**  
(This will fetch the dependencies for the projects) [1) Go get 2) GoABS deps.Fetch]  
to build it --> `$ go build fetch_dependencies.go`  
to run it --> give path of "project_commit_place.csv" from downloaderscript  
(you need to have projects already downloaded to your system for the dependency fetcher to work, see other repo: https://github.com/Zurichified/goprojectdownloader)  

*Example:*
`$ fetch_dependencies "...BenchmarkProjects\project_commit_place.csv"`  


- **Prophunt**  
to build it --> `$ go build -o Prophunt main.go data.go parser.go pathwalker.go csvhandler.go callgraph.go`  
to run it --> give the path of "project_commit_place.csv" from downloaderscript  
(you need to have projects already downloaded to your system for the dependency fetcher to work, see other repo: https://github.com/Zurichified/goprojectdownloader)  

*Example:*
`$ Prophunt "...BenchmarkProjects\project_commit_place.csv"`  

Output: csv1 folder --> every project's csv file, index.csv showing the files (name, project_csv, project_path)  
These can later be used for the Callgraph Analyzer or standalone analysis purposes



- **Callgraph Analyzer**  
to build it --> `$ go build -o CallgraphAnalyzer cg.go data.go parser.go pathwalker.go csvhandler.go callgraph.go`  
to run it --> just run CallgraphAnalyzer, it will read index.csv from csv1 and iterate through projects to cumulatively collect feature values  

*Example:*
`$ CallgraphAnalyzer`  

Output: csv2 folder --> every project's csv file, index.csv showing the files (name, project_csv, project_path)  
These can later be used for any analysis purposes
