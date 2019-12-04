
all: clean tag tag.push release

tag:
	git tag -a $(version)

tag.push:
	git push origin $(version)

release:
	source .env && goreleaser --rm-dist

clean:
	rm -rf gin-bin dist
