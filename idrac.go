package main

import (
	"bytes"
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func RebootViaIDRAC(username, password, ip string) error {
	client, err := ssh.Dial("tcp", ip+":22", &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(host string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var out bytes.Buffer
	session.Stdout = &out
	if err := session.Run("racadm serveraction hardreset"); err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "reset %s via idrac/racadm (output follows):%s\n", ip, out.String())
	return nil
}
