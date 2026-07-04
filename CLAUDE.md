# CLAUDE.md

本文件供 **Claude Code** 使用，只维护 Claude Code 专属规则；通用规则见下方引入的 `AGENTS.md`（Codex 直接读它）。

@AGENTS.md

## Claude Code 专属规则

- 检索代码用 Grep / Glob 或 Explore 子代理，不用 shell 的 `grep`、`find`；看文件用 Read，不用 `cat`。
- 较大或多步改动先进 plan mode 梳理方案，确认后再执行。
- 审查本地改动用 `/code-review`，涉及安全用 `/security-review`。
- 命中 `AGENTS.md`「顶层路由」触发规则时，用 Read 读取对应文档再改；Claude 专属路由（如某个 skill）在此按需补充。
