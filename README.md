# Templates

Templates allows you to bundle file template in your binary.

## Installation

To install simply go get it:

```
$ go get github.com/cyberdelia/templates
```

## Usage

Generate the package:

```
$ templates -source templates/ templates/templates.go
```

Use it in your code:

```
template = template.Must(templates.Parse(nil))
```
