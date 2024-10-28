clean:
	rm -rf build

zip:
	mkdir -p build/function
	cp go.mod build/function
	cp main.go build/function
	zip -j -r build/function.zip build/function