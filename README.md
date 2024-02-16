# dependency-analyzer
Analyzes dependencies in a project

# Process
Generate the dependency data files by running:
```shell
mvn dependency:tree -DoutputType=dot -DoutputFile=dependencies.dot -Dverbose=true
```