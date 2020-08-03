# jsgpack

The simplest CLI tool converts Message Pack <-> JSON(Newline Delimited JSON)

## Installation

```
$ go get -u github.com/syucream/jsgpack
```

## Usage

```
$ ./jsgpack
2020/08/04 00:17:08 call with subcommand, tojson or fromjson
$ ./jsgpack -h
Usage of ./jsgpack:
  -in string
        input file path
  -out string
        out file
```

```
# JSON -> MsgPack
$ jsgpack -in path/to/input.json -out path/to/output.bin fromjson
$ cat path/to/input.json | jsgpack fromjson > path/to/output.bin

# MsgPack -> JSON
$ jsgpack -in path/to/input.bin -out path/to/output.json tojson
$ cat path/to/input.bin | jsgpack tojson > path/to/output.json
```
