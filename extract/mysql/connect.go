package mysql

type Connect struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
}

// 打开数据库连接
func (c *Connect) OpenConnect() error {

	return nil
}

// 加载认证信息
func (c *Connect) LoadAuthInfo() error {

	return nil
}
