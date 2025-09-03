package main

import (
	"encoding/json"
	"os"
)

func irBuild() {
	bytes, _ := os.ReadFile("./.intext/cache/AST.json")
	nodes := []map[string]interface{}{}
	err := json.Unmarshal(bytes, &nodes)
	Check(err)

	irSlice := []map[string]interface{}{}

	for _, node := range nodes {
		meta := node["meta"].(map[string]interface{})

		switch node["type"] {
		case "output":
			irSlice = append(irSlice, map[string]interface{}{"type": "output", "val": node["value"], "val-type": meta["raw_type"]})
		}
	}

	ir, err := json.MarshalIndent(irSlice, "", "  ")
	Check(err)
	os.WriteFile("ir.json", ir, 0644)
}
