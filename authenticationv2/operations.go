package authenticationv2

type OperationName string

const (
	OperationGetAvailableExports OperationName = "GetAvailableExports"
)

//
// OperationResult
// @Description:
//
type OperationResult struct {
	UUID   string      //连接uuid
	Result interface{} //结果
	Error  string      //错误
}

//
// Operation
// @Description:
//
type Operation struct {
	UUID      string                 //连接uuid
	User      string                 //发起人
	Operation OperationName          //操作
	Params    map[string]interface{} //参数
}

//
// GetParams
// @Description:
// @receiver o
// @param key
//
func (o *Operation) GetParams(key string) interface{} {
	if o.Params != nil {
		return o.Params[key]
	}
	return nil
}
