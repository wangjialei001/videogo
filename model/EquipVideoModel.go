package model

import "time"

type LocalTime time.Time
type EquipVideoModel struct {
	EquipId    uint64
	ExpireTime LocalTime
}

type EquipVideoConfigModel struct {
	Url      string
	Port     uint16
	UserName string
	Pwd      string
	Num      uint16
	Id       uint32
	EquipId  uint64
}

type VideoProcessModel struct {
	EquipId   uint64
	ProcessId int
}
