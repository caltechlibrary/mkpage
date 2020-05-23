
# blogit

## USAGE

    blogit [OPTIONS] DOCUMENT_NAME [DATE]

## DESCRIPTION

Blogit provides a quick tool to add or replace blog content
organized around a date oriented file path. In addition to
placing documents it also will generate simple markdown documents
for inclusion in navigation.

## OPTIONS

```
    -V, -verbose         verbose output
    -e, -examples        display examples
    -generate-manpage    generate man page
    -generate-markdown   generate markdown documentation
    -h, -help            display help
    -l, -license         display license
    -p, -prefix          Set the prefix path before YYYY/MM/DD.
    -v, -version         display version
```


## EXAMPLES

I have a Markdown file called, "my-vacation-day.md". I want to
add it to my blog for the date July 1, 2020.  I've written
"my-vacation-day.md" in my home "Documents" folder and my blog
repository is in my "Sites" folder under "Sites/me.example.org".
Adding "my-vacation-day.md" to the blog me.example.org would
use the following command.

```
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

Option descripton for `-prefix`

```
   -p, -prefix    Set the prefix path before the YYYY/MM/DD path.
```

blogit v0.0.33
