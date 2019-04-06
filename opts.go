// Copyright (c) 2017 Renato Mastrulli. Tutti i diritti riservati. All rights reserved.

package clap

import (
	"fmt"
	"strings"
)

/*
Option rappresenta un'opzione della linea di comando, ovvero un valore bool.

Per creare un oggetto Option usa la funzione NewOption.

Un'opzione è una parola (word) che rappresenta un valore direttamente o precede il valore bool secondo la sintassi

  parola[=bool]

Il valore in maiuscole o minuscole è uno fra true, t, 1 per bool true e false, f, 0 per bool false.

  opzione=true opzione=t opzione=1
  opzione=false opzione=f opzione=0

Il valore acquisito con la sola parola è quello del campo MatchValue.

Il campo Value contiene il puntatore *bool che riceve il valore.

L'opzione può essere seguita da una sequenza di argomenti contenuti nella lista argomenti del campo Args.

*/
type Option struct {
	//id contiene l'identificativo dell'opzione.
	id string
	//word rappresenta l'opzione nella linea di comando.
	word string
	//Args contiene gli argomenti dell'opzione.
	Args ArgumentList
	//MatchValue è il valore assegnato all'opzione quando presente senza valore.
	MatchValue bool
	//Value contiene il valore.
	Value *bool
	//Help contiene la stringa che descrive l'argomento.
	Help string
}

/*
NewOption crea un'opzione con MatchValue impostato su true.
*/
func NewOption(optWord string, optID string, optValue *bool, optHelp string) *Option {
	optWord = strings.TrimSpace(optWord)
	if len(optWord) == 0 {
		return nil
	}
	return &Option{word: optWord, Value: optValue, MatchValue: true, Help: optHelp}
}

//ID restituisce l'identificativo dell'opzione.
func (opt *Option) ID() string {
	return opt.id
}

//Word restituisce la parola che rappresenta l'opzione nella linea di comando.
func (opt *Option) Word() string {
	return opt.word
}

//AppendArg aggiunge un argomento alla lista argomenti dell'opzione.
func (opt *Option) AppendArg(arg *Argument) {
	if arg != nil {
		if opt.Args == nil {
			opt.Args = make(ArgumentList, 0)
		}
		opt.Args = append(opt.Args, arg)
	}
}

//AppendNewArg crea un argomento nominale con la funzione NewArg, lo aggiunge alla lista argomenti dell'opzione e lo restituisce.
func (opt *Option) AppendNewArg(argWord string, argID string, argRequired bool, argValue interface{}, argHelp string) (arg *Argument) {
	arg = NewArg(argWord, argID, argRequired, argValue, argHelp)
	opt.AppendArg(arg)
	return
}

//Describe restituisce una stringa che descrive la sintassi dell'opzione.
func (opt *Option) Describe(showOptional bool, showValue bool, showArgs bool) string {
	var str string
	str = opt.word
	if showValue {
		str += "[=bool]"
	}
	if showArgs && len(opt.Args) > 0 {
		line, _ := opt.Args.HelpStrings(0)
		str += " "
		str += line
	}
	if showOptional {
		return ("[" + str + "]")
	}
	return str
}

//match verifica che la stringa passata corrisponde all'opzione e memorizza il valore.
func (opt *Option) match(s string) (result bool, err error) {
	//verifica presenza opzione
	if s == opt.word {
		opt.storeValue(opt.MatchValue)
		return true, nil
	}
	// argomento con verb e valore
	//verifica l'impostazione del valore per l'opzione
	//es. word=true
	if val := strings.TrimPrefix(s, (opt.word + "=")); val != s {
		result = true
		switch strings.ToLower(val) {
		case "true", "t", "1":
			opt.storeValue(true)
		case "false", "f", "0":
			opt.storeValue(false)
		default:
			err = fmt.Errorf("not valid value '%s' for option %s", val, opt.word)
		}
	}
	return
}

//storeValue memorizza il valore dell'opzione.
func (opt *Option) storeValue(b bool) {
	if opt.Value != nil {
		*opt.Value = b
	}
}

// =======================================================

// OptionList raccoglie una serie di opzioni.
type OptionList []*Option

/*
HelpStrings crea stringhe di help unendo in linea e in elenco l'help delle opzioni.

Il ritorno line contiene la descrizione su linea singola.
Il ritorno list contiene la descrizione su linee multiple.
*/
func (ol OptionList) HelpStrings(listIndent int) (line string, list string) {
	var indent string
	if listIndent > 0 {
		indent = strings.Repeat(" ", listIndent)
	}
	items := make([]string, 0, len(ol))
	for _, o := range ol {
		if o != nil {
			items = append(items, o.Describe(true, false, false))
			if len(o.Help) > 0 {
				list += fmt.Sprintf("\n%[1]s%[2]s\n  %[1]s%[3]s\n", indent, o.Describe(false, true, true), o.Help)
			} else {
				list += fmt.Sprintf("\n%s%s\n", indent, o.Describe(false, true, true))
			}
			if len(o.Args) > 0 {
				if _, alst := o.Args.HelpStrings(listIndent + 3); len(alst) > 0 {
					list += fmt.Sprintf("\n%s  Argomenti:\n", indent)
					list += alst
				}
			}
		}
	}
	// unisce gli elementi nella linea
	line = strings.Join(items, " ")
	return
}

// =======================================================

//optListEnum gestisce un oggetto OptionList per l'enumerazione.
type optListEnum struct {
	list     OptionList
	position int
}

//createOptListEnum crea un oggetto optListEnum.
func createOptListEnum(optList OptionList) *optListEnum {
	return &optListEnum{list: optList, position: -1}
}

//reset reimposta il cursore di enumerazione.
func (ole *optListEnum) reset() {
	ole.position = -1
}

//next sposta il cursore di enumerazione e restituisce l'opzione in quella posizione.
func (ole *optListEnum) next() *Option {
	ole.position++
	if ole.position < 0 || ole.position >= len(ole.list) {
		return nil
	}
	return ole.list[ole.position]
}
