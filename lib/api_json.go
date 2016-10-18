package godyne

import (
    "github.com/golang/glog"
    "encoding/json"
    "errors"
)

type UpdateRequest struct {
    Domain string `json:"domain"`
    IP string `json:"ip"`
}

func ParseJSON(msg string) (UpdateRequest, error) {
    var req UpdateRequest
    byt := []byte(msg)

    err := json.Unmarshal(byt, &req)
    if err != nil {
        return req, errors.New("Failed to parse message.")
    }

    return req, nil
}

func (req UpdateRequest) Handle() (error) {
    glog.Infof("Handle() called: Domain: %s, IP: %s", req.Domain, req.IP)
    return nil
}
