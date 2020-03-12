package main

import (
	"flag"
	"io"
	"io/ioutil"
	"os"

	"github.com/nikandfor/tlog"
)

type writer struct {
	w    io.Writer
	last byte
}

func (w *writer) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	w.last = p[len(p)-1]

	return w.w.Write(p)
}

func (w *writer) Newline() error {
	if w.last == '\n' {
		return nil
	}

	_, err := w.w.Write([]byte{'\n'})

	return err
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"-"}
	}

	w := &writer{
		w: os.Stdout,
	}

	for _, a := range args {
		func() {
			var r io.ReadCloser
			if a == "-" {
				r = ioutil.NopCloser(os.Stdin)
			} else {
				f, err := os.Open(a)
				if err != nil {
					tlog.Fatalf("open %q: %v", a, err)
				}

				defer func() {
					err := f.Close()
					if err != nil {
						tlog.Fatalf("close %q: %v", a, err)
					}
				}()

				r = f
			}

			_, err := io.Copy(w, r)
			if err != nil {
				tlog.Fatalf("copy %q: %v", a, err)
			}
		}()
	}

	err := w.Newline()
	if err != nil {
		tlog.Fatalf("add newline: %v", err)
	}
}
