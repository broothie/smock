binname=smock
cmdpath=cmd/cmd.go
binpath=$(GOPATH)/bin

all: install

install: build
	mv $(binname) $(binpath)
	chmod +x $(binpath)/$(binname)

build:
	go build -o $(binname) $(cmdpath)

clean:
	rm gin-bin $(binname)
