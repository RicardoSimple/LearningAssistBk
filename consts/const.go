package consts

const (
	Authotization_Header   = "Authorization"
	WebSocketAuthorization = "Sec-WebSocket-Protocol"

	OnLine  = "在线"
	OffLine = "离线"

	UserTypeTeacher = "teacher"
	UserTypeStudent = "student"
	UserTypeAdmin   = "admin"
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
)
