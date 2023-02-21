package errors

type ErrCode int

func (this ErrCode) Error() string {
	return CodeMap[this].Message
}

func (this ErrCode) GetMessage() string {
	return CodeMap[this].Message
}

func (this ErrCode) GetTranslate() string {
	return CodeMap[this].Translate
}

var (
	Null = ErrCode(0) //无错误,内部用
	// HTTP标准状态码定义
	Success       = ErrCode(200) //成功
	ErrBadRequest = ErrCode(400) //Http请求参数错误
	ErrForbidden  = ErrCode(403) //无访问权限

	//ip地址冲突
	ErrIPConflict = ErrCode(400951) //IP地址冲突
	// 通用状态码定义
	SuccessNeedReboot = ErrCode(200000) //成功,立即重启后生效
	SuccessAutoReboot = ErrCode(200001) //成功,自动重启后生效 //yj add
	ErrOprFailed      = ErrCode(400000) //操作失败
	ErrSystem         = ErrCode(400001) //系统异常
	ErrUploadFailed   = ErrCode(400002) //上传文件失败
	ErrDownloadFailed = ErrCode(400003) //下载文件失败
	ErrRedirectHttps  = ErrCode(400004) //重定向到https

	// 设备授权错误码定义
	// 设备授权错误码定义
	ErrNoAuthority      = ErrCode(400200) //无授权文件   //yjadd
	ErrInvalidAuthority = ErrCode(400201) //授权文件无效 //yjadd
	ErrGetChipIDFailed  = ErrCode(400202) //获取芯片ID失败 //yjadd

	ErrInvalidFilePath  = ErrCode(400203) //文件路径无效
	ErrFileIncomplete   = ErrCode(400204) //文件不完整
	ErrFileInavailable  = ErrCode(400205) //文件正在使用
	ErrFileTooLarge     = ErrCode(400206) //文件太大   //yjaddd
	ErrFileCreateFailed = ErrCode(400207) //文件创建失败   //iray add
	ErrFileOpenFailed   = ErrCode(400208) //文件打开失败   //iray add
	ErrFileReadFailed   = ErrCode(400209) //文件读取失败   //iray add
	ErrFileSaveFailed   = ErrCode(400210) //文件保存失败   //iray add
	ErrFileMD5Sum       = ErrCode(400211) //文件MD5值计算失败  //iray add
	ErrFileMD5Match     = ErrCode(400212) //文件MD5值匹配失败  //iray add

	// 配置相关错误码定义
	ErrGetCfgFailed    = ErrCode(400300) //获取配置失败
	ErrSetCfgFailed    = ErrCode(400301) //设置配置失败
	ErrGetDftCfgFailed = ErrCode(400302) //获取默认配置失败
	ErrTestCfgFailed   = ErrCode(400303) //测试配置失败
	ErrGetDBFailed     = ErrCode(400304) //读取数据库失败   //iray add
	ErrDelRegFailed    = ErrCode(400305) //删除测温区域失败 //iray add
	// 能力集相关错误码定义
	ErrGetCapsFailed = ErrCode(400401) //获取能力集失败
	ErrNoCaps        = ErrCode(400402) //设备不具有该项能力
	ErrOverCaps      = ErrCode(400403) //设备超出最大能力限制
	// 数据库相关错误码定义
	ErrSearchFailed = ErrCode(400501) //未查询到相关数据
	ErrInsertFailed = ErrCode(400502) //数据添加失败  //yj add
	ErrUpdateFailed = ErrCode(400503) //数据更新失败
	ErrDeleteFailed = ErrCode(400504) //数据删除失败 //yj add

	// 智能分析相关错误码定义
	ErrUnrecognized    = ErrCode(400601) //图片无法识别  //yj add
	ErrRecogFailed     = ErrCode(400602) //图片识别失败  //yj add
	ErrNoObjectRecog   = ErrCode(400603) //未识别出目标  //yj add
	ErrNoSimilarObject = ErrCode(400604) //无相似目标    //yj add

	// 用户管理错误码定义
	ErrInvalidToken        = ErrCode(400701) //用户Token无效
	ErrPasswordError       = ErrCode(400702) //用户密码错误
	ErrOverLoginLimit      = ErrCode(400703) //用户登录数量超出限制
	ErrUserNoExist         = ErrCode(400704) //用户不存在
	ErrGroupNotExist       = ErrCode(400705) //用户组不存在
	ErrCanNotCreSysAccount = ErrCode(400706) //无法创建系统账户
	ErrUserAlreadyExist    = ErrCode(400708) //用户已存在
	ErrGroupAlreadyExist   = ErrCode(400709) //用户组已存在
	ErrNotAuthorized       = ErrCode(400710) //没有该项操作权限
	ErrUsersInGroup        = ErrCode(400711) //组内还存在用户
	ErrUserCreateWrong     = ErrCode(400713) //创建用户失败
	ErrModifyUserFailed    = ErrCode(400716) //修改用户失败
	ErrPasswordSame        = ErrCode(400717) //修改密码与原密码相同
	ErrLoadUserFailed      = ErrCode(400720) //加载用户信息失败
	ErrLoadGroupsFailed    = ErrCode(400721) //加载用户组信息失败  //yj add
	ErrSaveUsersFailed     = ErrCode(400722) //保存用户信息失败   //yj add
	ErrSaveGroupsFailed    = ErrCode(400723) //保存用户组信息失败 //yj add
	ErrUserBlackList       = ErrCode(400724) //当前账户限制为黑名单 //yj add
	ErrUserLock            = ErrCode(400725) //当前账户已被锁定    //yj add
	ErrAccountLocked       = ErrCode(400730) //密码错误次数超过限制，账号锁定  //iray add
	ErrPwdWrongOnce        = ErrCode(400731) //密码错误1次，5次后账号将被锁定   //iray add
	ErrPwdWrongTwice       = ErrCode(400732) //密码错误2次，5次后账号将被锁定    //iray add
	ErrPwdWrongThird       = ErrCode(400733) //密码错误3次，5次后账号将被锁定     //iray add
	ErrPwdWrongFourth      = ErrCode(400734) //密码错误4次，5次后账号将被锁定     //iray add
	ErrCanNotDelSysAccount = ErrCode(400750) //无法删除系统账户   //iray add

	// 系统升级错误码定义
	ErrUpgradeFailed                = ErrCode(400801) //升级过程出错
	ErrUpgradeIncomplete            = ErrCode(400802) //升级包不完整
	ErrUpgradeMismatch              = ErrCode(400803) //升级包不匹配
	ErrFirmwareUpgradeFilesMismatch = ErrCode(400804) //固件升级文件不匹配
	ErrUpgrading                    = ErrCode(400805) //设备正在升级 //yj add
	ErrFirmwareUpgradeFPGAWrong     = ErrCode(400806) //固件升级FPGA状态异常

	// WiFi相关错误码定义
	ErrBeyondMaxWifiList = ErrCode(400901) //获取可用WiFi列表超出规定最大限制
	ErrWiFiPwd           = ErrCode(400902) //WiFi密码错误
	ErrWiFiConnect       = ErrCode(400903) //连接WiFi出现未知错误
	ErrDisconnectWiFi    = ErrCode(400904) //断开WiFi连接错误
	ErrAssignIP          = ErrCode(400905) //分配IP错误
	ErrConnectUnknown    = ErrCode(400906) //连接WiFi返回未定义错误
	ErrOpenHotSpot       = ErrCode(400907) //开启个人热点失败
	ErrCloseHotSpot      = ErrCode(400908) //关闭个人热点失败
	ErrGetWiFiStatus     = ErrCode(400909) //获取Wi-FI连接状态失败
	ErrGetHotSpotStatus  = ErrCode(400910) //获取个人热点功能状态失败

	//语音对讲错误码定义
	ErrTalkNotOpen          = ErrCode(401001) //音频功能没开启
	ErrTalkInavailable      = ErrCode(401002) //对讲被占用
	ErrVoiceFileInavailable = ErrCode(401003) //音频文件正在播放
	ErrAudioInit            = ErrCode(401004) //音频设备初始化失败 //iray add

	// 温度错误码定义
	ErrNotOpenTemp            = ErrCode(402001) //没有开启测温功能  //iray add
	ErrNotOpenColdHotTracking = ErrCode(402002) //没有开启冷热点跟踪 //iray add

	// 外设错误码定义
	ErrMQTTOpen   = ErrCode(403001) //MQTT开启     //iray add
	ErrModbusOpen = ErrCode(403002) //modbus开启   //iray add
	//请求太多次了，系统正忙
	ErrSystemIsBusy = ErrCode(404000) //系统正忙 //iray add

	// VISCA协议与宽动态互斥
	ErrWideDynamicMutexLowLight  = ErrCode(401101) //低照度慢快门与超宽动态互斥  //yj add
	ErrWideDynamicMutexBackLight = ErrCode(401102) //背光补偿与超宽动态互斥   //yj add
	ErrWideDynamicMutexHLC       = ErrCode(401103) //强光抑制与超宽动态互斥  //yj add
	ErrWideDynamicMutexDefog     = ErrCode(401104) //透雾与超宽动态互斥     //yj add
	ErrBackLightMutexHLC         = ErrCode(401105) //背光补偿与强光抑制互斥  //yj add
	ErrDefogMutexHLC             = ErrCode(401106) //透雾和强光抑制互斥    //yj add
	//云台相关错误码定义
	ErrNotInPtzLimitFailed            = ErrCode(401201) //当前坐标不在云台限位范围内  //yj add
	ErrBeforeSetLeftMarginPtzLimit    = ErrCode(401202) //设置左边界前将先清除右边界  //yj add
	ErrBeforeSetRightMarginPtzLimit   = ErrCode(401203) //设置右边界前必须先设置左边界 //yj add
	ErrRightLessLeftMarginPtzLimit    = ErrCode(401204) //请设置右边界大于左边界      //yj add
	ErrRightGreaterLeftMarginPtzLimit = ErrCode(401205) //请设置右边界小于左边界      //yj add
	ErrBeforeSetUpMarginPtzLimit      = ErrCode(401206) //设置上边界前将先清除下边界   //yj add
	ErrBeforeSetDownMarginPtzLimit    = ErrCode(401207) //设置下边界前必须先设置上边界 //yj add
	ErrDownLessUpMarginPtzLimit       = ErrCode(401208) //请设置下边界大于上边界      //yj add
	ErrDownGreaterUpMarginPtzLimit    = ErrCode(401209) //请设置下边界小于上边界      //yj add
	ErrPresetAssociatedScenario       = ErrCode(401210) //该预置点已绑定场景,请清除场景后删除  //yj add

	//摄像头设置错误码定义
	ErrCameraReadFailed  = ErrCode(401301) //摄像头参数读取错误  //yj add
	ErrCameraWriteFailed = ErrCode(401302) //摄像头参数写入错误  //yj add
)

