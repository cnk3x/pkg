package cmdx

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/cnk3x/pkg/cmdo"
	"github.com/cnk3x/pkg/errx"
	"github.com/cnk3x/pkg/jsonx"
	"github.com/cnk3x/pkg/logx"
	"github.com/cnk3x/pkg/x"
)

// 状态
const (
	statusStarting   = "starting"   //正在启动
	statusRunning    = "running"    //正在运行
	statusRestarting = "restarting" //正在重启
	statusStopped    = "stopped"    //已经停止
)

type Config struct {
	Path       string         `json:"path,omitempty"`
	Args       jsonx.Strings  `json:"args,omitempty"`
	Env        jsonx.Strings  `json:"env,omitempty"`
	InheritEnv bool           `json:"inherit_env,omitempty"`
	Dir        string         `json:"dir,omitempty"`
	Log        LogConfig      `json:"log"`
	Restart    RestartConfig  `json:"restart"`
	WaitDelay  jsonx.Duration `json:"wait_delay,omitempty"`
}

type Program struct {
	cfg       Config
	onPrepare func(cfg *Config) error
	log       *slog.Logger

	start   context.CancelFunc
	stop    context.CancelFunc
	restart context.CancelFunc

	status string
	done   <-chan struct{}
}

func Start(ctx context.Context, options ...Option) *Program {
	s := &Program{log: logx.With("运行")}
	done := make(chan struct{})
	s.done = done

	for _, option := range options {
		option.apply(s)
	}

	//方法：上报状态变更
	statusUp := func(status string) { s.status = status }
	//方法：状态判断
	statusIs := func(status string) bool { return s.status == status }

	//方法：主体执行
	directRun := func(stop_ctx context.Context) (err error) {
		if statusIs(statusRunning) {
			return errx.Errorf("cmdx: already running")
		}

		if s.restart != nil {
			s.restart()
		}

		restart_ctx, cancel := context.WithCancel(stop_ctx)
		defer cancel()
		s.restart = cancel

		if !statusIs(statusRestarting) {
			statusUp(statusStarting)
		}

		defer func() {
			if !statusIs(statusRestarting) {
				statusUp(statusStopped)
			}
		}()

		if s.onPrepare != nil {
			if err = s.onPrepare(&s.cfg); err != nil {
				return errx.Errorf("cmdx: %w", err)
			}
		}

		x := s.cfg
		c := exec.CommandContext(restart_ctx, x.Path, x.Args...)
		c.SysProcAttr = &syscall.SysProcAttr{}
		cmdo.PKill(c)
		c.Dir = x.Dir

		if x.InheritEnv {
			c.Env = append(c.Env, os.Environ()...)
		}
		c.Env = append(c.Env, x.Env...)

		c.WaitDelay = max(x.WaitDelay.Value(), time.Second*5) //调用cancel后等待退出，最低5s

		l0, l1, lc, le := x.Log.Open()
		if err = le; err != nil {
			return errx.Errorf("cmdx: %w", err)
		}
		c.Stdout, c.Stderr = l0, l1
		defer lc()

		if c.Dir != "" {
			if err = os.MkdirAll(c.Dir, 0777); err != nil {
				return errx.Errorf("cmdx: %w", err)
			}
		}

		s.log.Debug("启动", "cmdline", c.String())
		if err = c.Start(); err != nil {
			return errx.Errorf("cmdx: %w", err)
		}
		s.log.Debug("已启动", "pid", c.Process.Pid)

		statusUp(statusRunning)
		if err = c.Wait(); err != nil {
			return errx.Errorf("cmdx: %w", err)
		}
		return
	}

	//方法: 运行
	run := func(ctx context.Context) {
		if s.stop != nil {
			s.stop()
		}

		stop_ctx, cancel := context.WithCancel(ctx)
		s.stop = cancel

		for count := 1; ; count++ {
			select {
			case <-ctx.Done():
				return
			case <-stop_ctx.Done():
				return
			default:
			}

			err := directRun(stop_ctx)
			if err != nil {
				s.log.Debug("运行结果", "err", err.Error())
			}

			restart := s.cfg.Restart.CheckWait(ctx, stop_ctx, count, err)
			if !restart {
				return
			}

			//重启
			s.log.Debug("自动重启", "count", count)
		}
	}

	var initialized = make(chan struct{})
	setInitialized := sync.OnceFunc(func() { close(initialized) })

	//启动, 等待信号
	go func(ctx context.Context) {
		defer close(done)

		defer s.log.Debug("🔚 结束")
		s.log.Debug("初始化完成")

		for {
			startSignal := make(chan struct{})
			closeSignal := sync.OnceFunc(func() { close(startSignal) })
			s.start = closeSignal
			setInitialized()

			select {
			case <-ctx.Done():
				closeSignal()
				return
			case <-startSignal: //无限期等待, 直到... 调用了 s.start, 从而 start_ctx.Done!....
				s.log.Debug("启动")
				run(ctx)
			}
		}
	}(ctx)

	<-initialized
	x.Ig(s.Start())
	return s
}

// 启动
func (s *Program) Start() error { s.call(s.start, "启动"); return nil }

// 重启
func (s *Program) Restart() error { s.call(s.restart, "重启"); return nil }

// 停止
func (s *Program) Stop() error { s.call(s.stop, "停止"); return nil }

// 取得退出信号
func (s *Program) Done() <-chan struct{} { return s.done }

// 取得状态
func (s *Program) Status() string { return s.status }

func (s *Program) call(cancel context.CancelFunc, name string) {
	if cancel != nil {
		slog.Debug("请求命令: " + name)
		cancel()
	}
}
