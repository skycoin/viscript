package dbus

//register that a resource exists
func (di *DbusInstance) ResourceRegister(resId ResourceId, resType ResourceType) {
	x := ResourceMeta{}
	x.Id = resId
	x.Type = resType

	di.Resources = append(di.Resources, x)
}

//remove resource from list
func (di *DbusInstance) ResourceUnregister(ResourceID ResourceId) {
	println("<dbus/registry>.ResourceUnregister()")
	println("FIXME/TODO: THIS IS NOT CALLED ANYWHERE")

	for i, resourceMeta := range di.Resources {
		if resourceMeta.Id == ResourceID {
			di.Resources = append(di.Resources[:i], di.Resources[i+1:]...)
		}
	}
}
