package exco

import (
	"sync"
)

// 主-辅表达式对应关系
type MainSubMappingStruct struct {
	sync.RWMutex
	Mapping map[string]map[string]interface{} `json:"main_subs"`
}

func NewMainSubMappingStruct() *MainSubMappingStruct {
	return &MainSubMappingStruct{Mapping: make(map[string]map[string]interface{}, 0)}
}

func (this *MainSubMappingStruct) Set(m map[string]map[string]interface{}) {
	this.Lock()
	defer this.Unlock()
	this.Mapping = m
}

func (this *MainSubMappingStruct) Get(main string) (map[string]interface{}, bool) {
	this.RLock()
	this.RUnlock()
	m, found := this.Mapping[main]
	return m, found
}

func (this *MainSubMappingStruct) Exist(main string) bool {
	this.RLock()
	this.RUnlock()
	_, found := this.Mapping[main]
	return found
}

// 辅表达式列表
type SubMappingStruct struct {
	sync.RWMutex
	Mapping map[string]interface{} `json:"subs"`
}

func NewSubMappingStruct() *SubMappingStruct {
	return &SubMappingStruct{Mapping: make(map[string]interface{}, 0)}
}

func (this *SubMappingStruct) Set(m map[string]interface{}) {
	this.Lock()
	defer this.Unlock()
	this.Mapping = m
}

func (this *SubMappingStruct) Get() map[string]interface{} {
	this.RLock()
	this.RUnlock()
	return this.Mapping
}

func (this *SubMappingStruct) Exist(key string) bool {
	this.RLock()
	this.RUnlock()
	_, found := this.Mapping[key]
	return found
}

// 表达式 状态存储器
type ExpressionCounterStatus struct {
	sync.RWMutex
	StatusMap map[string]string `json:"status"`
}

func NewExpressionCounterStatus() *ExpressionCounterStatus {
	return &ExpressionCounterStatus{StatusMap: make(map[string]string, 0)}
}

func (this *ExpressionCounterStatus) Get(key string) (string, bool) {
	this.RLock()
	defer this.RUnlock()
	s, found := this.StatusMap[key]
	return s, found
}

func (this *ExpressionCounterStatus) Set(key, status string) {
	this.Lock()
	defer this.Unlock()
	this.set(key, status)
}

func (this *ExpressionCounterStatus) set(key, status string) {
	this.StatusMap[key] = status
}

func (this *ExpressionCounterStatus) Remove(key string) {
	this.Lock()
	defer this.Unlock()
	delete(this.StatusMap, key)
}
