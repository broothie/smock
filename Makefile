
all: tag tag.push release

tag:
	git tag -a $(version) -m $(message)

tag.push:
	git push origin $(version)

release:
	source .env && goreleaser --rm-dist
