OUTDIR = build

build: clean
	mkdir -p $(OUTDIR)
	GOOS=linux GOARCH=amd64 go build -o $(OUTDIR)/twitter_listener main.go

run:
	CGO_ENABLED=1 go run main.go xpeppers

test:
	go test ./listener

clean:
	rm -rf $(OUTDIR)
