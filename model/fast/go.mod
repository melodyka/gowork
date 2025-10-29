module main.go

go 1.20

replace work => ./patterns/work

require (
	github.com/tatsushid/go-fastping v0.0.0-20160109021039-d7bb493dee3e
	work v0.0.0-00010101000000-000000000000
)

require (
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
)
