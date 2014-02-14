package main

import (
	"flag"
	"os"
	"io/ioutil"
	"github.com/steveyen/gkvlite"
	"fmt"
)

func main() {
	flag.Parse()

	for _, fn := range flag.Args() {
		err := compress(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on file %q: %v", fn, err)
		}
	}
}

func compress(fn string) (err error) {
	f, err := os.OpenFile(fn, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer func() {
		err_ := f.Close()
		if err == nil {
			err = err_
		}
	}()

	tmpf, err := ioutil.TempFile("", "gkvlite")
	if err != nil {
		return err
	}
	defer os.Remove(tmpf.Name())
	defer tmpf.Close()

	store, err := gkvlite.NewStore(f)
	if err != nil {
		return err
	}

	store2, err := store.CopyTo(tmpf, 1000000)
	if err != nil {
		return err
	}

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = store2.CopyTo(f, 1000000)
	if err != nil {
		return err
	}

	return nil
}
