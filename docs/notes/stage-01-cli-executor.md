# 阶段 1：CLI 与任务执行器基础

> 时间：2026-07-05 ~ 2026-07-06
> 对应 tag：v0.1.0

## 1. 本阶段目标

先有一个能执行任务的核心，哪怕是串行的。专注打牢 Go 语言基础（struct、方法、接口、错误处理和测试），刻意不碰并发——并发留到阶段 3 独立深潜。

## 2. 学到的知识点

### 语言核心

- **可见性在领域模型里的应用**：
  阶段 0 已记过基本规则，这里落到 `task` 包：`Task`、`Status`、`Transition()` 要给 `runner` 编排执行流时调用，所以导出；`runOnce()` 是执行器内部细节，所以保持小写。
- **自定义类型 `type Status string`**：
  创建一个底层是 `string` 但不等于 `string` 的新类型，带来类型安全：不能把裸字符串误传进要求 `Status` 的地方。它和 `struct`（多字段打包）是两种建模手段；互转必须显式写 `Status(x)` 或 `string(x)`。
- **const 枚举**：
  Go 没有 enum，用一组 const 常量模拟。命名惯例用类型名做前缀（`StatusPending`），
  IDE 补全友好。
- **值接收者 vs 指针接收者**：
  值接收者拿到的是拷贝，改了不影响原对象；指针接收者拿到地址，改了会改到本体。
  要改字段、struct 大、追求一致性时，用指针接收者。拿不准就用指针。
  实例：`Transition()` 要改 `Status` 必须用 `*Task`；`IsTerminal()` 只读且类型小，用值接收者。
- **接口是隐式实现的**：
  类型只要有接口要求的所有方法，就自动满足该接口，不需要 implements 声明（鸭子类型）。
  小接口哲学（io.Reader 只有一个方法）——接口越小，实现它的类型越多，组合力越强。
  这是“换实现零改动”（阶段 4 换数据库）的底层原理：上层依赖接口，不依赖具体实现。
- **switch 的 Go 特色**：
  一个 case 可列多个值（`case A, B, C`）；默认不贯穿，每个 case 执行完自动跳出，不用写 break。

### 错误处理（重点）

- **error 是值，不是异常**：
  用返回值处理，`if err != nil` 显式检查。Go 几乎不用 panic。
- **`error` vs `panic` 的分界**：
  外部输入、IO 和运行环境导致的失败一律返回 `error`；`panic` 只用于“不变量被违反”这类代码 bug。像 `runner.RunAll()` 里 `Running` 到终态迁移失败这种理论不可达路径，可用 `unreachable:` 前缀标注断言式 `panic`。
- **哨兵错误做成包级 Err 变量**：
  `var ErrInvalidTransition = errors.New(...)`，命名 Err 开头。
  这样调用方能用 errors.Is 识别“是不是这种错误”。每次现造 error 就没法识别了。
- **`fmt.Errorf` 的 `%w` 包装**：
  既加上下文（pending 到 running），又保留原始错误身份，
  让 `errors.Is` 能穿透包装识别。用 `%v` 则丢失身份。`%w` 和 `errors.Is` 是配对使用的标准模式。
- **错误信息字符串惯例**：
  小写开头、结尾无标点（`"invalid status transition"`）。
  因为错误常被 `%w` 层层包装拼接，大写或句号会让拼接结果混乱。`go vet` 的 ST1005 会检查。
- **`errors.As` vs `errors.Is`**：
  Is 回答是否为某错误；As 回答是否匹配，并把 err 转成具体类型给你用。
  取 `exec.ExitError` 的退出码必须用 As。`errors.As(err, &target)` 传 `&target` 是因为 As 要“写入”该变量
  （函数改谁就传谁地址）；若 target 本身是指针，`&target` 就是二级指针。同 `json.Unmarshal(data, &v)` 的模式。

### 标准库与执行

- **`os/exec` 的 Cmd 模型**：
  `exec.Command(程序名, 参数...)` 不经过 shell、不按空格拆分、不认管道或重定向。
  想用 shell 特性要显式 `sh -c "命令"`。
- **`CombinedOutput()` 的 err 三种情况**：
  `nil`（退出码 0）、`*exec.ExitError`（命令跑了但退出码非 0，可掏退出码）
  和其它错误（命令根本没启动）。区分“失败”和“没跑起来”对调度器很重要。
- **[]byte ↔ string**：
  `CombinedOutput` 返回 `[]byte`，`string(output)` 转换。同类型转换家族还有 `Status("x")`。
- **time 包**：
  `time.Now()` 取当前时刻，`time.Since(start)` 算耗时得到 `time.Duration`（打印自带单位）。

### context 与 defer（重点）

- **context 超时入门**：
  `context.WithTimeout(父, 时长)` 派生一个到点自动取消的 ctx，返回 ctx 和 cancel 函数。
  配合 `exec.CommandContext` 可自动杀超时进程。`Background()` 是根 ctx。`DeadlineExceeded` 是超时哨兵错误。
- **defer 执行时机**：
  defer 后的调用推迟到“当前函数即将返回时”执行，不管从哪个 return 走、是否 panic 都保证执行。
  典型用于成对的获取和释放（开文件 defer 关、加锁 defer 解、WithTimeout defer cancel）。
