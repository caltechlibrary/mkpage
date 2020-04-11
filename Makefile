#
# Simple Makefile
#

PROJECT = mkpage

VERSION = $(shell grep -m1 "Version = " $(PROJECT).go | cut -d\` -f 2)

BRANCH = $(shell git branch | grep "* " | cut -d\   -f 2)

PKGASSETS = $(shell which pkgassets)

OS = $(shell uname)

EXT =
ifeq ($(OS),Windows)
	EXT = .exe
endif

build: bin/mkpage$(EXT) bin/mkslides$(EXT) bin/mkrss$(EXT) \
	bin/sitemapper$(EXT) bin/byline$(EXT) bin/titleline$(EXT) \
	bin/reldocpath$(EXT) bin/urlencode$(EXT) bin/urldecode$(EXT) \
	bin/ws$(EXT) bin/frontmatter$(EXT) bin/mkpongo

mkpage.go: assets.go codesnip.go

assets.go: defaults/templates/page.tmpl defaults/templates/slides.tmpl
	pkgassets -o assets.go -p mkpage Defaults defaults
	git add assets.go

bin/mkpage$(EXT): mkpage.go assets.go codesnip.go cmd/mkpage/mkpage.go
	go build -o bin/mkpage$(EXT) cmd/mkpage/mkpage.go

bin/mkslides$(EXT): mkpage.go cmd/mkslides/mkslides.go
	go build -o bin/mkslides$(EXT) cmd/mkslides/mkslides.go

bin/mkrss$(EXT): mkpage.go cmd/mkrss/mkrss.go
	go build -o bin/mkrss$(EXT) cmd/mkrss/mkrss.go

bin/sitemapper$(EXT): mkpage.go cmd/sitemapper/sitemapper.go
	go build -o bin/sitemapper$(EXT) cmd/sitemapper/sitemapper.go

bin/byline$(EXT): mkpage.go cmd/byline/byline.go
	go build -o bin/byline$(EXT) cmd/byline/byline.go

bin/titleline$(EXT): mkpage.go cmd/titleline/titleline.go
	go build -o bin/titleline$(EXT) cmd/titleline/titleline.go

bin/reldocpath$(EXT): cmd/reldocpath/reldocpath.go
	go build -o bin/reldocpath$(EXT) cmd/reldocpath/reldocpath.go

bin/urlencode$(EXT): cmd/urlencode/urlencode.go
	go build -o bin/urlencode$(EXT) cmd/urlencode/urlencode.go

bin/urldecode$(EXT): cmd/urldecode/urldecode.go
	go build -o bin/urldecode$(EXT) cmd/urldecode/urldecode.go

bin/ws$(EXT): mkpage.go cmd/ws/ws.go
	go build -o bin/ws$(EXT) cmd/ws/ws.go

bin/frontmatter$(EXT): mkpage.go cmd/frontmatter/frontmatter.go
	go build -o bin/frontmatter$(EXT) cmd/frontmatter/frontmatter.go

bin/mkpongo$(EXT): mkpage.go mkpongo.go cmd/mkpongo/mkpongo.go
	go build -o bin/mkpongo$(EXT) cmd/mkpongo/mkpongo.go

lint:
	golint mkpage.go
	golint mkpage_test.go
	golint cmd/mkpage/mkpage.go
	golint cmd/mkslides/mkslides.go
	golint cmd/mkrss/mkrss.go
	golint cmd/sitemapper/sitemapper.go
	golint cmd/byline/byline.go
	golint cmd/titleline/titleline.go
	golint cmd/reldocpath/reldocpath.go
	golint cmd/urlencode/urlencode.go
	golint cmd/urldecode/urldecode.go
	golint cmd/ws/ws.go
	golint cmd/frontmatter/frontmatter.go
	golint cmd/mkpongo/mkpongo.go

format:
	gofmt -w mkpage.go
	gofmt -w mkpage_test.go
	gofmt -w cmd/mkpage/mkpage.go
	gofmt -w cmd/mkslides/mkslides.go
	gofmt -w cmd/mkrss/mkrss.go
	gofmt -w cmd/sitemapper/sitemapper.go
	gofmt -w cmd/byline/byline.go
	gofmt -w cmd/titleline/titleline.go
	gofmt -w cmd/reldocpath/reldocpath.go
	gofmt -w cmd/urlencode/urlencode.go
	gofmt -w cmd/urldecode/urldecode.go
	gofmt -w cmd/ws/ws.go
	gofmt -w cmd/frontmatter/frontmatter.go
	gofmt -w cmd/mkpongo/mkpongo.go

test: bin/mkpage$(EXT) bin/mkslides$(EXT) bin/mkrss$(EXT) \
	bin/sitemapper$(EXT) bin/byline$(EXT) bin/titleline$(EXT) \
	bin/reldocpath$(EXT) bin/urlencode$(EXT) bin/urldecode$(EXT) \
	bin/ws$(EXT) bin/frontmatter$(EXT) bin/mkpongo$(EXT) FORCE
	go test
	bash test_cmd.bash

status:
	git status

save:
	if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

clean:
	if [ "$(PKGASSETS)" != "" ]; then rm assets.go; pkgassets -o assets.go -p mkpage Defaults defaults; git add assets.go; fi
	if [ -d bin ]; then rm -fR bin; fi
	if [ -d dist ]; then rm -fR dist; fi
	if [ -d man ]; then rm -fR man; fi

man: build
	mkdir -p man/man1
	bin/mkpage -generate-manpage | nroff -Tutf8 -man > man/man1/mkpage.1
	bin/mkslides -generate-manpage | nroff -Tutf8 -man > man/man1/mkslides.1
	bin/mkrss -generate-manpage | nroff -Tutf8 -man > man/man1/mkrss.1
	bin/sitemapper -generate-manpage | nroff -Tutf8 -man > man/man1/sitemapper.1
	bin/byline -generate-manpage | nroff -Tutf8 -man > man/man1/byline.1
	bin/titleline -generate-manpage | nroff -Tutf8 -man > man/man1/titleline.1
	bin/reldocpath -generate-manpage | nroff -Tutf8 -man > man/man1/reldocpath.1
	bin/urldecode -generate-manpage | nroff -Tutf8 -man > man/man1/urldecode.1
	bin/urlencode -generate-manpage | nroff -Tutf8 -man > man/man1/urlencode.1
	bin/ws -generate-manpage | nroff -Tutf8 -man > man/man1/ws.1
	bin/frontmatter -generate-manpage | nroff -Tutf8 -man > man/man1/frontmatter.1
	bin/mkpongo -generate-manpage | nroff -Tutf8 -man > man/man1/mkpongo.1

install: assets.go
	env GOBIN=$(GOPATH)/bin go install cmd/mkpage/mkpage.go
	env GOBIN=$(GOPATH)/bin go install cmd/mkslides/mkslides.go
	env GOBIN=$(GOPATH)/bin go install cmd/mkrss/mkrss.go
	env GOBIN=$(GOPATH)/bin go install cmd/sitemapper/sitemapper.go
	env GOBIN=$(GOPATH)/bin go install cmd/byline/byline.go
	env GOBIN=$(GOPATH)/bin go install cmd/titleline/titleline.go
	env GOBIN=$(GOPATH)/bin go install cmd/reldocpath/reldocpath.go
	env GOBIN=$(GOPATH)/bin go install cmd/urlencode/urlencode.go
	env GOBIN=$(GOPATH)/bin go install cmd/urldecode/urldecode.go
	env GOBIN=$(GOPATH)/bin go install cmd/ws/ws.go
	env GOBIN=$(GOPATH)/bin go install cmd/frontmatter/frontmatter.go
	env GOBIN=$(GOPATH)/bin go install cmd/mkpongo/mkpongo.go


dist/linux-amd64:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/mkpage cmd/mkpage/mkpage.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/mkslides cmd/mkslides/mkslides.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/mkrss cmd/mkrss/mkrss.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/sitemapper cmd/sitemapper/sitemapper.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/byline cmd/byline/byline.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/titleline cmd/titleline/titleline.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/reldocpath cmd/reldocpath/reldocpath.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/urlencode cmd/urlencode/urlencode.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/urldecode cmd/urldecode/urldecode.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/ws cmd/ws/ws.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/frontmatter cmd/frontmatter/frontmatter.go
	env  GOOS=linux GOARCH=amd64 go build -o dist/bin/mkpongo cmd/mkpongo/mkpongo.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* templates/*
	rm -fR dist/bin



dist/windows-amd64:
	mkdir -p dist/bin
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/mkpage.exe cmd/mkpage/mkpage.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/mkslides.exe cmd/mkslides/mkslides.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/mkrss.exe cmd/mkrss/mkrss.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/sitemapper.exe cmd/sitemapper/sitemapper.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/byline.exe cmd/byline/byline.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/titleline.exe cmd/titleline/titleline.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/reldocpath.exe cmd/reldocpath/reldocpath.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/urlencode.exe cmd/urlencode/urlencode.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/urldecode.exe cmd/urldecode/urldecode.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/ws.exe cmd/ws/ws.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/frontmatter.exe cmd/frontmatter/frontmatter.go
	env  GOOS=windows GOARCH=amd64 go build -o dist/bin/mkpongo.exe cmd/mkpongo/mkpongo.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* templates/*
	rm -fR dist/bin

dist/macosx-amd64:
	mkdir -p dist/bin
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/mkpage cmd/mkpage/mkpage.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/mkslides cmd/mkslides/mkslides.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/mkrss cmd/mkrss/mkrss.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/sitemapper cmd/sitemapper/sitemapper.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/byline cmd/byline/byline.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/titleline cmd/titleline/titleline.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/reldocpath cmd/reldocpath/reldocpath.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/urlencode cmd/urlencode/urlencode.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/urldecode cmd/urldecode/urldecode.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/ws cmd/ws/ws.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/frontmatter cmd/frontmatter/frontmatter.go
	env  GOOS=darwin GOARCH=amd64 go build -o dist/bin/mkpongo cmd/mkpongo/mkpongo.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-macosx-amd64.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* templates/*
	rm -fR dist/bin

dist/raspbian-arm7:
	mkdir -p dist/bin
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/mkpage cmd/mkpage/mkpage.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/mkslides cmd/mkslides/mkslides.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/mkrss cmd/mkrss/mkrss.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/sitemapper cmd/sitemapper/sitemapper.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/byline cmd/byline/byline.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/titleline cmd/titleline/titleline.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/reldocpath cmd/reldocpath/reldocpath.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/urlencode cmd/urlencode/urlencode.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/urldecode cmd/urldecode/urldecode.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/ws cmd/ws/ws.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/frontmatter cmd/frontmatter/frontmatter.go
	env  GOOS=linux GOARCH=arm GOARM=7 go build -o dist/bin/mkpongo cmd/mkpongo/mkpongo.go
	cd dist && zip -r $(PROJECT)-$(VERSION)-raspbian-arm7.zip README.md LICENSE INSTALL.md bin/* docs/* how-to/* templates/*
	rm -fR dist/bin

distribute_docs:
	mkdir -p dist/docs
	mkdir -p dist/how-to
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v docs/byline.md dist/docs/
	cp -v docs/mkpage.md dist/docs/
	cp -v docs/mkrss.md dist/docs
	cp -v docs/mkslides.md dist/docs/
	cp -v docs/reldocpath.md dist/docs/
	cp -v docs/sitemapper.md dist/docs/
	cp -v docs/titleline.md dist/docs/
	cp -v docs/urldecode.md dist/docs/
	cp -v docs/urlencode.md dist/docs/
	cp -v docs/ws.md dist/docs/
	cp -v docs/frontmatter.md dist/docs/
	cp -v docs/mkpongo.md dist/docs/
	cp -v how-to/go-template-recipes.md dist/how-to/
	cp -v how-to/the-basics.md dist/how-to/
	cp -vR templates dist/
	#FIXME: need to pull package versions from go.mod file.
	#./package-versions.bash > dist/package-versions.txt

release: clean website assets.go distribute_docs dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7

website:
	./mk-website.py

publish:
	./mk-website.py
	./publish.bash

FORCE:
