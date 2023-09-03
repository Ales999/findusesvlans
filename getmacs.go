package main

import (
	"fmt"
	"strings"

	"github.com/ales999/cisaccs"
)

// GetMacsFromCisco забрать mac-address тфблицу из cisco
func GetMacsFromCisco() {

	// Проверим что файл для записи задан.
	useoutfile := len(cli.Getmacs.Outfile) > 0

	fmt.Println("Use output file:", useoutfile)

	cmds := []string{"sh mac address-table | e CPU"}

	acc := cisaccs.NewCisAccount(cli.Getmacs.CisFileName, cli.Getmacs.PwdFileName)

	for _, cishost := range cli.Getmacs.Hosts {
		fmt.Println(cishost)
		fmt.Println("Host:", cishost, "Port:", cli.Getmacs.PortSsh)
		out, err := acc.OneCisExecuteSsh(cishost, cli.Getmacs.PortSsh, cmds)
		if err != nil {
			fmt.Printf("Error get data from host %s: %v\n", cishost, err)
			continue
		}
		// Перебираем полученные строки
		for _, line := range out {
			// Если есть необходимость пропускать что содержит исключенное
			if len(cli.Getmacs.ExclString) > 0 {
				// Тогда пропускаем
				if strings.Contains(line, cli.Getmacs.ExclString) {
					continue
				}
			}
			// Print
			fmt.Println(line)
		}

	}
}
