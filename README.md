# dependency-analyzer
Analyzes dependencies in a project

# Process
Generate the dependency data files by running:
```shell
mvn dependency:tree -DoutputType=dot -DoutputFile=dependencies.dot -Dverbose=true
```

To run the cli
```shell
go run main.go --path ..\..\data\ --filePattern *.dot --includedParents com.oracle --excludedDependencies com.oracle,javax. --mode notreferenced
```