package main

import "fmt"

func main() {
	fmt.Println("start")

	endpointMap := make(map[string]string)

	endpointMap["big.txt"] = "/uploads"
	endpointMap["big.txt"] = "/uploads/index.php"
	endpointMap["common.txt"] = "/admin"

	for key, value := range endpointMap {
		fmt.Println("Key: ", key, "Value: ", value)
	}

}
