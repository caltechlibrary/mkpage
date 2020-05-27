
# blogit

## USAGE

	blogit [OPTIONS] DOCUMENT_NAME [DATE]

## DESCRIPTION


Blogit provides a quick tool to add or replace blog content
organized around a date oriented file path. In addition to
placing documents it also will generate simple markdown documents
for inclusion in navigation.


## OPTIONS

Below are a set of options available.

```
    -C, -copyright      Set the blog copyright notice.
    -D, -description    Set the blog description
    -E, -ended          Set the blog ended date.
    -IT, -index-tmpl    Set index blog template
    -L, -language       Set the blog language.
    -License            Set the blog language license.
    -N, -name           Set the blog name.
    -P, -prefix         Set the prefix path before YYYY/MM/DD.
    -PT, -post-tmpl     Set index blog template
    -Q, -quip           Set the blog quip.
    -R, -refresh        Refresh blog.json for a given year
    -S, -started        Set the blog started date.
    -U, -url            Set blog's URL
    -V, -verbose        verbose output
    -a, -asset          Copy asset file to the blog path for provided date (YYYY-MM-DD)
    -e, -examples       display examples
    -generate-manpage   generate man page
    -generate-markdown  generate markdown documentation
    -h, -help           display help
    -l, -license        display license
    -v, -version        display version
```

## EXAMPLES

I have a Markdown file called, "my-vacation-day.md". I want to
add it to my blog for the date July 1, 2020.  I've written
"my-vacation-day.md" in my home "Documents" folder and my blog
repository is in my "Sites" folder under "Sites/me.example.org".
Adding "my-vacation-day.md" to the blog me.example.org would
use the following command.

```shell
   cd $HOME/Sites/me.example.org
   blogit $HOME/my-vacation-day.md 2020-07-01
```

The *blogit* command will copy "my-vacation-day.md", creating any
necessary file directories to "$HOME/Sites/me.example.org/2020/06/01".
It will also update article lists (index.md) at the year level, 
month, and day level and month level of the directory tree and
and generate/update a posts.json in the "$HOME/Sites/my.example.org"
that can be used in your home page template for listing recent
posts.

*blogit* includes an option to set the prefix path to
the blog posting.  In this way you could have separate blogs 
structures for things like podcasts or videocasts.

```
    # Add a landing page for the podcast
    blogit -prefix=podcast my-vacation.md 2020-07-01
    # Add an audio file containing the podcast
    blogit -prefix=podcast my-vacation.wav 2020-07-01
```

```
   -p, -prefix    Set the prefix path before the YYYY/MM/DD path.
```


If you have an existing blog paths in the form of
PREFIX/YYYY/MM/DD you can use blogit to create/update/recreate
the blog.json file.

```
    blogit -prefix=blog -year=2020
```

The option "-year" is what indicates you want to crawl
for blog posts for that year.


