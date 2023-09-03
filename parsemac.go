package main

import (
	"fmt"
	"regexp"
)

func ParseMacs(hn string) {
	fmt.Println("Macs file:", hn)
}

func parseArpLine(line string) MacLineData {

	/*
	     1    548a.ba01.50b3    DYNAMIC     Gi0/43
	     1    b022.7a2e.5561    DYNAMIC     Gi0/43
	    19    805e.c02d.4d50    DYNAMIC     Gi0/43
	   204    0000.aa8d.ada8    DYNAMIC     Gi0/43
	*/

	//fmt.Println(line)
	//tmp := `^\s?(\d+)\s{4}(\S+)\s{4}[D|S]\w+`

	re, _ := regexp.Compile(`^Internet  (\S+)\s+[0-9|-]+\s+(\S+)\s+ARPA\s+(\S+)$`)
	res := re.FindStringSubmatch(line)

	if len(res) > 0 {
		if res[2] != "Incomplete" {
			return *NewMacLineData(res[1], res[2], res[3], "")
		}
	}

	return MacLineData{}

}
