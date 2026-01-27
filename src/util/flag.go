package util

import "flag"

type flagStruct struct {
	UseCwd *bool
}

var Flag = flagStruct{}

func LoadFlag() {
	Flag.UseCwd = flag.Bool("usecwd", false, "")
	flag.Parse()
}
