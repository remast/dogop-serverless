clean:
	rm -rf build

zip:
	mkdir -p build/function
	cp go.mod build/function
	cp function.go build/function
	zip -j -r build/function.zip build/function

run-local:
	FUNCTION_TARGET=HandleQuote LOCAL_ONLY=true go run cmd/main.go