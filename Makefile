#
# Simple Makefile
#

PROJECT = mkpage

VERSION = $(shell grep -m1 "Version = " version.go | cut -d\` -f 2)

BRANCH = $(shell git branch | grep "* " | cut -d\   -f 2)

OS = $(shell uname)

EXT =
ifeq ($(OS),Windows)
	EXT = .exe
endif

build: bin/mkpage$(EXT) bin/mkrss$(EXT) \
	bin/sitemapper$(EXT) bin/byline$(EXT) bin/titleline$(EXT) \
	bin/reldocpath$(EXT) bin/urlencode$(EXT) bin/urldecode$(EXT) \
	bin/ws$(EXT) bin/frontmatter$(EXT) bin/blogit$(EXT)


bin/mkpage$(EXT): version.go mkpage.go codesnip.go pandoc.go cmd/mkpage/mkpage.go blogit.go pandoc.go generators.go
	go build -o bin/mkpage$(EXT) cmd/mkpage/mkpage.go

bin/mkrss$(EXT): version.go mkpage.go mkrss.go pandoc.go cmd/mkrss/mkrss.go
	go build -o bin/mkrss$(EXT) cmd/mkrss/mkrss.go

bin/sitemapper$(EXT): version.go mkpage.go cmd/sitemapper/sitemapper.go
	go build -o bin/sitemapper$(EXT) cmd/sitemapper/sitemapper.go

bin/byline$(EXT): version.go mkpage.go cmd/byline/byline.go
	go build -o bin/byline$(EXT) cmd/byline/byline.go

bin/titleline$(EXT): version.go mkpage.go cmd/titleline/titleline.go
	go build -o bin/titleline$(EXT) cmd/titleline/titleline.go

bin/reldocpath$(EXT): version.go cmd/reldocpath/reldocpath.go
	go build -o bin/reldocpath$(EXT) cmd/reldocpath/reldocpath.go

bin/urlencode$(EXT): version.go cmd/urlencode/urlencode.go
	go build -o bin/urlencode$(EXT) cmd/urlencode/urlencode.go

bin/urldecode$(EXT): version.go cmd/urldecode/urldecode.go
	go build -o bin/urldecode$(EXT) cmd/urldecode/urldecode.go

bin/ws$(EXT): version.go mkpage.go cmd/ws/ws.go
	go build -o bin/ws$(EXT) cmd/ws/ws.go

bin/frontmatter$(EXT): version.go mkpage.go pandoc.go cmd/frontmatter/frontmatter.go
	go build -o bin/frontmatter$(EXT) cmd/frontmatter/frontmatter.go

bin/blogit$(EXT): version.go blogit.go mkpage.go cmd/blogit/blogit.go
	go build -o bin/blogit$(EXT) cmd/blogit/blogit.go

lint:
	golint mkpage.go
	golint mkpage_test.go
	golint pandoc.go
	golint pandoc_test.go
	golint blogit.go
	golint blogit_test.go
	golint cmd/mkpage/mkpage.go
	golint cmd/mkrss/mkrss.go
	golint cmd/sitemapper/sitemapper.go
	golint cmd/byline/byline.go
	golint cmd/titleline/titleline.go
	golint cmd/reldocpath/reldocpath.go
	golint cmd/urlencode/urlencode.go
	golint cmd/urldecode/urldecode.go
	golint cmd/ws/ws.go
	golint cmd/frontmatter/frontmatter.go
	golint cmd/blogit/blotit.go

format:
	gofmt -w mkpage.go
	gofmt -w mkpage_test.go
	gofmt -w pandoc.go
	gofmt -w pandoc_test.go
	gofmt -w blogit.go
	gofmt -w blogit_test.go
	gofmt -w cmd/mkpage/mkpage.go
	gofmt -w cmd/mkrss/mkrss.go
	gofmt -w cmd/sitemapper/sitemapper.go
	gofmt -w cmd/byline/byline.go
	gofmt -w cmd/titleline/titleline.go
	gofmt -w cmd/reldocpath/reldocpath.go
	gofmt -w cmd/urlencode/urlencode.go
	gofmt -w cmd/urldecode/urldecode.go
	gofmt -w cmd/ws/ws.go
	gofmt -w cmd/frontmatter/frontmatter.go

test: .FORCE bin/mkpage$(EXT) bin/mkrss$(EXT) \
	bin/sitemapper$(EXT) bin/byline$(EXT) bin/titleline$(EXT) \
	bin/reldocpath$(EXT) bin/urlencode$(EXT) bin/urldecode$(EXT) \
	bin/ws$(EXT) bin/frontmatter$(EXT) bin/blogit$(ext)
	go test
	bash test_cmd.bash

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

clean:
	if [ -f assets.go ]; then rm assets.go; fi
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	if [ -d test ]; then rm -fR test/*; fi


install: 
	env GOBIN=$(GOPATH)/bin go install cmd/mkpage/mkpage.go
	env GOBIN=$(GOPATH)/bin go install cmd/mkrss/mkrss.go
	env GOBIN=$(GOPATH)/bin go install cmd/sitemapper/sitemapper.go
	env GOBIN=$(GOPATH)/bin go install cmd/byline/byline.go
	env GOBIN=$(GOPATH)/bin go install cmd/titleline/titleline.go
	env GOBIN=$(GOPATH)/bin go install cmd/reldocpath/reldocpath.go
	env GOBIN=$(GOPATH)/bin go install cmd/urlencode/urlencode.go
	env GOBIN=$(GOPATH)/bin go install cmd/urldecode/urldecode.go
	env GOBIN=$(GOPATH)/bin go install cmd/ws/ws.go
	env GOBIN=$(GOPATH)/bin go install cmd/frontmatter/frontmatter.go
	env GOBIN=$(GOPATH)/bin go install cmd/blogit/blogit.go


dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/mkpage cmd/mkpage/mkpage.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/mkrss cmd/mkrss/mkrss.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/sitemapper cmd/sitemapper/sitemapper.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/byline cmd/byline/byline.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/titleline cmd/titleline/titleline.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/reldocpath cmd/reldocpath/reldocpath.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/urlencode cmd/urlencode/urlencode.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/urldecode cmd/urldecode/urldecode.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/ws cmd/ws/ws.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/frontmatter cmd/frontmatter/frontmatter.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/blogit cmd/blogit/blogit.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* templates/*
	rm -fR dist/bin



dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/mkpage.exe cmd/mkpage/mkpage.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/mkrss.exe cmd/mkrss/mkrss.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/sitemapper.exe cmd/sitemapper/sitemapper.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/byline.exe cmd/byline/byline.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/titleline.exe cmd/titleline/titleline.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/reldocpath.exe cmd/reldocpath/reldocpath.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/urlencode.exe cmd/urlencode/urlencode.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/urldecode.exe cmd/urldecode/urldecode.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/ws.exe cmd/ws/ws.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/frontmatter.exe cmd/frontmatter/frontmatter.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/blogit.exe cmd/blogit/blogit.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* templates/*
	rm -fR dist/bin

dist/macos-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/mkpage cmd/mkpage/mkpage.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/mkrss cmd/mkrss/mkrss.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/sitemapper cmd/sitemapper/sitemapper.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/byline cmd/byline/byline.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/titleline cmd/titleline/titleline.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/reldocpath cmd/reldocpath/reldocpath.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/urlencode cmd/urlencode/urlencode.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/urldecode cmd/urldecode/urldecode.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/ws cmd/ws/ws.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/frontmatter cmd/frontmatter/frontmatter.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/blogit cmd/blogit/blogit.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macos-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* templates/*
	rm -fR dist/bin

dist/macos-arm64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/mkpage cmd/mkpage/mkpage.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/mkrss cmd/mkrss/mkrss.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/sitemapper cmd/sitemapper/sitemapper.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/byline cmd/byline/byline.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/titleline cmd/titleline/titleline.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/reldocpath cmd/reldocpath/reldocpath.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/urlencode cmd/urlencode/urlencode.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/urldecode cmd/urldecode/urldecode.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/ws cmd/ws/ws.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/frontmatter cmd/frontmatter/frontmatter.go
	env  GOOS=darwin GOARCH=arm64 go build -o dist/bin/blogit cmd/blogit/blogit.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macos-arm64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* templates/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/mkpage cmd/mkpage/mkpage.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/mkrss cmd/mkrss/mkrss.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/sitemapper cmd/sitemapper/sitemapper.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/byline cmd/byline/byline.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/titleline cmd/titleline/titleline.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/reldocpath cmd/reldocpath/reldocpath.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/urlencode cmd/urlencode/urlencode.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/urldecode cmd/urldecode/urldecode.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/ws cmd/ws/ws.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/frontmatter cmd/frontmatter/frontmatter.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/blogit cmd/blogit/blogit.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* templates/*
	rm -fR dist/bin

distribute_docs:
	mkdir -p dist/docs
	mkdir -p dist/how-to
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -vR docs/* dist/docs/
	cp -vR how-to/* dist/how-to/
	cp -vR templates dist/

release: clean website distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macos-amd64 dist/macos-arm64 dist/raspbian-arm7

website:
	./mk_website.py
	cd how-to/fountain-demo && make

publish: website
	./publish.bash

.FORCE:
