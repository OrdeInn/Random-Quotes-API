package models

type Quotation struct {
	Id        int    `json:"id"`
	Author    string `json:"author"`
	Quotation string `json:"quotation"`
}
