/*
Question: We have flow data which includes destination IP addresses and
bits per second (bps) for each flow, and need to aggregate over destination
prefixes (cidrs) based on bgp routing table. The flow table should be
added to all matching routes. i.e. a flow for 192.168.0.1 should be included
in the result for 192.168.0.0/24 and 192.168.0.0/20.

The result should be ordered by bps, highest to lowest. If there are
entries with equal bps then most specific route should be first,
i.e. /24 should come before /20. Routes without any traffic
shouldn't be included in the returned data.

*/

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sort"
)

type flow struct {
    address string
    bps     int
}

/*
Complete the 'aggregateOverPrefixes' function below by returning
an aggregated series of flows containing prefixes and bps.
This should be sorted descending by bps first, then prefix length.

Please explain your approach as you go. You are welcome to use online
resources (docs, etc) but please share your screen so that we can
assist you with questions you might have about what you're seeing.
*/
func aggregateOverPrefixes(flows []flow, routes []string) []flow {
    routeMap := make(map[string]int)

	for _, route := range routes {
		_, ipNet, err := net.ParseCIDR(route)
		if err != nil {
			continue
		}

		for _, f := range flows {
			ip := net.ParseIP(f.address)
			if ip != nil && ipNet.Contains(ip) {
				routeMap[route] += f.bps
			}
		}
	}

	var result []flow
	for route, totalBps := range routeMap {
		if totalBps > 0 {
			result = append(result, flow{address: route, bps: totalBps})
		}
	}

	// Sort the result by bps descending, then by prefix length ascending
	sort.Slice(result, func(i, j int) bool {
		if result[i].bps == result[j].bps {
			_, ipNetI, _ := net.ParseCIDR(result[i].address)
			_, ipNetJ, _ := net.ParseCIDR(result[j].address)
			onesI, _ := ipNetI.Mask.Size()
			onesJ, _ := ipNetJ.Mask.Size()
			return onesI > onesJ
		}
		return result[i].bps > result[j].bps
	})

	return result
}	


// Do not edit this function - your solution should be
// contained within aggregateOverPrefixes
func main() {
    stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
    checkError(err)

    defer stdout.Close()

    writer := bufio.NewWriterSize(stdout, 16*1024*1024)

    flows := []flow{
        {"192.168.0.1", 100},
        {"192.168.0.100", 200},
        {"10.100.100.1", 50},
        {"172.16.78.4", 100},
        {"172.16.100.1", 200},
        {"192.168.100.100", 400},
    }
    routes := []string{
        "192.168.0.0/24",
        "192.168.0.0/20",
        "192.168.0.0/16",
        "10.0.0.0/8",
        "10.0.0.0/16",
        "10.100.0.0/16",
        "10.100.0.0/24",
        "10.100.100.0/24",
        "172.16.64.0/18",
        "172.16.64.0/19",
    }
    result := aggregateOverPrefixes(flows, routes)

    for _, i := range result {
        fmt.Fprintf(writer, "%s:%d\n", i.address, i.bps)
    }

    writer.Flush()
}

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}