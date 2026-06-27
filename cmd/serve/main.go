// bug-serve — nginx kerak bo'lmagan oddiy HTTP dev server.
// Barcha so'rovlarni root'dagi index.bug ga yo'naltiradi (front controller).
//
//	go run ./cmd/serve --root examples/webapp --addr :8080
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/UzStack/bug-lang/internal/runner"
)

func main() {
	os.Setenv("BUG_FPM", "1") // web rejimi: ffi o'chiq

	root := flag.String("root", ".", "document root (index.bug shu yerda)")
	addr := flag.String("addr", ":8080", "tinglash manzili")
	flag.Parse()

	absRoot, _ := filepath.Abs(*root)
	index := filepath.Join(absRoot, "index.bug")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		runner.Run(index, absRoot+"/", r).Apply(w)
	})

	fmt.Printf("bug-serve 🚀  http://localhost%s  (root: %s)\n", *addr, absRoot)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}
