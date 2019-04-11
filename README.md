
## Package clap
### Libreria go per gestire gli argomenti della linea di comando

[![Go Report Card](https://goreportcard.com/badge/github.com/rmite/clap)](https://goreportcard.com/report/github.com/rmite/clap)  [![GoDoc](https://godoc.org/github.com/rmite/clap?status.svg)](https://godoc.org/github.com/rmite/clap)

#### Descrizione

Il package clap ti permette di **acquisire gli argomenti della linea di comando** dalla console definendo comandi per ogni azione della tua applicazione.

L'analisi degli argomenti con il package clap individua il comando in base ad una parola iniziale e acquisisce **i valori della linea di comando** per passarli ad **una funzione di elaborazione personale**.

La **sequenza** è questa:

```
comando argomenti [argomenti opzionali] [opzioni]
```

dove "comando" è la parola identificativa del comando, "argomenti" sono gli argomenti obbligatori a cui seguono quelli opzionali e poi le opzioni.

Gli argomenti catturano un valore. Possono essere nominali oppure no. Gli **argomenti nominali** sono identificati da una parola con la sintassi:

```
parola=valore
```

mentre gli **argomenti generici** prevedono solo il valore:

```
valore
```

Le **opzioni** invece sono sempre identificate da una parola a cui può seguire un valore booleano:

```
parola[=bool]
```

se non specificato, l'opzione acquisisce il valore default MatchValue impostato nella sua definizione.

Un'opzione può avere anche una serie di argomenti che devono essere indicati dopo essa nella riga di comando.

I valori acquisiti sono passati **in puntatori a variabili** impostati nella definizione di argomenti e opzioni.

Il package inoltre permette di:
- usare un comando default, cioè riconosciuto anche senza la parola identificativa
- definire argomenti e opzioni dei comandi sia in anticipo sia in maniera dinamica, cioè quando il comando deve essere eseguito
- evidenziare la mancanza di argomenti necessari
- mostrare la sintassi dei comandi
- mostrare un help generale

Per ulteriori informazioni, consulta la [documentazione](https://godoc.org/github.com/rmite/clap) del package.

#### Esempio

Definizione di comando e argomenti prima di chiamare Parse

```
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
```

e in maniera dinamica

```
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
```

Avendo specificato cmdFix come comando default nella chiamata a Parse, il comando può essere riconosciuto sia con la sua parola -fix

  -fix source.txt dest.txt -enc=UTF8 -a

sia senza, ovvero la linea di comando contiene direttamente i suoi argomenti

  source.txt dest.txt -enc=UTF8 -a

questo è l'unico modo possibile se cmdFix non fosse stato aggiunto al map commands.

Se cmdFix non fosse specificato come default, per essere riconosciuto dovrebbe essere aggiunto al map commands e la sua parola sarebbe obbligatoria.
