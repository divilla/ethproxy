package interfaces

type ErrorLogger interface {
	Error(i ...interface{})
	Errorf(format string, args ...interface{})
	Panic(i ...interface{})
	Panicf(format string, args ...interface{})
	//Info(i ...interface{})
	//Infof(format string, args ...interface{})
	//Infoj(j log.JSON)
	//Warn(i ...interface{})
	//Warnf(format string, args ...interface{})
	//Warnj(j log.JSON)
	//Errorj(j log.JSON)
	//Fatal(i ...interface{})
	//Fatalj(j log.JSON)
	//Fatalf(format string, args ...interface{})
	//Panicj(j log.JSON)
}
