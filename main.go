package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/vmihailenco/msgpack/v4"
)

const (
	subCommandToJson   = "tojson"
	subCommandFromJson = "fromjson"
)

func fromJson(data []byte) ([]byte, error) {
	lines := strings.Split(string(data), "\n")

	buf := bytes.NewBuffer(nil)
	encoder := msgpack.NewEncoder(buf)

	for _, l := range lines {
		var v interface{}
		err := json.Unmarshal([]byte(l), &v)
		if err != nil {
			return nil, err
		}

		err = encoder.EncodeValue(reflect.ValueOf(v))
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func toJson(data []byte) ([]byte, error) {
	d := msgpack.NewDecoder(bytes.NewReader(data))

	maps := make([]map[string]interface{}, 0)
	for {
		arr, err := d.DecodeInterface()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		m, ok := arr.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid record in msgpack binary")
		}

		maps = append(maps, m)
	}

	encoded, err := json.Marshal(maps)
	if err != nil {
		return nil, err
	}

	return encoded, nil
}

func main() {
	in := flag.String("in", "", "input file path")
	out := flag.String("out", "", "out file path")
	flag.Parse()

	var subcommand string
	if len(flag.Args()) >= 1 {
		subcommand = flag.Args()[0]
	}

	var inData []byte
	if *in != "" {
		d, err := ioutil.ReadFile(*in)
		if err != nil {
			log.Fatal(err)
		}
		inData = d
	} else {
		d, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		inData = d
	}

	var outData []byte
	switch subcommand {
	case subCommandToJson:
		d, err := toJson(inData)
		if err != nil {
			log.Fatal(err)
		}
		outData = d
	case subCommandFromJson:
		d, err := fromJson(inData)
		if err != nil {
			log.Fatal(err)
		}
		outData = d
	default:
		log.Fatalf("call with subcommand, tojson or fromjson")
	}

	if *out != "" {
		err := ioutil.WriteFile(*out, outData, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err := os.Stdout.Write(outData)
		if err != nil {
			log.Fatal(err)
		}
	}
}
