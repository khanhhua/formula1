package engine

type EngineError struct {
	error
	message string
}

func Error(message string) EngineError {
  return EngineError {
    message: message,
  }
}
