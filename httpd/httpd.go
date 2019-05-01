package httpd

import (
	"goNES/bus"
	"io"
	"net/http"
	"os"
	"strings"
)

type HttpServer struct {
	httpCh chan string
}

func NewHttpd(message chan string) HttpServer {
	return HttpServer{httpCh: message}
}

func (server *HttpServer) handleHello(w http.ResponseWriter, r *http.Request) {
	command := strings.TrimLeft(r.URL.String(), "/")
	server.httpCh <- command
	io.WriteString(w, command)
}

func (server *HttpServer) StartHttpd() {
	http.HandleFunc("/", server.handleHello)
	http.ListenAndServe(":8080", nil)
}

func GetKeyinput(httpCh chan string, buttons [8]bool) [8]bool {
	select {
	case key := <-httpCh:
		switch key {
		case "finish":
			os.Exit(0)
		case "start":
			buttons[bus.ButtonStart] = true
			break
		case "select":
			buttons[bus.ButtonSelect] = true
			break
		case "up":
			buttons[bus.ButtonUp] = true
			break
		case "down":
			buttons[bus.ButtonDown] = true
			break
		case "left":
			buttons[bus.ButtonLeft] = true
			break
		case "right":
			buttons[bus.ButtonRight] = true
			break
		case "a":
			buttons[bus.ButtonA] = true
			break
		case "b":
			buttons[bus.ButtonB] = true
			break
		}
		break
	default:
		break
	}
	return buttons
}
