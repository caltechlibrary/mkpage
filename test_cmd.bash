#!/bin/bash


function assert_exists() {
    if [ "$#" != "2" ]; then
        echo "wrong number of parameters for $1, $@"
        exit 1
    fi
    if [[ ! -f "$2" && ! -d "$2" ]]; then
        echo "$1: $2 does not exists"
        exit 1
    fi
}

function assert_equal() {
    if [ "$#" != "3" ]; then
        echo "wrong number of parameters for $1, $@"
        exit 1
    fi
    if [ "$2" != "$3" ]; then
        echo "$1: expected |$2| got |$3|"
        exit 1
    fi
}

function assert_empty() {
    if [ "$#" != "2" ]; then
        echo "wrong number of parameters for $1, $@"
        exit 1
    fi
    if [ "$2" != "" ]; then
        echo "$1: expected empty string got |$2|"
        exit 1
    fi
}


#
# Tests
#

function test_byline() {
    EXPECTED="By J. Q. Public 2018-12-04"
    # Test reading from file
    RESULT=$(bin/byline -i demo/byline/index.md)
    assert_equal "test_byline" "$EXPECTED" "$RESULT"
    # Test with standard input
    RESULT=$(cat demo/byline/index.md | bin/byline -i - )
    assert_equal "test_byline" "$EXPECTED" "$RESULT"
    echo "test_byline OK"
}

function test_mkpage() {
    # test basic markdown processing
    if [[ -f "temp.html" ]]; then rm temp.html; fi
    bin/mkpage content=demo/mkpage/helloworld.md page.tmpl > temp.html
    EXPECTED=""
    assert_exists "test_mkpage (simple)" "temp.html"
    RESULT=$(cmp demo/mkpage/helloworld.html temp.html)
    assert_equal "test_mkpage (simple)" "$EXPECTED" "$RESULT"

    # test codesnip support
    if [[ -f "temp.html" ]]; then rm temp.html; fi
    bin/mkpage content=demo/codesnip/index.md page.tmpl > temp.html
    EXPECTED=""
    assert_exists "test_mkpage (codesnip html)" "temp.html"
    RESULT=$(cmp demo/codesnip/index.html temp.html)
    assert_empty "test_mkpage (codesnip html)" "$RESULT"
    mkdir -p test/codesnip
    bin/mkpage -codesnip -i=demo/codesnip/index.md -o=test/codesnip/hello.bash
    assert_exists "test_mkpage (codesnip bash)" "test/codesnip/hello.bash"
    RESULT=$(cmp demo/codesnip/hello.bash test/codesnip/hello.bash)
    assert_empty "test_mkpage (codesnip bash)" "$RESULT"
    rm -fR test/codesnip
    echo "test_mkpage() OK"
}

function test_mkrss() {
    echo "test_mkrss() not implemented."
}

function test_mkslides() {
    echo "test_mkslides() not implemented."
    CWD=$(pwd)
    cd test
    ../bin/mkpage --mkslides 'content=../testdata/presentation_demo.md'
    cd $(CWD)
}

function test_reldocpath() {
    echo "test_reldocpath() not implemented."
}

function test_sitemapper() {
    echo "test_sitemapper() not implemented."
}

function test_titleline() {
    echo "test_titleline() not implemented."
}

function test_urldecode() {
    echo "test_urldecode() not implemented."
}

function test_urlencode() {
    echo "test_urlencode() not implemented."
}

echo "Testing command line tools"
test_byline
test_mkpage
test_mkrss
test_mkslides
test_reldocpath
test_sitemapper
test_titleline
test_urldecode
test_urlencode
echo 'Success!'
