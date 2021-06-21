# go-rectcp
A go tcp proxy application.

# How to use
```
go get github.com/mgr9525/go-rectcp

go-rectcp help
```

# Example
```
go-rectcp :80 localhost:8080
go-rectcp -d --timeout 20s :80 baidu.com:80
```