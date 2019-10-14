package basket

import (
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

func (this *Basket) Info(host, session_id string, db *sqlw.DB, currency_id int) string {
	if host == "" || session_id == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	// Load currency here

	this.Lock()
	defer this.Unlock()
	if h, ok := this.hosts[host]; ok == true {
		if s, ok := h.sessions[session_id]; ok == true {
			return s.String()
		} else {
			return (&dResponse{IsError: false, Msg: "basket_is_empty", Message: ""}).String()
		}
	} else {
		return (&dResponse{IsError: false, Msg: "basket_is_empty", Message: ""}).String()
	}
}

func (this *Basket) Plus(host, session_id string, db *sqlw.DB, product_id int) string {
	if host == "" || session_id == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[host]; ok == true {
		if s, ok := h.sessions[session_id]; ok == true {
			s.Plus(db, product_id)
		}
	} else {
		s := &session{}
		s.Products = map[int]*product{}
		s.Plus(db, product_id)
		h := &onehost{}
		h.sessions = map[string]*session{}
		h.sessions[session_id] = s
		this.hosts[host] = h
	}

	return (&dResponse{IsError: false, Msg: "basket_product_plus", Message: ""}).String()
}

func (this *Basket) Minus(host, session_id string, db *sqlw.DB, product_id int) string {
	if host == "" || session_id == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[host]; ok == true {
		if s, ok := h.sessions[session_id]; ok == true {
			s.Minus(db, product_id)
		}
	}

	return (&dResponse{IsError: false, Msg: "basket_product_minus", Message: ""}).String()
}

func (this *Basket) Remove(host, session_id string, db *sqlw.DB, product_id int) string {
	if host == "" || session_id == "" {
		return (&dResponse{IsError: true, Msg: "basket_host_or_session_not_set", Message: ""}).String()
	}

	this.Lock()
	defer this.Unlock()

	if h, ok := this.hosts[host]; ok == true {
		if s, ok := h.sessions[session_id]; ok == true {
			s.Remove(db, product_id)
		}
	}

	return (&dResponse{IsError: false, Msg: "basket_product_remove", Message: ""}).String()
}
