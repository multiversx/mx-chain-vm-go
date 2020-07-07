package async

import "errors"

var ErrAsyncCallGroupDoesNotExist = errors.New("async call group does not exist")

var ErrAsyncCallNotFound = errors.New("async call not found")

var ErrCallBackFuncCalledInRun = errors.New("calling callBack() directly is forbidden")

var ErrCallBackFuncNotExpected = errors.New("unexpected callback was received")
