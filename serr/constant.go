package serr

type Constant_ErrorCode uint32

const (
	Constant_SUCCESS Constant_ErrorCode = 0
	//public error 1~999
	Constant_ERROR_MULTI_CONFIG                    Constant_ErrorCode = 1
	Constant_ERROR_NO_ACTION                       Constant_ErrorCode = 2
	Constant_ERROR_INITIAL_TRANSITION_NOT_SUBSTATE Constant_ErrorCode = 3
	Constant_ERROR_PARAM_LEN                       Constant_ErrorCode = 4
	Constant_ERROR_PARAM_TYPE                      Constant_ErrorCode = 5

	Constant_ERROR_UNKNOW Constant_ErrorCode = 999
)
