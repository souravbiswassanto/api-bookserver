name:= bookserver
all: gorun

gorun: gobin
	./$(name) start
gobin:
	go build -o $(name) .
