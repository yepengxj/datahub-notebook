package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asiainfoLDP/datahub-client/cmd"
	"github.com/asiainfoLDP/datahub-client/daemon/daemonigo"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const g_strDpPath string = "/var/lib/datahub/"

type StoppableListener struct {
	*net.UnixListener          //Wrapped listener
	stop              chan int //Channel used only to indicate listener should shutdown
}

func New(l net.Listener) (*StoppableListener, error) {
	tcpL, ok := l.(*net.UnixListener)

	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}

	retval := &StoppableListener{}
	retval.UnixListener = tcpL
	retval.stop = make(chan int)

	return retval, nil
}

var StoppedError = errors.New("Listener stopped")
var sl = new(StoppableListener)

func (sl *StoppableListener) Accept() (net.Conn, error) {

	for {
		//Wait up to one second for a new connection
		sl.SetDeadline(time.Now().Add(time.Second))

		newConn, err := sl.UnixListener.Accept()

		//Check for the channel being closed
		select {
		case <-sl.stop:
			return nil, StoppedError
		default:
			//If the channel is still open, continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			//If this is a timeout, then continue to wait for
			//new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (sl *StoppableListener) Stop() {
	close(sl.stop)
}

func helloHttp(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Hello HTTP!\n")
}

func stopHttp(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Hello HTTP!\n")
	sl.Close()
	fmt.Println("connect close")
}

type formatDpCreate struct {
	Dpname  string `json:"dpname"`
	Dptype  string `json:"dptype"`
	Dpconn  string `json:"dpConn"`
	Dpquota string `json:"dpquota"`
}
type msgResp struct {
	Msg string `json:"msg"`
}

func dpHttp(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)

	if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		reqJson := formatDpCreate{}
		err := json.Unmarshal(result, &reqJson)
		if err != nil {
			fmt.Printf("%T\n%s\n%#v\n", err, err, err)
			fmt.Println(rw, "invalid argument.")
		}
		if len(reqJson.Dpname) == 0 {
			fmt.Fprintln(rw, "invalid argument.")
		} else {
			msg := &msgResp{}
			if err := os.Mkdir(reqJson.Dpname, 0755); err != nil {
				msg.Msg = err.Error()
			} else {
				msg.Msg = "OK."
			}
			resp, _ := json.Marshal(msg)
			respStr := string(resp)
			fmt.Fprintln(rw, respStr)
		}

	}

}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	panic("not reached")
}

func RunDaemon() {
	fmt.Println("run daemon..")
	// Daemonizing echo server application.
	switch isDaemon, err := daemonigo.Daemonize(); {
	case !isDaemon:
		return
	case err != nil:
		log.Fatalf("main(): could not start daemon, reason -> %s", err.Error())
	}
	fmt.Printf("server := http.Server{}\n")

	if false == isDirExists(g_strDpPath) {
		err := os.MkdirAll(g_strDpPath, 0755)
		if err != nil {
			fmt.Printf("mkdir %s error! %v ", g_strDpPath, err)
		}

	}
	os.Chdir(g_strDpPath)
	originalListener, err := net.Listen("unix", cmd.UnixSock)
	if err != nil {
		panic(err)
	}

	sl, err = New(originalListener)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", helloHttp)
	http.HandleFunc("/stop", stopHttp)
	http.HandleFunc("/datapool", dpHttp)
	http.HandleFunc("/Repository", repoHandler)
	//http.HandleFunc("/subscriptions", subHttp)

	server := http.Server{}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		server.Serve(sl)
	}()

	fmt.Printf("Serving HTTP\n")
	select {
	case signal := <-stop:
		fmt.Printf("Got signal:%v\n", signal)
	}
	fmt.Printf("Stopping listener\n")
	sl.Stop()
	fmt.Printf("Waiting on server\n")
	wg.Wait()
	daemonigo.UnlockPidFile()

}
