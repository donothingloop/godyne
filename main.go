/**
 * This file is part of the godyne project.
 * Copyright (c) 2016 Alexander Müller <donothingloop@gmail.com>
 */

package main

import (
    "os"
    "flag"
    "github.com/golang/glog"
    "./lib"
)

func usage() {
    flag.PrintDefaults()
    os.Exit(2)
}

func init() {
    flag.Parse()
    flag.Lookup("logtostderr").Value.Set("true")
}

func main() {
    glog.Info("Starting godyne")
    srv, err := godyne.NewServer()

    if err != nil {
        glog.Error("Failed to create server.")
        glog.Error(err)
        return
    }

    srv.Init()
}
