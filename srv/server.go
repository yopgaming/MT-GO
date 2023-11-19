package srv

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/goccy/go-json"
	"log"
	"net"
	"net/http"
	"sync/atomic"

	"MT-GO/data"
	"MT-GO/pkg"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgradeToWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(conn)

	sessionID := r.URL.Path[28:] //mongoID is 24 chars
	data.SetConnection(sessionID, conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			data.DeleteConnection(sessionID)
			return
		}
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

//const incomingRoute string = "[%s] %s on %s\n"

func inflateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if websocket.IsWebSocketUpgrade(r) {
			upgradeToWebsocket(w, r)
		} else {
			if r.Header.Get("Content-Length") == "" {
				next.ServeHTTP(w, r)
				return
			}

			buffer := pkg.Inflate(r)
			if buffer == nil || buffer.Len() == 0 {
				next.ServeHTTP(w, r)
				return
			}

			//TODO: Refactor to replace r.Body ((remove CTX))
			var parsedData map[string]any
			if err := json.Unmarshal(buffer.Bytes(), &parsedData); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), pkg.ParsedBodyKey, parsedData)))
		}
	})
}

var CW = &ConnectionWatcher{}

func startHTTPSServer(serverReady chan<- bool, certs *Certificate, mux *muxt) {
	mux.initRoutes(mux.mux)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(inflateRequest)

	/*httpsServer := &http.Server{
		Addr: mux.address,
		//ConnState: CW.OnStateChange,
		TLSConfig: &tls.Config{
			RootCAs:      nil,
			Certificates: []tls.Certificate{certs.Certificate},
		},
		Handler: logAndDecompress(mux.mux),
	}*/

	log.Println("Started " + mux.serverName + " HTTPS server on " + mux.address)
	serverReady <- true

	err := http.ListenAndServeTLS(mux.address, certs.CertFile, certs.KeyFile, r)
	if err != nil {
		log.Fatalln(err)
	}
}

func startHTTPServer(serverReady chan<- bool, mux *muxt) {
	mux.mux.Use(middleware.Logger)
	mux.mux.Use(inflateRequest)

	mux.initRoutes(mux.mux)

	fmt.Println("Started " + mux.serverName + " HTTP server on " + mux.address)
	serverReady <- true

	err := http.ListenAndServe(mux.address, mux.mux)
	if err != nil {
		log.Fatalln(err)
	}
}

type muxt struct {
	mux        *chi.Mux
	address    string
	serverName string
	initRoutes func(mux *chi.Mux)
}

func SetServer() {
	srv := data.GetServerConfig()
	muxers := []*muxt{
		{
			mux: chi.NewMux(), address: data.GetMainIPandPort(),
			serverName: "Main", initRoutes: setMainRoutes,
		},
		{
			mux: chi.NewMux(), address: data.GetTradingIPandPort(),
			serverName: "Trading", initRoutes: setTradingRoutes,
		},
		{
			mux: chi.NewMux(), address: data.GetMessagingIPandPort(),
			serverName: "Messaging", initRoutes: setMessagingRoutes,
		},
		{
			mux: chi.NewMux(), address: data.GetRagFairIPandPort(),
			serverName: "RagFair", initRoutes: setRagfairRoutes,
		},
		{
			mux: chi.NewMux(), address: data.GetLobbyIPandPort(),
			serverName: "Lobby", initRoutes: setLobbyRoutes,
		},
	}

	serverReady := make(chan bool)

	if srv.Secure {
		cert := GetCertificate(srv.IP)
		certs, err := tls.LoadX509KeyPair(cert.CertFile, cert.KeyFile)
		if err != nil {
			log.Fatalln(err)
		}
		cert.Certificate = certs

		for _, muxData := range muxers {
			go startHTTPSServer(serverReady, cert, muxData)
		}
	} else {
		for _, muxData := range muxers {
			go startHTTPServer(serverReady, muxData)
		}
	}

	for range muxers {
		<-serverReady
	}
	close(serverReady)

	pkg.SetDownloadLocal(srv.DownloadImageFiles)
	pkg.SetChannelTemplate()
	pkg.SetGameConfig()

}

type ConnectionWatcher struct {
	n int64
}

func (cw *ConnectionWatcher) OnStateChange(_ net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew: //Connection open
		cw.Add(1)
	case http.StateHijacked, http.StateClosed: //Connection Closed
		cw.Add(-1)
	}
}

func (cw *ConnectionWatcher) Count() int {
	return int(atomic.LoadInt64(&cw.n))
}

func (cw *ConnectionWatcher) Add(c int64) {
	atomic.AddInt64(&cw.n, c)
}
