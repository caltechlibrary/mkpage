
# publisher

## USAGE

    publisher [OPTIONS]

publisher is a tool to interact with a mkpage project published on AWS S3.
It support the basic CRUD operations in your S3 bucket.

## OPTIONS

```
	-create	create the named object in bucket
	-delete	delete an item in the bucket
	-h	display help
	-help	display help
	-l	display license
	-license	display license
	-local	local name of object, needed by create and update options
	-read	read the named object in bucket writing to stdout locally
	-update	update (delete,create) an object in a bucket
	-v	display version
	-version	display version
```

## EXAMPLES

Examples assume you've previously setup your AWS access via environment 
variables on configuration files.

Create an index.html file in your AWS Bucket. 

```shell
	publisher -create /index.html -local index.html
```

Read an item from the bucket (writes to stdout) 

```shell
	publisher -read /index.html
```

Update (actually a delete, followed by create) a file

```shell
	publisher -update /index.html -local index.html
```

Delete removes the copy in the bucket

```shell
    publisher -delete /index.html
```

