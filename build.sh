#/bin/bash
GOOS=linux GOARCH=amd64 go build -o ./bin/spaceelevator.amd64 main.go \
&& GOOS=darwin GOARCH=amd64 go build -o ./bin/spaceelevator.darwin main.go \
&& GOOS=windows GOARCH=amd64 go build -o ./bin/spaceelevator.exe main.go \
&& tar cvzf ./spaceelevator.tar.gz ./bin
