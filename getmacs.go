package main

import (
	"fmt"
	"strings"

	"github.com/ales999/cisaccs"
)

// GetMacsFromCisco забрать mac-address тфблицу из cisco
func GetMacsFromCisco() {

	// Проверим что файл для записи задан.
	useoutfile := len(cli.Getmacs.Outfile) > 1 // exlude "-" with name
	var outtofile []string

	fmt.Println("Use output file:", useoutfile)

	cmds := []string{"sh mac address-table | e CPU"}

	acc := cisaccs.NewCisAccount(cli.Getmacs.CisFileName, cli.Getmacs.PwdFileName)

	for _, cishost := range cli.Getmacs.Hosts {
		// Включим метку для парсинга:
		fmt.Println("hostgetmac:", cishost)
		//fmt.Println("Host:", cishost, "Port:", cli.Getmacs.PortSsh)
		out, err := acc.OneCisExecuteSsh(cishost, cli.Getmacs.PortSsh, cmds)
		if err != nil {
			fmt.Printf("Error get data from host %s: %v\n", cishost, err)
			continue
		}
		// Перебираем полученные строки
		for _, line := range out {
			//line = strings.TrimSpace(line)
			// Если есть необходимость пропускать что содержит исключенное
			if len(cli.Getmacs.ExclString) > 0 {
				// Тогда пропускаем
				if strings.Contains(line, cli.Getmacs.ExclString) {
					continue
				}
			}
			// If set save to file
			if useoutfile {
				outtofile = append(outtofile, line+"\n")

			} else {
				// Print
				fmt.Println(line)
			}
		}

	}
	if useoutfile {
		WriteFile(outtofile, cli.Getmacs.Outfile)
	}
}
