package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type APIServer struct {
	Config Config     `json:"-"`
	TFTP   TFTPServer `json:"-"`

	Flavor string            `json:"flavor"`
	Nodes  map[string]string `json:"nodes"`
}

func matches(r *http.Request, pat string) bool {
	ok, _ := regexp.MatchString(pat, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
	return ok
}

func bail(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprintf(w, "oops... %s\n", err)
}

func NewAPIServer(config Config) APIServer {
	srv := APIServer{
		Config: config,
		Flavor: "none",
	}

	srv.TFTP = NewTFTPServer(config.ListenTFTP, config.Root)

	srv.Nodes = make(map[string]string)
	srv.SyncNodes()
	return srv
}

func (srv *APIServer) Run() {
	go srv.TFTP.Run()

	http.Handle("/api/", srv)
	http.Handle("/", http.FileServer(http.Dir(srv.Config.Root)))
	http.ListenAndServe(srv.Config.ListenHTTP, nil)
}

func (srv *APIServer) SyncNodes() {
	seen := make(map[string]bool)

	for _, m := range srv.Config.Machines {
		if _, ok := srv.Nodes[m.Name]; !ok {
			srv.Nodes[m.Name] = "new"
		}
		seen[m.Name] = true
	}

	for k := range srv.Nodes {
		if !seen[k] {
			delete(srv.Nodes, k)
		}
	}
}

func (srv *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if matches(r, `^GET /api/status$`) {
		b, err := json.Marshal(srv)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "failed to jsonify: %s\n", err)
			return
		}
		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s\n", string(b))
		return
	}

	if matches(r, `^POST /api/install$`) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			bail(w, err)
			return
		}
		flavor := strings.TrimSuffix(string(b), "\n")

		if !srv.Config.ValidFlavor(flavor) {
			w.WriteHeader(400)
			fmt.Fprintf(w, "invalid flavor '%s'; must be one of %s\n", flavor, srv.Config.ValidFlavors())
			return
		}

		fmt.Fprintf(os.Stderr, "reconfiguring lab to install '%s'\n", flavor)
		for _, m := range srv.Config.Machines {
			fmt.Fprintf(os.Stderr, " - %s\n", m.Name)
			srv.TFTP.Install(m.MAC, flavor, m.Role)
			if err = m.Reboot(); err != nil {
				fmt.Fprintf(os.Stderr, "  %s (skipping)\n", err)
			}
			srv.Nodes[m.Name] = "installing"
		}

		srv.Flavor = flavor
		w.WriteHeader(204)
		return
	}

	if matches(r, `^POST /api/[^/]*/[^/]*$`) {
		srv.SyncNodes()

		re := regexp.MustCompile(`^/api/([^/]*)/([^/]*)$`)
		x := re.FindStringSubmatch(r.URL.Path)
		if _, ok := srv.Nodes[x[1]]; !ok {
			w.WriteHeader(404)
			fmt.Fprintf(w, "unrecognized lab node '%s'\n", x[1])
			return
		}

		found := false
		for _, m := range srv.Config.Machines {
			if m.Name == x[1] {
				found = true
				fmt.Fprintf(os.Stderr, "reseting TFTP configuration for %s (mac %s)\n", m.Name, m.MAC)
				srv.TFTP.Reset(m.MAC)
			}
		}
		if !found {
			w.WriteHeader(404)
			fmt.Fprintf(w, "unrecognized lab node '%s'\n", x[1])
			return
		}

		srv.Nodes[x[1]] = x[2]
		w.WriteHeader(204)
		return
	}

	w.WriteHeader(404)
	fmt.Fprintf(w, "%s is not a valid API endpoint\n", r.URL.Path)
}
