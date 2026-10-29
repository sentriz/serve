package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	flag := flag.CommandLine
	flag.Usage = func() {
		fmt.Fprintf(flag.Output(), "Usage:\n")
		fmt.Fprintf(flag.Output(), "  $ %s [<options>] download [<download options>] [<path>]\n", flag.Name())
		fmt.Fprintf(flag.Output(), "  $ %s [<options>] upload [<dir>]\n", flag.Name())
		fmt.Fprintf(flag.Output(), "\n")
		fmt.Fprintf(flag.Output(), "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.Output(), "\n")
		fmt.Fprintf(flag.Output(), "See also:\n")
		fmt.Fprintf(flag.Output(), "  $ %s download -h\n", flag.Name())
	}
}

func main() {
	addr := flag.String("addr", ":8000", "listen address")
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("no command provided")
	}

	var handler http.Handler
	switch command, args := flag.Arg(0), flag.Args()[1:]; command {
	case "download":
		flag := flag.NewFlagSet(command, flag.ExitOnError)
		attachment := flag.Bool("attachment", false, "set headers to prompt to download")
		flag.Parse(args)

		path := cmp.Or(flag.Arg(0), ".")
		info, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}

		switch {
		case info.Mode().IsRegular():
			handler = fileHandler(path, *attachment)
		case *attachment:
			log.Fatal("attachment needs a file")
		default:
			handler = http.FileServer(http.Dir(path))
		}

	case "upload":
		flag := flag.NewFlagSet(command, flag.ExitOnError)
		flag.Parse(args)

		dir := cmp.Or(flag.Arg(0), ".")
		info, err := os.Stat(dir)
		if err != nil {
			log.Fatal(err)
		}
		if !info.IsDir() {
			log.Fatal("upload needs a dir")
		}

		handler = uploadHandler(dir)

	default:
		log.Fatalf("unknown command %q", command)
	}

	log.Fatal(http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requests %q", r.RemoteAddr, r.URL)
		handler.ServeHTTP(w, r)
	})))
}

func fileHandler(path string, attachment bool) http.Handler {
	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": filepath.Base(path)})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if attachment {
			w.Header().Set("Content-Disposition", disposition)
			w.Header().Set("Content-Type", "application/octet-stream")
		}
		http.ServeFile(w, r, path)
	})
}

func uploadHandler(dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			io.WriteString(w, uploadForm)
			return
		}

		mr, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for {
			part, err := mr.NextPart()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if part.FileName() == "" {
				continue
			}

			dest := filepath.Join(dir, part.FileName())
			if err := saveFile(dest, part); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("%s uploaded %q", r.RemoteAddr, dest)
		}
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	})
}

const uploadForm = `<!doctype html>
<meta name="viewport" content="width=device-width">
<form method="post" enctype="multipart/form-data">
<input type="file" name="file" multiple>
<button>upload</button>
</form>
`

func saveFile(dest string, r io.Reader) error {
	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	if err := errors.Join(err, f.Close()); err != nil {
		os.Remove(dest)
		return err
	}
	return nil
}
