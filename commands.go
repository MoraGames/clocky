package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/MoraGames/clockyuwu/events"
	"github.com/MoraGames/clockyuwu/internal/app"
	"github.com/MoraGames/clockyuwu/pkg/types"
	"github.com/MoraGames/clockyuwu/pkg/utils"
	"github.com/MoraGames/clockyuwu/structs"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Command struct {
	Name             string
	ShortDescription string
	LongDescription  string
	Category         string
	Syntax           string
	AdminOnly        bool
	Execute          func(msg tgbotapi.Message) error
}

var Commands map[string]Command

func init() {
	Commands = map[string]Command{
		"alias": {
			Name:             "alias",
			ShortDescription: "Mostra la lista di tutti gli schemi e dei loro vecchi nomi",
			LongDescription:  "Mostra la lista di tutti gli schemi disponibili in gioco, per ciascuno di essi indica (se disponibile) il vecchio nome che rispecchiava il pattern dello schema.",
			Category:         "di gioco",
			Syntax:           "/alias",
			AdminOnly:        false,
			Execute: func(msg tgbotapi.Message) error {
				if len(events.Sets) == 0 {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Nessuno schema trovato."), msg.MessageID)
					return nil
				}
				rawText := "Nome nuovo => %%Nome vecchio%%\n\n"
				for _, set := range events.Sets {
					rawText += fmt.Sprintf("%v => %%%%%v%%%%\n", set.Name, set.Pattern)
				}
				entities, text := utils.ParseToEntities(rawText)
				respMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				respMsg.Entities = entities
				sendMessage(respMsg, msg.MessageID)
				return nil
			},
		},
		"credits": {
			Name:             "credits",
			ShortDescription: "Mostra informazioni sul progetto e il suo autore",
			LongDescription:  "Fornisce informazioni sul bot e il suo scopo, link utili dell'autore e del progetto oltre ai ringraziamenti speciali.",
			Category:         "generali",
			Syntax:           "/credits",
			AdminOnly:        false,
			Execute: func(msg tgbotapi.Message) error {
				entities, text := utils.ParseToEntities(ComposeMessage(
					[]string{
						"%v è un giochino perdi-tempo interamente gestito da questo bot.\n",
						"Una volta che il bot è stato aggiunto ad un gruppo, il gioco consiste principalmente (ma non esclusivamente) nell'inviare messaggi nel formato \"hh:mm\" a determinati orari del giorno in cambio di punti. ",
						"La persona che totalizza più punti alla fine del campionato viene proclamata Clocky Champion!\n\n",
						"Il codice sorgente, disponibile su GitHub, è scritto interamente in GoLang e fa uso della libreria \"telegram-bot-api\". Per ogni suggerimento o problema, riferisciti al progetto GitHub.\n\n",
						"- Telegram: MoraGames\n- Discord: moragames - Instagram: moragames.dev\n\n",
						"Un ringraziamento speciale ai primi beta tester (e giocatori) del minigioco, \"Vano\", \"Ale\" e \"Alex\".",
					},
					app.Name,
				))
				respMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				respMsg.Entities = entities
				sendMessage(respMsg, msg.MessageID)
				return nil
			},
		},
		"file": {
			Name:             "file",
			ShortDescription: "Gestisce le operazioni sui file di memorizzazione dati",
			LongDescription:  "Permette di ottenere, aggiornare o cancellare i file di gioco memorizzati dal bot, utili per backup o migrazione dati. In caso di aggiornamento, al messaggio del comando deve essere fornito un file in allegato. Se non si specifica un'operazione, verrà effettuata un'operazione di ottenimento per default.",
			Category:         "per admin",
			Syntax:           "/file [\"get\"|\"upd\"|\"del\"] <file_name> [\"overwrite\"]",
			AdminOnly:        true,
			Execute: func(msg tgbotapi.Message) error {
				go func(msg tgbotapi.Message) {
					if !isAdmin(msg.From.ID) {
						sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Non hai i permessi per eseguire questo comando."), msg.MessageID)
						logOutcome("file", fmt.Errorf("unauthorized user"))
						return
					}
					var args []string
					if cmdArgs := types.CommandArguments(&msg); cmdArgs != "" {
						args = strings.Split(cmdArgs, " ")
					}
					if len(args) < 1 || len(args) > 3 {
						sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Uso non valido del comando. Usa /help per maggiori informazioni."), msg.MessageID)
						logOutcome("file", fmt.Errorf("wrong number of arguments"))
						return
					}
					if (len(args) >= 2 && args[0] != "get" && args[0] != "upd" && args[0] != "del") || (len(args) == 3 && args[2] != "overwrite") || (len(args) == 3 && args[0] != "upd" && args[2] == "overwrite") {
						sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Argomento sconosciuto. Usa /help per maggiori informazioni."), msg.MessageID)
						logOutcome("file", fmt.Errorf("unknown argument value"))
						return
					}

					var operation, fileName string
					var overwrite bool
					switch len(args) {
					case 1:
						operation = "get"
						fileName = args[0]
						overwrite = false
					case 2:
						operation = args[0]
						fileName = args[1]
						overwrite = false
					case 3:
						operation = args[0]
						fileName = args[1]
						overwrite = true
					}
					if fileName != "sets.json" && fileName != "events.json" && fileName != "users.json" && fileName != "pinnedMessage.json" && fileName != "hints.json" && fileName != "championship.json" && fileName != "pinnedChampionshipMessage.json" && fileName != "logs/log.json" {
						sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "File non valido per l'operazione."), msg.MessageID)
						logOutcome("file", fmt.Errorf("invalid file"))
						return
					}

					switch operation {
					case "get":
						file, err := os.ReadFile("files/" + fileName)
						if err != nil {
							sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "File non trovato."), msg.MessageID)
							logOutcome("file", fmt.Errorf("file not found"))
							return
						}
						respMsg := tgbotapi.NewDocument(msg.Chat.ID, tgbotapi.FileBytes{Name: fileName, Bytes: file})
						respMsg.Caption = fileName + " recuperato con successo."
						sendDocument(respMsg, msg.MessageID)
					case "upd":
						if msg.Document == nil {
							sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Nessun file allegato al messaggio."), msg.MessageID)
							logOutcome("file", fmt.Errorf("attachment not found"))
							return
						}
						var err error
						switch fileName {
						case "sets.json":
							err = updateData(msg, "files/"+fileName, &events.Sets, events.AssignSetsFromSetsJson, overwrite)
						case "events.json":
							err = updateData(msg, "files/"+fileName, &events.Events, nil, overwrite)
						case "users.json":
							err = updateData(msg, "files/"+fileName, &Users, nil, overwrite)
						case "pinnedMessage.json":
							err = updateData(msg, "files/"+fileName, &events.PinnedResetMessage, nil, overwrite)
						case "hints.json":
							err = updateData(msg, "files/"+fileName, &events.HintRewardedUsers, nil, overwrite)
						case "championship.json":
							err = updateData(msg, "files/"+fileName, &events.CurrentChampionship, UpdateChampionshipResetCronjob, overwrite)
						case "pinnedChampionshipMessage.json":
							err = updateData(msg, "files/"+fileName, &structs.PinnedChampionshipResetMessage, nil, overwrite)
						default:
							sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "File non valido per l'operazione."), msg.MessageID)
							logOutcome("file", fmt.Errorf("invalid file"))
							return
						}
						logOutcome("file", err)
					case "del":
						err := os.Remove("files/" + fileName)
						if err != nil {
							sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "File non trovato."), msg.MessageID)
							logOutcome("file", fmt.Errorf("file not found"))
							return
						}
						sendMessage(tgbotapi.NewMessage(msg.Chat.ID, fileName+" cancellato con successo."), msg.MessageID)
					}
				}(msg)
				return nil
			},
		},
		"help": {
			Name:             "help",
			ShortDescription: "Mostra la lista dei comandi disponibili",
			LongDescription:  "Fornisce una lista di tutti i comandi disponibili per il bot accompagnati da una breve descrizione. Usando /help <command_name> si ottengono maggiori informazioni su un comando specifico.",
			Category:         "generali",
			Syntax:           "/help [command_name]",
			AdminOnly:        false,
			Execute: func(msg tgbotapi.Message) error {
				var entities []tgbotapi.MessageEntity
				var text string
				var args []string
				if cmdArgs := msg.CommandArguments(); cmdArgs != "" {
					args = strings.Split(cmdArgs, " ")
				}

				if len(args) > 1 {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Uso non valido del comando. Usa /help o /help <command_name>."), msg.MessageID)
					return fmt.Errorf("wrong number of arguments")
				}
				if len(args) == 0 {
					categoriezedCommands := slicefyCommands()
					rawText := fmt.Sprintf("Name: %v | Version: %v\n\nQuesta è una lista di tutti i comandi disponibili per il bot.\nPer maggiori informazioni usa /help <command_name>.\n\n", app.Name, app.Version)
					for _, cat := range categoriezedCommands {
						rawText += fmt.Sprintf("**Comandi %v:**\n", cat[0].Category)
						for _, cmd := range cat {
							rawText += fmt.Sprintf("| /%v - %%%%%v%%%%\n", cmd.Name, cmd.ShortDescription)
						}
						rawText += "\n"
					}
					entities, text = utils.ParseToEntities(rawText)
					respMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
					respMsg.Entities = entities
					sendMessage(respMsg, msg.MessageID)
					return nil
				}
				if cmd, exists := Commands[args[0]]; !exists {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Argomento sconosciuto. Usa /help per maggiori informazioni."), msg.MessageID)
					return fmt.Errorf("unknown argument")
				} else {
					entities, text := utils.ParseToEntities(fmt.Sprintf("**/%v - %v**\n\nDescrizione: %%%%%v%%%%\n\nCategoria: %%%%%v%%%%\nSintassi: %%%%%v%%%%\n", cmd.Name, cmd.ShortDescription, cmd.LongDescription, cmd.Category, cmd.Syntax))
					respMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
					respMsg.Entities = entities
					sendMessage(respMsg, msg.MessageID)
					return nil
				}
			},
		},
		"list": {
			Name:             "list",
			ShortDescription: "Mostra la lista degli schemi o degli effetti attualmente attivi",
			LongDescription:  "A seconda del parametro, \"sets\" o \"effects\", viene mostrata la lista corrispondente. Se non viene specificato alcun argomento, entrambe le liste verranno mostrate.",
			Category:         "di gioco",
			Syntax:           "/list [\"sets\"|\"effects\"]",
			AdminOnly:        false,
			Execute: func(msg tgbotapi.Message) error {
				var args []string
				if cmdArgs := msg.CommandArguments(); cmdArgs != "" {
					args = strings.Split(cmdArgs, " ")
				}
				if len(args) > 1 {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Uso non valido del comando. Usa /help per maggiori informazioni."), msg.MessageID)
					return fmt.Errorf("wrong number of arguments")
				}
				if len(args) == 1 && args[0] != "sets" && args[0] != "effects" {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Argomento sconosciuto. Usa /help per maggiori informazioni."), msg.MessageID)
					return fmt.Errorf("unknown argument")
				}
				var rawText string
				if len(args) == 0 || args[0] == "sets" {
					sets := events.Events.Stats.EnabledSets
					sort.Slice(sets, func(i, j int) bool {
						return sets[i] < sets[j]
					})

					rawText += fmt.Sprintf("**Schemi attivi (%v):**\n", len(sets))
					for _, setName := range sets {
						rawText += fmt.Sprintf(" | %q\n", setName)
					}
					rawText += "\n"
				}
				if len(args) == 0 || args[0] == "effects" {
					effects := slicefyEffects()

					rawText += fmt.Sprintf("**Effetti presenti (%v):**\n", len(effects))
					for _, effect := range effects {
						rawText += fmt.Sprintf(" | %q = %v\n", effect.Name, effect.Amount)
					}
					rawText += "\n"
				}
				entities, text := utils.ParseToEntities(rawText)
				respMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				respMsg.Entities = entities
				sendMessage(respMsg, msg.MessageID)
				return nil
			},
		},
		"ping": {
			Name:             "ping",
			ShortDescription: "Verifica se il bot è online",
			LongDescription:  "Qualora il sistema sia attivo e connesso, il bot risponderà con 'Pong!' confermando la sua operatività.",
			Category:         "generali",
			Syntax:           "/ping",
			AdminOnly:        false,
			Execute: func(msg tgbotapi.Message) error {
				sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Pong!"), msg.MessageID)
				return nil
			},
		},
		"ranking": {
			Name:             "ranking",
			ShortDescription: "Mostra la classifica dei giocatori partecipanti",
			LongDescription:  "Mostra la classifica dei giocatori partecipanti al campionato in corso. È possibile specificare l'ambito della classifica (giornaliero, del campionato o totale) e visualizzare la classifica dal punto di vista di un altro utente.",
			Category:         "di gioco",
			Syntax:           "/ranking [\"day\"|\"championship\"|\"total\"] [\"pov\" <username>]",
			AdminOnly:        false,
			Execute: func(msg tgbotapi.Message) error {
				var ranking []structs.Rank
				var povTelegramUserID int64
				var args []string
				if cmdArgs := msg.CommandArguments(); cmdArgs != "" {
					args = strings.Split(cmdArgs, " ")
				}
				if len(args) == 0 {
					ranking = structs.GetRanking(Users, structs.DefaultRankScope, true)
					povTelegramUserID = msg.From.ID
				} else {
					if len(args) > 3 {
						sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Uso non valido del comando. Usa /help per maggiori informazioni."), msg.MessageID)
						return fmt.Errorf("wrong number of arguments")
					} else if (len(args) == 1 && args[0] != string(structs.RankScopeDay) && args[0] != string(structs.RankScopeChampionship) && args[0] != string(structs.RankScopeTotal)) || (len(args) == 2 && args[0] != "pov") || (len(args) == 3 && ((args[0] != string(structs.RankScopeDay) && args[0] != string(structs.RankScopeChampionship) && args[0] != string(structs.RankScopeTotal)) || args[1] != "pov")) {
						sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Argomento sconosciuto. Usa /help per maggiori informazioni."), msg.MessageID)
						return fmt.Errorf("unknown argument")
					}
					if len(args) == 1 {
						switch args[0] {
						case string(structs.RankScopeDay):
							ranking = structs.GetRanking(Users, structs.RankScopeDay, true)
						case string(structs.RankScopeChampionship):
							ranking = structs.GetRanking(Users, structs.RankScopeChampionship, true)
						case string(structs.RankScopeTotal):
							ranking = structs.GetRanking(Users, structs.RankScopeTotal, false)
						}
						povTelegramUserID = msg.From.ID
					} else if len(args) == 2 {
						username := args[1]
						var userId int64
						var founded bool
						for userID, user := range Users {
							if user.UserName == username {
								founded = true
								userId = userID
								break
							}
						}
						if !founded {
							sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Utente non trovato."), msg.MessageID)
							return fmt.Errorf("unknown argument value")
						} else {
							ranking = structs.GetRanking(Users, structs.DefaultRankScope, true)
							povTelegramUserID = userId
						}
					} else if len(args) == 3 {
						username := args[2]
						var userId int64
						var founded bool
						for userID, user := range Users {
							if user.UserName == username {
								founded = true
								userId = userID
								break
							}
						}
						if !founded {
							sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Utente non trovato."), msg.MessageID)
							return fmt.Errorf("unknown argument value")
						} else {
							switch args[0] {
							case string(structs.RankScopeDay):
								ranking = structs.GetRanking(Users, structs.RankScopeDay, true)
							case string(structs.RankScopeChampionship):
								ranking = structs.GetRanking(Users, structs.RankScopeChampionship, true)
							case string(structs.RankScopeTotal):
								ranking = structs.GetRanking(Users, structs.RankScopeTotal, false)
							}
							povTelegramUserID = userId
						}
					}
				}

				// Generate the string to send
				rawText := "Ancora nessun utente ha partecipato agli eventi del campionato."
				if len(ranking) != 0 {
					var povPoints int
					for _, r := range ranking {
						if r.UserTelegramID == povTelegramUserID {
							povPoints = r.Points
							break
						}
					}
					rankingString := ""
					for i, r := range ranking {
						rankingString += fmt.Sprintf("**%v] %v:** %v %%%%(%+d)%%%%\n", i+1, r.Username, r.Points, r.Points-povPoints)
					}

					// Send the message
					rawText = fmt.Sprintf("__**La classifica è la seguente:**__\n\n%v", rankingString)
				}
				entities, text := utils.ParseToEntities(rawText)
				respMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				respMsg.Entities = entities
				sendMessage(respMsg, msg.MessageID)
				return nil
			},
		},
		"execute": {
			Name:             "execute",
			ShortDescription: "Esegue le funzioni per rigenerare i dati di gioco",
			LongDescription:  "Esegue la rigenerazione dei dati di gioco, sovrascrivendo i file attuali con nuovi dati calcolati. È possibile rigenerarli tutti oppure specificare singolarmente quali dati ricalcolare. Questa operazione è irreversibile.",
			Category:         "per admin",
			Syntax:           "/execute <\"all\"|\"events\"|\"championship\"> [\"all\"|\"reset\"|\"rewards\"]",
			AdminOnly:        true,
			Execute: func(msg tgbotapi.Message) error {
				if !isAdmin(msg.From.ID) {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Non hai i permessi per eseguire questo comando."), msg.MessageID)
					return fmt.Errorf("unauthorized user")
				}
				var args []string
				if cmdArgs := msg.CommandArguments(); cmdArgs != "" {
					args = strings.Split(cmdArgs, " ")
				}
				if len(args) < 1 || len(args) > 2 {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Uso non valido del comando. Usa /help per maggiori informazioni."), msg.MessageID)
					return fmt.Errorf("wrong number of arguments")
				}
				if (args[0] != "all" && args[0] != "events" && args[0] != "championship") || (len(args) == 2 && args[1] != "all" && args[1] != "reset" && args[1] != "rewards") {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Argomento sconosciuto. Usa /help per maggiori informazioni."), msg.MessageID)
					return fmt.Errorf("unknown argument value")
				}
				if len(args) == 1 {
					args = append(args, "all")
				}

				selectedCategories := make(map[string]bool)
				selectedActions := make(map[string]bool)
				switch args[0] {
				case "all":
					selectedCategories["championship"] = true
					selectedCategories["events"] = true
				case "championship", "events":
					selectedCategories[args[0]] = true
				}
				switch args[1] {
				case "all":
					selectedActions["reset"] = true
					selectedActions["rewards"] = true
				case "reset", "rewards":
					selectedActions[args[1]] = true
				}

				categories := []string{"championship", "events"}
				actions := []string{"reset", "rewards"}
				dailyEnabledEvents := events.Events.Stats.EnabledEventsNum
				for _, category := range categories {
					if _, exist := selectedCategories[category]; !exist {
						continue
					}
					for _, action := range actions {
						if _, exist := selectedActions[action]; !exist {
							continue
						}

						switch category {
						case "championship":
							switch action {
							case "reset":
								events.CurrentChampionship.Reset(
									structs.GetRanking(Users, structs.RankScopeChampionship, true),
									&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
									types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
								)
							case "rewards":
								ChampionshipUserRewardAndReset(
									Users,
									&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
									types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
								)
							}
						case "events":
							switch action {
							case "reset":
								events.Events.Reset(
									true,
									&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
									types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
								)
							case "rewards":
								DailyUserRewardAndReset(
									Users,
									dailyEnabledEvents,
									&types.WriteMessageData{Bot: App.BotAPI, ChatID: App.DefaultChatID, ReplyMessageID: -1},
									types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat},
								)
							}
						}
					}
				}

				return nil
			},
		},
		"stats": {
			Name:             "stats",
			ShortDescription: "Mostra le statistiche dei giocatori",
			LongDescription:  "Fornisce una panoramica delle statistiche di gioco dell'utente indicato in forma classica oppure in forma estesa. Se non viene specificato alcun utente, verranno mostrate le statistiche classiche dell'utente che ha inviato il comando.",
			Category:         "di gioco",
			Syntax:           "/stats [username [\"expand\"]]",
			AdminOnly:        false,
			Execute: func(msg tgbotapi.Message) error {
				var args []string
				if cmdArgs := msg.CommandArguments(); cmdArgs != "" {
					args = strings.Split(cmdArgs, " ")
				}
				if len(args) > 2 {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Uso non valido del comando. Usa /help per maggiori informazioni."), msg.MessageID)
					return fmt.Errorf("wrong number of arguments")
				}
				if len(args) == 2 && args[1] != "expand" {
					sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Argomento sconosciuto. Usa /help per maggiori informazioni."), msg.MessageID)
					return fmt.Errorf("unknown argument value")
				}

				var rawText string
				var user *structs.User
				if len(args) > 0 {
					username := args[0]
					var userId int64
					var founded bool
					for uID, u := range Users {
						if u.UserName == username {
							founded = true
							userId = uID
							break
						}
					}
					if !founded {
						sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Utente non trovato."), msg.MessageID)
						return fmt.Errorf("unknown argument value")
					}
					user = Users[userId]
					rawText = fmt.Sprintf("__**Le statistiche di %v sono:**__\n\n", user.UserName)
				} else {
					user = Users[msg.From.ID]
					rawText = "__**Le tue statistiche sono:**__\n\n"
				}

				if len(args) == 2 {
					rawText += ComposeMessage(
						[]string{
							"**Statistiche di oggi:**\n",
							"- Punti: %v\n- Partecipazioni: %v\n- Vittorie: %v\n- Sconfitte: %v\n\n",
							"- Punti/Partecipazioni: %.2f\n- Punti/Vittorie: %.2f\n- Vittorie/Partecipazioni: %.2f\n\n",
							"**Statistiche del campionato:**\n",
							"- Punti: %v\n- Partecipazioni: %v\n- Vittorie: %v\n- Sconfitte: %v\n\n",
							"- Punti/Partecipazioni: %.2f\n- Punti/Vittorie: %.2f\n- Vittorie/Partecipazioni: %.2f\n\n",
							"**Statistiche di sempre:**\n",
							"- Punti: %v\n- Partecipazioni: %v\n- Vittorie: %v\n- Sconfitte: %v\n\n",
							"- Punti/Partecipazioni: %.2f\n- Punti/Vittorie: %.2f\n- Vittorie/Partecipazioni: %.2f\n\n",
							"- Campionati svolti: %v\n- Campionati vinti: %v\n",
							"- Streak partecipazioni: %v\n- Streak attività: %v\n\n",
							"**Effetti attivi:**\n",
							"- %v",
						},
						user.DailyPoints, user.DailyEventPartecipations, user.DailyEventWins, user.DailyEventPartecipations-user.DailyEventWins,
						float64(user.DailyPoints)/float64(user.DailyEventPartecipations), float64(user.DailyPoints)/float64(user.DailyEventWins), float64(user.DailyEventWins)/float64(user.DailyEventPartecipations),
						user.ChampionshipPoints, user.ChampionshipEventPartecipations, user.ChampionshipEventWins, user.ChampionshipEventPartecipations-user.ChampionshipEventWins,
						float64(user.ChampionshipPoints)/float64(user.ChampionshipEventPartecipations), float64(user.ChampionshipPoints)/float64(user.ChampionshipEventWins), float64(user.ChampionshipEventWins)/float64(user.ChampionshipEventPartecipations),
						user.TotalPoints, user.TotalEventPartecipations, user.TotalEventWins, user.TotalEventPartecipations-user.TotalEventWins,
						float64(user.TotalPoints)/float64(user.TotalEventPartecipations), float64(user.TotalPoints)/float64(user.TotalEventWins), float64(user.TotalEventWins)/float64(user.TotalEventPartecipations),
						user.TotalChampionshipPartecipations, user.TotalChampionshipWins,
						user.DailyPartecipationStreak, user.DailyActivityStreak,
						user.StringifyEffects(false),
					)
				} else {
					rawText += ComposeMessage(
						[]string{
							"**Statistiche di oggi:**\n",
							"- Punti: %v\n- ~~Partecipazioni: %v\n- Vittorie~~: %v\n\n",
							"- ~~Punti/Vittorie~~: %.2f\n- Vittorie/Partecipazioni: %.2f\n\n",
							"**Statistiche del campionato:**\n",
							"- Punti: %v\n- Partecipazioni: %v\n- Vittorie: %v\n\n",
							"- Punti/Vittorie: %.2f\n- Vittorie/Partecipazioni: %.2f\n\n",
							"**Statistiche di sempre:**\n",
							"- Streak partecipazioni: %v\n- Streak attività: %v\n\n",
						},
						user.DailyPoints, user.DailyEventPartecipations, user.DailyEventWins,
						float64(user.DailyPoints)/float64(user.DailyEventWins), float64(user.DailyEventWins)/float64(user.DailyEventPartecipations),
						user.ChampionshipPoints, user.ChampionshipEventPartecipations, user.ChampionshipEventWins,
						float64(user.ChampionshipPoints)/float64(user.ChampionshipEventWins), float64(user.ChampionshipEventWins)/float64(user.ChampionshipEventPartecipations),
						user.DailyPartecipationStreak, user.DailyActivityStreak,
					)
				}

				entities, text := utils.ParseToEntities(rawText)
				respMsg := tgbotapi.NewMessage(msg.Chat.ID, text)
				respMsg.Entities = entities
				sendMessage(respMsg, msg.MessageID)
				return nil
			},
		},
	}
}

