/**
 * This file is part of the godyne project.
 * Copyright (c) 2016 Alexander Müller <donothingloop@gmail.com>
 */

package godyne

import (
    "strings"
    "fmt"
    "time"
    "github.com/golang/glog"
    "github.com/miekg/dns"
)

type DNSUpdater struct {
    config *Config
    client *dns.Client
    idCounter uint16
}

func NewDNSUpdater(cfg *Config) (*DNSUpdater, error) {
    dnsupd := DNSUpdater{}
    dnsupd.config = cfg
    dnsupd.client = &dns.Client{}
    return &dnsupd, nil
}

func (updater *DNSUpdater) makeDeleteRR(rec string, req *Request) (dns.RR, error) {
    ustr := fmt.Sprintf("%s. 0 IN A 0.0.0.0", rec)
    glog.Infof("Delete: %s", ustr)

    return dns.NewRR(ustr)
}

func (updater *DNSUpdater) makeCreateRR(rec string, req *Request) (dns.RR, error) {
    ustr := fmt.Sprintf("%s. %d IN A %s", rec, updater.config.DnsConfig.TTL, req.IP)
    glog.Infof("Create: %s", ustr)

    return dns.NewRR(ustr)
}

func (updater *DNSUpdater) UpdateRecord(req *Request) error {
    glog.Infof("Updating %s on nameserver %s", req.Record, updater.config.DnsConfig.Nameserver)

    record, err := updater.config.FindRecord(req.Record)

    if err != nil {
        glog.Error(err)
        return err
    }

    glog.Infof("Found record: %s", req.Record)
    zone, err := updater.config.FindZone(record.Zone)

    if err != nil {
        glog.Error(err)
        return err
    }

    glog.Infof("Found zone: %s", record.Zone)

    var key []string = strings.Split(zone.Key, ":")
    glog.Infof("Using update key: %s", key[0])

    // preparing dns update request
    updater.client.TsigSecret = map[string]string{key[0]: key[1]}
    msg := &dns.Msg{}
    msg.SetUpdate(record.Zone + ".")
    msg.Id = updater.idCounter

    updater.idCounter++

    rr, err := updater.makeCreateRR(req.Record, req)

    if err != nil {
        glog.Error("Failed to make RR for main record.")
        return err
    }

    drr, err := updater.makeDeleteRR(req.Record, req)

    if err != nil {
        glog.Error("Failed to make delete RR for main record.")
        return err
    }

    crrs := []dns.RR{rr}
    ddrs := []dns.RR{drr}

    for _, rec := range record.Aliases {
        crec, err := updater.makeCreateRR(rec, req)
        drec, err2 := updater.makeDeleteRR(rec, req)

        if err != nil || err2 != nil {
            glog.Errorf("Failed to update alias: %s", rec)
            continue
        }

        crrs = append(crrs, crec)
        ddrs = append(ddrs, drec)
    }

    msg.RemoveName(ddrs)
    msg.Insert(crrs)

    msg.SetTsig(key[0], dns.HmacSHA512, 300, time.Now().Unix())

    // send the request
    _, _, err = updater.client.Exchange(msg, updater.config.DnsConfig.Nameserver + ":53")

    if err != nil {
        glog.Error(err)
        return err
    }

    return nil
}
