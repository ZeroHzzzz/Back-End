package utils

type MyError struct {
	Code int
	Msg  string
	Data interface{}
}

// 考虑分为用户，数据格式，数据库
var (
	LOGIN_ERROR              = NewError(101, "账号或密码错误")
	UNAUTHORIZED             = NewError(102, "没有足够权限")
	NOACCESS                 = NewError(103, "未到接口开放时间")
	CONTEXT_ERROR            = NewError(104, "上下文错误")
	PARAM_ERROR              = NewError(201, "参数错误")
	FILE_ERROR               = NewError(203, "文件错误")
	DECODE_ERROR             = NewError(204, "解码错误")
	MONGODB_INIT_ERROR       = NewError(301, "MongoDB Client 初始化失败")
	RMQ_INIT_ERROR           = NewError(302, "RabbitMQ 初始化失败")
	DATABASE_OPERATION_ERROR = NewError(303, "数据库操作失败")
	QUEUE_OPERATION_ERROR    = NewError(304, "队列操作失败")
	INNER_ERROR              = NewError(501, "服务器内部错误")
	CONNECT_ERROR            = NewError(502, "连接异常")
	LIMIT_ERROR              = NewError(503, "当前访问人数过多，服务器限流")
)

func (e *MyError) Error() string {
	return e.Msg
}

func NewError(code int, msg string) *MyError {
	return &MyError{
		Msg:  msg,
		Code: code,
	}
}

func GetError(e *MyError, data interface{}) *MyError {
	return &MyError{
		Msg:  e.Msg,
		Code: e.Code,
		Data: data,
	}
}
