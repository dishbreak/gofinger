.PHONY = build

build: finger/finger fingerd/fingerd

finger/finger: finger/finger.go
	cd finger && go build

fingerd/fingerd: fingerd/fingerd.go
	cd fingerd && go build

clean:
	rm finger/finger
	rm fingerd/fingerd
