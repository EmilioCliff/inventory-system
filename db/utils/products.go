package utils

import (
	"math/rand"
	"time"
)

var products = []map[string]interface{}{
	{
		"Name":      "HIV 1&2 Test Monopack cassette",
		"Pack Size": 10,
		"Price":     1100,
	},
	{
		"Name":      "HIV 1&2 Test Monopack cassette(strips)",
		"Pack Size": 30,
		"Price":     2550,
	},
	{
		"Name":      "HIV 1&2 Test Monopack cassette(strips)",
		"Pack Size": 50,
		"Price":     4250,
	},
	{
		"Name":      "H. Pylori Ab Device",
		"Pack Size": 30,
		"Price":     2100,
	},
	{
		"Name":      "H. Pylori Ag",
		"Pack Size": 25,
		"Price":     2500,
	},
	{
		"Name":      "hCG (Pregnancy) Strip - Urine",
		"Pack Size": 50,
		"Price":     400,
	},
	{
		"Name":      "Gonorrhea Test Cassette(swap/Urine)",
		"Pack Size": 25,
		"Price":     2500,
	},
	{
		"Name":      "LH (Ovulation) Device - Urine",
		"Pack Size": 50,
		"Price":     2500,
	},
	{
		"Name":      "Malaria Anigen Pf/Pan Device",
		"Pack Size": 50,
		"Price":     3500,
	},
	{
		"Name":      "Malaria Anigen Pf/Pv Device",
		"Pack Size": 50,
		"Price":     3500,
	},
	{
		"Name":      "Syphillis Device",
		"Pack Size": 50,
		"Price":     3500,
	},
	{
		"Name":      "Syphillis Strip",
		"Pack Size": 100,
		"Price":     3000,
	},
	{
		"Name":      "Toxo IgG/IgM Device",
		"Pack Size": 30,
		"Price":     2700,
	},
	{
		"Name":      "Troponin - I Device",
		"Pack Size": 30,
		"Price":     7500,
	},
	{
		"Name":      "Tsutsugamushi Ab Device",
		"Pack Size": 30,
		"Price":     8100,
	},
	{
		"Name":      "Typhoid IgG/IgM Device",
		"Pack Size": 50,
		"Price":     4300,
	},
	{
		"Name":      "Salmonella Ag",
		"Pack Size": 25,
		"Price":     2800,
	},
}

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
}

func RandomProduct() map[string]interface{} {
	k := len(products)
	return products[rand.Intn(k+1)]
}
