package sora_errors

type NotSupportedError struct {
}

func (e NotSupportedError) Error() string {
	return "Không hỗ trợ"
}

type LogicError struct {
	Code    string
	Message string
}

func (e LogicError) Error() string {
	return "Có lỗi logic"
}
