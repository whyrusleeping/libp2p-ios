framework:
	mkdir -p frameworks
	rm -rf frameworks/Libp2p.framework
	gomobile bind -v -target ios -o frameworks/Libp2p.framework github.com/whyrusleeping/libp2p-ios/go/libp2p
