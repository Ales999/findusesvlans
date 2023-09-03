package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var cli struct {
	Getmacs struct {
		Hosts       []string `arg:"" name:"hosts"`
		Outfile     string   `help:"Output file" type:"path"`
		ExclString  string   `name:"exclude" help:"Exclude string if exist from mac-address line" short:"e"`
		PortSsh     int      `help:"use ssh port" default:"22" short:"p" env:"CISPORT"`
		CisFileName string   `help:"Patch to cis.yaml" type:"existingfile" env:"CISFILE"`
		PwdFileName string   `help:"Patch to passw.json" type:"existingfile" env:"CISPWDS"`
	} `cmd:"" help:"Get mac-address table from cisco hosts."`

	Parsemac struct {
		MacsFileName string `arg:"" name:"macsfile" type:"existingfile"`
	} `cmd:"" help:"Parsing file with mac-address table"`
}

func main() {

	ctx := kong.Parse(&cli,
		kong.Name("findusesvlans"),
		kong.Description("Get and Parse mac-address table for find used vlans"),
		kong.UsageOnError(),
	)

	switch ctx.Command() {
	case "getmacs <hosts>":
		GetMacsFromCisco()
	case "parsemac <macsfile>":
		fmt.Println("Use parsemac")
		ParseMacs(cli.Parsemac.MacsFileName)
	default:
		panic(ctx.Command())
	}
}
