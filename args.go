// Copyright (c) 2017 Renato Mastrulli. Tutti i diritti riservati. All rights reserved.

package clap

import (
	"fmt"
	"strconv"
	"strings"
)

/*
Argument rappresenta un argomento del comando.

Puo essere nominale, creato con NewArg, oppure generico, creato con NewGenericArg.

Un argomento nominale possiede una parola (word) che precede il valore secondo la sintassi

  parola=valore

mentre un argomento generico intercetta un valore senza alcuna parola e il suo metodo Word restituisce una stringa vuota.

Value è l'oggetto che riceve il valore, cioè StoreValue, ArgumentStore o un puntatore a variabile.
I tipi variabile riconosciuti sono:
  string
  bool
  int
  int8
  int16
  int32
  int64
  uint
  uint8
  uint16
  uint32
  uint64
  float32
  float64
e il parsing del valore avviene con i metodi Parse di strconv.

L'identificativo dell'argomento (id) non è strettamente necessario se non per gli argomenti generici, può essere impostato su una stringa vuota ma non può essere modificato dopo l'inizializzazione.

*/
type Argument struct {
	// id contiene l'identificativo dell'argomento.
	id string
	// word rappresenta l'argomento nella linea di comando.
	word string
	// Required indica se l'argomento deve essere specificato.
	Required bool
	// Value contiene il valore.
	Value interface{}
	// Help contiene la stringa che descrive l'argomento.
	Help string
}

/*
StoreValue è il tipo funzione che memorizza il valore str nell'argomento parent.

La funzione deve restituire nil se il valore è stato memorizzato con successo, altrimenti l'errore che sarà restituito dalla funzione Parse.
*/
type StoreValue func(str string, parent Argument) (err error)

/*
ArgumentStore è l'interfaccia per memorizzare il valore di un argomento.

Il metodo StoreValue deve agire come il tipo funzione omonimo, vedi la sua documentazione.
*/
type ArgumentStore interface {
	// StoreValue memorizza il valore di un argomento e restituisce l'eventuale errore.
	StoreValue(str string, parent Argument) (err error)
}

/*
NewGenericArg crea un argomento generico.

Il metodo Word restituirà una stringa vuota.
*/
func NewGenericArg(argID string, argRequired bool, argValue interface{}, argHelp string) *Argument {
	return NewArg("", argID, argRequired, argValue, argHelp)
}

// NewArg crea un argomento nominale.
func NewArg(argWord string, argID string, argRequired bool, argValue interface{}, argHelp string) *Argument {
	if len(argID) == 0 {
		return nil
	}
	argWord = strings.TrimSpace(argWord)
	return &Argument{id: argID, word: argWord, Required: argRequired, Value: argValue, Help: argHelp}
}

// ID restituisce l'identificativo dell'argomento.
// L'identificativo è usato al posto della parola per rappresentare gli argomenti generici nelle stringhe di help e di errore.
func (arg *Argument) ID() string {
	return arg.id
}

// Word restituisce la parola che rappresenta l'argomento nella linea di comando.
func (arg *Argument) Word() string {
	return arg.word
}

// IsGeneric restituisce true se l'argomento non ha una parola rappresentativa.
func (arg *Argument) IsGeneric() bool {
	return (len(arg.word) == 0)
}

// Describe restituisce una stringa che descrive la sintassi dell'argomento.
// L'identificativo dell'argomento è usato al posto della parola per rappresentare gli argomenti generici.
func (arg *Argument) Describe(showOptional bool, showValue bool) string {
	var str string
	if arg.IsGeneric() {
		str = ("<" + arg.id + ">")
	} else {
		str = arg.word
		if showValue {
			str += "=value"
		}
	}
	if (arg.Required == false) && showOptional {
		return ("[" + str + "]")
	}
	return str
}

// match verifica che la stringa passata corrisponde all'argomento e memorizza il valore.
func (arg *Argument) match(s string) (result bool, err error) {
	// argomento generico
	if arg.IsGeneric() {
		result = true
		err = arg.storeValue(s)
		return
	}
	// argomento con word e valore
	if sv := strings.TrimPrefix(s, (arg.word + "=")); sv != s {
		result = true
		err = arg.storeValue(sv)
		return
	}
	// argomento non corrispondente
	if arg.Required {
		if arg.IsGeneric() {
			err = NewExpectedArgError(arg.id)
		} else {
			err = NewExpectedArgError(arg.word)
		}
	}
	return
}

