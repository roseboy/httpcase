package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type demoCmd struct {
	cmd  *cobra.Command
	opts demoOpts
}

type demoOpts struct {
	port string
}

func newDemoCmd() *demoCmd {
	root := &demoCmd{}

	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Run an api demo server",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("httpcase demo api server start on:", root.opts.port)
			srv := &http.Server{Addr: ":" + root.opts.port}
			return srv.ListenAndServe()
		},
	}
	cmd.Flags().StringVarP(&root.opts.port, "port", "p", "8080", "server port")

	root.cmd = cmd
	return root
}

var (
	userList = make(map[string]string)
)

func init() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/user/", userHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/upload", fileHandler)
	http.HandleFunc("/callback", callbackHandle)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "httpcase api demo")
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("token") != "123456" {
		w.WriteHeader(401)
		return
	}

	log.Println("USERS-DB:", userList)

	id := getUrlArgs(r.URL.Path, 2)

	switch r.Method {
	case "GET":
		log.Println("get user by id:", id)
		if id == "" {
			w.WriteHeader(500)
			fmt.Fprint(w, `{"status":"fail"}`)
			return
		}
		if _, ok := userList[id]; !ok {
			w.WriteHeader(404)
			fmt.Fprint(w, `{"status":"fail"}`)
			return
		}
		fmt.Fprint(w, userList[id])
		return
	case "POST":
		body, _ := ioutil.ReadAll(r.Body)
		bodyStr := string(body)
		bodyStr = strings.Trim(bodyStr, "\n")
		bodyStr = strings.Trim(bodyStr, "\r")
		if id == "" {
			log.Println("add user:", bodyStr)
			id := fmt.Sprintf("uid-%d", time.Now().Unix())
			bodyStr = fmt.Sprintf(`{"id":"%s",%s`, id, bodyStr[1:])
			userList[id] = bodyStr
		} else {
			log.Println("update user:", bodyStr)
			userList[id] = bodyStr
		}
		fmt.Fprint(w, bodyStr)
		return
	case "DELETE":
		log.Println("delete user by id:", id)
		if id == "" {
			w.WriteHeader(500)
			fmt.Fprint(w, `{"status":"fail"}`)
			return
		}
		if _, ok := userList[id]; !ok {
			w.WriteHeader(404)
			fmt.Fprint(w, `{"status":"fail"}`)
			return
		}

		delete(userList, id)
		fmt.Fprint(w, `{"status":"success"}`)
		return
	}

}

func fileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var buf bytes.Buffer
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(&buf, file)

	buf.Reset()
	fmt.Println(r.PostForm)
	fmt.Println(header)
	fmt.Fprintln(w, "ok")
}

func callbackHandle(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return
	}
	fmt.Println(string(body))
}

func getUrlArgs(url string, index int) string {
	args := strings.Split(url, "/")
	if len(args) <= index {
		return ""
	}
	return args[index]
}
