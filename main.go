package main

import (
	"github.com/alecthomas/kong"
)

var cli struct {
	Getmacs struct {
		Hosts          []string `arg:"" name:"hosts"`
		Outfile        string   `help:"Output file" short:"o" type:"path"`
		ForceOverwrite bool     `name:"force" help:"Force overwrite output file file" short:"f"`
		ExclString     []string `name:"exclude" help:"Exclude string if exist from mac-address line" short:"e"`
		PortSsh        int      `help:"use ssh port" default:"22" short:"p" env:"CISPORT"`
		CisFileName    string   `help:"Patch to cis.yaml" type:"existingfile" env:"CISFILE"`
		PwdFileName    string   `help:"Patch to passw.json" type:"existingfile" env:"CISPWDS"`
	} `cmd:"" help:"Get mac-address table from cisco hosts."`

	Parsemac struct {
		MacsFileName   string `arg:"" name:"macsfile" type:"existingfile"`
		Outfile        string `help:"Output file" type:"path"`
		ForceOverwrite bool   `name:"force" help:"Force overwrite output file file" short:"f"`
		Reportfile     string `help:"Reporting output file" type:"path"`
		UseMaxi        bool   `help:"Использовать максимальное количество что найдем" short:"m"`
		CisFileName    string `help:"Patch to cis.yaml" hidden:"" type:"existingfile" env:"CISFILE"`
		PwdFileName    string `help:"Patch to passw.json" hidden:"" type:"existingfile" env:"CISPWDS"`
	} `cmd:"" help:"Parsing file with mac-address table"`
}

var skipVlans []string

func main() {

	//	52, 170, 246, 248(нет запросов), 242(нет запросов),  620(нет запросов), 8, 16, 6, 7, 204, 19(voice domain), 172

	skipVlans = append(skipVlans, "6", "7", "8", "16", "19", "52", "170", "172", "204", "242", "246", "248", "620")

	ctx := kong.Parse(&cli,
		kong.Name("findusesvlans"),
		kong.Description("Get and Parse mac-address table for find used vlans"),
		kong.UsageOnError(),
	)

	switch ctx.Command() {
	case "getmacs <hosts>":
		GetMacsFromCisco()
	case "parsemac <macsfile>":
		ParseMacs(cli.Parsemac.MacsFileName)
	default:
		panic(ctx.Command())
	}
}
