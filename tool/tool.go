package tool

type ToolInterface interface {
	Configure(interface{})
	Info(string)
	Run(string)
}
