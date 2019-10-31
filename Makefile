GOOS=linux
GOARCH=amd64
GOARM=6


.PHONY: clean
clean:
ifneq ("$(wildcard go-player)","")
	@rm go-player
endif
ifneq ("$(wildcard go-player-arm)","")
	@rm go-player-arm
endif


.PHONY: build
build:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o go-player . || (echo "build failed $$?"; exit 1)
	@echo 'Build suceeded... done'


.PHONY: buildarm
buildarm: GOARCH=arm GOARM=6
buildarm:
	@GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -o go-player-arm . || (echo "build failed $$?"; exit 1)
	@echo 'Build suceeded... done'
