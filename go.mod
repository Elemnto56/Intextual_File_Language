module intext

go 1.24.5

replace intext/eval-engine => ./eval-engine

require intext/eval-engine v0.0.0

require (
	github.com/expr-lang/expr v1.17.6
	github.com/spf13/cobra v1.10.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)
