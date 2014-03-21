OUTDIR = build

build: clean
	mkdir -p $(OUTDIR)
	GOOS=linux GOARCH=amd64 go build -o $(OUTDIR)/twitter_listener main.go

clean:
	rm -rf $(OUTDIR)
