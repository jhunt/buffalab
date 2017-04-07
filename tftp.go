package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pin/tftp"
)

type TFTPServer struct {
	Listen string
	Root   string
}

func NewTFTPServer(listen, root string) TFTPServer {
	srv := TFTPServer{
		Listen: listen,
		Root:   root,
	}
	return srv
}

func (srv TFTPServer) Run() {
	s := tftp.NewServer(func(path string, rf io.ReaderFrom) error {
		file, err := os.Open(srv.Root + "/" + path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: read failed: %s\n", path, err)
			return err
		}
		n, err := rf.ReadFrom(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: read failed: %s\n", path, err)
			return err
		}
		fmt.Printf("%s: %d bytes sent\n", path, n)
		return nil
	}, nil)

	err := s.ListenAndServe(srv.Listen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tftp server died: %s\n", err)
		os.Exit(1)
	}
}

func (srv TFTPServer) Install(mac, flavor, role string) {
	in, err := os.Open(srv.Root + "/tpl/" + flavor + "-" + role)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to install pxe config for flavor '%s' / role '%s', on behalf of client with mac %s: %s\n", flavor, role, mac, err)
		return
	}
	defer in.Close()

	out, err := os.Create(srv.Root + "/pxelinux.cfg/01-" + mac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to install pxe config for flavor '%s' / role '%s', on behalf of client with mac %s: %s\n", flavor, role, mac, err)
		return
	}
	defer out.Close()

	io.Copy(out, in)
}

func (srv TFTPServer) Reset(mac string) {
	err := os.Remove(srv.Root + "/pxelinux.cfg/01-" + mac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to remove pxe config for client with mac %s: %s\n", mac, err)
		return
	}
}
