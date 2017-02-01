package dbus

/*
	ID generation (should eventually be per dbus instance)
*/
var ChannelIdGlobal ChannelId = 2 //sequential

func GetChannelId() ChannelId {
	print("(dbus/dbus.go).GetChannelId(): ")
	ChannelIdGlobal += 1
	println(ChannelIdGlobal)
	return ChannelIdGlobal
	//return (ProccesId)(rand.Int63())
}
