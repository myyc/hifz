حفژ
===

A stupid webapp with Arabic flashcards scraped from [this website](http://arabic.desert-sky.net/vocab.html).

```
$ go get github.com/labstack/echo
$ go run main.go
```

Supports SSL if you have certificates but it's not necessary.

```
$ go run main.go --help
Usage of /var/folders/vd/s1dsrtzx16n73l7y7gb1v9pr0000gn/T/go-build768780365/command-line-arguments/_obj/exe/main:
  -port int
        server port (default 2000)
  -tls-cert string
        TLS certificate (.pem) location (default "certs/cert.pem")
  -tls-key string
        TLS key (.pem) location (default "certs/key.pem")
exit status 2
```
