build:
	go generate
	go build

clean:
	$(RM) pingo pilang.go pilang.output

tests: build
	python3 tests/runtests.py 
