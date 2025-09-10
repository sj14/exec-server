package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// TODO:
// - add auth
// - allowlist commands
// - denylist commands

func lookupEnvString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func main() {
	var (
		listen     = flag.String("listen", lookupEnvString("LISTEN", ":8080"), "listen for connections")
		workingDir = flag.String("working-dir", lookupEnvString("WORKING_DIR", ""), "the default directory the server should run commands in")
		// allowList  = flag.String("allow-list", lookupEnvString("ALLOW_LIST", ""), "only allow the configured commands, e.g. 'cmd1,cmd2' (ignored when empty)")
		// denyList   = flag.String("deny-list", lookupEnvString("DENY_LIST", ""), "deny the configured commands, e.g. 'cmd1,cmd2'")
	)
	flag.Parse()
	os.Chdir(*workingDir)

	http.HandleFunc("POST /exec", execHandler(nil, nil))

	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.Fatalln(err)
	}
}

func execHandler(allowList, denylist []string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Split the input to separate the command and the arguments.
		args := strings.Split(string(body), " ")

		// Prepare the command to execute.
		cmd := exec.Command(args[0], args[1:]...)

		// Command output is redirected to the http response.
		cmd.Stderr = w
		cmd.Stdout = w

		if err = cmd.Start(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if req.Header.Get("fire-and-forget") != "true" {
			if err = cmd.Wait(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
}
