<p align="center"><a href="#readme"><img src="https://gh.kaos.st/go-z7.svg"/></a></p>

<p align="center"><a href="#installation">Installation</a> • <a href="#compatibility-and-os-support">Compatibility and OS support</a> • <a href="#license">License</a></p>

<p align="center">
  <a href="https://godoc.org/pkg.re/essentialkaos/z7.v7"><img src="https://godoc.org/pkg.re/essentialkaos/z7.v7?status.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/z7"><img src="https://goreportcard.com/badge/github.com/essentialkaos/z7"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-z7"><img alt="codebeat badge" src="https://codebeat.co/badges/7d5b1210-a853-4d1d-a34a-4afcf574861e" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

`z7` package provides methods for working with 7z archives (`p7zip` wrapper).

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (_reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)_):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.8+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get pkg.re/essentialkaos/z7.v7
```

If you want to update `z7` to latest stable release, do:

```
go get -u pkg.re/essentialkaos/z7.v7
```

### Compatibility and OS support

|      Version |      1.x |    2.x  | 3.x-8.x |
|--------------|----------|---------|---------|
|  `p7zip 9.x` |    Full  | Partial | Partial |
| `p7zip 15.x` |  Partial |    Full |    Full |
| `p7zip 16.x` |  Partial | Partial |    Full |

| OS       | Support            |
|----------|--------------------|
| Linux    | :heavy_check_mark: |
| Mac OS X | :heavy_check_mark: |
| FreeBSD  | :heavy_check_mark: |
| Windows  | :x:                |

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
