package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type HstVl struct {
	hstname string
	iface   string
}

func ParseMacs(macFileName string) {

	mlds, err := ParseMacFile(macFileName)
	if err != nil {
		panic(err)
	}

	var ciscoHostName string
	lastHostName := "nonehost"

	//hl := make(map[string]HstVl)

	for _, mld := range mlds {
		ciscoHostName = mld.hostname
		if len(ciscoHostName) > 0 {
			if !strings.EqualFold(ciscoHostName, lastHostName) {
				fmt.Printf("- %s:\n", ciscoHostName)
				lastHostName = ciscoHostName
			}
			// Если VLAN в списке пропускаемы - смотрим далее
			if FindSkip(mld.vlan, &skipVlans) {
				continue
			}
			fmt.Printf("Hst: %s, vlan: %s, iface: %s\n", ciscoHostName, mld.vlan, mld.iface)

			//hl[mld.vlan] = HstVl{hstname: mld.hostname, iface: mld.iface}
		}
	}
	/*
		for hli, hle := range hl {
			fmt.Printf("Vlan%s, Last Iface: %s, Host: %s\n", hli, hle.iface, hle.hstname)
		}
	*/

}

// FindSkip - вернуть true если vl есть в массиве skip
func FindSkip(vl string, skip *[]string) bool {

	for _, s := range *skip {
		if strings.EqualFold(vl, s) {
			return true
		}
	}
	return false
}

func ParseMacFile(macFileName string) ([]MacLineData, error) {

	fmt.Println("Parse MAC file:", macFileName)

	MacLines := []MacLineData{}

	// Читаем ACL файл
	aclFile, err := os.OpenFile(macFileName, os.O_RDONLY, 0644)
	if err != nil {
		return MacLines, fmt.Errorf("ошибка открытия файла: %s", err)
	}
	defer aclFile.Close()

	scanner := bufio.NewScanner(aclFile)
	scanner.Split(bufio.ScanLines)

	// Строки ACL файла
	var aclFileLines []string

	for scanner.Scan() {
		aclFileLines = append(aclFileLines, scanner.Text())
	}
	aclFile.Close()

	var hostName string
	for _, s := range aclFileLines {
		tr := strings.TrimSpace(s)
		if len(tr) > 0 {
			if strings.Contains(tr, "hostgetmac:") {
				hostName = strings.TrimPrefix(tr, "hostgetmac: ")
			}
			a := parseArpLine(tr, hostName)
			if len(a.vlan) > 0 {
				MacLines = append(MacLines, a)
			}
			// Добавим разделитель между .
			// if strings.Contains(tr, "--") {
			// MacLines = append(MacLines, AclLineData{original: tr, iface: "original"})
			// }

		}
	}

	return MacLines, nil

}

func parseArpLine(line string, hostName string) MacLineData {

	/*
	   1    548a.ba01.50b3    DYNAMIC     Gi0/43
	   1    b022.7a2e.5561    DYNAMIC     Gi0/43
	   19    805e.c02d.4d50    DYNAMIC     Gi0/43
	   204    0000.aa8d.ada8    DYNAMIC     Gi0/43
	*/

	re, _ := regexp.Compile(`^^(\d+)\s{4}(\S+)\s{4}([D|S]\S+)\s{4,6}(\S+)$`)
	res := re.FindStringSubmatch(line)

	if len(res) > 0 {
		return *NewMacLineData(res[1], res[2], res[3], res[4], hostName)
	}

	return MacLineData{}

}
