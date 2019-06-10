/*
functionality:
- perform forward lookup returning all ips for host
- perform reverse lookup returning all hosts for ip
- recursive mode where the host/ip looking is performed on the result of the initial request
- return output in json, -oJ
eventually:
- perform arin queries
- perform whois queries
- let recursive take a number for the number of times to recurse?
- don't think ipv4/6 is being respected
*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

type Addr struct {
	Addr    string  `json:"addr"`
	Lookups []*Addr `json:"lookups"`
	// Type     string  `json:"type"` // host, ip
	resolver *net.Resolver
}

func (i *Addr) typeIP() bool {
	if ip := net.ParseIP(i.Addr); ip != nil {
		return true
	}
	return false
}

func (i *Addr) typeHost() bool {
	return !i.typeIP()
}

// may need to remove duplicates
// func flatten(addrs []*Addr) []*Addr {
// 	var flat []*Addr
// 	for _, addr := range addrs {
// 		if len(addr.Lookups) > 0 {
// 			// flat = append(flat, addr) need to blank out the lookups in this case
// 			flat = append(flat, flatten(addr.Lookups)...)
// 		} else {
// 			flat = append(flat, addr)
// 		}
// 	}
// 	return flat
// }

func (i *Addr) resolve(ipv4, ipv6 bool) {
	if i.typeIP() {
		names, err := i.resolver.LookupAddr(context.Background(), i.Addr)
		// typecheck error ???
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}
		for _, n := range names {
			i.Lookups = append(i.Lookups, &Addr{Addr: n, resolver: i.resolver})
		}
	}
	if i.typeHost() {
		cname, err := net.LookupCNAME(i.Addr)
		if err == nil {
			i.Lookups = append(i.Lookups, &Addr{Addr: cname, resolver: i.resolver})
		}
		// net.LookupCNAME
		ips, err := net.LookupHost(i.Addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}
		for _, ip := range ips {
			typedIP := net.ParseIP(ip)
			if typedIP == nil {
				continue
			}
			isIPv4 := !(typedIP.To4() == nil)
			isIPv6 := typedIP != nil && !isIPv4
			if (ipv4 && isIPv4) || (ipv6 && isIPv6) {
				i.Lookups = append(i.Lookups, &Addr{Addr: ip, resolver: i.resolver})
			}
		}
	}
}

// automatically detect Addr as hostname or ip
func main() {
	var ipv4, ipv6, recursive, oj, ojp bool
	flag.BoolVar(&ipv4, "4", false, "ipv4 responses only")
	flag.BoolVar(&ipv6, "6", false, "ipv6 responses only")
	flag.BoolVar(&recursive, "r", false, "recursive (resolve lookups)")
	// flag.BoolVar(&flat, "f", false, "flatten output") // json only?
	flag.BoolVar(&oj, "j", false, "output json")
	flag.BoolVar(&ojp, "jp", false, "output json (pretty)")
	// _ = flat // not sure what flat really means right now...

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "error: provide at least one ip/hostname as an argument")
		os.Exit(1)
	}

	if !ipv4 && !ipv6 {
		ipv4 = true
		ipv6 = true
	}

	resolver := &net.Resolver{PreferGo: true}

	var addrs []*Addr
	for _, arg := range flag.Args() {
		addrs = append(addrs, &Addr{Addr: arg, resolver: resolver})
	}

	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr *Addr) {
			defer wg.Done()
			addr.resolve(ipv4, ipv6)
			if recursive {
				var wgr sync.WaitGroup
				for _, addr := range addr.Lookups {
					wgr.Add(1)
					go func(addr *Addr) {
						defer wgr.Done()
						addr.resolve(ipv4, ipv6)
					}(addr)
				}
				wgr.Wait()
			}
		}(addr)
	}
	wg.Wait()

	if oj {
		b, _ := json.Marshal(addrs)
		fmt.Printf("%s\n", string(b))
		return
	}
	if ojp {
		b, _ := json.MarshalIndent(addrs, "", "  ")
		fmt.Printf("%s\n", string(b))
		return
	}

	for i, addr := range addrs {
		if i > 0 {
			fmt.Println()
		}
		fmt.Println(addr.Addr)
		for _, addr := range addr.Lookups {
			fmt.Printf("\t%s\n", addr.Addr)
			if recursive {
				for _, addr := range addr.Lookups {
					fmt.Printf("\t\t%s\n", addr.Addr)
				}
			}
		}
	}
}
