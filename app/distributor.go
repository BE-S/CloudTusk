package app

import (
	"CloudTusk/lib/config"
	"CloudTusk/lib/log"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"sync"
	"time"
)

type Server struct {
	Url    *url.URL
	Proxy  *httputil.ReverseProxy
	Health bool
	Load   int
}

type Distributor struct {
	servers      map[string]*Server
	mutex        sync.Mutex
	ConfigParams *config.ConfigParams
}

func (d *Distributor) getServersFromConfig() map[string]bool {
	serversConfig := d.ConfigParams.Get("config", "hosts").SearchValue

	servers := make(map[string]bool)

	sc := reflect.ValueOf(serversConfig)

	if sc.Kind() == reflect.Slice {
		for _, value := range serversConfig.([]interface{}) {
			v := reflect.ValueOf(value)

			servers[v.String()] = true
		}
	} else {
		log.Fatal("Ошибка счетения серверов.")
	}

	return servers
}

func LoadServers() *Distributor {
	d := &Distributor{
		ConfigParams: &config.ConfigParams{},
	}

	d.servers = make(map[string]*Server)

	servers := d.getServersFromConfig()

	for urlConfig, _ := range servers {
		status := d.addServer(urlConfig)

		log.Info("Первичное подключение к серверу. Статус " + urlConfig + " " + status)
	}

	return d
}

func (d *Distributor) CheckLifeServers() {
	for {
		d.mutex.Lock()

		serversConfig := d.getServersFromConfig()

		/*
		 * Удаление серверов основываясь на конфиге
		 */
		for serverUrl, _ := range d.servers {
			_, exist := serversConfig[serverUrl]

			if !exist {
				delete(d.servers, serverUrl)

				log.Info("Изменение конфига: удалён сервер " + serverUrl)
			}
		}

		/*
		 * Анализ серверов
		 */
		for _, server := range d.servers {
			url := server.Url

			_, err := http.Head(url.String())

			if err != nil {
				log.Error("Сервер " + url.String() + " не активен")

				server.Health = false
			} else {
				server.Health = true
			}
		}

		/*
		 * Добавление новых серверов из конфига
		 */
		for urlConfig, _ := range serversConfig {
			_, exist := d.servers[urlConfig]

			if !exist {
				status := d.addServer(urlConfig)

				log.Info("Изменение конфига: добавлен сервер. Статус " + urlConfig + " " + status)
			}
		}

		d.mutex.Unlock()

		time.Sleep(time.Second * 10)
	}
}

func (d *Distributor) addServer(urlConfig string) string {
	status := "активен"
	urlAdded := false

	urlParsed, err := url.Parse(urlConfig)

	if err == nil {
		resp, err := http.Head(urlConfig)

		if err == nil {
			serverStatus := false

			proxy := httputil.NewSingleHostReverseProxy(urlParsed)

			if resp.StatusCode == 200 {
				serverStatus = true
			}

			d.servers[urlConfig] = &Server{
				Url:    urlParsed,
				Proxy:  proxy,
				Health: serverStatus,
				Load:   0,
			}

			urlAdded = true
		}
	}

	if err != nil {
		log.Error(err.Error())
	}

	if !urlAdded {
		status = fmt.Sprintf("не %s", status)
	}

	return status
}

/*
 * Получить сервер с наименьшей нагрузкой
 */
func (d *Distributor) getServer() *Server {
	var bestServer *Server

	d.mutex.Lock()

	for _, server := range d.servers {
		if server.Health && (bestServer == nil || bestServer.Load > server.Load) {
			bestServer = server
		}
	}

	if bestServer != nil {
		bestServer.Load += 1
	}

	d.mutex.Unlock()

	return bestServer
}
