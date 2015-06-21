# consulkv

CUI frontend for [Consul](https://consul.io/) Key/value store.

## Install

```
$ go get github.com/spiritloose/consulkv
```

## Synopsis

```
$ consulkv list
$ consulkv cat foo
$ consulkv delete foo
$ consulkv edit foo
$ consulkv flags foo
$ consulkv flags foo 42
$ consulkv put foo bar
$ consulkv put foo < /path/to/file
```

## Environment Variables

* `CONSUL_HTTP_ADDR`
* `CONSUL_HTTP_TOKEN`
* `CONSUL_HTTP_AUTH`
* `CONSUL_HTTP_SSL`
* `CONSUL_HTTP_SSL_VERIFY`

## Author

Jiro Nishiguchi <<jiro@cpan.org>>
