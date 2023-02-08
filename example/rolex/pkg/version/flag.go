package version

import (
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"
)

type versionValue bool

const (
	VersionFalse versionValue = false
	VersionTrue  versionValue = true
)

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

func (v *versionValue) String() string {
	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

func (v *versionValue) Set(s string) error {
	boolVar, err := strconv.ParseBool(s)
	if boolVar {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}

	return err
}

func (v *versionValue) Type() string {
	return "version"
}

const versionFlagName = "version"

var versionFLag = Version(versionFlagName, VersionFalse, "Print Version information and quit.")

func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)

	*p = value
	flag.Var(p, name, usage)
	flag.Lookup(name).NoOptDefVal = "true"

	return p
}

func PrintAndExitIfRequested() {
	if *versionFLag == VersionTrue {
		fmt.Printf("%s\n", Get())
		os.Exit(0)
	}
}
