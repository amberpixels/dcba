// Package dcba is the canonical DCBA layering spec: four sibling layers under
// internal/ and one dependency rule, A -> {C, D} -> B.
//
// There is no code here to call. The package exists so the spec is an ordinary,
// versioned Go dependency rather than a URL: a project that follows DCBA
// imports it for its side effect on the build graph, which pins the spec in
// go.mod, checksums it in go.sum, and keeps `go mod tidy` from dropping it.
//
//	import _ "github.com/amberpixels/dcba"
//
// That is what lets a standardgo project write `presets: [github.com/amberpixels/dcba]`
// and get the enforcement below without copying it.
package dcba

import _ "embed"

// Preset is the layering, expressed as depguard rules for standardgo to merge.
//
// ${MODULE} in it stands for the module path of the project being linted;
// depguard matches deny entries as import-path prefixes with no globbing, so
// the rules cannot name the packages they forbid until they are applied.
//
//go:embed standardgo-preset.yml
var Preset []byte
