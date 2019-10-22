package main

import(
	"app/gate"
	"app"
)

func main()  {
	var myapp = gate.GateApp{}
	app.Run(&myapp)
}
