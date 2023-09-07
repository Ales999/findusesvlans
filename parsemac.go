package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ales999/cisaccs"
)

/*
	type HostData struct {
		vlan  string
		iface string
	}

	type MacDb struct {
		ciscoName string
		hstdat    []HostData
	}

	func NewMacDB(ciscoName string) *MacDb {
		return &MacDb{ciscoName: ciscoName}
	}
*/

func ParseMacs(macFileName string) {

	//var ReportOut []MacDb

	mlds, err := ParseMacFile(macFileName)
	if err != nil {
		panic(err)
	}

	// Проверим что файл для записи задан.
	useoutfile := len(cli.Parsemac.Outfile) > 1    // type bool -->  true если задан выходной файл.
	usedreport := len(cli.Parsemac.Reportfile) > 1 // type bool --> true если задан файл отчета.
	var outtofile []string                         // массив вывода в файл отчета
	var reprtfile []string                         // массив вывода в файл репорта

	if useoutfile {
		fmt.Println("Use output file:", cli.Parsemac.Outfile)
	}
	if usedreport {
		fmt.Println("Use mac-report file:", cli.Parsemac.Reportfile)
	}
	//hl := make(map[string]HstVl)
	var vlans []string
	var firstVlan bool

	for _, hmld := range mlds {
		firstVlan = true
		var outstr string  // Строка для вывода готового набора  vlan-ов
		vlans = []string{} // Масив vlan-ов

		acc := cisaccs.NewCisAccount(cli.Parsemac.CisFileName, cli.Parsemac.PwdFileName)

		iface, err := acc.GetIfaceByHost(hmld.HostName)

		if err != nil {
			fmt.Println("Error:", err.Error())
			continue
		}

		if useoutfile {
			outtofile = append(outtofile, fmt.Sprintf("!--- Host: %s\n", hmld.HostName))
		} else { // Выводим на экран
			fmt.Println("-----------------------")
			fmt.Println("Host:", hmld.HostName)
		}
		if usedreport { // Для записи в файл отчета.
			reprtfile = append(reprtfile, "#-----------------------\n")
			reprtfile = append(reprtfile, fmt.Sprintf("Host: %s\n", hmld.HostName))

		}

		for _, mld := range hmld.mld {

			// Если VLAN в списке пропускаемы - смотрим далее
			if FindSkip(mld.vlan, &skipVlans) {
				continue
			}
			vlans = append(vlans, mld.vlan)
			// Debug output
			if usedreport {
				reprtfile = append(reprtfile, fmt.Sprintf("Vlan: %s,\tMac: %s\tIface: %s\n", mld.vlan, mld.mac, mld.iface))
			}

		}
		// Добавим в список VLAN-ы которые обязательно должны быть.
		vlans = append(vlans, skipVlans...)

		// Удалим дубликаты
		vlans = RemoveDuplicate(vlans)
		// Конвертируем номера vlan-ов в Int
		vlints := IntedStringToInts(vlans)
		// Сортируем
		sort.Ints(vlints)

		for _, v := range vlints {
			if firstVlan {
				if useoutfile { // Если указано вывод в файл
					outstr = fmt.Sprintf(" interface %s\n switchport trunk allowed vlan %d", iface, v)
				} else {
					fmt.Printf(" switchport trunk allowed vlan %d", v)
				}
				firstVlan = false

			} else {
				if useoutfile {
					outstr += fmt.Sprintf(",%d", v)
					//outtofile = append(outtofile, fmt.Sprintf(",%d", v))
				} else {
					fmt.Printf(",%d", v)
				}
			}
		}
		if useoutfile {
			outstr += "\n exit\n" // Выходим из редактирования интерфейса
			outtofile = append(outtofile, outstr)
		} else {
			fmt.Println()
		}
		firstVlan = true
	}

	if useoutfile {
		WriteOutFile(outtofile, cli.Parsemac.Outfile)
	}
	if usedreport {
		WriteOutFile(reprtfile, cli.Parsemac.Reportfile)
	}
}

// IntedStringToInts - сконвертировать массив string (в которых только числа) в массив Integer
func IntedStringToInts(strarr []string) []int {
	var out []int
	for _, v := range strarr {
		nmbr, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		out = append(out, nmbr)
	}
	return out
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

// ParseMacFile - открыть файл с MAC-адресами, и распасрсить его в массив хостов с данными.
func ParseMacFile(macFileName string) ([]HostMacLineData, error) {

	fmt.Println("Parse MAC file:", macFileName)

	MacLines := []MacLineData{}  // Временный массив
	var output []HostMacLineData // Исходящие данныы

	// Читаем ACL файл
	aclFile, err := os.OpenFile(macFileName, os.O_RDONLY, 0644)
	if err != nil {
		return output, fmt.Errorf("ошибка открытия файла: %s", err)
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
	var hmld HostMacLineData
	for _, s := range aclFileLines {
		tr := strings.TrimSpace(s)
		if len(tr) > 0 {

			if strings.Contains(tr, "hostgetmac:") {
				hostName = strings.TrimPrefix(tr, "hostgetmac: ")
				if len(MacLines) > 0 {
					hmld.mld = MacLines
					output = append(output, hmld)
				}
				// Новая
				hmld = *NewHostMacLineData(hostName)

			} else {
				a := parseArpLine(tr)
				if len(a.vlan) > 0 {
					MacLines = append(MacLines, a)
				}
			}
		}
	}
	// Добавляем последний проверяемый в выходной массив
	if len(MacLines) > 0 {
		hmld.mld = MacLines
		output = append(output, hmld)
	}

	return output, nil

}

func parseArpLine(line string) MacLineData {

	/*
	   1    548a.ba01.50b3    DYNAMIC     Gi0/43
	   1    b022.7a2e.5561    DYNAMIC     Gi0/43
	   19    805e.c02d.4d50    DYNAMIC     Gi0/43
	   204    0000.aa8d.ada8    DYNAMIC     Gi0/43
	*/

	re, _ := regexp.Compile(`^^(\d+)\s{4}(\S+)\s{4}([D|S]\S+)\s{4,6}(\S+)$`)
	res := re.FindStringSubmatch(line)

	if len(res) > 0 {
		return *NewMacLineData(res[1], res[2], res[4])
	}

	return MacLineData{}

}
