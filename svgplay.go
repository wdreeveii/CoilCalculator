// svgplay: sketch with SVGo, (derived from the old misc/goplay), except:
// (1) only listen on localhost, (default port 1999)
// (2) always render html,
// (3) SVGo default code,
package main

import (
	"flag"
	//"fmt"
	"github.com/ajstarks/svgo"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	httpListen = flag.String("port", "9004", "port to listen on")
)

var (
	frontPage []byte
)

func init() {
	var err error

	frontPage, err = ioutil.ReadFile("client/frontPage.html")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/", FrontPage)
	http.HandleFunc("/view", GenSVG)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("client"))))

	log.Fatal(http.ListenAndServe(":"+*httpListen, nil))
}

// FrontPage is an HTTP handler that renders the svgoplay interface.
// If a filename is supplied in the path component of the URI,
// its contents will be put in the interface's text area.
// Otherwise, the default "hello, world" program is displayed.
func FrontPage(w http.ResponseWriter, req *http.Request) {
	w.Write(frontPage)
}

// Compile is an HTTP handler that reads Go source code from the request,
// runs the program (returning any errors),
// and sends the program's output as the HTTP response.
func GenSVG(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(500, 500)
	s.Circle(250, 250, 125, "fill:none;stroke:black")
	s.End()
}
