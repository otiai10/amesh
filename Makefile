
.DEFAULT_GOAL := build

build:
	rm -rf ./release && mkdir ./release
	gox -output="./release/{{.Dir}}_{{.OS}}_{{.Arch}}"
