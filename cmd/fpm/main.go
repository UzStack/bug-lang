package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path/filepath"
	"strings"

	"github.com/UzStack/bug-lang/internal/runner"
)

type Job struct {
	File     string
	Response chan runner.Response
	Request  *http.Request
	BaseDir  string
}

func Worker(jobs <-chan Job) {
	for job := range jobs {
		job.Response <- processJob(job)
		close(job.Response)
	}
}

// processJob bitta so'rovni ishlaydi (runner.Run panic'ni 500 ga aylantiradi —
// shunday qilib yomon skript butun worker goroutinesini o'ldirmaydi).
func processJob(job Job) runner.Response {
	return runner.Run(job.File, job.BaseDir, job.Request)
}

// handler so'rovni mos faylga moslab worker'ga uzatadi.
func handler(w http.ResponseWriter, r *http.Request, jobs chan<- Job) {
	params := fcgi.ProcessEnv(r)
	root := filepath.Clean(params["DOCUMENT_ROOT"])
	// filepath.Join "../" larni tozalaydi; natija DOCUMENT_ROOT ichida qolishi shart
	file := filepath.Join(root, params["DOCUMENT_URI"])
	if file != root && !strings.HasPrefix(file, root+string(os.PathSeparator)) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 Forbidden"))
		return
	}
	result := make(chan runner.Response)
	jobs <- Job{File: file, Response: result, Request: r, BaseDir: root + "/"}
	(<-result).Apply(w)
}

func main() {
	// web rejimida ffi (ixtiyoriy .so yuklash) o'chiriladi
	os.Setenv("BUG_FPM", "1")

	jobs := make(chan Job)
	for range 4 { // 4 worker
		go Worker(jobs)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, jobs)
	})

	sockPath := "/tmp/bug-fpm.sock"
	os.Remove(sockPath)
	listen, err := net.Listen("unix", sockPath)
	if err != nil {
		panic(err.Error())
	}
	defer listen.Close()

	fmt.Println("Started server " + sockPath + " 🚀")

	go func() {
		if err := fcgi.Serve(listen, http.DefaultServeMux); err != nil {
			fmt.Println("Xatolik:", err)
		}
	}()

	select {}
}
