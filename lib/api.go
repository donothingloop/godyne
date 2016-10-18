/**
 * This file is part of the godyne project.
 * Copyright (c) 2016 Alexander Müller <donothingloop@gmail.com>
 */

package godyne

import (
    "errors"
    "github.com/golang/glog"
)

type Api struct {
    Updater *DNSUpdater
    Config *Config
}

type Request struct {
    User string
    Record string
    IP string
}

func (api *Api) handle(req *Request) error {
    glog.Infof("Handling request: %s -> %s for %s", req.Record, req.IP, req.User)

    record, err := api.Config.FindRecord(req.Record)

    if err != nil {
        glog.Error(err)
        return err
    }

    allowed := false

    for _, b := range record.Users {
        if b == req.User {
            allowed = true
            break
        }
    }

    if !allowed {
        glog.Warning("User is not allowed to update record.")
        return errors.New("User is not allowed to update record.")
    }

    return api.Updater.UpdateRecord(req)
}
