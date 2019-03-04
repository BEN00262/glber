package server

// // ErrInvalidBalancingAlgorithm is error of invalid balancing algorithm
// var ErrInvalidBalancingAlgorithm = errors.New("invalid balancing algorithm")
//
// // LB represents load balancer
// type LB struct {
// 	*http.Server
// 	// balancing *b.Balancing
// 	lf lockfree.LockFree
// }
//
// // NewLB returns LB object
// func NewLB(addr string) *LB {
// 	return &LB{
// 		Server: &http.Server{
// 			Addr: addr,
// 		},
// 		lf: lockfree.New(),
// 	}
// }
//
// // Build builds config of load balancer
// func (lb *LB) Build(conf config.Config) *LB {
// 	switch conf.Balancing {
// 	case "ip-hash":
// 		ih, err := iphash.New(conf.Servers.GetAddresses())
// 		if err != nil {
// 			glg.Fatalln(errors.Wrap(err, "faild to load ip-hash algorithm"))
// 		}
//
// 		// lb.balancing = b.New(ih)
// 		lb.Handler = http.HandlerFunc(lb.ipHashBalancing)
// 	case "round-robin":
// 		rr, err := roundrobin.New(conf.Servers.GetAddresses())
// 		if err != nil {
// 			glg.Fatalln(errors.Wrap(err, "faild to load round-robin algorithm"))
// 		}
//
// 		// lb.balancing = b.New(rr)
// 		lb.Handler = http.HandlerFunc(lb.roundRobinBalancing)
// 	case "least-connections":
// 		lc, err := leastconnections.New(conf.Servers.GetAddresses())
// 		if err == nil {
// 			glg.Fatalln(errors.Wrap(err, "faild to load least-connections algorithm"))
// 		}
//
// 		// lb.balancing = b.New(lc)
// 		lb.Handler = http.HandlerFunc(lb.ipHashBalancing)
// 	default:
// 		glg.Fatalln(errors.Wrap(ErrInvalidBalancingAlgorithm, conf.Balancing))
// 	}
//
// 	return lb
// }
//
// // ServeTLS runs load balancer with TLS
// func (lb *LB) ServeTLS(tlsConfig *tls.Config, certFile, keyFile string) error {
// 	lisner, err := net.Listen("tcp", lb.Addr)
// 	if err != nil {
// 		return errors.Wrap(err, "faild to listen")
// 	}
//
// 	lb.TLSConfig = tlsConfig
//
// 	glg.Success("Load Balancer starting on " + lb.Addr)
// 	err = lb.Server.ServeTLS(lisner, certFile, keyFile)
// 	if err != nil {
// 		return errors.Wrap(err, "faild to serve with tls")
// 	}
//
// 	return nil
// }
//
// // Serve runs load balancer
// func (lb *LB) Serve() error {
// 	lisner, err := net.Listen("tcp", lb.Addr)
// 	if err != nil {
// 		return errors.Wrap(err, "faild to listen")
// 	}
//
// 	glg.Success("Load Balancer starting on " + lb.Addr)
// 	err = lb.Server.Serve(lisner)
// 	if err != nil {
// 		return errors.Wrap(err, "faild to serve")
// 	}
//
// 	return nil
// }
