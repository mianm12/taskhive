# 仓库指南

> 本文件是 Claude Code 与 Codex 共用的**通用规则**。Codex 直接读取本文件（含子目录的 `AGENTS.md`）；Claude Code 通过 `CLAUDE.md` 引入，并在其中维护自身专属规则。通用规则只写在这里，不在两处重复。

## 项目结构与模块组织

TaskHive 是一个 Go CLI 学习项目。CLI 入口位于 `cmd/taskhive/`。内部包放在 `internal/` 下，目前 `internal/version` 负责构建期版本信息和格式化。测试文件与被测代码放在同一包内，命名为 `*_test.go`。项目笔记和阶段复盘放在 `docs/notes/`；新增阶段笔记时使用 `docs/notes/stage-template.md`。CI 配置位于 `.github/workflows/`。

## 构建、测试与开发命令

- `make build`：使用版本 `ldflags` 编译 `./cmd/taskhive` 到 `bin/taskhive`。
- `make run`：构建并运行 `./bin/taskhive version`。
- `make test`：运行 `go test -race ./...`。
- `make lint`：运行 `golangci-lint run ./...`。
- `make fmt`：通过 `golangci-lint fmt ./...` 运行已配置的格式化器。
- `make clean`：删除 `bin/` 下的本地构建产物。

## 代码风格与命名约定

使用标准 Go 包布局：CLI 组装逻辑放在 `cmd/taskhive`，可复用逻辑放在 `internal/`。Go 代码使用仓库工具格式化；`.golangci.yml` 已配置 `gofumpt` 和 `goimports`。Go 文件缩进以格式化器输出为准。导出标识符需要有有效注释，因为 `revive` 会检查导出文档。优先编写短小、可测试的函数；除构建期配置外，避免新增全局状态。

## 测试指南

使用 Go 标准库 `testing` 包。测试文件放在被测包旁边，文件名使用 `*_test.go`。纯函数优先使用表驱动测试，并用 `t.Run` 写清楚用例名。提交行为变更前运行 `make test`；修改 Go 代码或 CI 配置时同时运行 `make lint`。

## Commit 与 Pull Request 指南

近期提交使用类似 Conventional Commit 的前缀，例如 `feat:`、`docs:`、`ci:` 和 `chore:`。提交标题保持简洁，并说明动作。前缀后的描述统一使用中文，例如 `docs: 添加 Claude Code 与 Codex 共用的 agent 指南`，与历史提交保持一致的格式。Pull Request 需要描述改动内容、列出已运行的验证命令，并说明对文档、CI 或配置的影响；如有关联 issue，应一并链接。

## 顶层路由

某些目录或任务有专属规则，登记在下表。**动手前先判断是否命中触发规则；命中则必须先读完“读取文档”列的文件再修改。** 新增规则时写进对应文档并在此登记一行，不在本文件展开。

| 用途     | 读取文档                          | 触发规则                                 |
| -------- | --------------------------------- | ---------------------------------------- |
| 文档写作 | `docs/AGENTS.md`                  | 新增、修改或润色 `docs/**` 下的 Markdown |

## 安全与配置提示

不要提交密钥、token、私钥或本地 `.env` 文件。生成的二进制文件和覆盖率输出不要纳入 git；这些内容已经由 `.gitignore` 覆盖。
