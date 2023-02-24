package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type writer struct {
	w     io.Writer
	count int
	last  byte
}

func (w *writer) Write(p []byte) (n int, err error) {
	n, err = w.w.Write(p)

	w.count += n
	if n < len(p) {
		w.last = p[n-1]
	}

	return
}

func (w *writer) Newline() error {
	if w.last == '\n' {
		return nil
	}

	_, err := w.w.Write([]byte{'\n'})

	return err
}

var force = flag.Bool("f", false, "force. add newline only on empty output")

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
					log.Fatalf("open %q: %v", a, err)
				}

				defer func() {
					err := f.Close()
					if err != nil {
						log.Fatalf("close %q: %v", a, err)
					}
				}()

				r = f
			}

			_, err := io.Copy(w, r)
			if err != nil {
				log.Fatalf("copy %q: %v", a, err)
			}
		}()
	}

	if !*force && w.count == 0 {
		return
	}

	err := w.Newline()
	if err != nil {
		log.Fatalf("add newline: %v", err)
	}
}
