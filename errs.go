// Copyright (c) 2017 Renato Mastrulli. Tutti i diritti riservati. All rights reserved.

package clap

// argError rappresenta un errore relativo a un argomento.
type argError struct {
	argName string
	str     string
}

// NewParseArgError crea un errore generico di parsing dell'argomento. Il parametro text descrive l'errore.
func NewParseArgError(argument string, text string) error {
	return &argError{argName: argument, str: text}
}

// NewExpectedArgError crea un errore per un argomento necessario.
func NewExpectedArgError(argument string) error {
	return &argError{argName: argument, str: "expected"}
}

// NewUnknownArgError crea un errore per un argomento sconosciuto.
func NewUnknownArgError(argument string) error {
	return &argError{argName: argument, str: "unknown"}
}

// NewTooMuchArgError crea un errore per un argomento oltre la fine della sequenza.
func NewTooMuchArgError(argument string) error {
	return &argError{argName: argument, str: "too much arguments"}
}

func (e *argError) Error() string {
	if len(e.argName) > 0 {
		return (e.str + ": " + e.argName)
	}
	return e.str
}
