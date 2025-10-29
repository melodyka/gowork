module main.go

go 1.23.0

toolchain go1.23.4

replace work => ./patterns/work

require work v0.0.0-00010101000000-000000000000

require (
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.33.0 // indirect
)
