# consulkv [![Build Status](https://travis-ci.org/spiritloose/consulkv.svg)](https://travis-ci.org/spiritloose/consulkv)

CUI frontend for [Consul](https://consul.io/) Key/value store.

## Installation

Download latest release binary from [release page](https://github.com/spiritloose/consulkv/releases).

Extract the archive and put `consulkv` binary into the directory that is in your `PATH`.

Or install using `go get` if that's what you want.

```
$ go get github.com/spiritloose/consulkv
```

## Usage

```
$ consulkv list
$ consulkv cat foo
$ consulkv delete foo
$ consulkv edit foo
$ consulkv flags foo
$ consulkv flags foo 42
$ consulkv put foo bar
$ consulkv put foo < /path/to/file
$ consulkv dump > /path/to/dump.txt
$ consulkv load < /path/to/dump.txt
```

## Environment Variables

* `CONSUL_HTTP_ADDR`
* `CONSUL_HTTP_TOKEN`
* `CONSUL_HTTP_AUTH`
* `CONSUL_HTTP_SSL`
* `CONSUL_HTTP_SSL_VERIFY`

## Author

Jiro Nishiguchi <<jiro@cpan.org>>
