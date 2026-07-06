# Changelog

## v0.1.0 - 2026-07-06

### Added

- `internal/task`：任务领域模型与状态机，覆盖合法迁移校验。
- `internal/executor`：基于 `os/exec` 的执行器，支持超时控制与失败重试。
- `internal/runner`：从 JSON 加载任务并串行执行，输出汇总报表。
- `cmd/taskhive`：基于 `cobra` 的 CLI，提供 `run` 和 `version` 子命令。

### Known Issues

- 超时杀进程在 Linux 上需杀整个进程组，待阶段 3 修复（见 `docs/notes/TODO.md`）。

## v0.0.1 - 2026-07-03

### Added

- 项目骨架与 Go module。
- `version` 包和 CLI `version` 子命令，支持 `ldflags` 注入。
- `Makefile`（`build`、`test`、`lint`、`fmt`、`run`）。
- `golangci-lint` 配置。
- GitHub Actions CI（`lint`、`test -race`、`build`）。
