# Octopus HTTP client

## Description

Octopus is an http client that allows you to make HTTP requests described in the file.
Query templates can contain variables, and you can set them when you run the program.

## Templates

The **defaults** key describes default values of a template. This values will be set for **all** templates without defined fields. If you don't need to set any default values, leave it blank. 

Templates should be described in the json file in the key **templates**.
Each request must contain the mandatory **target** field with the URL of the target host.
There are also optional fields that allow you to specify request fields, such as headers or request type.

### Fields description

|Key|Description|
|-|--------|
| target | target URL
| description | description of the request
| method | HTTP request method (get/post/etc.). **GET** by default
| headings | list of headers
| cookies | list of cookies
| delay | delay in milliseconds before the request
| timeout | request timeout in milliseconds
| log | message to display on request
| response.log | message to display on response

Variables can be used in any text field. The variable name must have the format `${name}`.
You can set variables at program startup using the command line arguments `-v="variable name:the value of the variable"`

You can use variables `respCode` and `respBody` to get data from the response.

See [example.json](assets/example.json)


## Build

`go build -o oct cmd/main`

### Arguments

|Key|Description|
|-|--------|
| f | Templates file path
| v | Define a variable
| p | Parallelism level

## Run

`./oct -v "host:http://192.168.1.29:3000 " -v "token:123" -f example.json`

or

`run cmd/main.go -v="host:http://192.168.1.29:3000" -v="token:123" -f="example.json"`