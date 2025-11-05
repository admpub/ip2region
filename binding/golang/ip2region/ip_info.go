package ip2region

import (
	"strings"
)

type IpInfo struct {
	Continent      string       `json:",omitempty"` // 所属的七大洲
	Country        string       `json:",omitempty"` // 国家
	Province       string       `json:",omitempty"` // 省份
	City           string       `json:",omitempty"` // 城市
	District       string       `json:",omitempty"` // 区县
	ISP            string       `json:",omitempty"` // 互联网供应商 如：移动、联通、电信
	Coordinates    *Coordinates `json:",omitempty"` // 坐标
	Currency       string       `json:",omitempty"` // 货币 如：CNY、USD
	TimeZone       string       `json:",omitempty"` // 时区 如：Asia/Shanghai
	Codes          *Codes       `json:",omitempty"`
	ASN            string       `json:",omitempty"` // ASN号 如：AS4134
	Scene          string       `json:",omitempty"` // 应用场景 如：MOB
	Elevation      string       `json:",omitempty"` // 海拔
	WeatherStation string       `json:",omitempty"` // 气象站 如：CHXX0120
	parser         func(*IpInfo, string)
}

func (a *IpInfo) SetParser(parser func(*IpInfo, string)) *IpInfo {
	a.parser = parser
	return a
}

func (a *IpInfo) Parse(result string) *IpInfo {
	if a.parser != nil {
		a.parser(a, result)
	} else {
		DefaultParser(a, result)
	}
	return a
}

type Coordinates struct {
	Longitude string // 经度
	Latitude  string // 纬度
}

type Codes struct {
	Area string // 行政区码
	City string // 电话和区号
	Zip  string // 邮编
}

func (ip IpInfo) String() string {
	return ip.Country + "|" + ip.Province + "|" + ip.City + "|" + ip.ISP
}

func ParseResult(result string) IpInfo {
	r := IpInfo{}
	r.Parse(result)
	return r
}

var DefaultParser = func(ipInfo *IpInfo, line string) {
	lineSlice := strings.Split(line, "|")
	length := len(lineSlice)
	if length < 5 {
		for i := 0; i <= 5-length; i++ {
			lineSlice = append(lineSlice, "")
		}
	} else if length >= 15 {
		/* document: https://ip2region.net/doc/data/ipv4_base
		[
		    0: "亚洲",
		    1: "中国",
		    2: "广东",
		    3: "深圳",
		    4: "宝安",
		    5: "电信",
		    6: "113.88311",
		    7: "22.55371",
		    8: "440306",
		    9: "0755",
		    10: "518100",
		    11: "Asia/Shanghai",
		    12: "CNY",
		    13: "11",
		    14: "CHXX0120"
		];
		*/
		ipInfo.Continent = lineSlice[0]
		ipInfo.Country = lineSlice[1]
		ipInfo.Province = lineSlice[2]
		ipInfo.City = lineSlice[3]
		ipInfo.District = lineSlice[4]
		ipInfo.ISP = lineSlice[5]
		ipInfo.Coordinates = &Coordinates{}
		ipInfo.Coordinates.Longitude = lineSlice[6]
		ipInfo.Coordinates.Latitude = lineSlice[7]
		ipInfo.Currency = lineSlice[12]
		ipInfo.TimeZone = lineSlice[11]
		ipInfo.Codes = &Codes{}
		ipInfo.Codes.Area = lineSlice[8]
		ipInfo.Codes.City = lineSlice[9]
		ipInfo.Codes.Zip = lineSlice[10]
		switch length {
		case 15:
			ipInfo.Elevation = lineSlice[13]
			ipInfo.WeatherStation = lineSlice[14]
		case 16:
			ipInfo.ASN = lineSlice[13]
			ipInfo.Elevation = lineSlice[14]
			ipInfo.WeatherStation = lineSlice[15]
		case 17:
			ipInfo.ASN = lineSlice[13]
			ipInfo.Scene = lineSlice[14]
			ipInfo.Elevation = lineSlice[15]
			ipInfo.WeatherStation = lineSlice[16]
		}
		return
	}

	ipInfo.Country = lineSlice[0]
	ipInfo.Province = lineSlice[1]
	ipInfo.City = lineSlice[2]
	ipInfo.ISP = lineSlice[3]
	return
}
