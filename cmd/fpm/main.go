package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"

	"github.com/UzStack/bug-lang/internal/lexar"
	"github.com/UzStack/bug-lang/internal/parser"
	"github.com/UzStack/bug-lang/internal/runtime"
	"github.com/UzStack/bug-lang/internal/runtime/enviroment"
	"github.com/UzStack/bug-lang/internal/runtime/std"
	"github.com/UzStack/bug-lang/internal/runtime/types"
)

type Job struct {
	File     string
	Response chan Result
}

type Result struct {
	Body    string
	Headers []Header
}

type Header struct {
	Key   string
	Value string
}

func Worker(jobs <-chan Job) {
	for job := range jobs {
		code, err := os.ReadFile(job.File)
		if err != nil {
			job.Response <- Result{
				Body: fmt.Sprintf("Error reading file %s: %v", job.File, err),
			}
			continue
		}
		tokenize := lexar.NewTokenize()
		tokens := tokenize.Tokenize(string(code))
		parser := parser.NewParser(tokens)
		ast := parser.CreateAST()
		env := enviroment.NewGlobalEnv()
		std.Load(env)
		var buf bytes.Buffer
		var headers []Header
		env.AssignmenVariable("print", &types.NativeFunctionDeclaration{
			Call: func(values ...any) {
				std.Pprint(&buf, values)
			},
		}, -1)
		env.AssignmenVariable("header", &types.NativeFunctionDeclaration{
			Call: func(key any, value any) {
				headers = append(headers, Header{
					Key:   key.(*types.StringValue).Value,
					Value: value.(*types.StringValue).Value,
				})
			},
		}, -1)
		runtime.Interpreter(ast, env)
		job.Response <- Result{
			Body:    buf.String(),
			Headers: headers,
		}
		close(job.Response)
	}
}

// handler funksiyasida jobs kanalini parametr sifatida uzatamiz
func handler(w http.ResponseWriter, r *http.Request, jobs chan<- Job) {
	params := fcgi.ProcessEnv(r)
	file := params["DOCUMENT_ROOT"] + params["DOCUMENT_URI"]
	result := make(chan Result)
	// Faylni workerga yuborish
	jobs <- Job{File: file, Response: result}
	res := <-result
	for _, header := range res.Headers {
		w.Header().Set(header.Key, header.Value)
	}
	w.Write([]byte(res.Body))
}

func main() {
	// jobs va results kanallarini yaratamiz
	jobs := make(chan Job)

	// Workerlarni ishga tushiramiz
	for i := 0; i < 4; i++ { // 4 worker ishga tushirilsin
		go Worker(jobs)
	}

	// Handlerga jobs va results kanallarini uzatamiz
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, jobs)
	})

	// UNIX socketni ochamiz
	sockPath := "/tmp/bug-fpm.sock"
	os.Remove(sockPath)
	listen, err := net.Listen("unix", sockPath)
	if err != nil {
		panic(err.Error())
	}
	defer listen.Close()

	fmt.Println("Started server " + sockPath + " 🚀")

	// FastCGI serverni ishga tushiramiz
	go func() {
		err := fcgi.Serve(listen, http.DefaultServeMux)
		if err != nil {
			fmt.Println("Xatolik:", err)
		}
	}()

	select {}
}
