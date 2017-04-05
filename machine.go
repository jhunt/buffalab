package main

import (
	"fmt"
)

type Machine struct {
	Name     string `yaml:"name"`
	MAC      string `yaml:"mac"`
	Role     string `yaml:"role"`
	Type     string `yaml:"type"`
	IP       string `yaml:"ip"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (m Machine) Reboot() error {
	switch m.Type {
	case "idrac":
		if m.IP == "" {
			return fmt.Errorf("cannot reboot '%s': no IP specified in configuration.", m.Name)
		}
		if m.Username == "" {
			return fmt.Errorf("cannot reboot '%s': no iDRAC Username specified in configuration.", m.Name)
		}
		return RebootViaIDRAC(m.Username, m.Password, m.IP)

	case "":
		return fmt.Errorf("cannot reboot '%s': no remote management type set in configuration.", m.Name)

	default:
		return fmt.Errorf("cannot reboot '%s': unrecognized remote management type '%s'.", m.Name, m.Type)
	}
}
