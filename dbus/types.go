package dbus

import ()

type ChannelId int32

type ResourceType int32

type ResourceId int32

//resource types
const (
	ResourceChannel  = 1
	ResourceViewport = 2
	ResourceTerminal = 3
	ResourceProcess  = 4
)

//channel type
const (
	ChannelPubsub = 1
)
