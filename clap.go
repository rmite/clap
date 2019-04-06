// Copyright (c) 2017 Renato Mastrulli. Tutti i diritti riservati. All rights reserved.

//Command Line Arguments Processing

/*
Package clap implementa l'elaborazione di argomenti in linea di comando.

La funzione principale è Parse che analizza una sequenza di stringhe cercando un comando, i suoi argomenti e infine eventuali opzioni.

La sequenza è questa:

  comando argomenti [argomenti opzionali] [opzioni]

Le strutture Command, Argument e Option rappresentano rispettivamente comando, argomento e opzione. Sono contraddistinti nella linea di comando con una parola loro assegnata e restituita dal metodo Word.

Gli oggetti Command devono essere creati prima di chiamare Parse, mentre gli oggetti Argument e Option possono essere definiti sia prima con i metodi Append sia in maniera dinamica usando il metodo Prepare di Command.

Il campo Data di Command può contenere i dati del comando come ad esempio il puntatore a una struttura i cui campi sono stati referenziati come puntatori nel campo Value di Argument e Option.

Gli argomenti catturano un valore. Possono essere nominali oppure generici. Gli argomenti nominali sono contraddistinti da una parola con la sintassi:

  parola=valore

mentre gli argomenti generici prevedono solo il valore:

  valore

Le opzioni sono sempre contraddistinte da una parola a cui può seguire un valore booleano:

  parola[=bool]

se il valore non è specificato, l'opzione acquisisce il valore default MatchValue impostato nella sua definizione.

Un'opzione può avere anche una serie di argomenti che devono essere indicati dopo essa nella riga di comando.

I valori degli argomenti sono convertiti nei tipi standard o attribuiti con StoreValue e ArgumentStore, vedi la descrizione di Argument, mentre il valore di un'opzione è bool, vedi la descrizione di Option.

L'analisi permette di verificare se sono stati specificati un comando valido e gli argomenti necessari.
L'esempio qui in basso mostra come preparare e usare i comandi, inoltre evidenzia la funzionalità che trascrive l'help dei comandi.

Esempio

Definizione di comando e argomenti prima di chiamare Parse

  var commands clap.CommandMap

  func main() {

    cmdHelp := clap.NewCommand("-h", "Help", "Mostra l'help generale o per un comando.")
    cmdHelp.Data = new(string)
    cmdHelp.Exec = showHelp
    cmdHelp.AppendNewArg("", "cmd", false, cmdHelp.Data, "comando da descrivere")

    cmdVersion := clap.NewCommand("-v", "Versione", "Mostra la versione.")
    cmdVersion.Exec = showVersion

    commands = clap.NewCommandMap(2)
    commands.Insert(cmdHelp)
    commands.Insert(cmdVersion)

    var cmd *clap.Command
    var err error

    if cmd, err = clap.Parse(os.Args[1:], commands, nil); err != nil {
      fmt.Println(err)
      fmt.Println("")
      fmt.Println("Usa -h per vedere l'help.")
      return
    }

    if cmd != nil {
      cmd.Execute()
    } else {
      fmt.Println("Specifica un comando. Usa -h per vedere l'help.")
    }

  }

  func showHelp(data interface{}) {
    var str string
    if data != nil {
      str = *data.(*string)
    }
    if str != "" {
      if cmd, ok := commands[str]; ok {
        cmd.ShowFullHelp()
      } else {
        fmt.Printf("Comando inesistente: %s\n", str)
      }
      return
    }
    fmt.Println("Help generale")
    commands.ShowHelp()
  }

  func showVersion(data interface{}) {
    fmt.Println("Versione 1.0")
  }

e in maniera dinamica

  type fixArgs struct {
    input      string
    output     string
    encoding   string
    overwrite  bool
    appendText bool
  }

  var commands clap.CommandMap

  func main() {

    [...]

    cmdFix := clap.NewCommand("-fix", "Fix", "Sistema la codifica.")
    cmdFix.Prep = prepareFixArgs
    cmdFix.Exec = doFix

    commands = clap.NewCommandMap(3)
    commands.Insert(cmdHelp)
    commands.Insert(cmdVersion)
    commands.Insert(cmdFix)  // non inserire il comando cmdFix nel map se vuoi usarlo solo come default con Parse

    var cmd *clap.Command
    var err error

    if cmd, err = clap.Parse(os.Args[1:], commands, cmdFix); err != nil {
      fmt.Println(err)
      fmt.Println("")
      fmt.Println("Usa -h per vedere l'help.")
      return
    }

    if cmd != nil {
      cmd.Execute()
    } else {
      fmt.Println("Specifica un comando. Usa -h per vedere l'help.")
    }

  }

  [...]

  // esempio comando
  // -fix input output [-enc=value] [-w] [-a]

  func prepareFixArgs(cmd *clap.Command) {
    fixdata := fixArgs{input: "none", output: "none", encoding: "none"}
    cmd.AppendNewArg("", "input", true, &fixdata.input, "file sorgente")
    cmd.AppendNewArg("", "output", true, &fixdata.output, "file destinazione")
    cmd.AppendNewArg("-enc", "encoding", false, &fixdata.encoding, "codifica")
    cmd.AppendNewOpt("-w", "overwrite", &fixdata.overwrite, "sovrascrivi")
    cmd.AppendNewOpt("-a", "append", &fixdata.appendText, "aggiungi")
    cmd.Data = &fixdata
  }

  func doFix(data interface{}) {
    if d, ok := data.(*fixArgs); ok {
      fmt.Printf("Input: '%s'\n", d.input)
      fmt.Printf("Output: '%s'\n", d.output)
      fmt.Printf("Encoding: '%s'\n", d.encoding)
      fmt.Printf("Overwrite: %v\n", d.overwrite)
      fmt.Printf("Append: %v\n", d.appendText)
      //codice...
    }
  }

Avendo specificato cmdFix come comando default nella chiamata a Parse, il comando può essere riconosciuto sia con la sua parola -fix

  -fix source.txt dest.txt -enc=UTF8 -a

sia senza, ovvero la linea di comando contiene direttamente i suoi argomenti

  source.txt dest.txt -enc=UTF8 -a

questo è l'unico modo possibile se cmdFix non fosse stato aggiunto al map commands.

Se cmdFix non fosse specificato come default, per essere riconosciuto dovrebbe essere aggiunto al map commands e la sua parola sarebbe obbligatoria.

*/
package clap

