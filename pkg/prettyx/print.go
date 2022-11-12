package prettyx

import (
	"encoding/json"
	"fmt"
)

func printValue(val interface{}) {
	b, err := json.MarshalIndent(val, "", "	")
	if err != nil {
		fmt.Printf("failed to print val. err: %v \n", err)
		return
	}
	fmt.Println(string(b))
}

func Print(val interface{}) {
	printValue(val)
}

func Printf(comment string, val interface{}) {
	fmt.Printf("[PRINT] %s \n", comment)
	printValue(val)
}
