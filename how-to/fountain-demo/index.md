{
    "has_code": true
}

# Fountain Demo

**mkpage** supports processing [Fountain](https://fountain.io) markup.
This is helpful if you're transcribing and scriptting an interview or 
presentation. Of course Fountain is also used to write screenplays and 
plays.

**mkpage** will recognize fountain documents by their extensions, 
e.g. .spmd and .fountain. It will then invoke the fountain markup 
engine rather than the Markdown engine to turn the Fountain text into 
HTML. You can set the options for the fountain processor in the 
front matter of the document.

You can see the fountain source at [interview-with-a-dog.fountain](interview-with-a-dog.fountain) and an HTML rendering in [interview-with-a-dog.html](interview-with-a-dog.html).


The TOML front matter configures how the fountain document is rendered.

```TOML
    +++
    [fountain]
    page = false
    link_css = true
    css = "fountain.css"
    +++
```

The "+++" indicate the front matter is TOML. The fountain processor is looking for a heading of "[fountain]" then for the various settings, e.g. page (rendering full HTML page), link_css (using a link element to bring the CSS into the page), inline_css (use an inline style element for style the content), css (the CSS file to reference, either for inline content or as a link).

**mkpage** support three types of front matter. TOML, JSON and YAML. TOML and JSON are preferred and the best tested.

Processing our document can be done with the following command.

```shell
    mkpage content=interview-with-a-dog.fountain page.tmpl \
      >interview-with-a-dog.html
```

