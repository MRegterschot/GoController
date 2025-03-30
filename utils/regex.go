package utils

import "regexp"

var CommandRegex = regexp.MustCompile(`^/{1,2}`)
var FakePlayerRegex = regexp.MustCompile(`^\*fakeplayer\d+\*$`)
var MapFileRegex = regexp.MustCompile(`(?i)[.](Map|Challenge).*[.]Gbx$`)