import "strings"

/*
Parse esamina uno slice di stringhe per individuare un comando da eseguire e lo restituisce dopo avergli attribuito i valori dei vari argomenti.

Una volta individuato un comando, la funzione chiama il suo metodo Prepare e prosegue con la ricerca dei valori.

La ricerca prosegue in sequenza comando -> argomenti -> opzioni.
Le opzioni sono tutte opzionali, se un'opzione ha degli argomenti questi vengono cercati dopo aver trovato l'opzione.
Dopo gli argomenti, che siano del comando o di un'opzione, vengono cercate altre eventuali opzioni fino ad esaurimento delle stringhe passate alla funzione Parse.

Gli argomenti sono cercati nell'ordine in cui sono stati inseriti nelle rispettive liste argomenti, le opzioni invece possono trovarsi in qualsiasi ordine dopo gli argomenti.

Se la prima stringa in args non è un comando fra quelli contenuti in cmdMap, la funzione cerca di attribuire i valori al comando default cmdDefault se diverso da nil.

La funzione si interrompe quando ha analizzato tutte le stringhe oppure se c'è un errore nell'attribuzione di un valore, se c'è un argomento sconosciuto o se un argomento obbligatorio non è stato specificato.
La funzione restituisce l'oggetto Command che ha individuato insieme all'eventuale errore.

Se args è vuoto o cmdMap non contiene il comando e cmdDefault è nil, la funzione restituisce nil come comando. Se cmdDefault è nil, restituisce anche un errore di comando sconosciuto.

*/
func Parse(args []string, cmdMap CommandMap, cmdDefault *Command) (cmd *Command, err error) {
	if len(args) == 0 {
		return nil, nil
	}
	var result bool
	var ale *argListEnum
	var curArg *Argument
	var ole *optListEnum
	var curOpt *Option
	count := 0
	// cerca il comando
	if cmd, result = cmdMap[args[count]]; result {
		// comando trovato
		count++
	} else {
		// comando non trovato
		if cmdDefault == nil {
			return nil, NewUnknownArgError(args[count])
		}
		// usa il comando default
		cmd = cmdDefault
	}
	// prepara il comando se necessario
	cmd.Prepare()
	// imposta e avvia l'enumeratore degli argomenti
	ale = createArgListEnum(cmd.Args)
	ale.reset()
	curArg = ale.next()
	// imposta l'enumeratore delle opzioni
	ole = createOptListEnum(cmd.Opts)
	ole.reset()
	// ciclo di analisi argomenti
	for {
		if count < len(args) {
			if curArg == nil {
				if cmd.Opts == nil {
					err = NewTooMuchArgError(args[count])
					return
				}
				// cerca fra le opzioni
				oldcnt := count
				curOpt = ole.next()
				for curOpt != nil {
					if result, err = curOpt.match(args[count]); err == nil {
						if result {
							// opzione corrispondente
							count++ // passa alla stringa successiva
							// avvia l'enumeratore degli argomenti dell'opzione
							ale = createArgListEnum(curOpt.Args)
							ale.reset()
							curArg = ale.next()
							// reimposta l'enumeratore delle opzioni
							ole.reset()
							break
						}
					} else {
						// errore, es. valore non valido
						return
					}
					curOpt = ole.next()
				} // for curOpt != nil
				if oldcnt == count {
					// nessuna corrispondenza
					err = NewUnknownArgError(args[count])
					return
				}
			} else { // curArg != nil
				// verifica corrispondenza dell'argomento corrente
				if result, err = curArg.match(args[count]); err == nil {
					if result {
						// argomento corrispondente
						count++ // passa alla stringa successiva
					}
				} else {
					// errore, es. argomento mancante o valore non valido
					return
				}
				curArg = ale.next()
			}
		} else {
			// la linea di comando non contiene altri argomenti
			skipped := make([]string, 1)
			for curArg != nil {
				if curArg.Required {
					// ci deve essere un argomento
					skipped = append(skipped, curArg.Describe(true, false))
					err = NewExpectedArgError(strings.Join(skipped, " "))
				} else {
					skipped = append(skipped, curArg.Describe(true, false))
				}
				curArg = ale.next()
			}
			return
		}
	} // ciclo di analisi argomenti

}
