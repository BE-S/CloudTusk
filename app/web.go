package app

import (
	"CloudTusk/lib/config"
	"CloudTusk/lib/log"
	"html/template"
	"net/http"
)

type Web struct {
	port        string
	distributor *Distributor
}

func New(distributor *Distributor) *Web {
	var configParams config.ConfigParams

	port := configParams.Get("config", "web->port").String()

	if port == "" {
		log.Fatal("Не возможно установить порт из конфигурационного файла")
	}

	return &Web{
		port:        port,
		distributor: distributor,
	}
}

func (web *Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := web.distributor.getServer()

	if server != nil {
		server.IncrementLoad()

		server.Proxy.ServeHTTP(w, r)

		server.DecrementLoad()
	} else {
		tpl, _ := outputCustomHtml("template/failed.html")

		tpl.Execute(w, map[string]string{
			"Path": r.URL.Path,
		})
	}
}

func (web *Web) Start() {
	err := http.ListenAndServe(":"+web.port, web)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func outputCustomHtml(filePath string) (*template.Template, error) {
	tpl, err := template.ParseFiles(filePath)

	if err != nil {
		log.Fatal(err.Error())
	}

	return tpl, err
}
