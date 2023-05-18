# Homemade Container in Go

This repository contains a basic implementation of a container using Go. The aim is to illustrate the principles behind containerization technology by creating a simplified version from scratch. The container offers a segregated environment to execute specific commands with its own user namespace, UTS namespace, and PID namespace.

## Getting Started

Before running the project, ensure you have Go installed on your machine. If you haven't, you can download and install it from the official [Go website](https://golang.org/dl/).

## Structure of the Project

The project is a single Go file `main.go` that creates a simple container. It operates in two modes: 'run' mode to create a new container, and 'child' mode to run commands within that container.

## Running the Container

To build the project, navigate to the project directory and use the Go build command:

```bash
go build main.go
```

Then, you can run the container as follows:

```bash
./main run /bin/bash
```
For a detailed walkthrough on how the container was built and how each component works, refer to the blog post on my website.