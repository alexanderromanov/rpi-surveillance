PI_USER = pi
PI_HOST = 192.168.178.22

.PHONY: deploy
deploy:
	GOOS=linux GOARCH=arm GOARM=7 go build \
        && scp rpi-surveillance $(PI_USER)@$(PI_HOST):/home/$(PI_USER)/surveillance

.PHONY: run
run:
	ssh -t $(PI_USER)@$(PI_HOST) "/home/$(PI_USER)/surveillance/rpi-surveillance"
