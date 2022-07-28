package main

import (
	"api_bank/momo"
	"fmt"
)

func main() {
	momoService := momo.NewMomo("sdt", "pass")
	//fmt.Println(momoService.RegNewDevice()) // Gửi otp
	//fmt.Println(momoService.VerifyDevice("5940")) // Xác thực otp
	momoService.UserLogin("phash", "key")                           // điền phash và key từ cái xác thực otp nó trả về
	data, err := momoService.GetHistory("1/6/2022", "2/6/2022", 10) // get lsgd
	fmt.Println(data, err)                                          // in
}
