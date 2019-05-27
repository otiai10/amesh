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

create_release:
	curl -XPOST \
		https://api.github.com/repos/otiai10/amesh/releases \
		-H "Content-Type: application/json" \
		-H "Authorization: token $(GITHUB_ACCESS_TOKEN)" \
		-d '{"tag_name": "$(CIRCLE_TAG)", "target_commitish": "$(CIRCLE_SHA1)"}'

publish: $(addsuffix .log.json, $(TARGET_FILES))

release/%.log.json:
	$(eval RELEASE_ID := $(shell curl -s https://api.github.com/repos/otiai10/amesh/releases | jq 'select(.[].tag_name == "$(CIRCLE_TAG)") | .[].id'))
	curl -s -X POST \
		-H "Authorization: token $(GITHUB_ACCESS_TOKEN)" \
		-H "Content-Type: application/zip" \
		--data-binary @release/$* \
		"https://uploads.github.com/repos/otiai10/amesh/releases/$(RELEASE_ID)/assets?name=$*" >$@
