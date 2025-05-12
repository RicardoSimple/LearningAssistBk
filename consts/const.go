package consts

const (
	Authotization_Header   = "Authorization"
	WebSocketAuthorization = "Sec-WebSocket-Protocol"

	OnLine  = "在线"
	OffLine = "离线"

	UserTypeTeacher = "teacher"
	UserTypeStudent = "student"
	UserTypeAdmin   = "admin"

	ClassNumLength = 6

	DefaultPassword = "123123"
)

var (
	UserTypeToIntMap = map[string]uint{
		UserTypeAdmin:   0,
		UserTypeTeacher: 1,
		UserTypeStudent: 2,
	}
	UserTypeToStringMap = map[uint]string{
		0: UserTypeAdmin,
		1: UserTypeTeacher,
		2: UserTypeStudent,
	}
	GradeGroups = map[string][]string{
		"小学":  {"一年级", "二年级", "三年级", "四年级", "五年级", "六年级"},
		"初中":  {"初一", "初二", "初三"},
		"高中":  {"高一", "高二", "高三"},
		"大学":  {"大一", "大二", "大三", "大四"},
		"研究生": {"研究生"},
	}

	GradeOptions = []string{
		"一年级",
		"二年级",
		"三年级",
		"四年级",
		"五年级",
		"六年级",
		"初一",
		"初二",
		"初三",
		"高一",
		"高二",
		"高三",
		"大一",
		"大二",
		"大三",
		"大四",
		"研究生",
	}
)
