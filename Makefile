mac:
	env GOOS=darwin GOARCH=amd64 go build -o lexbelt.mac

linux:
	env GOOS=linux GOARCH=amd64 go build -o lexbelt.linux
	gzip -fk9 lexbelt.linux

clean:
	rm lexbelt.linux lexbelt.linux.gz
