# Go Serverless - Code Beispiel Container Based

## Projekt aufsetzen

Go Modul erstellen mit `go mod init crossnative.com/dogop-serverless`.

## Beispiel Last erzeugen

Mit hey:
```
hey -n 200 -m POST -d '{ "age": 8, "breed": "chow" }' https://dogop-serverless-746651650023.europe-west10.run.app/api/quote
```

Mit vegeta:
```
vegeta attack -duration=30s -rate=3 -targets=target-az-container.list -output=report_raw.bin
vegeta plot -title=Quote%20Results report_raw.bin > results.html
```