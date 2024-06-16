package business

import "github.com/hanselacn/banking-transaction/repo"

type Business interface {
}

type business struct {
	Repo repo.Repo
}

func NewBusiness(repo repo.Repo) Business {
	return business{Repo: repo}
}

func (b business) Create() {

}
