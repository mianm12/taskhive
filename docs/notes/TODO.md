# TODO

散落在代码里的 TODO 汇总总账。阶段末或对应阶段回来逐条清理。
代码里就地用 `// TODO(stage-N): ...` 标注，这里记详细版。

## 阶段 1 已清

- [x] 用 cobra 替换 `cmd/taskhive/main.go` 里手写的 `os.Args[1] == "version"` 判断（做 CLI 那一站）

## 阶段 3（并发）回来清

- [ ] **executor 超时杀进程的跨平台 bug**
      现象：`sh -c "sleep 5"` 设置 200ms 超时，本地 mac 通过，CI Linux 跑满 5 秒失败。
      根因：`CommandContext` 只杀直接子进程 `sh`，`sh` fork 出的 `sleep` 在 Linux 上不被连带杀掉，变孤儿进程睡完。
      修复方向：`cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}` 让进程自成一组，超时时 `syscall.Kill(-pid, SIGKILL)` 杀整组；配合 `goroutine` 和 `select` 控制超时。
      相关：`internal/executor/executor.go` 的 `runOnce`；`TestRun_Timeout`（已 Skip）。

- [ ] **重试退避改为指数退避**
      现状：阶段 1 用固定 `time.Sleep(RetryDelay)` 退避。
      目标：改成指数退避（每次等待翻倍，可加抖动）。阶段 3 并发改造 executor 时一起做。
