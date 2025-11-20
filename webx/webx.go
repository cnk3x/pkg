package webx

// type Server struct {
// 	server
// 	done        <-chan struct{}
// 	log         *slog.Logger
// 	name        string
// 	shutTimeout time.Duration
// 	err         error
// }

// // hide for webx.Server
// type server = http.Server

// func (s *Server) Init(options ...ServerOption) {
// 	s.name = "webx"
// 	s.log = logx.With("webx")
// 	s.server = server{Addr: "127.0.0.1:23432"}
// 	for _, apply := range options {
// 		apply(s)
// 	}
// }

// func (s *Server) Handle(handler http.Handler) { s.Handler = handler }

// func (s *Server) Start(ctx context.Context) error {
// 	done := make(chan struct{})
// 	s.done = done

// 	started := make(chan net.Addr, 1)

// 	s.server.BaseContext = func(l net.Listener) context.Context { return ctx }

// 	go func() {
// 		defer close(done)
// 		defer close(started)

// 		addr := s.Addr
// 		if addr == "" {
// 			addr = ":http"
// 		}

// 		s.log.Info("正在启动", "addr", s.Addr)
// 		ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", addr)
// 		if err != nil {
// 			s.err = err
// 			close(started)
// 		} else {
// 			s.err = nil
// 			started <- ln.Addr()
// 		}

// 		if err == nil {
// 			err = s.Serve(ln)
// 			s.err = err
// 		}

// 		if err != nil && err != http.ErrServerClosed {
// 			s.log.Warn("已停止", "err", err)
// 		} else {
// 			s.log.Info("已停止")
// 		}
// 	}()

// 	go func() {
// 		//当ctx done的时候关闭服务
// 		select {
// 		case <-done:
// 			return
// 		case <-ctx.Done():
// 			shutdown_ctx, cancel := context.WithTimeout(context.Background(), max(min(time.Second, s.shutTimeout), time.Second*3))
// 			_ = cancel
// 			errx.Ig(s.Shutdown(shutdown_ctx))
// 		}
// 	}()

// 	if addr, ok := <-started; ok {
// 		s.log.InfoContext(ctx, "服务已启动")
// 		s.printEndpoint(ctx, addr)
// 	}

// 	return nil
// }

// func (s *Server) Stop(ctx context.Context) error { return s.server.Shutdown(ctx) }

// func (s *Server) Done() <-chan struct{} { return s.done }

// func (s *Server) printEndpoint(ctx context.Context, addr net.Addr) {
// 	switch t := addr.(type) {
// 	case *net.TCPAddr:
// 		na2s := func(at net.Addr, _ int) (ip string, ok bool) {
// 			if nip, yes := at.(*net.IPNet); yes {
// 				if ok = len(nip.IP) == net.IPv4len || len(nip.IP) == net.IPv6len; ok {
// 					ip = nip.IP.String()
// 					if strings.Contains(ip, ":") {
// 						ip = `[` + ip + `]`
// 					}
// 					return
// 				}
// 			}
// 			return "", false
// 		}

// 		port := lo.Ternary(t.Port == 80, "", ":"+strconv.Itoa(t.Port))

// 		if t.IP.IsUnspecified() {
// 			allIp := lo.Flatten(lo.Map(
// 				errx.May(net.Interfaces()),
// 				func(ifi net.Interface, _ int) []string {
// 					return lo.FilterMap(errx.May(ifi.Addrs()), na2s)
// 				},
// 			))
// 			slices.Sort(allIp)
// 			lo.ForEach(allIp, func(ip string, i int) {
// 				s.log.InfoContext(ctx, fmt.Sprintf("访问端点: %2d: http://%s%s", i, ip, port))
// 			})
// 			return
// 		}

// 		if ip, ok := na2s(t, 0); ok {
// 			s.log.InfoContext(ctx, fmt.Sprintf("访问端点: http://%s%s", ip, port))
// 		}
// 	case *net.UnixAddr:
// 		//unix:/var/run/some.sock
// 		s.log.InfoContext(ctx, fmt.Sprintf("访问端点: %s:%s", t.Net, t.Name))
// 	}
// }
