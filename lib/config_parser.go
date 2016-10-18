package godyne

import (
    "errors"
    "encoding/json"
    "io/ioutil"
    "github.com/golang/glog"
)

type ConfigParser struct {}

func NewConfigParser() *ConfigParser {
    return &ConfigParser{}
}

type DnsConfig struct {
    Nameserver string `json:"nameserver"`
    TTL uint16 `json:"ttl"`
}

type ZoneConfig struct {
    Key string `json:"key"`
}

type RecordConfig struct {
    Zone string `json:"zone"`
    Users []string `json:"users"`
    Aliases []string `json:"aliases"`
}

type UserConfig struct {
    Password string `json:"password"`
}

type ServerConfig struct {
    Port uint `json:"port"`
}

type Config struct {
    Server *ServerConfig `json:"server"`
    DnsConfig *DnsConfig `json:"dns"`
    Zones map[string]*ZoneConfig `json:"zones"`
    Records map[string]*RecordConfig `json:"records"`
    Users map[string]*UserConfig `json:"users"`
}

func (config *Config) FindZone(zone string) (*ZoneConfig, error) {
    val, ok := config.Zones[zone]

    if !ok {
        return nil, errors.New("dns_updater: Zone not found.")
    }

    return val, nil;
}

func (config *Config) FindRecord(domain string) (*RecordConfig, error) {
    val, ok := config.Records[domain]

    if !ok {
        return nil, errors.New("dns_updater: Record not found.")
    }

    return val, nil
}

func (config *Config) FindUser(user string) (*UserConfig, error) {
    val, ok := config.Users[user]

    if !ok {
        return nil, errors.New("dns_updater: User not found.")
    }

    return val, nil
}

func (conf *ConfigParser) ParseConfig(path string) (*Config, error) {
    glog.Infof("Parsing config file %s", path)

    file, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    err = json.Unmarshal(file, &cfg)

    if err != nil {
        return nil, err
    }

    return &cfg, nil
}
