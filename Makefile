build:
	go generate
	go build

clean:
	$(RM) pingo pilang.go pilang.output
