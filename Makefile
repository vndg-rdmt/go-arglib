.DEFAULT_GOAL := example

clean:
	@rm ./app

build:
	@go build -o app example/example.go

run: build
	@./app

example_help:
	@./app -h

example_missing_required_params:
	@./app -c 1

example_wrong_type1:
	@./app -c test

example_wrong_type2:
	@./app -c 1 2 3 -r -n testName -i a a a

example_wrong_type3:
	@./app -c 1 -r a -n testName -i a a a

example: build example_help clean