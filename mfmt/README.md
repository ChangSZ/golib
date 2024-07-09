# 使用方法

```golang
package main

import "github.com/ChangSZ/golib/mfmt"

func main() {
	mfmt.Run("./")
}

```

使用效果类似于
```bash
gofmt -w ./
goimports -w -local github.com/ChangSZ/golib ./
```