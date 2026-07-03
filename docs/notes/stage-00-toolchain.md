# 阶段 0：工具链与项目骨架

> 时间：2026-07-01 ~ 2026-07-03
> 对应 tag：v0.0.1

## 1. 本阶段目标

不写任何业务逻辑，只把“工程闭环”跑通：写代码 → 格式化 → lint → 测试 → 构建 → CI，全链路一键运转。之后每阶段只往骨架里填肉，不推翻重来。

## 2. 学到的知识点

- **GOPATH / GOBIN 在 modules 时代的定位**：
  GOPATH（默认 `~/go`）不再是项目必须放的地方，只剩两个用途：模块缓存（`~/go/pkg/mod`）和 `go install` 装的工具（`~/go/bin`）。GOBIN 留空即用默认。关键动作是把 `~/go/bin` 加进 PATH，否则 go install 装的工具调不动。
- **可见性由首字母大小写决定**：
  Go 没有 public/private 关键字。大写开头 = 导出（包外可见），小写 = 未导出（仅包内）。所以 version 包里 `String()` 大写能被 main 调，`format()` 小写只能包内用。
- **零值（zero value）**：
  变量声明后不赋值也有确定的零值（数字 `0` / 字符串 `""` / 指针·slice·map 为 `nil`），绝不会是未初始化的随机内存。很多类型“零值可用”（如 sync.Mutex 声明即能用）。
- **ldflags 注入版本号**：
  `-ldflags "-X 包路径.变量=值"` 在链接阶段替换包级 string 变量的值。同一份源码，构建命令不同就能带上精确的 git commit 和构建时间，不用硬编码。
- **表驱动测试**：
  用一个 struct 切片列出所有用例，一个循环 + `t.Run(子测试名, ...)` 跑完。加用例只加一行。Go 标准测试不用断言库，老实写 `if got != want { t.Errorf(...) }`。
- **Makefile 的坑**：
  命令行（target 下面）必须 Tab 开头；但变量定义的续行（`\` 换行）不能用 Tab，用空格或顶格。`cat -A Makefile` 看 `^I` 可自查。`.PHONY` 声明这些名字是命令不是文件。

## 3. 关键决策（为什么这么做）

- **golangci-lint 用 brew 装，不用 go install**：官方推荐二进制安装，版本更稳、更快。
- **配置用 v2 格式（`version: "2"`）**：brew 装的是新版 golangci-lint，老教程的 v1 格式会报错。
  linter 从 `default: standard` 起步——官方精选、既有用又不满屏误报，适合新手。
- **CI 跑在 ubuntu-latest，本地是 mac**：故意的。跨平台能帮忙抓“只在我机器上能跑”的问题。
- **测试从阶段 0 就带 `-race`**：现在没并发抓不到东西，但习惯要现在养成，阶段 3 它会真正救命。
- **tag 用带注释/签名版**：设了 `git config tag.gpgSign`，签名必须存在 annotated tag 里，
  所以 tag 自动升级为带注释模式。用 `git tag -a vX -m "..."` 可避免每次弹 vim。

## 4. 踩过的坑

- **git tag 打开了 vim**：因为配了 gpgSign，签名需要 annotated tag，于是弹编辑器写说明。
  解决：`i` 写说明 → `Esc` → `:wq`；或以后直接用 `-m` 给说明。
- **（在这里补充你自己实际遇到的报错）**

## 5. 遗留 TODO

- [ ] 阶段 1 用 cobra 替换掉 main 里手写的 `os.Args[1] == "version"` 判断
