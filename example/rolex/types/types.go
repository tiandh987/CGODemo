package types

type KV struct {
	Key   int32
	Value string
}

type Direction struct {
	// 方向的起始点A，两点相同表示双向
	PointA Point `json:"PointA"`
	// 方向的结束点B，两点相同表示双向
	PointB Point `json:"PointB"`
	// 箭头指向,0表示双向,1表示A->B,2表示B->A
	PointTo int32 `json:"PointTo"`
}

type Point struct {
	// 点的X轴坐标
	X int `json:"X" description:"点的X轴坐标"`
	// 点的Y轴坐标
	Y int `json:"Y" description:"点的Y轴坐标"`
}

type Line struct {
	// 线的起始点
	Point1 Point `json:"Point1" description:"线的起始点"`
	// 线的结束点
	Point2 Point `json:"Point2" description:"线的结束点"`
}

type Rect struct {
	// 矩形区域右下角Y轴坐标
	Bottom int `json:"Bottom" description:"矩形区域右下角Y轴坐标"`
	// 矩形区域左上角X轴坐标
	Left int `json:"Left" description:"矩形区域左上角X轴坐标"`
	// 矩形区域右下角X轴坐标
	Right int `json:"Right" description:"矩形区域右下角X轴坐标"`
	// 矩形区域左上角Y轴坐标
	Top int `json:"Top" description:"矩形区域左上角Y轴坐标"`
}

type Region struct {
	Points []Point `description:"多边形区域的顶点坐标"`
}

type SystemTime struct {
	Day    int `description:"日"`
	Hour   int `description:"时"`
	Minute int `description:"分"`
	Month  int `description:"月"`
	Second int `description:"秒"`
	Year   int `description:"年"`
}

type TimeSlice struct {
	// 时间片段所对应的以秒计时间，内部使用，便于做时间比对
	TimeSec []int ` json:"TimeSec" description:"时间片段所对应的以秒计时间，内部使用，便于做时间比对"`
	// 时间片段,数组0为开始时间，1为结束时间, 时间格式:2006-01-02 15:04:05
	TimeStr []string `json:"TimeStr" description:"时间片段, 时间格式:15:04:05"`
}

type TimeSection struct {
	// 每日时间段,定义每日最多6个时间片段
	Section []TimeSlice `json:"Section" description:"每日时间段,定义每日最多6个时间片段"`
}

type WeekSchedule struct {
	// 每周时间表
	WeekDay []TimeSection `json:"WeekDay" description:"每周时间表"`
}

type LinkChannel struct {
	// 使能开关
	Enable bool `json:"Enable" description:"使能开关"`
	// 联动通道
	Channel []int `json:"Channel" description:"联动通道"`
	// 联动延时
	Delay int `json:"Delay" description:"联动延时"`
}

type LinkLight struct {
	// 使能开关
	Enable bool `json:"Enable" description:"使能开关"`

	// 模式，0-闪烁；1-常亮
	Mode int `json:"Mode" description:"模式，0-闪烁；1-常亮"`

	// 闪烁频率，0-高；1-中；2低
	Frequency int `json:"Frequency" description:"闪烁频率，0-高；1-中；2低"`

	// 停留时间
	StayTime int `json:"StayTime" description:"停留时间"`
}

type LinkVoice struct {
	// 使能开关
	Enable bool `json:"Enable" description:"使能开关"`
	// 语音提示文件路径
	FilePath string `json:"FilePath" description:"语音提示文件路径"`
	// 播放次数
	Times int `json:"Times" description:"播放次数"`
}

type Linkage struct {
	// 联动 I/O 通道.
	AlarmOut LinkChannel `json:"AlarmOut" description:"联动I/O通道"`
	// 联动蜂鸣器.
	Buzzer bool `json:"Buzzer" description:"联动蜂鸣器"`
	// 联动上传 FTP.
	FTP bool `json:"FTP" description:"联动上传FTP"`
	// 联动白光灯闪光.
	Flash LinkLight `json:"Flash" description:"联动白光灯闪光"`
	// 联动云台操作.
	PTZ bool `json:"PTZ" description:"联动云台操作"`
	// 联动录像通道.
	Record LinkChannel `json:"Record" description:"联动录像通道"`
	// 联动发送邮件.
	SMTP bool `json:"SMTP" description:"联动发送邮件"`
	// 联动抓图通道.
	Snap LinkChannel `json:"Snap" description:"联动抓图通道"`
	// 联动语音提示.
	VoicePrompt LinkVoice `json:"VoicePrompt" description:"联动语音提示"`
}

