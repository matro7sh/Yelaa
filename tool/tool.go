package tool

type ToolInterface interface {
	Configure(interface{})
	Run(string)
}
