/**
 * This file is part of the godyne project.
 * Copyright (c) 2016 Alexander Müller <donothingloop@gmail.com>
 */

package godyne

import (
    "io"
    "net/http"
    "github.com/golang/glog"
    "github.com/abbot/go-http-auth"
)

type ApiHTTP struct {
    Api
}

func (api ApiHTTP) update(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
    query := r.URL.Query()

    ip, ok := query["myip"]

    if !ok {
        io.WriteString(w, "INVALID_ARGS")
        return
    }

    hostname, ok := query["hostname"]

    if !ok {
        io.WriteString(w, "INVALID_ARGS")
        return
    }

    glog.Infof("Handling request for user %s: %s -> %s", r.Username, hostname[0], ip[0])

    req := Request{User: r.Username, Record: hostname[0], IP: ip[0]}
    err := api.handle(&req)

    if err != nil {
        io.WriteString(w, "ERROR")
    } else {
        io.WriteString(w, "OK")
    }
}

func (api *ApiHTTP) handleAuth(user, realm string) string {
    glog.Infof("HTTP Basic Auth: %s/%s", user, realm)

    uconf, err := api.Config.FindUser(user)

    if err != nil {
        glog.Error(err)
        return ""
    }

    return uconf.Password
}

func httpLog(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        glog.Infof("HTTP request: %s %s %s", r.RemoteAddr, r.Method, r.URL)
        handler.ServeHTTP(w, r)
    })
}

func (api *ApiHTTP) Start() {
    glog.Info("Listening on port 33009")

    authenticator := auth.NewBasicAuthenticator("godyne", api.handleAuth)
    http.HandleFunc("/api/update", authenticator.Wrap(api.update))
    http.ListenAndServe(":33009", httpLog(http.DefaultServeMux))
}

func NewHttpApi(dns *DNSUpdater, cfg *Config) (*ApiHTTP, error) {
    api := ApiHTTP{Api{dns, cfg}}

    return &api, nil
}
