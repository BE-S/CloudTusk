package app

import (
	"CloudTusk/lib/config"
	"html/template"
	"log"
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
		server.Proxy.ServeHTTP(w, r)
	} else {
		tpl, _ := outputCustomHtml("template/failed.html")

		tpl.Execute(w, map[string]string{
			"Path": r.URL.Path,
		})
	}
}

func (web *Web) Start() {
	http.ListenAndServe(":"+web.port, web)

	log.Println("Сервер запущен на хосте " + web.port)
}

func outputCustomHtml(filePath string) (*template.Template, error) {
	tpl, err := template.ParseFiles(filePath)

	if err != nil {
		log.Fatal("Отсутствует указанный html-шаблон: " + filePath)
	}

	return tpl, err
}
