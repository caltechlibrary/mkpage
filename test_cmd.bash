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

function test_blogit() {
    mkdir -p test/blog
    if ! bin/blogit --prefix=test/blog README.md '2015-01-03'; then
        echo "Failed: bin/blogit --prefix=test/blog README.md 2015-01-03"
        exit 1
    fi 
    if ! bin/blogit --prefix=test/blog INSTALL.md '2015-01-03'; then
        echo "bin/blogit --prefix=test/blog INSTALL.md 2015-01-03"
        exit 1
    fi
    if ! bin/blogit --prefix=test/blog LICENSE '2015-01-03'; then
        echo "bin/blogit --prefix=test/blog LICENSE '2015-01-03'"
        exit 1
    fi
    if ! bin/blogit --prefix=test/blog Pandoc-Integration.md '2020-05-19'; then
        echo "bin/blogit --prefix=test/blog Pandoc-Integration.md '2020-05-19'"
        exit 1
    fi
    if ! bin/blogit --prefix=test/blog DEVELOPERS.md '2018-07-01'; then
        echo "bin/blogit --prefix=test/blog DEVELOPERS.md '2018-07-01'"
        exit 1
    fi
    echo "test_blogit OK"
}

function test_byline() {
    EXPECTED="By J. Q. Public 2018-12-04"
    # Test reading from file
    RESULT=$(bin/byline -i examples/byline/index.md)
    assert_equal "test_byline" "$EXPECTED" "$RESULT"
    # Test with standard input
    RESULT=$(cat examples/byline/index.md | bin/byline -i - )
    assert_equal "test_byline" "$EXPECTED" "$RESULT"
    echo "test_byline OK"
}

function test_mkpage() {
    # test basic markdown processing
    if [[ -f "temp.html" ]]; then rm temp.html; fi
    if [[ -f "temp.md" ]]; then rm temp.md; fi
    bin/mkpage "title=text:Hello World" \
        content=examples/helloworld.md page.tmpl > temp.html
    EXPECTED=""
    assert_exists "test_mkpage (simple)" "temp.html"
    RESULT=$(cmp examples/helloworld.html temp.html)
    assert_equal "test_mkpage (simple)" "$EXPECTED" "$RESULT"

    bin/mkpage -templates=testdata/codemeta.tmpl \
        "codemeta=codemeta.json" > temp.md
    EXPECTED=""
    assert_exists "test_mkpage (markdown doc)" "temp.md"
    RESULT=$(cmp examples/codemeta.md temp.md)
    assert_equal "test_mkpage (markdown doc)" "$EXPECTED" "$RESULT"

    echo "test_mkpage() OK"
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
test_blogit
#test_reldocpath
#test_sitemapper
#test_titleline
#test_urldecode
#test_urlencode
echo 'Success!'
