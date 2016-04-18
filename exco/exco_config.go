package exco

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/toolkits/file"

	"github.com/open-falcon/alarm/g"
)

const (
	exCoCfgFile = "./exco_config.json"
)

func ReadConfig() (*EcConfig, error) {
	cfg := exCoCfgFile
	if !file.IsExist(cfg) {
		return nil, fmt.Errorf("config file %s not found", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return nil, fmt.Errorf("read config file %s fail: %v", cfg, err)
	}

	var c EcConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return nil, fmt.Errorf("parse config file %s fail: %v", cfg, err)
	}

	if g.Config().Debug {
		log.Println("read exco config:", cfg, "successfully")
		log.Println(c.String())
	}

	return &c, nil
}

// raw config struct
type EcConfig struct {
	Ecs []*EcItem `json:"ecs"`
}

func (this *EcConfig) String() string {
	ret := "[ "
	for _, ec := range this.Ecs {
		ret += fmt.Sprintf("{main:%s [", ec.Main)
		for _, ex := range ec.Exs {
			ret += fmt.Sprintf("{id:%s counter:%s}", ex.Id, ex.Counter)
		}
		ret += "]}, "
	}
	return ret[0 : len(ret)-2]
}

type EcItem struct {
	Main string    `json:"main"`
	Exs  []*ExItem `json:"exs"`
}

func (this *EcItem) MainEx() (*ExItem, error) {
	for _, ex := range this.Exs {
		if ex.Id == this.Main {
			return ex, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

type ExItem struct {
	Id      string `json:"id"`
	Counter string `json:"counter"`
}

func (this *ExItem) Key() string {
	return fmt.Sprintf("%s_%s", this.Id, this.Counter)
}
