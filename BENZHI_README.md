本项目为 Go 库/组件（含配套管理页）：gRPC 流式代理与帧编解码，配套 proxyd。

# go-grpcproxy

gRPC 流式代理组件，提供 Frame Codec、Backpressure Window、Handler Chain 与双向 Stream Bridge。

## 构建

```bash
go build -o proxyd ./cmd/proxyd
```
