# go-torrent-cli

Go 语言实现 BT 下载器。博客地址：https://tiandiyijian.top/posts/go-bt/

## 使用
```
go run main.go path/to/torrent output/path
```

注意事项:
- 仅支持种子文件不支持磁力链接
- 仅支持 HTTP tracker
- 仅支持单文件种子
- 不支持上传功能