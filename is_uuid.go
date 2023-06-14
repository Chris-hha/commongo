package main

import (
	"fmt"
	// "regexp"
	"github.com/google/uuid"
)

func IsUuid(inputStr string) bool {
	/*
	uuidReg := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	if !uuidReg.MatchString(inputStr) {
		return false
	}
	*/
	parseRes, err := uuid.Parse(inputStr)
	fmt.Println("parseRes: ", parseRes)

	return err == nil
}

func main(){
	inputStr := "cn-chengdu-sdv"
	// inputStr := "87d1a587-7bdc-481b-8309-63931e1740eb"
	// inputStr := "87d1a5877bdc481b830963931e1740eb"
	res := IsUuid(inputStr)
	fmt.Println("inputStr: ", inputStr, "res: ", res)
}
