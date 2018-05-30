package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
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
		port:       Config.Port,
		staticpath: Config.StaticPath,
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
	if Config.UseIPv6 {
		nettype = "tcp6"
		addrString = "[" + Config.ListenAddress + "]:" + server.port
	} else {
		nettype = "tcp4"
		addrString = Config.ListenAddress + ":" + server.port
	}

	address, err := net.ResolveTCPAddr(nettype, addrString)
	if err != nil {
		log.Fatalf("Error resolving address %s (%s)", addrString, err.Error())
	}

	http.Handle("/", server.router)
	log.Notice("Starting HTTP Server on ", addrString)

	var srv *http.Server
	if Config.TLSHost != "" {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(Config.TLSHost),
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
	var resp mdalResponse
	var query Query

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
	query.Composition = inq.Composition
	for _, s := range inq.Selectors {
		query.Selectors = append(query.Selectors, Selector(s))
	}
	query.Variables = inq.Variables
	t0, err := time.Parse("2006-01-02 15:04:05", inq.Time.T0)
	if err != nil {
		handleErr(errors.Wrapf(err, "Could not parse T0 (%s)", inq.Time.T0))
		return
	}
	t1, err := time.Parse("2006-01-02 15:04:05", inq.Time.T1)
	if err != nil {
		handleErr(errors.Wrapf(err, "Could not parse T1 (%s)", inq.Time.T1))
		return
	}
	query.Time.T0 = t0
	query.Time.T1 = t1
	//query.Time.WindowSize = inq.Time.WindowSize
	query.Time.Aligned = inq.Time.Aligned

	log.Info("Serving query", query)

	ts, err := srv.c.HandleQuery(context.Background(), &query)
	if err != nil {
		handleErr(errors.Wrap(err, "Could not run query"))
		return
	}

	// serialize the result
	packed, err := ts.msg.MarshalPacked()
	if err != nil {
		handleErr(errors.Wrap(err, "Error marshalling results"))
		return
	}
	log.Debug(len(packed))
	log.Debugf("%+v", query)

	resp.Rows = query.uuids
	resp.Data = packed
	resp.Nonce = inq.Nonce

	enc := json.NewEncoder(rw)
	if err := enc.Encode(resp); err != nil {
		handleErr(errors.Wrap(err, "Error encoding results as json"))
		return
	}
	return
}
