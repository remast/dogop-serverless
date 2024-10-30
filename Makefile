clean:
	rm -rf build

build-dist:
	GOOS=linux GOARCH=amd64 go build -o build/function main.go