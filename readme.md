# z7 [![GoDoc](https://godoc.org/pkg.re/essentialkaos/z7.v6?status.svg)](https://godoc.org/pkg.re/essentialkaos/z7.v6) [![Go Report Card](https://goreportcard.com/badge/essentialkaos/z7)](https://goreportcard.com/report/essentialkaos/z7) [![codebeat badge](https://codebeat.co/badges/7d5b1210-a853-4d1d-a34a-4afcf574861e)](https://codebeat.co/projects/github-com-essentialkaos-z7) [![License](https://gh.kaos.io/ekol.svg)](https://essentialkaos.com/ekol)

`z7` package provides methods for working with 7z archives (`p7zip` wrapper).

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.7+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get pkg.re/essentialkaos/z7.v6
```

If you want to update `z7` to latest stable release, do:

```
go get -u pkg.re/essentialkaos/z7.v6
```

### Compatibility

|      Version |      1.x |    2.x  | 3.x-6.x |
|--------------|----------|---------|---------|
|  `p7zip 9.x` |    Full  | Partial | Partial |
| `p7zip 15.x` |  Partial |    Full |    Full |
| `p7zip 16.x` |  Partial | Partial |    Full |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.io/ekgh.svg"/></a></p>
