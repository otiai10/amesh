TARGET_ARCHS=linux/amd64 darwin/amd64 windows/amd64
TARGET_DIRS=$(subst /,_,$(TARGET_ARCHS))
TARGET_FILES=$(TARGET_DIRS:%=release/%.zip)

all: $(TARGET_FILES)

release/%.zip:
	$(eval OSARCH := $(subst _,/,$(notdir $*)))
	gox -output="release/{{.OS}}_{{.Arch}}/amesh" --osarch="$(OSARCH)"
	cd ./release && zip -r $*.zip $*/* && cd ..

clean:
	rm -rf ./release/

publish: $(addsuffix .log.json, $(TARGET_FILES))

release/%.log.json:
	$(eval LATEST_TAG := $(shell git tag | tail -1))
	$(eval RELEASE_ID := $(shell curl -s https://api.github.com/repos/otiai10/amesh/releases | jq 'select(.[].tag_name == "v1.0.0") | .[].id'))
	$(eval FILE_SIZE := $(shell stat -f%z release/$*))
	curl -s -X POST \
		-H "Authorization: token $(GITHUB_ACCESS_TOKEN)" \
		-H "Content-Type: application/zip" \
		-H "Content-Length: $(FILE_SIZE)" \
		--data-binary @release/$* \
		"https://uploads.github.com/repos/otiai10/amesh/releases/$(RELEASE_ID)/assets?name=$*" >$@
