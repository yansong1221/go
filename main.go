package main

import(
	"app/gate"
	"app"
	"fmt"
)

func main()  {
	var myapp = gate.GateApp{}
	app.Run(&myapp)

	fmt.Println("-------------")
}
