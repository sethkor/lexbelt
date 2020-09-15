mac:
	env GOOS=darwin GOARCH=amd64 go build -o lexbelt

linux:
	env GOOS=linux GOARCH=amd64 go build -o lexbelt.linux
	gzip -fk9 lexbelt.linux

clean:
	rm lexbelt.linux lexbelt.linux.gz

publish-test:
	goreleaser --snapshot --skip-publish --rm-dist

publish:
	goreleaser --rm-dist --skip-validate