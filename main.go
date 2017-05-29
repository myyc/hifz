package main

import (
   "github.com/labstack/echo"
   "net/http"
   "html/template"
   "io"
   "fmt"
   "flag"
   "os"
   "encoding/json"
   "math/rand"
   mw "github.com/labstack/echo/middleware"
   "github.com/labstack/gommon/log"
   "strconv"
)

var PORT = 2000
var CERT_PATH = "certs/cert.pem"
var KEY_PATH = "certs/key.pem"

var wl []WL

type WL struct {
   // required
   English      string `json:"en"`
   SArabic      string `json:"ar"`
   EArabic      string `json:"eg"`
   SArabicTrans string `json:"ar-tr"`
   EArabicTrans string `json:"eg-tr"`
}

type Template struct {
   templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
   return t.templates.ExecuteTemplate(w, name, data)
}

func getWl() []WL {
   wlFile, _ := os.Open("wl.json")

   jd := json.NewDecoder(wlFile)

   if err := jd.Decode(&wl); err != nil {
      log.Fatal(fmt.Sprintf("Error parsing the file (%s)", err))
      return nil
   }

   return wl
}

func getCard(c echo.Context, n int) error {
   jo, _ := json.Marshal(wl[n])

   return c.Render(http.StatusOK, "main", struct {
      Js template.JS
      N  int
   }{template.JS(jo), n})
}

func main() {
   flag.IntVar(&PORT, "port", PORT, "server port")
   flag.StringVar(&CERT_PATH, "tls-cert", CERT_PATH, "TLS certificate (.pem) location")
   flag.StringVar(&KEY_PATH, "tls-key", KEY_PATH, "TLS key (.pem) location")

   flag.Parse()

   e := echo.New()
   e.Use(mw.Logger())

   e.Static("/static", "st")
   e.HideBanner = true

   getWl()

   t := &Template{
      templates: template.Must(template.ParseGlob("tpl/html/*")),
   }

   e.Renderer = t

   e.GET("/", func(c echo.Context) error {
      n := rand.Intn(len(wl))

      return getCard(c, n)
   })

   e.GET("id/:num", func(c echo.Context) error {
      n, err := strconv.Atoi(c.Param("num"))

      if err != nil {
         return http.ErrBodyNotAllowed
      }

      return getCard(c, n)
   })

   log.Info("Listening on port ", PORT)

   _, err1 := os.Stat(CERT_PATH)
   _, err2 := os.Stat(KEY_PATH)

   if err1 != nil || err2 != nil {
      ge := func(e error) string {
         if e == nil {
            return "ø"
         } else {
            return e.Error()
         }
      }

      log.Warn(fmt.Sprintf("Some error with the certificates, not using SSL – (%s, %s)", ge(err1), ge(err2)))
      e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", PORT)))
   } else {
      log.Info("Using SSL")
      e.Logger.Fatal(e.StartTLS(fmt.Sprintf(":%d", PORT), CERT_PATH, KEY_PATH))
   }
}
