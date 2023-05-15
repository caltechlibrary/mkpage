#
# Simple Makefile for Golang based Projects.
#
PROJECT = mkpage

PROGRAMS = $(shell ls -1 cmd/)

PACKAGE = $(shell ls -1 *.go)

VERSION = $(shell grep '"version":' codemeta.json | cut -d\"  -f 4)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

OS = $(shell uname)

#PREFIX = /usr/local/bin
PREFIX = $(HOME)

ifneq ($(prefix),)
        PREFIX = $(prefix)
endif

EXT =
ifeq ($(OS), Windows)
        EXT = .exe
endif

PANDOC = $(shell which pandoc)

build: version.go $(PROGRAMS) CITATION.cff about.md installer

CITATION.cff: codemeta.json .FORCE
	cat codemeta.json | sed -E   's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	if [ -f $(PANDOC) ]; then echo "" | $(PANDOC) --metadata title="Citation $(PROJECT)" --metadata-file=_codemeta.json --template=codemeta-cff.tmpl >CITATION.cff; fi

about.md: codemeta.json .FORCE
	cat codemeta.json | sed -E   's/"@context"/"at__context"/g;s/"@type"/"at__type"/g;s/"@id"/"at__id"/g' >_codemeta.json
	if [ -f $(PANDOC) ]; then echo "" | $(PANDOC) --metadata title="About $(PROJECT)" --metadata-file=_codemeta.json --template=codemeta-md.tmpl >about.md; fi

version.go: .FORCE
	@echo "package $(PROJECT)" >version.go
	@echo '' >>version.go
	@echo 'const Version = "$(VERSION)"' >>version.go
	@echo '' >>version.go
	@git add version.go
	@if [ -f bin/codemeta ]; then ./bin/codemeta; fi

$(PROGRAMS): $(PACKAGE)
	@mkdir -p bin
	go build -o "bin/$@$(EXT)" cmd/$@/*.go

test: $(PACKAGE)
	go test

website: build about.md
	./mk_website.py

status:
	git status

save:
	@if [ "$(msg)" != "" ]; then git commit -am "$(msg)"; else git commit -am "Quick Save"; fi
	git push origin $(BRANCH)

refresh:
	git fetch origin
	git pull origin $(BRANCH)

publish:
	./mk_website.py
	./publish.bash

installer.sh: .FORCE
	echo '' | pandoc --metadata-file codemeta.json --template codemeta-installer.tmpl >installer.sh
	chmod 775 installer.sh
	git add -f installer.sh

clean:
	@if [ -f version.go ]; then rm version.go; fi
	@if [ -d bin ]; then rm -fR bin; fi
	@if [ -d dist ]; then rm -fR dist; fi
	@if [ -d man ]; then rm -fR man; fi

install: build
	@echo "Installing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f "./bin/$${FNAME}$(EXT)" ]; then mv -v "./bin/$${FNAME}$(EXT)" "$(PREFIX)/bin/$${FNAME}$(EXT)"; fi; done
	@echo ""
	@echo "Make sure $(PREFIX)/bin is in your PATH"

uninstall: .FORCE
	@echo "Removing programs in $(PREFIX)/bin"
	@for FNAME in $(PROGRAMS); do if [ -f "$(PREFIX)/bin/$${FNAME}$(EXT)" ]; then rm -v "$(PREFIX)/bin/$${FNAME}$(EXT)"; fi; done


dist/linux-amd64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env  GOOS=linux GOARCH=amd64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-$(VERSION)-linux-amd64.zip LICENSE codemeta.json CITATION.cff *.md bin/* docs/* how-to/*
	@rm -fR dist/bin


dist/macos-amd64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=amd64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-$(VERSION)-macos-amd64.zip LICENSE codemeta.json CITATION.cff *.md bin/* docs/* how-to/*
	@rm -fR dist/bin


dist/macos-arm64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=darwin GOARCH=arm64 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-$(VERSION)-macos-arm64.zip LICENSE codemeta.json CITATION.cff *.md bin/* docs/* how-to/*
	@rm -fR dist/bin


dist/windows-amd64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=amd64 go build -o "dist/bin/$${FNAME}.exe" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-$(VERSION)-windows-amd64.zip LICENSE codemeta.json CITATION.cff *.md bin/* docs/* how-to/*
	@rm -fR dist/bin

dist/windows-arm64: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=windows GOARCH=arm64 go build -o "dist/bin/$${FNAME}.exe" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-$(VERSION)-windows-arm64.zip LICENSE codemeta.json CITATION.cff *.md bin/* docs/* how-to/*
	@rm -fR dist/bin


dist/raspberry_pi_os-arm7: $(PROGRAMS)
	@mkdir -p dist/bin
	@for FNAME in $(PROGRAMS); do env GOOS=linux GOARCH=arm GOARM=7 go build -o "dist/bin/$${FNAME}" cmd/$${FNAME}/*.go; done
	@cd dist && zip -r $(PROJECT)-$(VERSION)-raspberry_pi_os-arm7.zip LICENSE codemeta.json CITATION.cff *.md bin/* docs/* how-to/*
	@rm -fR dist/bin

distribute_docs:
	@mkdir -p dist/
	@cp -v codemeta.json dist/
	@cp -v CITATION.cff dist/
	@cp -v README.md dist/
	@cp -v LICENSE dist/
	@cp -v INSTALL.md dist/
	@cp -v installer.sh dist/
	@cp -vR docs dist/
	@cp -vR how-to dist/

release: build installer.sh distribute_docs dist/linux-amd64 dist/macos-amd64 dist/macos-arm64 dist/windows-amd64 dist/windows-arm64 dist/raspberry_pi_os-arm7


.FORCE:
