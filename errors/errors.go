package errors

import "log"

func CheckErr(err error) bool {

	if err != nil {
		isHandled := HandleError(err)
		if !isHandled {
			log.Fatal(err)
		}
	}
	return false
}

func HandleError(err error) bool {
	//TODO Build some errors handling logic
	return false
}
