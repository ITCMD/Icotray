package main

import (
	_ "github.com/josephspurrier/goversioninfo"
	cmd "icotray/pkg/cmd/root"
)

//go:generate goversioninfo -icon=../../assets/image/icotray.ico -manifest=../../assets/data/icotray.exe.manifest

func main() {
	cmd.Execute()
}
