// Copyright (c) 2018, IBM
// Author(s): Brandon Lum, Ricardo Koller, Dan Williams
//
// Permission to use, copy, modify, and/or distribute this software for
// any purpose with or without fee is hereby granted, provided that the
// above copyright notice and this permission notice appear in all
// copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL
// WARRANTIES WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE
// AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL
// DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE, DATA
// OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
// TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.

// +build linux

package runnc_cont

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/retrage/runnc/nabla-lib/network"
)

type lklArgsNetwork struct {
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Param   string `json:"param,omitempty"`
	Mtu     string `json:"mtu,omitempty"`
	Ip      string `json:"ip,omitempty"`
	Ipv6    string `json:"ipv6,omitempty"`
	IfGw    string `json:"ifgateway,omitempty"`
	IfGw6   string `json:"ifgateway6,omitempty"`
	Mac     string `json:"mac,omitempty"`
	Mask    string `json:"masklen,omitempty"`
	Mask6   string `json:"masklen6,omitempty"`
	Neigh   string `json:"neigh,omitempty"`
	Qdisc   string `json:"qdisc,omitempty"`
    Offload string `json:"offload,omitempty"`
}

type lklArgs struct {
	If      []lklArgsNetwork `json:"interfaces"`
	Gw      string `json:"gateway,omitempty"`
	Gw6     string `json:"gateway6,omitempty"`
	Debug   string `json:"debug,omitempty"`
	Mount   string `json:"mount,omitempty"`
	Cpu     string `json:"singlecpu,omitempty"`
	Sysctl  string `json:"sysctl,omitempty"`
	Cmdline string `json:"boot_cmdline,omitempty"`
	Dump    string `json:"dump,omitempty"`
	Delay   string `json:"delay_main,omitempty"`
}

// CreateLklArgs returns the config string for lkl (a json)
func CreateLklArgs(ip net.IP, mask net.IPMask, gw net.IP, mac string) (string, error) {

	// XXX: Due to bug in: https://github.com/nabla-containers/runnc/issues/40
	// If we detect a /32 mask, we set it to 1 as a "fix", and hope we are in
	// the same subnet... (working on a fix for mask:0)
	cidr := strconv.Itoa(network.MaskCIDR(mask))
	if cidr == "32" {
		fmt.Printf("WARNING: Changing CIDR from 32 to 1 due to Issue https://github.com/nabla-containers/runnc/issues/40\n")
		cidr = "1"
	}

	net := lklArgsNetwork{
		Type: "rumpfd",
		Name: "tap",
		Ip:   ip.String(),
		Mask: cidr,
		IfGw: gw.String(),
		Mac:  mac,
	}

	la := &lklArgs{
		If:       []lklArgsNetwork{net},
		Debug:    "1",
		Cpu:      "1",
		Delay:    "50000",
		Sysctl:   "net.ipv4.tcp_wmem=4096 87380 2147483647",
	}

	b, err := json.Marshal(la)
	if err != nil {
		return "", fmt.Errorf("error with lkl config json: %v", err)
	}

    s := string(b)
	s = strings.Replace(s, "\"", "\\\"", -1)
	s = strings.Replace(s, " ", "\\ ", -1)

	return s, nil
}
