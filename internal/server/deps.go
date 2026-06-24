package server

// calc-lib is the external arithmetic library this service exposes. The
// foundation sub adds it as a direct, resolvable dependency in go.mod; the
// endpoint subs that add /add, /subtract, /multiply, and /divide import its
// calc package and call Add/Subtract/Multiply/Divide. The blank import keeps
// the dependency wired (and `go mod tidy` stable) until those handlers land,
// without implementing any arithmetic here.
import _ "github.com/strategylippo/calc-lib/calc"
