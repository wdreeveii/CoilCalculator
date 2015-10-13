// svgplay: sketch with SVGo, (derived from the old misc/goplay), except:
// (1) only listen on localhost, (default port 1999)
// (2) always render html,
// (3) SVGo default code,
package main

import (
	"flag"
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	httpListen = flag.String("port", "9004", "port to listen on")
	dynamic    = flag.Bool("dynamic", false, "Dynamically reload the front page")
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
	if *dynamic {
		page, err := ioutil.ReadFile("client/frontPage.html")
		if err != nil {
			log.Println(err)
		}
		w.Write(page)
	} else {

		w.Write(frontPage)
	}
}

func formString(req *http.Request, name string) (string, error) {
	var s string
	var err error
	if len(req.Form[name]) > 0 {
		s = req.Form[name][0]
	} else {
		err = fmt.Errorf("Form field %s not provided.", name)
	}
	return s, err
}

// Compile is an HTTP handler that reads Go source code from the request,
// runs the program (returning any errors),
// and sends the program's output as the HTTP response.
func GenSVG(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Println(err)
	}
	style, err := formString(req, "s")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if style == "rect" {
		diameter, err := formInt(req, "bd")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		genrect(w, diameter)
	} else {
		w.Header().Set("Content-Type", "image/svg+xml")
		gencir(w)
	}
}

var (
	canvas_width  = 500
	canvas_height = 500
	middle_width  = canvas_width / 2
)

func genrect(w io.Writer) {
	s := svg.New(w)
	s.Start(canvas_width, canvas_height)

	s.Rect(middle_width-125, canvas_height/4-75, 100, 100, "fill:none;stroke:black")
	s.Rect(middle_width-125, 2*canvas_height/3, 250, canvas_height/3, "fill:none;stroke:black")

	s.End()
}

func gencir(w io.Writer) {
	s := svg.New(w)
	s.Start(canvas_width, canvas_height)

	s.Circle(middle_width, canvas_height/3, 125, "fill:none;stroke:black")
	s.Rect(middle_width-125, 2*canvas_height/3, 250, canvas_height/3, "fill:none;stroke:black")

	s.End()
}
