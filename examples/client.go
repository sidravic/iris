package main

import (
	"github.com/supersid/iris/client"
	"fmt"
)


func SendMessage(c *client.Client, msg string){
	c.SendMessage("echo", msg)
}



func main() {


	_c := client.Start("tcp://127.0.0.1:5555")





	retry_count := 0
        seq := 1
	for {
		msg := fmt.Sprintf("Hello %d", seq)
		SendMessage(_c, msg)

		err, m := _c.ReceiveMessage()

		if err != nil {
			break;
		}

		if len(m) == 0 {
			SendMessage(_c, msg)
			if retry_count == 3 {
				seq++
				continue
			}else{
				retry_count++
				continue
			}


		}

		fmt.Println("------------------------")
		fmt.Println(m)
		fmt.Println("------------------------")
		seq++
	}





}
