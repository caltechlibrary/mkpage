Bottler
=======

> `bottler` is an Python "Bottle" application generator

Bottler is an application factory or generator intended to be a tool to create and maintain a specific, narrow type of web application based on Bottle.py and Pee Wee ORM.

Features
--------

- Simple
- Highly targetted
- Useful set of minimal definition required for a running application
- Generated applications have common features for working with programmer defined objects

In the past building web applications have involved frameworks, platforms and stacks. Today web applications are commonly built as light weight "micro services". Micro services as jargon for small, often simple, applications that are highly focused cooprating with the rest of the web ecosystem. This is large part is a result of core services being broken off such as authenictation, storage, messaging.

In a library and archives setting there is a common need to manage objects, these objects often representing simple things (e.g. a person's name and their email address) to complex things (e.g. a bibliographic record). But whether simple or complex the objects have a small common set of operations that can be performed on one or more objects.  These are the function of list (a.k.a. query in DB termonology or search engine terminology) combined with the standard database management functions of create, read, update, delete.

In the past library and archives applications have had to be all in one system providing services like user management, authentication, object storage and management, website rendering as well as the core curation functions of list, create, read, update and delete.  Today we can be simplier web applications off loading user management and authentication to central campus services like Shibboleth, Open ID/OAuth or even the web's venerable Basic Auth.  Likewise we've learned to re-use content in other systems and can focus solely on curation as a task, web site generation as another, full text search as another.

At Caltech Library we've taken this simplity to heart. We've settled a small core of Python libraries to build on. The web service is built as a "Bottle" application, data management and persistence relies on the Pee Wee ORM.  Comparing three very different applications I see a core set of functionality where the parts of the application can be easily predicted even calculated in advance. 

From the web broswer point of view these applications require you to login, that login is external to our application but do provide does provide us with a consistent "REMOTE_USER" user identifity.  The applications tend to "list" objects being created with links to pages for creating and object, display an object in detail (read), updating an object and when appropriate deleting an object.  These can generally be described using common routes (URL paths) regardless of the objects being curated.

### Common paths

- `/list?q=...` or `/search?q=...` lists objects, provides filter and sorting of object lists
- `/create` and `/edit/<ObjectID>` provides a means of creating and updating an object
- `/delete/<ObjectID>` provides for object deletion

In addition any static content can be provided in the applications own "htdocs" directory. These paths can be calculated based on a simple walk of the file system.

### What varies

There are a couple variences in application. Some applications require user and role relations to determine what the user is allowed to do in the application. Some applications only require an authenticated user.  Yet others may need to trigger additional services (e.g. email users, respond to file system changes).Finally the objects that are manage very greately from the simple object id submission triggering a processing request to the complex records found in systems like EPrints. All these can be used 

### Strengths

In the age of micro applications and the web of 2021 a cluster of simple applications can be combined to perform more complex problems.  The is especially true of applications that curate objects.

Bottler results is a limited number of specific files being generated based on a known, specified, object schema.  

Bottler can be used along side GitHub templates to create a common organization for simple web applications base on Bottle and Pee Wee.

Bottler provides a simple layer on top of Python Bottle and Pee Wee using standards based data formats like JSON-LD and TOML ini file.

### What bottler doesn't do

Bottler builds a single, simple type of application based on Python Bottle and Pee Wee packages. It isn't a general purpose attempt to build complex framework, configure web servers based on XML or other structured markup. While it does use JSON to describe an object's schema it's simple limit vocabulary in a specific conext of JSON-LD. The object's attribute data types are limited to those that are supported in Pee Wee.

Bottler isn't a GUI generator. It does create Bottle templates based on the object schema described to does not render a complete GUI application.  The developer using `bottler` still still need to provide JavaScript, CSS and media assets to complete the web site experience.

Bottler isn't a build system though you can use bottler to re-generate the files it can derive from the provided schema.


Core Concepts
-------------

- `bottler` functions as a factory which builds customizable "Bottle" Python applications for curating objects.
- `bottler` applications provide a common set of "verbs" or operations for working with objects you want to curated.
- `bottler` provides a mechanism to update a `bottler` generated application based on changes in either a "bottler.ini" file or in the schema pointed to by the "bottler.ini" file.
- `bottler` can generate a default "bottler.ini" and schema to start from
- Object schema are defined in JSON and can leverage existing JSON-LD metadata schema (e.g. Person as defined at https://schema.org)
- Object schema are based on the simple types provided by the Python Pee Wee ORM
- Object schema can guide the generation of Bottle templates used to manage objects and lists objects
- Object schema can tie specific Web Components to specific fields editable using the generated templates


Verbs
-----

__list__
: would include listing, sorting and filter an objects being curated

__create__
: the ability to create a new object (e.g. submit a Thesis, Photo or PDF record)

__read__
: the ability to display details of an object

__update__
: the ability to update (edit) and object

__delete__ 'd'
: the ability to remove an object from a collection of object


Base technology
---------------

- Python Bottle
- Python Pee Wee
- SQLite3 database
- Extenal systems like Shibboleth or Basic AUTH for authenticating users
- A known function for authorization based on user name and verb permission

Life cycle of a Bottler application
-----------------------------------

- Application Creation
- Updating the application based on updated schema
- Running and testing a Bottler application
- Extending an Bottler application

### Application creation

This is performed with the command `bottler init [APP_NAME]` where "APP_NAME" is the name of your appliction. If APP_NAME is not provided then `bottler` assumes it is the name of the parent folder.

The initialization Bottler creates the following files (assuming they haven't been created previously)

- bottler.ini if it does not exist, should be updated by the programmer
- schema/schema.json defining the object managed by the application, should be updated by the programmer, it can be renamed and specified in the "bottler.ini" file
- run-server
- adapter.wsgi
- APP_NAME sub directory which holds your Bottle Py aplication
- APP_NAME/routes.py holding the known routes manager by Bottler
- APP_NAME/models.py 
- templates/ holding the application templates
- templates/page.tpl holds outer HTML markup used by all pages
- templates/nav.tpl holds the navigation links, included by page.tpl
- templates/list.tpl is the "content body" of the list page, usually is a table
- templates/create.tpl is the "content body" to use to create a new object defined by schema/schema.json
- templates/edit.tpl is the "content body" to use to create the edit page for an object
- componants holds the Web Components used by the application which can be included automagically in the bottler managed templates

### Example bottler.ini

The bottler.ini file contains enough metadata about the project to know where to find the schema for the application's object and also can hold some switch about what is managed by bottler.

Let's create a very simple app that accepts a object key like our Distillery
application at Caltech Library.

#### Generate the skeleton using `bottler init`

In this example we have a simple application that curates a name, an email address and slogan.

Command the programmer types in

```shell
bottler init hithere
```

This would generate a "bottler.ini" file 

```
[bottler]
app_name = hithere
verbs = [ list, create, read, update, delete ]
schema = schema/schema.json
```

It would also generate a `schema/schema.json` file hold our models' schema.

This file would look like

```json
[
    {
        "@context": "https://apps.library.caltech.edu/schema/v0",
        "@type": "ObjectSchema",
        "class": {
            "name": "hithere",
            "attributes": {
            }
        }
    }
]
```

At this point the programmer would modidify the schema.json file defining the object's attributes we need, i.e. name, email and slogan.

The programmer modified file would look like

```json
[
    {
        "@context": "https://apps.library.caltech.edu/schema/v0",
        "@type": "ObjectSchema",
        "class": {
            "name": "hithere",
            "attributes": {
                "name": { "type": "CharField(required=True)", "editable": true },
                "email": { "type": "CharField(unique=True)", "editable": true },
                "slogan": { "type": "TextField()", "editable": true },
                "created": {
                    "type": "DateTimeField()",
                    "editable": false
                }
            }
        }
    }
]
```

Notice the type corresponds the the types described in the [Pee Wee documentation](http://docs.peewee-orm.com/en/latest/peewee/models.html), the "editable" value indicates whether or not an editable field is generated in the "edit.tpl". If it
is not editable then the value will be displayed.

Now the programmer can run the "build" option.

```shell
bottler build
```

This will generate the following bottler managed files.

- run-server (generated the first time build run, does not get overwritten in subsequent builds)
- adapter.wsgi (generated the first time build is run, does not get overwritten in subsequent builds)
- hithree (directory, if not previously create)
- hithere/models.py (will be overwritten on each "build" if the comment "bottler: true" is present)
- hithere/routes.py (will be overwritten on each "build" if the comment "bottler:true" is present)
- templates (directory, if not previously created)
- templates/page.tpl (will be overwritten on each "build" if the comment "bottler:true" is present in the template)
- templates/nav.tpl (will be overwritten on "build" if the comment "bottler:true" present)
- templates/create.tpl (will be overwritten on "build" if the comment "bottler:true" present)
- templates/edit.tpl (will be overwritten on "build" if the comment "bottler:true" present)

Note if you removed one or more verbs from the "bottler.ini" file then the template page and routes will not be generated.


Files managed by Bottler
------------------------
- The list verb display either the list of sorta object keys or results from a search query (e.g. SQLite3 table, Open Search, LunrJS depending on how app is configured)
- Authentication is external, e.g. Shibboleth, Basic Auth
- Authorization contrained a function which takes a username and a verb and returns TRUE if it should be allowed or FALSE otherwise
    - This function is written outside the bottle generated ones but could be as simple as a look up table
- Authorization function name is specificed in the build.ini
- Supported verbs are identified in the build.ini file
- `bottler` will evaluate the routes to static content and generates functions defs for them
- `bottler` needs to support additional customization that will not be overwritten when you run `bottler build` (e.g. hooks to email, file system flags, etc)
- A settings.ini will be used for run time configuration


`bottler` is driven by declaring data structures in JSON-LD and providing a web component for any data type described.

The directory layout and file layout of a `bottler` application is.

- schema (a JSON-LD describing the schema for the objects curated by the application)
- components (HTML templates and JavaScript used to build the data entry screens - htdocs (for static content)
- bottler.ini describes how to build the application.
    - the name of the application's object directory
    - The "verbs" supported by the Bottle application generated
    - mapping schema elements to conponents and templates

- settings.ini is used to run the bottle app generated

Developer workflow
------------------

Creating a new web UI

1. `bottler init` would create a default directory structure, files and codemeta.json file.
2. complete our object schema
3. complete our web form markup
4. `bottler build` would build the application, e.g. generate fresh Python code for object modules and routes
5. `run-server` to test the resulting "Bottle" Python application


Updating an `bottler` app.

1. Update our `bottler` metadata files
    - Add/update/remove any schema needed by the app
    - Add/update/remove any webforms in the app
    - Add/update/remove any custom Web Components used in the app
    - Add/Update/remove any custom Python code (e.g. class extensions)
    - Update `bottler.ini` if needed
2. Run `bottler build` to rebuild the site
3. Run `run-server` to test the updated "Bottle" Python application

Static content
--------------

Static content should be placed in the htdocs directory. The directy structure will be used assign the appropirate routesin the "Bottle" application generated by `bottler`.

Pandoc is used to render Markdown content and can also be used to generating things like menus and other navigation elements used in the applicatin.

The htdocs directory should also contain any custom Javascript, media assets, CSS or  other materials needed to render the website outside of the generated Python code, templates and webforms.



Reference Material
------------------

- [Tutorial on Bottle app](https://bottlepy.org/docs/dev/tutorial_app.html)
- [Open Search](https://www.opensearch.org/)
- [Web Components](https://developer.mozilla.org/en-US/docs/Web/Web_Components)
    - [Examples](https://github.com/mdn/web-components-examples)
    - ECMAScript 2015 [Classes](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes)
    - [slots](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/slot)
    - [template](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/template)
        - [The truth about templates](https://developer.mozilla.org/en-US/docs/Web/Web_Components/Using_templates_and_slots#the_truth_about_templates)
    - [content](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/content)
- [Web forms](https://developer.mozilla.org/en-US/docs/Learn/Forms)
    - [HTML5 input types](https://developer.mozilla.org/en-US/docs/Learn/Forms/HTML5_input_types)
    - [Basic Native Controls](https://developer.mozilla.org/en-US/docs/Learn/Forms/Basic_native_form_controls)
    - [How to structure your web form](https://developer.mozilla.org/en-US/docs/Learn/Forms/How_to_structure_a_web_form)
- [HTML form element](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/form)
- [HTML forms guide](https://developer.mozilla.org/en-US/docs/Learn/Forms)
- Other elements that are used when creating forms: [&lt;button&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/button), [&lt;datalist&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/datalist), [&lt;fieldset&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/fieldset), [&lt;input&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input), [&lt;label&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/label), [&lt;legend&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/legend), [&lt;meter&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meter), [&lt;optgroup&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/optgroup), [&lt;option&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/option), [&lt;output&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/output), [&lt;progress&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/progress), [&lt;select&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/select), [&lt;textarea&gt;](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/textarea).
- Getting a list of the elements in the form: [HTMLFormElement.elements](https://developer.mozilla.org/en-US/docs/Web/API/HTMLFormElement/elements)
    - [ARIA: Form role](https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Roles/Form_Role)
    - [ARIA: Search role](https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Roles/Search_role)
