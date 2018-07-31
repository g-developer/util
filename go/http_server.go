package util

import (
	"fmt"
	"net/http"
	"sync"
	"errors"
	"time"
)

const (
	MAX_HTTP_SERVER_NUM = 10
)

const (
	STOP      = -1
	NEW       = 0
	AVAILABLE = 1
	SERVING   = 2
)

type httpServer struct {
	port         int
	readTimeout  int
	writeTimeout int
	ins          *http.Server
	status       int
}

type httpServerManager struct {
	servers map[int]*httpServer
	size    int
	mutex   *sync.Mutex
}

type HttpHandler func(http.ResponseWriter, *http.Request)


func newHttpServer(port int) *httpServer {
	tmp := &httpServer{
		port,
		10,
		10,
		&http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			Handler: http.NewServeMux(),
			ReadTimeout:    600 * time.Second,
			WriteTimeout:   600 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		NEW,
	}
	return tmp
}

var httpServerMgrIns *httpServerManager
var once sync.Once

func GetHttpMgrInstance() *httpServerManager {
	once.Do(func() {
		httpServerMgrIns = &httpServerManager{map[int]*httpServer{}, 0, &sync.Mutex{}}
	})
	return httpServerMgrIns
}

func (self *httpServerManager) NewHttpServer(port int) *httpServer{
	if _, ok := self.servers[port]; !ok {
		tmp := newHttpServer(port)
		self.mutex.Lock()
		tmp.status = AVAILABLE
		self.servers[port] = tmp
		self.size += 1
		self.mutex.Unlock()
		return tmp
	} else {
		return nil
	}
}

func (self *httpServerManager) GetHttpServer(port int) (*httpServer) {
	if value, ok := self.servers[port]; ok {
		return value
	} else {
		logIns.Fatal("No Such HttpServer Linstend!")
		return nil
	}
}

func (self *httpServer) AddFileServer (pattern string, dir string) error {
	//self.ins.Handle(pattern, http.StripPrefix(pattern, http.FileServer(http.Dir(dir))))
	if nil != self {
		if mux, ok := self.ins.Handler.(*http.ServeMux); ok {
			fmt.Println("AddHandler---", pattern)
			mux.Handle(pattern, http.StripPrefix(pattern, http.FileServer(http.Dir(dir))))
			return nil
		} else {
			return errors.New("Handler Is Not Type *http.ServeMux")
		}
	} else {
		return errors.New("Self Is nil in AddHandler")
	}
}

func (self *httpServer) AddHandler(pattern string, handler HttpHandler) error {
	if nil != self {
		if mux, ok := self.ins.Handler.(*http.ServeMux); ok {
			fmt.Println("AddHandler---", pattern)
			mux.HandleFunc(pattern, handler)
			return nil
		} else {
			return errors.New("Handler Is Not Type *http.ServeMux")
		}
	} else {
		return errors.New("Self Is nil in AddHandler")
	}
}

func (self *httpServerManager) StopHttpServer(port int) error {
	return self.GetHttpServer(port).stop()
}

func (self *httpServer) start() error {
	if AVAILABLE != self.status {
		return errors.New(fmt.Sprintf("Port:%v Status Error! Status=%v", self.port, self.status))
	}
	err := self.ins.ListenAndServe()
	if nil == err {
		self.status = SERVING
	} else {
		fmt.Println("Server Port : ", self.port, "Start Failed! ", err)
	}
	return err
}

func (self *httpServer) restart () error {
	if STOP != self.status {
		return errors.New(fmt.Sprintf("Port:%v Status Error! Status=%v", self.port, self.status))
	}
	err := self.ins.ListenAndServe()
	if nil == err {
		self.status = SERVING
	}
	return err

}

func (self *httpServer) stop() error {
	if SERVING == self.status {
		err := self.ins.Close()
		if nil == err {
			self.status = STOP
		}
		return err
	} else {
		self.status = STOP
		return nil
	}
}

func (self *httpServerManager) DeleteHttpServer(port int) error {
	err := self.GetHttpServer(port).stop()
		if nil == err {
			self.mutex.Lock()
			delete(self.servers, port)
			self.size -= 1
			self.mutex.Unlock()
			return nil
		} else {
			return err
		}
}

func (self *httpServerManager) Start () {
	for _, s := range self.servers {
		if s.status != SERVING {
			err := s.start()
			if nil != err {
				logIns.Fatal(fmt.Sprintf("Port[%v] start Failed!", s.port))
			}
		}
	}
}

func (self *httpServerManager) Stop () {
	for _, s := range self.servers {
		if s.status == SERVING {
			err := s.stop()
			if nil != err {
				logIns.Fatal(fmt.Sprintf("Port[%v] stop Failed!", s.port))
			}
		}
	}
}