func manageCommands(update tgbotapi.Update) {
	// fmt.Printf("\n>>> DEBUG <<<\n |- %q\n |- %q (%v)\n |- %v (%v)\n\n", types.Command(update.Message), types.CommandArguments(update.Message), utf8.RuneCountInString(types.CommandArguments(update.Message)), strings.Split(types.CommandArguments(update.Message), " "), len(strings.Split(types.CommandArguments(update.Message), " ")))
	if cmd, exists := Commands[types.Command(update.Message)]; exists {
		err := cmd.Execute(*update.Message)
		// Since /file it's executed in a goroutine, is the only command that manage it own outcome logging. The check avoid double logging.
		if cmd.Name != "file" {
			logOutcome(cmd.Name, err)
		}
	} else {
		sendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Comando sconosciuto. Usa /help per vedere la lista dei comandi disponibili."), update.Message.MessageID)
	}
}

func logOutcome(cmdName string, err error) {
	if err != nil {
		App.Logger.WithFields(logrus.Fields{
			"cmd": cmdName,
			"err": err.Error(),
		}).Error("Error while executing command")
	} else if cmdName != "file" {
		App.Logger.WithField("cmd", cmdName).Debug("Command executed successfully")
	}
}

func sendMessage(msg tgbotapi.MessageConfig, replyTo int) {
	msg.ReplyToMessageID = replyTo
	_, err := App.BotAPI.Send(msg)
	if err != nil {
		App.Logger.WithField("err", err).Error("Error while sending message")
	}
}

