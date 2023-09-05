package main

import (
	"fmt"
	"strings"

	"github.com/ales999/cisaccs"
)

// GetMacsFromCisco забрать mac-address тфблицу из cisco
func GetMacsFromCisco() {

	// Проверим что файл для записи задан.
	useoutfile := len(cli.Getmacs.Outfile) > 1 // type bool -->  exlude "-" with name
	var outtofile []string

	if useoutfile {
		fmt.Println("Use output file:", cli.Getmacs.Outfile)
	}

	cmds := []string{"sh mac address-table | e CPU"}

	acc := cisaccs.NewCisAccount(cli.Getmacs.CisFileName, cli.Getmacs.PwdFileName)

	for _, cishost := range cli.Getmacs.Hosts {
		// Включим метку для парсинга:
		fmt.Println("hostgetmac:", cishost)

		// Получаем массив строк что вернет нам cisco
		sshouts, err := acc.OneCisExecuteSsh(cishost, cli.Getmacs.PortSsh, cmds)
		if err != nil {
			fmt.Printf("Error get data from host %s: %v\n", cishost, err)
			continue
		}

		// Перебираем полученные строки
		for _, line := range sshouts {
			usedLine := true // Признак что что строка не содержит одну из исключенных, если они есть.
			line = strings.TrimSpace(line)
			// Если есть необходимость пропускать что содержит исключенное
			if len(cli.Getmacs.ExclString) > 0 {
				// Тогда пропускаем
				for _, exclstr := range cli.Getmacs.ExclString {
					if strings.Contains(line, exclstr) {
						usedLine = false // Запрет использовать данную строку
					}
				}
			}

			if !usedLine {
				continue // go get new line
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
		WriteOutFile(outtofile, cli.Getmacs.Outfile)
	}
}