type Color struct {
	// 颜色类型[R(0-255),G(0-255),B(0-255),A(0-100,Alpha透明度)]
	RGBA [4]int `json:"RGBA" description:"颜色类型[R(0-255),G(0-255),B(0-255),A(0-100,Alpha透明度)]"`
}

// 初始化报警联动结构体
func InitLinkage(avChn, alarmOutChn int, linkage *Linkage) bool {
	linkage.AlarmOut.Channel = make([]int, 0)
	linkage.AlarmOut.Delay = 10
	linkage.AlarmOut.Enable = false
	linkage.Record.Enable = false
	linkage.Record.Delay = 10
	linkage.Record.Channel = make([]int, 0)
	linkage.Snap.Enable = false
	linkage.Snap.Delay = 10
	linkage.Snap.Channel = make([]int, 0)
	linkage.Buzzer = false
	linkage.FTP = false
	linkage.PTZ = false
	linkage.SMTP = false
	linkage.Flash.Enable = false
	linkage.Flash.Mode = 0
	linkage.Flash.Frequency = 1
	linkage.Flash.StayTime = 5
	linkage.VoicePrompt.Enable = false
	linkage.VoicePrompt.FilePath = "alarm.wav"
	linkage.VoicePrompt.Times = 5
	return true
}

func InitWeekSchedule(schedule *WeekSchedule) bool {
	schedule.WeekDay = make([]TimeSection, 0)
	for weekday := 0; weekday < 7; weekday++ {
		var timesection TimeSection
		var timeslice TimeSlice
		timeslice.TimeSec = make([]int, 0)
		timeslice.TimeSec = append(timeslice.TimeSec, 0)
		timeslice.TimeSec = append(timeslice.TimeSec, 86399)
		timeslice.TimeStr = make([]string, 0)
		timeslice.TimeStr = append(timeslice.TimeStr, "00:00:00")
		timeslice.TimeStr = append(timeslice.TimeStr, "23:59:59")
		timesection.Section = make([]TimeSlice, 0)
		timesection.Section = append(timesection.Section, timeslice)
		schedule.WeekDay = append(schedule.WeekDay, timesection)

	}
	return true
}

////检验schedule的正确性
//func CheckWeekSchedule(schedule *WeekSchedule) (bool, string) {
//	for _, weekdayData := range schedule.WeekDay {
//		sectionArr := weekdayData.Section
//		if len(sectionArr) == 0 { //sectionArr 长度为0
//			return false, "sectionarr length is 0"
//		}
//
//		for _, tmpSeciton := range sectionArr {
//			timeStrArr := tmpSeciton.TimeStr
//			if len(timeStrArr) != 2 {
//				return false, "TimeStr length is not 2"
//			}
//
//			ret := validator.Ishhmmss(timeStrArr[0])
//			if !ret {
//				return false, "TimeStr format is invalid"
//			}
//
//			ret = validator.Ishhmmss(timeStrArr[1])
//			if !ret {
//				return false, "TimeStr format is invalid"
//			}
//			timeSecArr := tmpSeciton.TimeSec
//			if len(timeSecArr) != 2 {
//				return false, "TimeSec length is not 2"
//			}
//			timeSecStr0 := strconv.Itoa(timeSecArr[0])
//			timeSecStr1 := strconv.Itoa(timeSecArr[1])
//
//			ret = validator.IstimeNum5(timeSecStr0)
//			if !ret {
//				return false, "TimeSec format is invalid"
//			}
//			ret = validator.IstimeNum5(timeSecStr1)
//			if !ret {
//				return false, "TimeSec format is invalid"
//			}
//
//		}
//
//	}
//
//	return true, ""
//}

type Hwid struct {
	Magic          string `description:"“IRAY”"`
	Dev_type       string `description:"设备类型	“IPC”、“NVR”or “TPC” "`
	Dev_os         string `description:"设备操作系统 如Linux or rtthread"`
	Cpu_arch       string `description:"设备cpu类型 如“arm”“mips”"`
	Flash_type     string `description:"设备中的flash类型:“emmc”，“spi nor”"`
	Flash_size     string `description:"单位固定为MB，如“8MB”“4096MB”"`
	Dev_oem        string `description:"设备OEM类型,相同硬件针对不同厂商的定制需求"`
	Dev_language   string `description:"设备支持的语言种类"`
	Dev_plat_major string `description:"设备固件主类型"`
	Dev_plat_minor string `description:"设备固件次类型: dev_plat_major和dev_plat_minor需要使用这两个就能指定一个产品固件型号，结合上述其他参数可以指定唯一的产品固件型号"`
	Res            string `description:"预留字节，用于后面扩展使用，比如国家码等"`
}
