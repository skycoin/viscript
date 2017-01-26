package dbus

import ()

type ChannelId uint32

type ResourceType uint32

type ResourceId uint32

//resource types
const (
	ResourceTypeChannel  = 1
	ResourceTypeViewport = 2
	ResourceTypeTerminal = 3
	ResourceTypeProcess  = 4
)

//channel type
const (
	ChannelTypePubsub = 1
)

//ID generation

/*
	Id gen should eventually be per dbus instance
*/

var ChannelIdGlobal ChannelId = 2 //sequential

func RandChannelId() ChannelId {
	ChannelIdGlobal += 1
	return ChannelIdGlobal
	//return (ProccesId)(rand.Int63())
}
