
all: clean tag tag.push release

tag.latest:
	git tag | cat | sort | tail -n 5

tag:
	git tag -a $(version)

tag.push:
	git push origin $(version)

release:
	source .env && goreleaser --rm-dist

clean:
	rm -rf gin-bin dist
