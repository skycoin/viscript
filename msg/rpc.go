package msg

import "github.com/corpusc/viscript/hypervisor/dbus"

type TermAndAttachedProcessID struct {
	TerminalId        TerminalId
	AttachedProcessId ProcessId
}

type ChannelInfo struct {
	ChannelId          dbus.ChannelId
	Owner              dbus.ResourceId
	OwnerType          dbus.ResourceType
	ResourceIdentifier string

	Subscribers []dbus.PubsubSubscriber
}
