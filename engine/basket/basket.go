package basket

import (
	"net/http"
	"sync"

	"golang-fave/engine/sqlw"
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
			s.Preload(p.R, p.DB)
			return s.String(p.DB, (*p.Config).Shop.Price.Format, (*p.Config).Shop.Price.Round)
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
			s.Preload(p.R, p.DB)
			s.Plus(p.DB, product_id, (*p.Config).Shop.Price.Format, (*p.Config).Shop.Price.Round)
		}
	} else {
		s := &session{}
		s.listCurrencies = map[int]*currency{}
		s.Products = map[int]*product{}
		s.Preload(p.R, p.DB)
		s.Plus(p.DB, product_id, (*p.Config).Shop.Price.Format, (*p.Config).Shop.Price.Round)
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
			s.Preload(p.R, p.DB)
			s.Minus(p.DB, product_id, (*p.Config).Shop.Price.Format, (*p.Config).Shop.Price.Round)
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
			s.Preload(p.R, p.DB)
			s.Remove(p.DB, product_id, (*p.Config).Shop.Price.Format, (*p.Config).Shop.Price.Round)
		}
	}

	return (&dResponse{IsError: false, Msg: "basket_product_remove", Message: ""}).String()
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
