/**
 * This file is part of the godyne project.
 * Copyright (c) 2016 Alexander Müller <donothingloop@gmail.com>
 */

package godyne

import (
    "github.com/golang/glog"
)

type Server struct {
    HTTP *ApiHTTP
    Updater *DNSUpdater
}

func (srv Server) Init() {
    glog.Info("Starting server...")

    srv.HTTP.Start()
}

func NewServer() (*Server, error) {
    parser := NewConfigParser()
    cfg, err := parser.ParseConfig("config.json")

    if err != nil {
        return nil, err
    }

    srv := new(Server)

    // initialize DNS Updater
    srv.Updater, err = NewDNSUpdater(cfg)

    if err != nil {
        return nil, err
    }

    // initialize HTTP API
    srv.HTTP, err = NewHttpApi(srv.Updater, cfg)

    if err != nil {
        return nil, err
    }

    return srv, nil
}