func sendDocument(msg tgbotapi.DocumentConfig, replyTo int) {
	msg.ReplyToMessageID = replyTo
	_, err := App.BotAPI.Send(msg)
	if err != nil {
		App.Logger.WithField("err", err).Error("Error while sending message")
	}
}

func isAdmin(userID int64) bool {
	env, exists := os.LookupEnv("TELEGRAM_ADMIN_ID")
	if !exists {
		App.Logger.WithField("env", "TELEGRAM_ADMIN_ID").Warn("Environment variable not set")
	}
	for _, idStr := range strings.Split(env, ", ") {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil && id == userID {
			return true
		}
	}
	return false
}

func updateData(msg tgbotapi.Message, filePath string, dataStructure any, ifOkay func(utils types.Utils), overwrite bool) error {
	file, err := App.BotAPI.GetFile(tgbotapi.FileConfig{FileID: msg.Document.FileID})
	if err != nil {
		sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Errore nel recupero del file allegato."), msg.MessageID)
		return fmt.Errorf("error retrieving attachment")
	}

	resp, err := http.Get("https://api.telegram.org/file/bot" + App.BotAPI.Token + "/" + file.FilePath)
	if err != nil {
		sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Errore nel download del file allegato."), msg.MessageID)
		return fmt.Errorf("error downloading attachment: %v", err)
	}
	defer resp.Body.Close()

	// Read the entire response body (file content)
	downloadedFile, err := io.ReadAll(resp.Body)
	if err != nil {
		sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Errore nella lettura del file scaricato."), msg.MessageID)
		return fmt.Errorf("error reading downloaded file: %v", err)
	}

	err = json.Unmarshal(downloadedFile, dataStructure)
	if err != nil {
		sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Errore nella conversione del file allegato."), msg.MessageID)
		return fmt.Errorf("error unmarshalling attachment")
	}

	if ifOkay != nil {
		ifOkay(types.Utils{Config: App.Config, Logger: App.Logger, TimeFormat: App.TimeFormat})
	}
	sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Struttura dati sovrascritta con successo."), msg.MessageID)

	if overwrite {
		fileToWrite, err := json.MarshalIndent(dataStructure, "", " ")
		if err != nil {
			sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Errore nella conversione del file allegato."), msg.MessageID)
			return fmt.Errorf("error marshalling attachment")
		}
		err = os.WriteFile(filePath, fileToWrite, 0644)
		if err != nil {
			sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "Errore nella scrittura del file allegato."), msg.MessageID)
			return fmt.Errorf("error writing file")
		}
		sendMessage(tgbotapi.NewMessage(msg.Chat.ID, "File dati sovrascritto dall'allegato con successo."), msg.MessageID)
	}
	return nil
}

func slicefyCommands() [][]Command {
	// Group by category
	byCategory := make(map[string][]Command)
	for _, cmd := range Commands {
		byCategory[cmd.Category] = append(byCategory[cmd.Category], cmd)
	}

	// Sort the categories between them
	categories := make([]string, 0, len(byCategory))
	for c := range byCategory {
		categories = append(categories, c)
	}
	sort.Strings(categories)

	// Build output, sorting each category
	out := make([][]Command, 0, len(categories))
	for _, cat := range categories {
		cmds := byCategory[cat]
		sort.Slice(cmds, func(i, j int) bool {
			return cmds[i].Name < cmds[j].Name
		})
		out = append(out, cmds)
	}

	return out
}

type EffectPresence struct {
	Name   string
	Amount int
}

func slicefyEffects() []EffectPresence {
	out := make([]EffectPresence, 0, len(events.Events.Stats.EnabledEffects))
	for effectName, effectNum := range events.Events.Stats.EnabledEffects {
		out = append(out, EffectPresence{
			Name:   effectName,
			Amount: effectNum,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func ComposeMessage(subMessages []string, args ...any) string {
	msg := ""
	for _, subMessage := range subMessages {
		msg += subMessage
	}
	return fmt.Sprintf(msg, args...)
}
