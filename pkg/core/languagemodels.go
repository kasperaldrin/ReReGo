package core

type LanguageModelAPI struct {
	// A running model to call.
	call func()
}
