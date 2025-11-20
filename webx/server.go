package webx

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/cnk3x/gopkg/errx"
	"github.com/cnk3x/gopkg/logx"
	"github.com/samber/lo"
)

func ServerStart(ctx context.Context, log *slog.Logger, options ...ServerOption) <-chan struct{} {
	if log == nil {
		log = logx.With("接口")
	}

	s := &http.Server{BaseContext: func(l net.Listener) context.Context { return ctx }}
	for _, apply := range options {
		apply(s)
	}

	done := make(chan struct{})
	started := make(chan net.Addr, 1)
	go func() {
		defer close(started)
		defer close(done)

		addr := s.Addr
		if addr == "" {
			addr = ":http"
		}

		log.Info("正在启动", "addr", s.Addr)
		ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", addr)
		if err == nil {
			started <- ln.Addr()
		}

		if err == nil {
			err = s.Serve(ln)
		}

		if err != nil && err != http.ErrServerClosed {
			log.Warn("已停止", "err", err)
		} else {
			log.Info("已停止")
		}
	}()

	//nolint: errcheck
	go func() {
		select {
		case <-done:
			return
		case <-ctx.Done():
			shutdown_ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			_ = cancel
			errx.Ig(s.Shutdown(shutdown_ctx))
		}
	}()

	var printEndpoint = func(addr net.Addr) {
		switch t := addr.(type) {
		case *net.TCPAddr:
			na2s := func(at net.Addr, _ int) (ip string, ok bool) {
				if nip, yes := at.(*net.IPNet); yes {
					if ok = len(nip.IP) == net.IPv4len || len(nip.IP) == net.IPv6len; ok {
						ip = nip.IP.String()
						if strings.Contains(ip, ":") {
							ip = `[` + ip + `]`
						}
						return
					}
				}
				return "", false
			}

			port := lo.Ternary(t.Port == 80, "", ":"+strconv.Itoa(t.Port))

			if t.IP.IsUnspecified() {
				allIp := lo.Flatten(lo.Map(
					errx.May(net.Interfaces()),
					func(ifi net.Interface, _ int) []string {
						return lo.FilterMap(errx.May(ifi.Addrs()), na2s)
					},
				))
				slices.Sort(allIp)
				lo.ForEach(allIp, func(ip string, i int) {
					log.InfoContext(ctx, fmt.Sprintf("访问端点: %2d: http://%s%s", i, ip, port))
				})
				return
			}

			if ip, ok := na2s(t, 0); ok {
				log.InfoContext(ctx, fmt.Sprintf("访问端点: http://%s%s", ip, port))
			}
		case *net.UnixAddr:
			//unix:/var/run/some.sock
			log.InfoContext(ctx, fmt.Sprintf("访问端点: %s:%s", t.Net, t.Name))
		}
	}

	select {
	case <-done:
	case listen, ok := <-started:
		if ok {
			log.InfoContext(ctx, "服务已启动")
			printEndpoint(listen)
		}
	}

	return done
}

type ServerOption func(s *http.Server)

func ServerPort(port uint16, internal ...bool) ServerOption {
	return func(s *http.Server) {
		var addr string
		if cmp.Or(internal...) {
			addr = "127.0.0.1:"
		} else {
			addr = "[::]:"
		}
		s.Addr = addr + strconv.Itoa(int(port))
	}
}

func ServerSocket(path string, perm int) ServerOption {
	return func(s *http.Server) {
		s.Addr = "unix:" + path
	}
}

func ServerHandle(handler http.Handler) ServerOption {
	return func(s *http.Server) { s.Handler = handler }
}
