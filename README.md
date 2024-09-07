#  ReReGo

```go
go test ./test
```

```go
go run ./cmd/main.go
```




Representation based automatic Research in Golang

I intend to write a minimal easy-to-use agent engine to be used with LLM's and other deep models to conduct agentic (recursing, iterating) research through web search. Unlike other frameworks out there, the sole and only interest is to do just that. The agent will follow a set of standards and formats including the DER files and the ESSENCE files, both which will have datastructures in Golang, JSON and SQL, since these files may become very large they may be divided into parts, more on that later.

## Components of ReReGo

Knowledge management and Collection
- Searching the web, API's, collection of files or asking the user to find knowledge.
- Retrieving knowledge.


Reasoning
- Traversing through knowledge to form connections and "understanding"
- Guiding knowledge collection and management.


ReReGo's core consists of a open and breezy framework for the most fundamental and clutterless way of orchestrating the Agents.

Basically there are these operations:

1. Search for new facts(essence)
2. Iterate these facts to create structure
3. Evaluate
4. Learn how to research - give the agent opurtunities to change it's own code/prompts to test wether alterations may enable it to achieve it's goal.


Given a query, the Agent will have to decompose it into subqueries, and then augment them with instructions on how to answer them.

## Agents

WikiData - given a task or set of instructions, the WikiData will provide back-context based on WikiData content.





## Prediction
