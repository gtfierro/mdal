package version

import (
	"fmt"
)

var Commit = "unset"
var Release = "unset"
var LOGO = fmt.Sprintf("MDAL: Commit: %s Release: %s", Commit, Release)
