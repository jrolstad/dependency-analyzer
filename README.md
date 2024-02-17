# Dependency Analyzer
When building a Java project, there are many third party dependencies that used to enable functionality that would otherwise be difficult to build on your own.  It can be difficult to determine all dependencies that a projecdt has and how they rely on each other.  The dependency analyzer intends to assist in this problem.

# How it works
The dependency analyzer is a cli tool that analyzes maven generated dot files and identifies all third party dependencies.  To run, used the following steps:
1. At the root of the Java project being analyzed generate the dependency data files by running:
```shell
mvn dependency:tree -DoutputType=dot -DoutputFile=dependencies.dot -Dverbose=true
```
2. Once files are generated, cd to the cmd\cli directory and run the main.go application.  Assuming the generated files from step 1 are in the data directory two levels up, a sample command is:
```shell
go run main.go --path ..\..\data\ --filePattern *.dot --includedParents com.oracle --excludedDependencies com.oracle,javax. --mode notreferenced
```
3. The command will output the dependency names and versions requested

# Requirements
* golang 1.18 or higher

# Components
## Applications

|Name|Location| Purpose                                                     |
|---|---|-------------------------------------------------------------|
|Command Line Interface|cmd/cli| Command line tool is the primary interface for the analyzer |

## Libraries

|Name|Location|Purpose|
|---|---|---|
|core|internal/core|Extensions to the native types in golang|
|models|internal/models|Data models / structs used in the components|
|orchestration|internal/orchestration|The main entry point for all exposed functionality. Contains business logic, processing, and flow implementations|


# License
This projects is made available under the [MIT License](LICENSE).
