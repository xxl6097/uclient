package openwrt

func (this *openWRT) webUpdateAll(obj []*DHCPLease) {
	if obj == nil {
		return
	}
	if this.fnEvent != nil {
		this.fnEvent(0, obj)
	}
}

func (this *openWRT) webUpdateOne(obj *DHCPLease) {
	if obj == nil {
		return
	}
	if this.fnEvent != nil {
		this.fnEvent(1, obj)
	}
}

func (this *openWRT) webNotify(obj *DHCPLease) {
	if obj == nil {
		return
	}
	if this.fnEvent != nil {
		this.fnEvent(2, obj)
	}
}
