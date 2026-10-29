### serve

I use `python -m http.server` a lot, but it likes to complain about broken pipes. This seems to be working a bit better

###### installation

`$ go install go.senan.xyz/serve@latest`

###### usage

```
$ serve [-addr [host]:<port>] download [-attachment] [path]
$ serve [-addr [host]:<port>] upload [dir]
```
