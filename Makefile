clean:
	rm -rf build

build-local:
	go build -o function main.go

build-dist:
	GOOS=linux GOARCH=amd64 go build -o build/function/function main.go
	cp host.json build/function
	cp -r quote build/function

start-local: build-local
	func start