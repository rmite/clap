// Copyright (c) 2017 Renato Mastrulli. Tutti i diritti riservati. All rights reserved.

package clap

import (
	"fmt"
	"sort"
	"strings"
)

/*
Command rappresenta un comando, ovvero il primo argomento nella linea di comando che determina l'azione da svolgere.

Per creare un oggetto Command usa la funzione NewCommand.

Il comando è composto da una parola (word) che deve apparire obbligatoriamente nella linea di comando.
Tranne se il comando è usato per default nella funzione Parse come spiegato nell'introduzione.

L'identificativo del comando (id) non è necessario, può essere impostato su una stringa vuota ma non può essere modificato dopo l'inizializzazione.

La preparazione del comando può essere fatta in anticipo o impostando la funzione del campo Prep.
Questa funzione è richiamata da Parse con il metodo Prepare quando il comando è individuato nella linea di comando.

In ogni caso, la sequenza da seguire è questa:

  argomenti [argomenti opzionali] [opzioni]

vedi l'introduzione per informazioni al riguardo.

*/
type Command struct {
	// id contiene l'identificativo del comando.
	id string
	// word restituisce la parola che rappresenta il comando.
	word string
	// Args contiene gli argomenti del comando.
	Args ArgumentList
	// Opts contiene le opzioni del comando.
	Opts OptionList
	// Data contiene i dati relativi al comando.
	Data interface{}
	// Prep contiene la funzione di preparazione del comando.
	Prep PrepareCommand
	//Exec contiene la funzione di esecuzione del comando.
	Exec ExecuteCommand
	//Help contiene la stringa che descrive il comando.
	Help string
}

// PrepareCommand è il tipo funzione per la preparazione di un comando.
type PrepareCommand func(cmd *Command)

// ExecuteCommand è il tipo funzione per l'esecuzione di un comando.
type ExecuteCommand func(data interface{})

// NewCommand crea un nuovo oggetto Command.
func NewCommand(cmdWord string, cmdID string, cmdHelp string) *Command {
	cmdWord = strings.TrimSpace(cmdWord)
	if len(cmdWord) == 0 {
		return nil
	}
	return &Command{id: cmdID, word: cmdWord, Help: cmdHelp}
}

// ID restituisce l'identificativo del comando.
func (cmd *Command) ID() string {
	return cmd.id
}

// Word restituisce la parola che rappresenta il comando.
func (cmd *Command) Word() string {
	return cmd.word
}

// Prepare è un metodo di comodo, chiama la funzione del campo Prep, se impostata, passando il comando stesso.
func (cmd *Command) Prepare() {
	if cmd.Prep != nil {
		cmd.Prep(cmd)
	}
}

// Execute è un metodo di comodo, chiama la funzione del campo Exec, se impostata, passando i dati (campo Data) del comando.
func (cmd *Command) Execute() {
	if cmd.Exec != nil {
		cmd.Exec(cmd.Data)
	}
}

// AppendArg aggiunge un argomento alla lista argomenti del comando.
func (cmd *Command) AppendArg(arg *Argument) {
	if arg != nil {
		if cmd.Args == nil {
			cmd.Args = make(ArgumentList, 0)
		}
		cmd.Args = append(cmd.Args, arg)
	}
}

// AppendNewArg crea un argomento nominale con la funzione NewArg, lo aggiunge alla lista argomenti del comando e lo restituisce.
func (cmd *Command) AppendNewArg(argWord string, argID string, argRequired bool, argValue interface{}, argHelp string) (arg *Argument) {
	arg = NewArg(argWord, argID, argRequired, argValue, argHelp)
	cmd.AppendArg(arg)
	return
}

// AppendOpt aggiunge un'opzione alla lista opzioni del comando.
func (cmd *Command) AppendOpt(opt *Option) {
	if opt != nil {
		if cmd.Opts == nil {
			cmd.Opts = make(OptionList, 0)
		}
		cmd.Opts = append(cmd.Opts, opt)
	}
}

// AppendNewOpt crea un'opzione con la funzione NewOption, la aggiunge alla lista opzioni del comando e la restituisce.
func (cmd *Command) AppendNewOpt(optWord string, optID string, optValue *bool, optHelp string) (opt *Option) {
	opt = NewOption(optWord, optID, optValue, optHelp)
	cmd.AppendOpt(opt)
	return
}

// ShortHelp restituisce l'help breve del comando, ovvero word e help su una stessa riga.
func (cmd *Command) ShortHelp() string {
	return fmt.Sprintf("%s		%s", cmd.word, cmd.Help)
}

// ShowShortHelp mostra l'help breve del comando nella console.
func (cmd *Command) ShowShortHelp() {
	fmt.Println(cmd.ShortHelp())
}

// FullHelp restituisce l'help esteso del comando, con sequenza argomenti e loro help.
func (cmd *Command) FullHelp() string {
	var hlp string
	var line string
	var arglist string
	var optlist string
	// prepara il comando
	cmd.Prepare()
	// usa l'id come titolo
	if len(cmd.id) > 0 {
		hlp = fmt.Sprintf("%s\n\n", cmd.id)
	}
	// comando
	hlp += cmd.word
	// argomenti
	line, arglist = cmd.Args.HelpStrings(0)
	if len(line) > 0 {
		hlp += " "
		hlp += line
	}
	// opzioni
	line, optlist = cmd.Opts.HelpStrings(0)
	if len(line) > 0 {
		hlp += " "
		hlp += line
	}
	// descrizione del comando
	hlp += fmt.Sprintf("\n\n%s\n", cmd.Help)
	// dettagli argomenti
	if len(arglist) > 0 {
		hlp += "\nARGOMENTI\n"
		hlp += arglist
	}
	// dettagli opzioni
	if len(optlist) > 0 {
		hlp += "\nOPZIONI\n"
		hlp += optlist
	}
	return hlp
}

// ShowFullHelp mostra l'help esteso del comando nella console.
func (cmd *Command) ShowFullHelp() {
	fmt.Println(cmd.FullHelp())
}

// =======================================================

// CommandMap raccoglie e organizza i comandi.
type CommandMap map[string]*Command

// NewCommandMap restituisce un oggetto CommandMap inizializzato con la lunghezza specificata.
func NewCommandMap(length int) CommandMap {
	if length < 1 {
		length = 0
	}
	return make(CommandMap, length)
}

// Insert inserisce un comando nella lista.
func (cm CommandMap) Insert(cmd *Command) {
	if cmd != nil {
		cm[cmd.word] = cmd
	}
}

// Remove rimuove un comando dalla lista.
func (cm CommandMap) Remove(cmd *Command) {
	if cmd != nil {
		delete(cm, cmd.word)
	}
}

// ShowHelp elenca i comandi nella console mostrando l'help accanto a ciascuno di essi.
func (cm CommandMap) ShowHelp() {
	// ordina le chiavi
	keys := make(sort.StringSlice, len(cm))
	i := 0
	for k := range cm {
		keys[i] = k
		i++
	}
	sort.Sort(keys)
	// mostra l'help breve per ogni comando
	var cmd *Command
	var ok bool
	for i = 0; i < len(keys); i++ {
		if cmd, ok = cm[keys[i]]; ok {
			cmd.ShowShortHelp()
		}
	}
}
