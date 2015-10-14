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
	"strconv"
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

func formFloat(req *http.Request, name string) (float64, error) {
	str, err := formString(req, name)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(str, 64)
	return f, err
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

	height, err := formFloat(req, "bh")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if style == "rect" {
		width, err := formFloat(req, "bw")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		length, err := formFloat(req, "bl")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		genrect(w, width, length, height)
	} else {
		diameter, err := formFloat(req, "bd")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		gencir(w, diameter, height)
	}
}

var (
	canvas_width  = 500
	canvas_height = 500

	middle_width = canvas_width / 2

	xy_center = canvas_height / 4
	xz_center = 3 * canvas_height / 4
)

func genrect(w io.Writer, width, length, height float64) {
	s := svg.New(w)
	s.Start(canvas_width, canvas_height)

	drawbackground(s)

	s.Rect(middle_width-int(width/2), xy_center-int(length/2), int(width), int(length), "fill:none;stroke:black")
	s.Rect(middle_width-int(width/2), xz_center-int(height/2), int(width), int(height), "fill:none;stroke:black")

	s.End()
}

func gencir(w io.Writer, diameter, height float64) {
	s := svg.New(w)
	s.Start(canvas_width, canvas_height)

	drawbackground(s)

	s.Circle(middle_width, xy_center, int(diameter/2), "fill:none;stroke:black")
	s.Rect(middle_width-int(diameter/2), xz_center-int(height/2), int(diameter), int(height), "fill:none;stroke:black")

	s.End()
}

func drawbackground(s *svg.SVG) {
	s.Def()
	s.Marker("uparrow", 5, 10, 20, 20)
	s.Path("M5 0 L0 10 L10 10 Z", "fill:black")
	s.MarkerEnd()

	s.Marker("rightarrow", 10, 5, 20, 20)
	s.Path("M0 0 L10 5 L0 10 Z", "fill:black")
	s.MarkerEnd()
	s.DefEnd()

	// middle divider
	s.Line(middle_width, 0, middle_width, canvas_height, "stroke:black")

	// xy_center line
	s.Line(0, xy_center, canvas_width, xy_center, "stroke:black")

	// xz_center line
	s.Line(0, xz_center, canvas_width, xz_center, "stroke:black")

	s.Line(10, 30, 10, canvas_height/2, `fill="none"`, `stroke="black"`, `marker-start="url(#uparrow)"`)
	s.Line(10, canvas_height/2, middle_width, canvas_height/2, `stroke="black"`, `marker-end="url(#rightarrow)"`)
}
