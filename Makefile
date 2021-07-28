OUTPUT_EXE=goGit

all: build

build:
	go build -o $(OUTPUT_EXE) cmd/main.go
