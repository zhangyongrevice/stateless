package serr

type Constant_ErrorCode int32

const (
	Constant_SUCCESS Constant_ErrorCode = 0
	//public error 1~999
	Constant_ERROR_MULTI_CONFIG                    = 1
	Constant_ERROR_NO_ACTION                       = 2
	Constant_ERROR_CONFIC_INIT                     = 3
	Constant_ERROR_INITIAL_TRANSITION_NOT_SUBSTATE = 4
	Constant_ERROR_PARAM_LEN                       = 5
	Constant_ERROR_PARAM_TYPE                      = 6

	Constant_ERROR_UNKNOW Constant_ErrorCode = 999
)
