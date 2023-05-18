package callorder

import (
	"fmt"

	"go.uber.org/fx"
)

type Config struct {
	data string
}

type (
	A struct{}
	B struct{}
	C struct{}
	D struct{}
)

func provide0(cfg *Config) *A {
	fmt.Println("Config.data = ", cfg.data)
	return &A{}
}

func provide1(a *A) *B {
	fmt.Println("provide1 is executing...")
	return &B{}
}
func provide2(b *B) *C {
	fmt.Println("provide2 is executing...")
	return &C{}
}
func invoker1(c *C) {
	fmt.Println("invoker1 is executing...")
}
func invoker2(c *C) {
	fmt.Println("invoker2 is executing...")
}

func printBanner(cfg *Config) {
	fmt.Println("Print banner here ...", "data=", cfg.data)
}

// How to run this demo
// go build
// .\go-learning-fx.exe (make sure callorder.Main() is uncommented in main.go)
// we can the output in this order:
// [Fx] INVOKE             hi-supergirl/go-learning-fx/callOrder.printBanner()
// Print banner here ... data= hello world
// [Fx] INVOKE             hi-supergirl/go-learning-fx/callOrder.invoker1()
// Config.data =  hello world
// provide1 is executing...
// provide2 is executing...
// invoker1 is executing...
// [Fx] INVOKE             hi-supergirl/go-learning-fx/callOrder.invoker2()
// invoker2 is executing...
// [Fx] INVOKE             hi-supergirl/go-learning-fx/callOrder.Main.func2()
// Invoke2.....  func(*C) is executing...
// [Fx] RUNNING
func Main() {
	config := &Config{data: "hello world"}

	app := fx.New(
		fx.Supply(config),      // provide the initial variables used in Provide
		fx.Invoke(printBanner), // this is the first invoke to be executed
		fx.Provide(
			provide0,
			provide1,
			provide2,
		),
		fx.Invoke(
			invoker1, // this is the second invoke to be executed
			invoker2, // this is the third invoke to be executed
			func(*C) { // this is the fourth invoke to be executed
				fmt.Println("Invoke2.....  func(*C) is executing...")
			}),
	)
	app.Run()
}
