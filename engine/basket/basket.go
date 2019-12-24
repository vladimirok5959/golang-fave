package basket

import (
	"net/http"
	"os"
	"sync"

	"golang-fave/engine/consts"
	"golang-fave/engine/sqlw"
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper/config"
)

type SBParam struct {
	R         *http.Request
	DB        *sqlw.DB
	Host      string
	Config    *config.Config
	SessionId string
}

type Basket struct {
	sync.RWMutex
	hosts map[string]*onehost
}

func New() *Basket {
	b := Basket{}
	b.hosts = map[string]*onehost{}
	return &b
}

func (this *Basket) Info(p *SBParam) string {
	if p.Host == "" || p.SessionId == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[p.Host]; ok == true {
		if s, ok := h.sessions[p.SessionId]; ok == true {
			s.Preload(p)
			return s.String(p)
		} else {
			return (&dResponse{IsError: false, Msg: "basket_is_empty", Message: ""}).String()
		}
	} else {
		return (&dResponse{IsError: false, Msg: "basket_is_empty", Message: ""}).String()
	}
}

func (this *Basket) Plus(p *SBParam, product_id int) string {
	if p.Host == "" || p.SessionId == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[p.Host]; ok == true {
		if s, ok := h.sessions[p.SessionId]; ok == true {
			s.Preload(p)
			s.Plus(p, product_id)
		} else {
			s := &session{}
			s.listCurrencies = map[int]*currency{}
			s.Products = map[int]*product{}
			s.Preload(p)
			s.Plus(p, product_id)
			h.sessions[p.SessionId] = s
		}
	} else {
		s := &session{}
		s.listCurrencies = map[int]*currency{}
		s.Products = map[int]*product{}
		s.Preload(p)
		s.Plus(p, product_id)
		h := &onehost{}
		h.sessions = map[string]*session{}
		h.sessions[p.SessionId] = s
		this.hosts[p.Host] = h
	}

	return (&dResponse{IsError: false, Msg: "basket_product_plus", Message: ""}).String()
}

func (this *Basket) Minus(p *SBParam, product_id int) string {
	if p.Host == "" || p.SessionId == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[p.Host]; ok == true {
		if s, ok := h.sessions[p.SessionId]; ok == true {
			s.Preload(p)
			s.Minus(p, product_id)
		}
	}

	return (&dResponse{IsError: false, Msg: "basket_product_minus", Message: ""}).String()
}

func (this *Basket) Remove(p *SBParam, product_id int) string {
	if p.Host == "" || p.SessionId == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[p.Host]; ok == true {
		if s, ok := h.sessions[p.SessionId]; ok == true {
			s.Preload(p)
			s.Remove(p, product_id)
		}
	}

	return (&dResponse{IsError: false, Msg: "basket_product_remove", Message: ""}).String()
}

func (this *Basket) ClearBasket(p *SBParam) {
	if p.Host == "" || p.SessionId == "" {
		return
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[p.Host]; ok == true {
		if s, ok := h.sessions[p.SessionId]; ok == true {
			s.Preload(p)
			s.ClearBasket(p)
		}
	}
}

func (this *Basket) ProductsCount(p *SBParam) int {
	if p.Host != "" && p.SessionId != "" {
		this.Lock()
		defer this.Unlock()

		if h, ok := this.hosts[p.Host]; ok == true {
			if s, ok := h.sessions[p.SessionId]; ok == true {
				return s.ProductsCount()
			}
		}
	}

	return 0
}

func (this *Basket) GetAll(p *SBParam) *utils.MySql_basket {
	if p.Host != "" && p.SessionId != "" {
		this.Lock()
		defer this.Unlock()

		if h, ok := this.hosts[p.Host]; ok == true {
			if s, ok := h.sessions[p.SessionId]; ok == true {
				return s.GetAll(p)
			}
		}
	}

	return nil
}

func (this *Basket) Cleanup() {
	this.Lock()
	defer this.Unlock()

	for host_name, host_data := range this.hosts {
		var remove []string
		for session_id, _ := range host_data.sessions {
			session_file := consts.ParamWwwDir + string(os.PathSeparator) + host_name + string(os.PathSeparator) + "tmp" + string(os.PathSeparator) + session_id
			if !utils.IsFileExists(session_file) {
				remove = append(remove, session_id)
			}
		}
		for _, session_id := range remove {
			delete(host_data.sessions, session_id)
		}
	}
}
