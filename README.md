# Go Serverless - Code Beispiel Azure Function App

Lokale Ausführung mit Azure Function App Core Tools

## Nützliche Befehle

- Tests ausführen mit: `go test -v ./...`

- Anwendung bauen: `go build -o build/dogop .`

- Anwendung ausführen: `go run .`

- Ausführen mit Hot Reload über [air](https://github.com/cosmtrek/air): `air`

- Go Dokumentation lesen: `go doc http.HandlerFunc`

- Docker Container bauen: `docker build . -t crossnative/dogop`

## Projekt aufsetzen

Go Modul erstellen mit `go mod init crossnative/dogop`.

Erste Dependency einbinden mit `go get github.com/go-chi/chi/v5`.

## Docker Container bauen

### Cloudnative Buildpacks nutzen

    pack build dogop-cnb --buildpack paketo-buildpacks/go --builder paketobuildpacks/builder-jammy-base

    docker run --network host dogop-cnb


hey -n 200 -m POST -d '{ "age": 8, "breed": "chow" }' http://localhost:8080/api/quote

https://alphasec.io/how-to-deploy-a-github-container-image-to-google-cloud-run/

europe-west10-docker.pkg.dev/dogop-serverless/ghcr/remast/dogop-serverless:latest

hey -n 200 -m POST -d '{ "age": 8, "breed": "chow" }' https://dogop-serverless-746651650023.europe-west10.run.app/api/quote


hey -n 200 -m POST -d '{ "age": 8, "breed": "chow" }' https://europe-west10-dogop-serverless.cloudfunctions.net/quote


vegeta attack -duration=10s -rate=2 -targets=target.list -output=report_raw.bin
vegeta attack -duration=1m -rate=2 -targets=target.list -output=report_raw.bin