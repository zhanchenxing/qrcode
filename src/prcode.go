package qrcode

import (
    "fmt"
    "net/http"
    "html/template"

    "appengine"
    "appengine/user"
)

func init() {
    http.HandleFunc("/base.html", base)
    http.HandleFunc("/", QR)
}

func handler(w http.ResponseWriter, r *http.Request){
    c := appengine.NewContext(r)
    u := user.Current(c)
    if u != nil {
        fmt.Fprintf(w, "Hello, %v!", u)
        return
    }

    url, err := user.LoginURL(c, r.URL.String())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Location", url)
    w.WriteHeader(http.StatusFound)
}

// var templ = template.Must(template.New("qr").Parse(templateStr))
var templ = template.Must(template.ParseFiles("templates/qr.html"))
func QR(w http.ResponseWriter, req *http.Request){
    templ.Execute(w, req.FormValue("s"))
}

var baseTempl = template.Must(template.ParseFiles("templates/base.html"))
func base(w http.ResponseWriter, req *http.Request){
    tc := make(map[string]string)
    tc["title"] = "Base Html"
    tc["content"] = "This is just a demo content."
    baseTempl.Execute(w, tc)
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/" name=f method="GET"><input maxLength=1024 size=70
name=s value="" title="Text to QR Encode"><input type=submit
value="Show QR" name=qr>
</form>
</body>
</html>
`
