module github.com/WakkeWang/EmbedFlow/serial-client

go 1.27

require (
	github.com/WakkeWang/EmbedFlow/pkg/protocol v0.0.0
	github.com/coder/websocket v1.8.15
)

replace github.com/WakkeWang/EmbedFlow/pkg/protocol => ../pkg/protocol
