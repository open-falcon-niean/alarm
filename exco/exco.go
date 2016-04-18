package exco

import (
	"log"

	"github.com/open-falcon/alarm/g"
)

var (
	MainSubMap = NewMainSubMappingStruct()
	SubMap     = NewSubMappingStruct()
)

func Start() {
	// config
	rawcfg, err := ReadConfig()
	if err != nil {
		log.Fatalln("exco.Start, read config failed", err)
	}

	main_sub, sub, err := convert(rawcfg)
	if err != nil {
		log.Fatalln("exco.Start convert config failed", err)
	}
	MainSubMap.Set(main_sub)
	SubMap.Set(sub)
	if g.Config().Debug {
		log.Println("MainSubMap:", MainSubMap)
		log.Println("SubMap:", SubMap)
	}

	log.Println("exco.Start, ok")
}

func convert(ecc *EcConfig) (main_sub_mapping map[string]map[string]interface{},
	sub_mapping map[string]interface{}, err error) {
	// ret
	main_sub_mapping = map[string]map[string]interface{}{}
	sub_mapping = map[string]interface{}{}

	for _, ec := range ecc.Ecs {
		main, err := ec.MainEx()
		if err != nil {
			return nil, nil, err
		}

		key := main.Key()
		subsMap, exist := main_sub_mapping[key]
		if !exist {
			subsMap = map[string]interface{}{}
			main_sub_mapping[key] = subsMap
		}

		for _, ex := range ec.Exs {
			if ex.Id == main.Id {
				continue
			}
			subsMap[ex.Key()] = true
			// sub
			sub_mapping[ex.Key()] = true
		}
	}

	return main_sub_mapping, sub_mapping, nil
}