type CodeDesc struct {
	Message   string
	Translate string
}

type CodeKey struct {
	Code      ErrCode
	Translate string
}

var (
	CodeMap = make(map[ErrCode]CodeDesc)

	CodeMapCN = map[ErrCode]CodeDesc{
		Null:              {"msg: Null", "无错误"},
		Success:           {"msg: success", "操作成功"},
		ErrBadRequest:     {"msg: bad request", "参数错误"},
		ErrIPConflict:     {"msg: ip conflict", "ip 冲突"},
		ErrForbidden:      {"msg: forbidden", "无访问权限"},
		SuccessNeedReboot: {"msg: success, need to reboot device.", "操作成功,需要立即重启设备"},
		ErrOprFailed:      {"msg: operation failed", "操作失败"},
		ErrSystem:         {"msg: system exception", "系统异常"},
		ErrUploadFailed:   {"msg: upload file failed", "上传文件失败"},
		ErrDownloadFailed: {"msg: download file failed", "下载文件失败"},
		ErrRedirectHttps:  {"msg: need redirect https", "重定向到https"},

		ErrGetCapsFailed: {"msg: get capability failed", "获取能力集失败"},
		ErrNoCaps:        {"msg: no capability", "设备不具有该项能力"},
		ErrOverCaps:      {"msg: over capability", "设备超出最大能力限制"},

		ErrInvalidFilePath: {"msg: Invalid file path", "文件路径无效"},
		ErrFileIncomplete:  {"msg: file is incomplete", "文件不完整"},
		ErrFileInavailable: {"msg: file is being used", "文件正在使用"},
		ErrFileTooLarge:    {"msg: file is too large", "文件过大"},

		ErrFileCreateFailed: {"msg: file creation failed", "文件创建失败"},
		ErrFileOpenFailed:   {"msg: file open failed", "文件打开失败"},
		ErrFileReadFailed:   {"msg: file read failed", "文件读取失败"},
		ErrFileSaveFailed:   {"msg: file save failed", "文件保存失败"},
		ErrFileMD5Sum:       {"msg: file MD5 value calculation failed", "文件MD5值计算失败"},
		ErrFileMD5Match:     {"msg: file MD5 value matching failed", "文件MD5值匹配失败"},

		ErrGetCfgFailed:    {"msg: get config failed", "获取配置失败"},
		ErrSetCfgFailed:    {"msg: set config failed", "设置配置失败"},
		ErrGetDftCfgFailed: {"msg: get default config failed", "获取默认配置失败"},
		ErrDelRegFailed:    {"msg: Failed to delete the temperature measurement area, please delete the area in the area comparison first", "删除测温区域失败，请先在区域比较中删除该区域"},

		ErrSearchFailed: {"msg: data search failed", "未查询到相关数据"},
		ErrUpdateFailed: {"msg: data update failed", "数据更新失败"},

		ErrInvalidToken:   {"msg: the token is invalid", "Token失效"},
		ErrPasswordError:  {"msg: user password error", "用户密码错误"},
		ErrOverLoginLimit: {"msg: The number of user logins exceeded the limit", "用户登录数量超出限制"},

		ErrUpgradeFailed:     {"msg: error in upgradeing", "升级过程出错"},
		ErrUpgradeIncomplete: {"msg: upgrade package is not complete", "升级包不完整"},
		ErrUpgradeMismatch:   {"msg: upgrade package is not match", "升级包不匹配"},

		ErrTalkNotOpen:          {"msg: audio talk not open", "音频功能没开启"},
		ErrTalkInavailable:      {"msg: audio talk is not availiable", "对讲被占用"},
		ErrVoiceFileInavailable: {"msg: audio file is playing", "音频文件正在播放"},
		ErrUserNoExist:          {"msg: the user does not exist", "用户不存在"},
		ErrGroupNotExist:        {"msg: the user group does not exist", "用户组不存在"},
		ErrCanNotCreSysAccount:  {"msg: can not create admin account", "不能创建admin账户"},
		ErrUserAlreadyExist:     {"msg: user already exist", "用户已经存在"},
		ErrGroupAlreadyExist:    {"msg: group already exist", "用户组已经存在"},
		ErrNotAuthorized:        {"msg: not authorized", "没有授权"},
		ErrUsersInGroup:         {"msg: some users in group", "用户组中有用户存在"},

		ErrCanNotDelSysAccount:          {"msg: system user can not be delete", "无法删除系统账户"},
		ErrTestCfgFailed:                {"msg: test config failed", "测试配置失败"},
		ErrGetDBFailed:                  {"msg: get database failed", "读取数据库失败"},
		ErrUserCreateWrong:              {"msg: create new user failed", "创建用户失败"},
		ErrModifyUserFailed:             {"msg: modify user failed", "修改用户失败"},
		ErrLoadUserFailed:               {"msg: load user information failed", "加载用户信息失败"},
		ErrBeyondMaxWifiList:            {"msg: exceed the maximum capacity of wifi list", "获取可用WiFi列表超出规定最大限制"},
		ErrWiFiPwd:                      {"msg: connect wifi failed due to password incorrect", "WiFi密码错误"},
		ErrWiFiConnect:                  {"msg: connect wifi failed due to unknown reason", "连接WiFi出现未知错误"},
		ErrDisconnectWiFi:               {"msg: disconnect wifi failed", "断开WiFi连接错误"},
		ErrAssignIP:                     {"msg: assign ip failed", "分配IP错误"},
		ErrConnectUnknown:               {"msg: Connect WiFi return unknown error code", "连接WiFi返回未定义错误"},
		ErrOpenHotSpot:                  {"msg: Open HotSpot Failed", "开启个人热点失败"},
		ErrCloseHotSpot:                 {"msg: Close HotSpot Failed", "关闭个人热点失败"},
		ErrGetWiFiStatus:                {"msg: Get WiFi Status Failed", "获取Wi-FI连接状态失败"},
		ErrGetHotSpotStatus:             {"msg: Get HotSpot Status Failed", "获取个人热点功能状态失败"},
		ErrAudioInit:                    {"msg: Audio Init Failed", "音频设备初始化失败"},
		ErrPasswordSame:                 {"msg: The password is the same as the original password", "修改密码与原密码相同"},
		ErrAccountLocked:                {"msg: The account is locked", "密码错误次数超过限制，账号锁定"},
		ErrPwdWrongOnce:                 {"msg: The password is wrong once, the account will be locked after 5 times", "密码错误1次，5次后账号将被锁定"},
		ErrPwdWrongTwice:                {"msg: The password is wrong twice, the account will be locked after 5 times", "密码错误2次，5次后账号将被锁定"},
		ErrPwdWrongThird:                {"msg: The password is wrong three times, the account will be locked after 5 times", "密码错误3次，5次后账号将被锁定"},
		ErrPwdWrongFourth:               {"msg: The password is wrong four times, the account will be locked after 5 times", "密码错误4次，5次后账号将被锁定"},
		ErrFirmwareUpgradeFilesMismatch: {"msg: Firmware upgrade files do not match", "固件升级文件不匹配"},
		ErrFirmwareUpgradeFPGAWrong:     {"msg: Firmware Upgrade FPGA Status Abnormal", "固件升级FPGA状态异常"},
		ErrNotOpenTemp:                  {"msg: The temp function is not open", "测温功能未开启"},
		ErrNotOpenColdHotTracking:       {"msg: The cold hot tracking is not open", "冷热点跟踪未开启"},
		ErrMQTTOpen:                     {"msg: The MQTT is open", "MQTT已开启，不允许关闭测温功能"},
		ErrModbusOpen:                   {"msg: The Modbus is open", "Modbus已开启，不允许关闭测温功能"},
	}

	CodeMapEN = map[ErrCode]CodeDesc{
		Null:              {"msg: Null", "null"},
		Success:           {"msg: success", "success"},
		ErrIPConflict:     {"msg: ip conflict", "ip conflict"},
		ErrBadRequest:     {"msg: bad request", "bad request"},
		ErrForbidden:      {"msg: forbidden", "forbidden"},
		SuccessNeedReboot: {"msg: success, need to reboot device.", "success, need to reboot device"},
		ErrOprFailed:      {"msg: operation failed", "operation failed"},
		ErrSystem:         {"msg: system exception", "system exception"},
		ErrUploadFailed:   {"msg: upload file failed", "upload file failed"},
		ErrDownloadFailed: {"msg: download file failed", "download file failed"},
		ErrRedirectHttps:  {"msg: need redirect https", "need redirect https"},

		ErrGetCapsFailed: {"msg: get capability failed", "get capability failed"},
		ErrNoCaps:        {"msg: no capability", "no capability"},
		ErrOverCaps:      {"msg: over capability", "over capability"},

		ErrSearchFailed: {"msg: data search failed", "data search failed"},
		ErrUpdateFailed: {"msg: data update failed", "data update failed"},

		ErrInvalidFilePath:  {"msg: invalid file path", "invalid file path"},
		ErrFileIncomplete:   {"msg: file is incomplete", "file is incomplete"},
		ErrFileInavailable:  {"msg: file is being used", "file is being used"},
		ErrFileTooLarge:     {"msg:file is too large", "file is too large"},
		ErrFileCreateFailed: {"msg: file creation failed", "file creation failed"},
		ErrFileOpenFailed:   {"msg: file open failed", "file open failed"},
		ErrFileReadFailed:   {"msg: file read failed", "file read failed"},
		ErrFileSaveFailed:   {"msg: file save failed", "file save failed"},
		ErrFileMD5Sum:       {"msg: file MD5 value calculation failed", "file MD5 value calculation failed"},
		ErrFileMD5Match:     {"msg: file MD5 value matching failed", "file MD5 value matching failed"},

		ErrGetCfgFailed:    {"msg: get config failed", "get config failed"},
		ErrSetCfgFailed:    {"msg: set config failed", "set config failed"},
		ErrGetDftCfgFailed: {"msg: get default config failed", "get default config failed"},
		ErrDelRegFailed:    {"msg: Failed to delete the temperature measurement area, please delete the area in the area comparison first", "Failed to delete the temperature measurement area, please delete the area in the area comparison first"},

		ErrInvalidToken:        {"msg: the token is invalid", "the token is invalid"},
		ErrPasswordError:       {"msg: user password error", "user password error"},
		ErrOverLoginLimit:      {"msg: the number of user logins exceeded the limit", "the number of user logins exceeded the limit"},
		ErrUserNoExist:         {"msg: the user does not exist", "the user does not exist"},
		ErrGroupNotExist:       {"msg: the user group does not exist", "the user group does not exist"},
		ErrCanNotCreSysAccount: {"msg: can not create admin account", "can not create admin account"},
		ErrUserAlreadyExist:    {"msg: user already exist", "user already exist"},
		ErrGroupAlreadyExist:   {"msg: group already exist", "group already exist"},
		ErrNotAuthorized:       {"msg: not authorized", "not authorized"},
		ErrUsersInGroup:        {"msg: some users in group", "some users in group"},

		ErrUpgradeFailed:     {"msg: error in upgradeing", "error in upgradeing"},
		ErrUpgradeIncomplete: {"msg: upgrade package is not complete", "upgrade package is not complete"},
		ErrUpgradeMismatch:   {"msg: upgrade package is not match", "upgrade package is not match"},

		ErrTalkNotOpen:          {"msg: audio talk not open", "audio talk not open"},
		ErrTalkInavailable:      {"msg: audio talk is not availiable", "audio talk is not availiable"},
		ErrVoiceFileInavailable: {"msg: audio file is playing", "audio file is playing"},

		ErrCanNotDelSysAccount:          {"msg: system user can not be delete", "system user can not be delete"},
		ErrTestCfgFailed:                {"msg: test config failed", "test config failed"},
		ErrGetDBFailed:                  {"msg: get database failed", "get database failed"},
		ErrUserCreateWrong:              {"msg: create new user failed", "create new user failed"},
		ErrModifyUserFailed:             {"msg: modify user failed", "modify user failed"},
		ErrLoadUserFailed:               {"msg: load user information failed", "load user information failed"},
		ErrBeyondMaxWifiList:            {"msg: exceed the maximum capacity of wifi list", "exceed the maximum capacity of wifi list"},
		ErrWiFiPwd:                      {"msg: connect wifi failed due to password incorrect", "connect wifi failed due to password incorrect"},
		ErrWiFiConnect:                  {"msg: connect wifi failed due to unknown reason", "connect wifi failed due to unknown reason"},
		ErrDisconnectWiFi:               {"msg: disconnect wifi failed", "disconnect wifi failed"},
		ErrAssignIP:                     {"msg: assign ip failed", "assign ip failed"},
		ErrConnectUnknown:               {"msg: Connect WiFi return unknown error code", "Connect WiFi return unknown error code"},
		ErrOpenHotSpot:                  {"msg: Open HotSpot Failed", "Open HotSpot Failed"},
		ErrCloseHotSpot:                 {"msg: Close HotSpot Failed", "Close HotSpot Failed"},
		ErrGetWiFiStatus:                {"msg: Get WiFi Status Failed", "Get WiFi Status Failed"},
		ErrGetHotSpotStatus:             {"msg: Get HotSpot Status Failed", "Get HotSpot Status Failed"},
		ErrAudioInit:                    {"msg: Audio Init Failed", "Audio Init Failed"},
		ErrPasswordSame:                 {"msg: The password is the same as the original password", "The password is the same as the original password"},
		ErrAccountLocked:                {"msg: The account is locked", "The account is locked"},
		ErrPwdWrongOnce:                 {"msg: The password is wrong once, the account will be locked after 5 times", "The password is wrong once, the account will be locked after 5 times"},
		ErrPwdWrongTwice:                {"msg: The password is wrong twice, the account will be locked after 5 times", "The password is wrong twice, the account will be locked after 5 times"},
		ErrPwdWrongThird:                {"msg: The password is wrong three times, the account will be locked after 5 times", "The password is wrong three times, the account will be locked after 5 times"},
		ErrPwdWrongFourth:               {"msg: The password is wrong four times, the account will be locked after 5 times", "The password is wrong four times, the account will be locked after 5 times"},
		ErrFirmwareUpgradeFilesMismatch: {"msg: Firmware upgrade files do not match", "Firmware upgrade files do not match"},
		ErrFirmwareUpgradeFPGAWrong:     {"msg: Firmware Upgrade FPGA Status Abnormal", "Firmware Upgrade FPGA Status Abnormal"},
		// 温度相关
		ErrNotOpenTemp:            {"msg: The temp function is not open", "The temp function is not open"},
		ErrNotOpenColdHotTracking: {"msg: The cold hot tracking is not open", "The cold hot tracking is not open"},
		ErrMQTTOpen:               {"msg: The MQTT is open", "The MQTT is open"},
		ErrModbusOpen:             {"msg: The Modbus is open", "The Modbus is open"},
	}

	CodeMsgMap = make(map[string]CodeKey)
)
