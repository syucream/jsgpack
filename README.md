# jsgpack

The simplest CLI tool converts msgpack&lt;->JSON

## Installation

```
$ go get -u github.com/syucream/jsgpack
```

## Usage

```
# JSON -> MsgPack
$ jsgpack -in path/to/input.json -out path/to/output.bin fromjson
$ cat path/to/input.json | jsgpack fromjson > path/to/output.bin

# MsgPack -> JSON
$ jsgpack -in path/to/input.bin -out path/to/output.json tojson
$ cat path/to/input.bin | jsgpack tojson > path/to/output.json
```
