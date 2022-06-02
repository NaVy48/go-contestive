#! /bin/bash
go build -o build/server  contestive/cmd/server
cp config.json build/config.json
go build -o build/judge contestive/cmd/judge
