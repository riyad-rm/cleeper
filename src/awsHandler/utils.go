package awsHandler

//import (
//	"fmt"
//)

func stringInList(list *[]string, str string) bool{
	for _, elem := range *list {
		//fmt.Println("list elem: ", elem)
		//fmt.Println("looking for: ", str)
		if elem == str {
			return true
		}
	}
	return false
}