// storeValue memorizza il valore dell'argomento.
func (arg *Argument) storeValue(str string) (err error) {
	if arg.Value != nil {
		switch v := arg.Value.(type) {
		case *string:
			*v = str
		case *bool:
			var bval bool
			bval, err = strconv.ParseBool(strings.ToLower(str))
			if err == nil {
				*v = bval
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *int:
			var uval int64
			uval, err = strconv.ParseInt(str, 10, 0)
			if err == nil {
				*v = int(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *int8:
			var uval int64
			uval, err = strconv.ParseInt(str, 10, 8)
			if err == nil {
				*v = int8(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *int16:
			var uval int64
			uval, err = strconv.ParseInt(str, 10, 16)
			if err == nil {
				*v = int16(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *int32:
			var uval int64
			uval, err = strconv.ParseInt(str, 10, 32)
			if err == nil {
				*v = int32(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *int64:
			var uval int64
			uval, err = strconv.ParseInt(str, 10, 64)
			if err == nil {
				*v = int64(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *uint:
			var uval uint64
			uval, err = strconv.ParseUint(str, 10, 0)
			if err == nil {
				*v = uint(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *uint8:
			var uval uint64
			uval, err = strconv.ParseUint(str, 10, 8)
			if err == nil {
				*v = uint8(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *uint16:
			var uval uint64
			uval, err = strconv.ParseUint(str, 10, 16)
			if err == nil {
				*v = uint16(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *uint32:
			var uval uint64
			uval, err = strconv.ParseUint(str, 10, 32)
			if err == nil {
				*v = uint32(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *uint64:
			var uval uint64
			uval, err = strconv.ParseUint(str, 10, 64)
			if err == nil {
				*v = uint64(uval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case StoreValue:
			err = v(str, *arg)
		case ArgumentStore:
			err = v.StoreValue(str, *arg)
		case *float32:
			var fval float64
			fval, err = strconv.ParseFloat(str, 32)
			if err == nil {
				*v = float32(fval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		case *float64:
			var fval float64
			fval, err = strconv.ParseFloat(str, 64)
			if err == nil {
				*v = float64(fval)
			} else {
				err = fmt.Errorf("not valid value '%s' for %s", str, arg.Describe(false, false))
			}
		}
	}
	return
}

// =======================================================

// ArgumentList contiene gli argomenti nell'ordine di sequenza di analisi.
type ArgumentList []*Argument

/*
HelpStrings crea stringhe di help unendo in linea e in elenco la descrizione degli argomenti.

Il ritorno line contiene la descrizione su linea singola.
Il ritorno list contiene la descrizione su linee multiple.
*/
func (al ArgumentList) HelpStrings(listIndent int) (line string, list string) {
	var indent string
	if listIndent > 0 {
		indent = strings.Repeat(" ", listIndent)
	}
	items := make([]string, 0, len(al))
	for _, a := range al {
		if a != nil {
			items = append(items, a.Describe(true, true))
			if len(a.Help) > 0 {
				list += fmt.Sprintf("\n%[1]s%[2]s\n  %[1]s%[3]s\n", indent, a.Describe(false, false), a.Help)
			} else {
				list += fmt.Sprintf("\n%s%s\n", indent, a.Describe(false, false))
			}
		}
	}
	// unisce gli elementi nella linea
	line = strings.Join(items, " ")
	return
}

// =======================================================

// argListEnum gestisce un oggetto ArgumentList per l'enumerazione.
type argListEnum struct {
	list     ArgumentList
	position int
}

// createargListEnum crea un oggetto argListEnum.
func createArgListEnum(argList ArgumentList) *argListEnum {
	return &argListEnum{list: argList, position: -1}
}

// reset reimposta il cursore di enumerazione.
func (ale *argListEnum) reset() {
	ale.position = -1
}

// next sposta il cursore di enumerazione e restituisce l'argomento in quella posizione.
func (ale *argListEnum) next() *Argument {
	ale.position++
	if ale.position < 0 || ale.position >= len(ale.list) {
		return nil
	}
	return ale.list[ale.position]
}
