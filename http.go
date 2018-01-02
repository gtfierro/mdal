package main

import (
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/acme/autocert"
)

type httpServer struct {
	c          *Core
	port       string
	staticpath string
	router     *httprouter.Router
}

func RunHTTP(c *Core) error {
	server := &httpServer{
		c:          c,
		port:       Config.HTTP.Port,
		staticpath: Config.HTTP.StaticPath,
	}

	log.Info("Static Path", server.staticpath)
	server.router = httprouter.New()
	server.router.POST("/query", server.handleQuery)

	// configure server
	var (
		addrString string
		nettype    string
	)

	// check if ipv6
	if Config.HTTP.UseIPv6 {
		nettype = "tcp6"
		addrString = "[" + Config.HTTP.ListenAddress + "]:" + server.port
	} else {
		nettype = "tcp4"
		addrString = Config.HTTP.ListenAddress + ":" + server.port
	}

	address, err := net.ResolveTCPAddr(nettype, addrString)
	if err != nil {
		log.Fatalf("Error resolving address %s (%s)", addrString, err.Error())
	}

	http.Handle("/", server.router)
	log.Notice("Starting HTTP Server on ", addrString)

	var srv *http.Server
	if Config.HTTP.TLSHost != "" {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(Config.HTTP.TLSHost),
			Cache:      autocert.DirCache("certs"),
		}
		srv = &http.Server{
			Addr:      address.String(),
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		}
		go func() {
			log.Warning(srv.ListenAndServeTLS("", ""))
		}()
	} else {
		srv = &http.Server{
			Addr: address.String(),
		}
		go func() {
			log.Warning(srv.ListenAndServe())
		}()
	}
	return nil
}

func (srv *httpServer) handleQuery(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var inq mdalQuery
	//var query Query

	handleErr := func(err error) {
		log.Error(err)
		http.Error(rw, err.Error(), 500)
	}

	defer req.Body.Close()
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&inq); err != nil {
		handleErr(err)
		return
	}

	log.Debugf("%+v", inq)

}
