package entity

import "time"

type InputEvent struct {
	Prompt string `json:"prompt"`
}

type OutputEvent struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TxInfo struct {
	BlockHash     string    `json:"block_hash"`
	BlockHeight   int       `json:"block_height"`
	BlockIndex    int       `json:"block_index"`
	Hash          string    `json:"hash"`
	Addresses     []string  `json:"addresses"`
	Total         int       `json:"total"`
	Fees          int       `json:"fees"`
	Size          int       `json:"size"`
	Vsize         int       `json:"vsize"`
	Preference    string    `json:"preference"`
	Confirmed     time.Time `json:"confirmed"`
	Received      time.Time `json:"received"`
	Ver           int       `json:"ver"`
	DoubleSpend   bool      `json:"double_spend"`
	VinSz         int       `json:"vin_sz"`
	VoutSz        int       `json:"vout_sz"`
	OptInRbf      bool      `json:"opt_in_rbf"`
	Confirmations int       `json:"confirmations"`
	Confidence    int       `json:"confidence"`
	Inputs        []Inputs  `json:"inputs"`
	Outputs       []Outputs `json:"outputs"`
}

type Inputs struct {
	PrevHash    string   `json:"prev_hash"`
	OutputIndex int      `json:"output_index"`
	OutputValue int      `json:"output_value"`
	Sequence    int64    `json:"sequence"`
	Addresses   []string `json:"addresses"`
	ScriptType  string   `json:"script_type"`
	Age         int      `json:"age"`
	Witness     []string `json:"witness"`
}

type Outputs struct {
	Value      int      `json:"value"`
	Script     string   `json:"script"`
	Addresses  []string `json:"addresses"`
	ScriptType string   `json:"script_type"`
}
