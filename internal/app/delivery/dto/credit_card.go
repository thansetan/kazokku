package dto

type CreditCardResponse struct {
	Type    string `json:"type"`
	Number  string `json:"number"`
	Name    string `json:"name"`
	Expired string `json:"expired"`
}