- **defer 循环陷阱（面试高频）**：
  defer 是函数级不是循环级。在 for 里 defer 会全部攒到函数返回才一起执行，
  可能导致资源（文件句柄等）堆积耗尽。解法：把循环体里“需要 defer 清理的单次操作”抽成独立函数。
  本项目把单次执行抽成 `runOnce()`，`defer cancel` 在其中，每次调用返回即释放，天然规避了陷阱。

### CLI、JSON 与数据流

- **`encoding/json` 与结构体标签**：
  结构体标签（struct tag）`json:"id"` 做 JSON 和字段名的映射；`Unmarshal(data, &v)` 传地址（同 `errors.As` 模式）；只有导出字段能被序列化。
  自定义类型实现 `UnmarshalJSON([]byte) error` 后，`json.Unmarshal` 会调用它；内部再把原始 `[]byte` 解成 `string`，可复用标准库的去引号、转义和类型校验。
- **`cobra` 命令树**：
  `cobra` 用命令树组织 CLI；`Run` 和 `RunE` 区分是否返回 error；flag 用 `StringVarP` 定义。这是项目第一个第三方依赖，`go.sum` 首次生成用于锁定校验和。
- **slice 与地址语义**：
  `make(T, 0, n)` 预分配容量，减少 `append` 扩容；`for i := range` 配合 `&slice[i]`，避开“range 循环变量是拷贝”的坑。

### 测试

- **表驱动测试的用例组织**：
  阶段 0 已记过基本形态，这一阶段重点是用一张表覆盖正向、反向和边界状态。匿名 struct 切片列用例，加上一个循环和 `t.Run(子测试名)` 跑完。加用例只加一行。
  正反用例都要有（既测“该成功的成功”，也测“该失败的失败”），再加边界用例。
- **白盒 vs 黑盒**：
  `package task`（同包，能测未导出成员）与 `package task_test`（独立包，只能测导出 API）。
  同目录允许 xxx 和 xxx_test 两个包并存，这是 Go 唯一的“一目录两包”例外。现阶段用白盒。
- **t.Fatalf vs t.Errorf**：
  Fatalf 立即终止当前测试（后续检查依赖前面结果时用）；Errorf 报告后继续（独立检查用）。
- **t.Skip**：
  主动跳过测试并显示 SKIP 和理由，测试仍“活着”提醒你待办，优于注释掉（那等于删除和遗忘）。
- **go build 无产物**：
  非 main 包 build 不产出文件，只当“编译检查器”。日常用 go test 即可（它会先编译）。
- **测试包与临时目录**：
  Go 测试只认包，不按文件名筛选；`t.TempDir()` 提供自动清理的临时目录；`filepath.Join` 用于跨平台拼路径。
- **`errcheck` 的错误态度**：
  `errcheck` 逼你对每个 error 表态：检查，或显式用 `_` 忽略。判断标准是“失败会不会让后续失去正确性”。

## 3. 关键决策（为什么这么做）

- **状态机用 map 表达规则，不用一堆 if**：
  `validTransitions map[Status][]Status` 把“规则即数据”，
  加状态或改规则只动数据不动逻辑。且利用 nil slice 可安全 range（零值可用）自然处理非法 from，不需特判。
- **校验放数据入口，核心函数信任已校验数据**：
  CanTransition 不特判非法 from（靠零值返回 false 即可）；
  真正防脏数据是构造或反序列化 Task 时的边界校验的职责。避免每个函数都重复防御。
- **执行结果用 Result struct 打包返回**：
  不返回一堆零散值，可扩展、语义清晰。
  ExitCode 记退出码、Err 只记“没跑起来”，两种失败语义分开存。
- **状态迁移放在 runner 编排层**：
  `executor` 只负责运行命令并返回 `Result`；`runner.RunAll()` 根据执行结果推进 `Task.Status`。
  这样任务生命周期编排和命令执行细节保持分离。
- **单次执行抽成 runOnce**：
  与重试循环 Run 分离，各司其职；也天然规避 defer 循环陷阱。
- **`sh -c` 执行命令的安全权衡**：
  见 [ADR-0001](../adr/0001-sh-c-execution.md)。

## 4. 踩过的坑

- **revive 报 exported const should have comment，但我明明写了行尾注释**。
  原因：Go 的“文档注释”专指对象**上方独占一行**、以对象名开头的注释；行尾注释（`//` 跟在代码后面）不算文档注释，`go doc` 不提取，revive 也不认。
  解决：给 const 块**上方**加一条独占行注释即可覆盖整块；或逐个常量在上方写。行尾注释可作为补充保留。

- **本地过、CI 挂：超时测试在 Linux 上失败**：`sh -c "sleep 5"` 设置 200ms 超时，本地 mac 通过，CI Linux 跑满 5 秒。
  根因：CommandContext 只杀直接子进程 sh，sh fork 出的 sleep 在 Linux 上不被连带杀掉、变孤儿进程睡完，超时失效；
  mac 因信号/管道行为不同蒙混过关。
  教训：碰真实进程的行为务必以 Linux（CI 或生产）为准，本地“能跑”可能是假象。
  已用 t.Skip 推迟，正确修复（进程组和并发控制）留待阶段 3，见 docs/notes/TODO.md。

## 5. 遗留 TODO

- [ ] 阶段 3 回来修 executor 超时杀进程的跨平台 bug（需进程组 `Setpgid`、`syscall.Kill(-pid)`，以及 `goroutine` 和 `select` 控制）
- [ ] 阶段 3 把 time.Sleep 退避改成指数退避
