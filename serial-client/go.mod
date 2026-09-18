module github.com/WakkeWang/EmbedFlow/serial-client

go 1.27

require (
	github.com/WakkeWang/EmbedFlow/pkg/protocol v0.0.0
	github.com/coder/websocket v1.8.15
	go.bug.st/serial v1.6.4
)

require (
	github.com/creack/goselect v0.1.2 // indirect
	golang.org/x/sys v0.19.0 // indirect
)

replace github.com/WakkeWang/EmbedFlow/pkg/protocol => ../pkg/protocol
