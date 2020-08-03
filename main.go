package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/vmihailenco/msgpack/v4"
)

const (
	// overwrite size; default value is 64k
	maxScanTokenSize = 4 * 1024 * 1024

	subCommandToJson   = "tojson"
	subCommandFromJson = "fromjson"
)

func fromJson(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	buf := make([]byte, bufio.MaxScanTokenSize)
	s.Buffer(buf, maxScanTokenSize)

	encoder := msgpack.NewEncoder(w)

	for s.Scan() {
		var v interface{}
		err := json.Unmarshal(s.Bytes(), &v)
		if err != nil {
			return err
		}

		err = encoder.EncodeValue(reflect.ValueOf(v))
		if err != nil {
			return err
		}
	}

	return s.Err()
}

func toJson(r io.Reader, w io.Writer) error {
	d := msgpack.NewDecoder(r)

	for {
		arr, err := d.DecodeInterface()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		m, ok := arr.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid record in msgpack binary")
		}

		marshaled, err := json.Marshal(m)
		if err != nil {
			return err
		}
		if _, err := w.Write(marshaled); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	in := flag.String("in", "", "input file path")
	out := flag.String("out", "", "out file path")
	flag.Parse()

	var subcommand string
	if len(flag.Args()) >= 1 {
		subcommand = flag.Args()[0]
	}

	var r io.Reader
	if *in != "" {
		rr, err := os.Open(*in)
		defer rr.Close()
		if err != nil {
			log.Fatal(err)
		}
		r = rr
	} else {
		r = os.Stdin
	}

	var w io.WriteCloser
	if *out != "" {
		ww, err := os.Create(*out)
		if err != nil {
			log.Fatal(err)
		}
		w = ww
	} else {
		w = os.Stdout
	}
	defer w.Close()

	switch subcommand {
	case subCommandToJson:
		if err := toJson(r, w); err != nil {
			log.Fatal(err)
		}
	case subCommandFromJson:
		if err := fromJson(r, w); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("call with subcommand, tojson or fromjson")
	}

}
