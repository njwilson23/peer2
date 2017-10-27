all: peer peerbib

peer: peer.go search.go
	go build -o $@ $^

peerbib: bibtex/peerbib.go bibtex/bibtex.go
	go build -o $@ $^

install:
	cp peer $(GOPATH)/bin/
	cp peerbib $(GOPATH)/bin/

clean:
	rm peer peerbib