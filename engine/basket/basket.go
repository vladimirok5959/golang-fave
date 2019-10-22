package basket

import (
	"net/http"
	"sync"

	"golang-fave/engine/sqlw"
)

type Basket struct {
	sync.RWMutex
	hosts map[string]*onehost
}

func New() *Basket {
	b := Basket{}
	b.hosts = map[string]*onehost{}
	return &b
}

func (this *Basket) Info(r *http.Request, host, session_id string, db *sqlw.DB, currency_id int) string {
	if host == "" || session_id == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	// Load currency here

	this.Lock()
	defer this.Unlock()
	if h, ok := this.hosts[host]; ok == true {
		if s, ok := h.sessions[session_id]; ok == true {
			s.Preload(r, db)
			return s.String(db)
		} else {
			return (&dResponse{IsError: false, Msg: "basket_is_empty", Message: ""}).String()
		}
	} else {
		return (&dResponse{IsError: false, Msg: "basket_is_empty", Message: ""}).String()
	}
}

func (this *Basket) Plus(r *http.Request, host, session_id string, db *sqlw.DB, product_id int) string {
	if host == "" || session_id == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[host]; ok == true {
		if s, ok := h.sessions[session_id]; ok == true {
			s.Preload(r, db)
			s.Plus(db, product_id)
		}
	} else {
		s := &session{}
		s.listCurrencies = map[int]*currency{}
		s.Products = map[int]*product{}
		s.Preload(r, db)
		s.Plus(db, product_id)
		h := &onehost{}
		h.sessions = map[string]*session{}
		h.sessions[session_id] = s
		this.hosts[host] = h
	}

	return (&dResponse{IsError: false, Msg: "basket_product_plus", Message: ""}).String()
}

func (this *Basket) Minus(r *http.Request, host, session_id string, db *sqlw.DB, product_id int) string {
	if host == "" || session_id == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[host]; ok == true {
		if s, ok := h.sessions[session_id]; ok == true {
			s.Preload(r, db)
			s.Minus(db, product_id)
		}
	}

	return (&dResponse{IsError: false, Msg: "basket_product_minus", Message: ""}).String()
}

func (this *Basket) Remove(r *http.Request, host, session_id string, db *sqlw.DB, product_id int) string {
	if host == "" || session_id == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[host]; ok == true {
		if s, ok := h.sessions[session_id]; ok == true {
			s.Preload(r, db)
			s.Remove(db, product_id)
		}
	}

	return (&dResponse{IsError: false, Msg: "basket_product_remove", Message: ""}).String()
}

func (this *Basket) ProductsCount(r *http.Request, host, session_id string) int {
	if host != "" && session_id != "" {
		this.Lock()
		defer this.Unlock()

		if h, ok := this.hosts[host]; ok == true {
			if s, ok := h.sessions[session_id]; ok == true {
				return s.ProductsCount()
			}
		}
	}

	return 0
}
