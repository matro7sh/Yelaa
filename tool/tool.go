package tool

type ToolInterface interface {
	Configure(interface{})
	Info()
	Run(string)
}
