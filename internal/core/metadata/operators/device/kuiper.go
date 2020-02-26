package device

import (
    "bytes"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/edgexfoundry/edgex-go/internal/core/metadata/config"
    "github.com/edgexfoundry/go-mod-core-contracts/clients"
    "io/ioutil"
    "net/http"
    "strconv"
)

type KuiperClient struct {
    Info config.KuiperInfo
}

func NewKuiperClient(info config.KuiperInfo) KuiperClient {
    return KuiperClient{
        Info: info,
    }
}

// TODO Clean up
func (k *KuiperClient) AddRule(name string) (string, error) {
    encoded := base64.StdEncoding.EncodeToString([]byte(name))
    r := RuleRequest{
        Name:        name,
        For:         "message.publish",
        RawSql:      "select * from \"new_device\" where payload.Payload = '" + encoded +"'",
        Description: "",
        Actions:     []RuleAction {
            {
                Name:   "republish",
                Params: RepublishAction{
                    TargetTopic: "blacklist_device",
                    TargetQos:   -1,
                    PayloadTmpl: "${payload}",
                },
            },
        },
    }
    b,_ := json.Marshal(r)
    url := k.Info.Protocol + "://" + k.Info.Host + ":" + strconv.Itoa(k.Info.Port) + "/api/v4/rules"
    kr, err := http.NewRequest("POST", url, bytes.NewReader(b))
    if err != nil {
        fmt.Errorf("error: %s", err.Error())
    }

    creds := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", k.Info.Username, k.Info.Password)))
    kr.Header.Set("Authorization", "Basic " + creds)
    kr.Header.Set(clients.ContentType, clients.ContentTypeJSON)

    res, err := http.DefaultClient.Do(kr)
    if err != nil {
        fmt.Errorf("err %s", err.Error())
    }

    rb, err := ioutil.ReadAll(res.Body)
    if err != nil {
        fmt.Errorf("Err: %s", err.Error())
    }
    var respBody AddRuleResponse
    json.Unmarshal(rb, &respBody)

    return respBody.Data.Id, nil
}

// TODO Clean up
func (k *KuiperClient) DeleteRule(n string) error {
    url := k.Info.Protocol + "://" + k.Info.Host + ":" + strconv.Itoa(k.Info.Port) + "/api/v4/rules/" + n

    kr, err := http.NewRequest("DELETE", url, nil)

    if err != nil {
        fmt.Errorf("error: %s", err.Error())
    }
    creds := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", k.Info.Username, k.Info.Password)))
    kr.Header.Set("Authorization", "Basic " + creds)
    kr.Header.Set(clients.ContentType, clients.ContentTypeJSON)

    _, err = http.DefaultClient.Do(kr)
    return  err
}

// TODO Lets return the proper response body
type AddRuleResponse struct {
    Data RuleResponse   `json:"data"`
    Code int            `json:"code"`
}


type RuleResponse struct {
    Id string `json:"id"`
    Rawsql string `json:"-"`
    Metrics map[string]interface{} `json:"-"`
    For []string    `json:"-"`
    Enabled string  `json:"-"`
    Description string  `json:"-"`
    Actions []interface{}   `json:"-"`
}

type RuleAction struct {
    Name    string  `json:"name"`
    Params  RepublishAction `json:"params"`
}

type RepublishAction struct {
    TargetTopic  string  `json:"target_topic"`
    TargetQos    int     `json:"target_qos"`
    PayloadTmpl  string  `json:"payload_tmpl"`
}

type Rule struct {
    Name string `json:"name"`
}

type RuleRequest struct {
    Name        string          `json:"name"`
    For         string          `json:"for"`
    RawSql      string          `json:"rawsql"`
    Description string          `json:"description"`
    Actions     []RuleAction    `json:"actions"`
}

func (r RuleRequest) MarshalJSON() ([]byte, error) {
    test := struct {
        Name        string  `json:"name"`
        For         string  `json:"for"`
        RawSql      string  `json:"rawsql"`
        Description string  `json:"description"`
        Actions     []RuleAction `json:"actions"`

    }{
        Name:           r.Name,
        For:            r.For,
        RawSql:         r.RawSql,
        Description:    r.Description,
        Actions:        r.Actions,
    }

    return json.Marshal(test)
}

func (r *RuleRequest) UnmarshalJSON(data []byte) error {
    var err error
    type Alias struct {
        Name        string          `json:"name"`
        For         string          `json:"for"`
        RawSql      string          `json:"rawSql"`
        Description string          `json:"description"`
        Actions     []RuleAction    `json:"actions"`
    }

    a := Alias{}

    if err = json.Unmarshal(data, &a); err != nil {
        return err
    }

    r.Name = a.Name
    r.For = a.For
    r.RawSql = a.RawSql
    r.Description = a.Description
    r.Actions = a.Actions

    return err
}


