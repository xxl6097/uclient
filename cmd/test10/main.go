package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type ValuesData struct {
	Values struct {
		Het struct {
			Anonymous        bool   `json:".anonymous"`
			Type             string `json:".type"`
			Name             string `json:".name"`
			Index            int    `json:".index"`
			Model            string `json:"model"`
			Manufacturer     string `json:"manufacturer"`
			HardwareVersion  string `json:"hardwareVersion"`
			Telnetd          string `json:"telnetd"`
			Mesh             string `json:"mesh"`
			Role             string `json:"role"`
			DeviceSn         string `json:"deviceSn"`
			Cmei             string `json:"cmei"`
			FirewareVersion  string `json:"firewareVersion"`
			BuildVersion     string `json:"buildVersion"`
			Opmode           string `json:"opmode"`
			OpmodeConfigFlag string `json:"opmode_config_flag"`
			FirstUpFlag      string `json:"first_up_flag"`
			FirstLogin       string `json:"first_login"`
		} `json:"het"`
		Andlink struct {
			Anonymous       bool   `json:".anonymous"`
			Type            string `json:".type"`
			Name            string `json:".name"`
			Index           int    `json:".index"`
			DeviceType      string `json:"deviceType"`
			DeviceId        string `json:"deviceId"`
			ProductToken    string `json:"productToken"`
			FirewareVersion string `json:"firewareVersion"`
			HardwareVersion string `json:"hardwareVersion"`
			VendorCode      string `json:"vendorCode"`
			DeviceMode      string `json:"deviceMode"`
			FlashSize       string `json:"flashSize"`
			RamSize         string `json:"ramSize"`
			CpuType         string `json:"cpuType"`
			Reserve         string `json:"reserve"`
			If5Enable       string `json:"if5enable"`
			DebugEnable     string `json:"debugEnable"`
			Disabled        string `json:"disabled"`
			MeshRole        string `json:"meshRole"`
			DeviceBrand     string `json:"deviceBrand"`
			DeviceModel     string `json:"deviceModel"`
			WebJkyxEnable   string `json:"web_jkyx_enable"`
			DeviceSn        string `json:"deviceSn"`
			CMEI            string `json:"CMEI"`
			DeviceKey       string `json:"deviceKey"`
			ProvinceCode    string `json:"provinceCode"`
			SoftwareVersion string `json:"softwareVersion"`
		} `json:"andlink"`
		Led struct {
			Anonymous bool   `json:".anonymous"`
			Type      string `json:".type"`
			Name      string `json:".name"`
			Index     int    `json:".index"`
			Status    string `json:"status"`
			Ontime    string `json:"ontime"`
			Offtime   string `json:"offtime"`
			Colour    string `json:"colour"`
		} `json:"led"`
		Wifi struct {
			Anonymous       bool   `json:".anonymous"`
			Type            string `json:".type"`
			Name            string `json:".name"`
			Index           int    `json:".index"`
			Bandone         string `json:"bandone"`
			Roamswitch      string `json:"roamswitch"`
			Lowrssi2G       string `json:"lowrssi2g"`
			Lowrssi5G       string `json:"lowrssi5g"`
			RSSIThreshold   string `json:"RSSIThreshold"`
			RSSIThreshold5G string `json:"RSSIThreshold5G"`
			Wpswitch        string `json:"wpswitch"`
			Limittime       string `json:"limittime"`
			Preferredsta5G  string `json:"preferredsta5g"`
			Packtloss       string `json:"packtloss"`
			Retryratio      string `json:"retryratio"`
			Lowrrsiperiod   string `json:"lowrrsiperiod"`
			Detectperiod    string `json:"detectperiod"`
			Dismisstime     string `json:"dismisstime"`
			Stopdetecttime  string `json:"stopdetecttime"`
			Igmpsnenable    string `json:"igmpsnenable"`
			Macaddress      string `json:"macaddress"`
		} `json:"wifi"`
		Web struct {
			Anonymous   bool     `json:".anonymous"`
			Type        string   `json:".type"`
			Name        string   `json:".name"`
			Index       int      `json:".index"`
			Username    string   `json:"username"`
			Loginuser   string   `json:"loginuser"`
			Timeout     string   `json:"timeout"`
			Read        []string `json:"read"`
			Write       []string `json:"write"`
			Password    string   `json:"password"`
			Md5Password string   `json:"md5_password"`
		} `json:"web"`
		Macfilter struct {
			Anonymous bool   `json:".anonymous"`
			Type      string `json:".type"`
			Name      string `json:".name"`
			Index     int    `json:".index"`
			Policy    string `json:"policy"`
			Enable    string `json:"enable"`
		} `json:"macfilter"`
		Ipv6 struct {
			Anonymous bool   `json:".anonymous"`
			Type      string `json:".type"`
			Name      string `json:".name"`
			Index     int    `json:".index"`
			Enable    string `json:"enable"`
		} `json:"ipv6"`
		Lock struct {
			Anonymous bool   `json:".anonymous"`
			Type      string `json:".type"`
			Name      string `json:".name"`
			Index     int    `json:".index"`
			Enable    string `json:"enable"`
		} `json:"lock"`
		Iptv struct {
			Anonymous       bool     `json:".anonymous"`
			Type            string   `json:".type"`
			Name            string   `json:".name"`
			Index           int      `json:".index"`
			Enable          string   `json:"enable"`
			Ethport         []string `json:"ethport"`
			Port2G4         []string `json:"port_2g4"`
			Port5G          []string `json:"port_5g"`
			Port5G2         []string `json:"port_5g_2"`
			TransparentFlag string   `json:"transparent_flag"`
		} `json:"iptv"`
		Pmf struct {
			Anonymous bool   `json:".anonymous"`
			Type      string `json:".type"`
			Name      string `json:".name"`
			Index     int    `json:".index"`
			Enable    string `json:"enable"`
		} `json:"pmf"`
		KernelCount struct {
			Anonymous bool   `json:".anonymous"`
			Type      string `json:".type"`
			Name      string `json:".name"`
			Index     int    `json:".index"`
			Enable    string `json:"enable"`
		} `json:"kernelCount"`
		Log struct {
			Anonymous bool   `json:".anonymous"`
			Type      string `json:".type"`
			Name      string `json:".name"`
			Index     int    `json:".index"`
			LogSwitch string `json:"logSwitch"`
			LogLevel  string `json:"logLevel"`
			LogTime   string `json:"logTime"`
		} `json:"log"`
		Multicast struct {
			Anonymous       bool   `json:".anonymous"`
			Type            string `json:".type"`
			Name            string `json:".name"`
			Index           int    `json:".index"`
			MulticastEnable string `json:"multicastEnable"`
		} `json:"multicast"`
	} `json:"values"`
}

func GetDeviceSn() string {
	out, _ := exec.Command("ubus", "call", "uci", "get", `{"config":"hetsystem"}`).Output()
	var res ValuesData
	//fmt.Println(string(out))
	json.Unmarshal(out, &res)
	return res.Values.Het.DeviceSn
}
func main() {
	fmt.Println(string(GetDeviceSn()))
}
