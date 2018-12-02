package types

import (
	"fmt"
)

const (
	FormatErrorMessage    = "命令格式错误，检查后输入"
	ArgErrorMessage       = "缺少参数，检查参数列表及格式后输入"
	ArgFormatErrorMessage = "参数格式错误"
)

var (
	FormatError    = fmt.Errorf(FormatErrorMessage)
	ArgError       = fmt.Errorf(ArgErrorMessage)
	ArgFormatError = fmt.Errorf(ArgFormatErrorMessage)
)
