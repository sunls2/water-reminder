package wechatwork

type Company interface {
	NewApplication(secret string, agentId int) (Application, error)
}

type company struct {
	Id string
}

func NewCompany(id string) Company {
	return &company{Id: id}
}

func (c company) NewApplication(secret string, agentId int) (Application, error) {
	return NewApplication(c.Id, secret, agentId)
}
