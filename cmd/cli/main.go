package main

import (
	"database/sql"
	"flag"
	"fmt"
	"runtime/debug"
	"strings"
	"time"
	"unsafe"

	"github.com/asaskevich/EventBus"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/timdevlet/todo/internal/cards"
	configs "github.com/timdevlet/todo/internal/configs"
	"github.com/timdevlet/todo/internal/dic"
	"github.com/timdevlet/todo/internal/eventbus"
	"github.com/timdevlet/todo/internal/helpers"
	"github.com/timdevlet/todo/internal/messenger"
	"github.com/timdevlet/todo/internal/notifications"
	"github.com/timdevlet/todo/internal/profile"
	"github.com/timdevlet/todo/internal/todo"
	"github.com/timdevlet/todo/pkg/postgres"

	. "github.com/dave/jennifer/jen"
)

//

type Hook struct {
}

// Add a hook to an instance of logger. This is called with
// `log.Hooks.Add(new(MyHook))` where `MyHook` implements the `Hook` interface.
func (hooks Hook) Add(hook Hook) {

}

// Fire all the hooks for the passed level. Used by `entry.log` to fire
// appropriate hooks for a log entry.
func (hooks Hook) Fire(level log.Level, entry *log.Entry) error {

	return nil
}

//

type dbHook struct {
	*sql.DB
}

func (db dbHook) Fire(e *logrus.Entry) error {
	strs := []string{}

	x, y, z := 0, 0, 0
	res := ""
	for i, c := range e.Message[:len(e.Message)-1] {
		if c == '[' {
			x = i
		}

		if c == ']' {
			z = i
			y = i
			strs = append(strs, e.Message[x:y])

			tag := e.Message[x+1 : y]

			if strings.Contains(tag, ":") {
				tagWithValue := strings.Split(tag, ":")
				e.Data[tagWithValue[0]] = tagWithValue[1]
			} else {
				e.Data["tag"] = tag
			}
		}
	}

	if z != 0 {
		res += strings.Trim(e.Message[z+1:], " ")
	}

	if res == "" {
		res = e.Message
	}

	e.Message = res

	return nil
}

func (dbHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func initLog(format string) {
	log.SetLevel(log.DebugLevel)

	log.SetFormatter(&log.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if format == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}

	h := dbHook{}
	log.AddHook(h)

	//

	// Dsn := "https://21189a97e1712a0986a779a1857fadda@o244099.ingest.sentry.io/4505647049080832"
	// hook, err := logrus_sentry.NewSentryHook(Dsn, []logrus.Level{
	// 	logrus.PanicLevel,
	// 	logrus.FatalLevel,
	// 	logrus.ErrorLevel,
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.AddHook(hook)
}

func main() {
	//

	// func (repo *DicRepository) FetchCities() ([]Citi, error) {
	// 	items := []Citi{}

	// 	result := repo.gorm.DB.Table("cities").Scan(&items)

	// 	if result.Error != nil {
	// 		return []Citi{}, result.Error
	// 	}

	// 	return items, nil
	// }

	// init struct with name and age

	inp := [][]string{[]string{"Cities", "Citi"}}

	f := NewFile("dic")

	for _, v := range inp {
		one := v[1]
		many := v[0]

		f.Func().
			Params(
				Id("repo").Op("*").Id("DicRepository"),
			).
			Id("Fetch"+one).Params().
			Params(
				Index().Id(one),
				Id("error"),
			).Block(
			Id("items").Op(":=").Index().Id(one).Values(),
			Id("result").Op(":=").Id("repo").Dot("gorm").Dot("DB").Dot("Table").Call(Lit(many)).Dot("Scan").Call(Op("&").Id("items")),

			If(Id("result.Error").Op("!=").Nil()).Block(
				Return(Index().Id(one).Values(), Id("result").Dot("Error")),
			),

			Return(Id("items"), Nil()),
		)

	}

	f.Save("../../internal/dic/dic.generated.go")

	return

	//

	// Init.
	logFlag := ""
	version := false
	action := ""

	flag.StringVar(&logFlag, "l", "plain", "log (plain|json)")
	flag.BoolVar(&version, "version", false, "version")
	flag.StringVar(&action, "action", "", "action")

	flag.Parse()

	//
	initLog(logFlag)

	log.WithField("f", "f").Error("asdasd 2")

	log.Error("[module:test][myvar:123] test error msg with tags")

	defer func() {
		if r := recover(); r != nil {
			log.Fatal("exception", string(debug.Stack()), r)
		}
	}()

	opt := configs.NewConfigsFromEnv()

	if version {
		log.Info("👋 Version: " + opt.VERSION)
		return
	}

	// Debug opt.
	optDebug := opt.GetFieldsWithValues()
	for key, value := range optDebug {
		log.Debug(key + ": " + value)
	}
	log.Debug("--------------------")

	log.Debug("[yyy] aasdasd[xxxx] ")

	addMessage(opt)
}

func addMessage(opt *configs.Configs) {
	p, err := postgres.NewPDB(postgres.Config{
		Host:     opt.DB_HOST,
		Port:     opt.DB_PORT,
		User:     opt.DB_USER,
		Password: opt.DB_PASSWORD,
		DBName:   opt.DB_NAME,
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal(err)
	}

	//

	gdb, err := postgres.NewGDB(postgres.Config{
		Host:     opt.DB_HOST,
		Port:     opt.DB_PORT,
		User:     opt.DB_USER,
		Password: opt.DB_PASSWORD,
		DBName:   opt.DB_NAME,
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal(err)
	}

	{
		err = gdb.DB.AutoMigrate(
			notifications.SMS{},
			cards.Card{},
			notifications.Email{},
		)
		if err != nil {
			log.Fatal(err)
		}

	}

	{
		//cardsService := cards.NewCardsService(cards.NewCardsRepository(gdb))

		cards.NewCardsRepository(gdb).Search()

	}

	sms2db := notifications.NewSms2DB(gdb)
	sms2db.SetBefore(func(s notifications.SMS) {
		log.WithField("text", s.Text).Info("sending sms...")
	})
	sms2db.SetAfter(func(s notifications.SMS) {
		log.WithField("text", s.Text).Info("sms send")
	})

	{

		eb := EventBus.New()

		eb.Subscribe("user_add", func(e string) {
			log.WithField("email", e).Info("user_add event")
		})

		event := eventbus.NewUpdateProfileEvent("")

		eb.Subscribe(event.Channel, func(e eventbus.ProfileEvent) {

			log.WithField("channel", event.Channel).Info("new event")

			sms := &notifications.SMS{
				Text: "sddsd",
				To:   12332,
			}

			err := sms2db.Set(sms).Send()

			if err != nil {
				log.Fatal(err)
			}

		})

		profileService := profile.NewProfileService(profile.NewProfileRepository(gdb), &eb)
		_, err := profileService.CreateRandomUser()

		if err != nil {
			log.Fatal(err)
		}

		iam := profile.IAM{
			UUID:  "asdasdasd",
			Email: "asdasd",
			Actions: []string{
				"users/add",
				"users/patch/fio",
			},
		}

		err = profileService.AddUser(iam, "adasdasd@gmail.com")

		if err != nil {
			log.Fatal(err)
		}

		action := &profile.ActionChangeName{
			IAM:  iam,
			UUID: "7a752619-ece6-4f89-9986-53181b0486bc",
			Name: "asdasd",
		}

		err = profileService.ActionChangeName(action)

		if err != nil {
			log.Fatal(err)
		}
	}

	return

	{
		dic := dic.NewDicService(dic.NewDicGRepository(gdb))
		err := dic.Sync()

		if err != nil {
			log.Fatal(err)
		}
	}

	return

	repo := messenger.NewMessageRepository(p)
	msg := messenger.NewMessengerService(repo)

	usersUuids := []string{}
	for i := 0; i < 300; i++ {
		usersUuids = append(usersUuids, helpers.UuidByHash(fmt.Sprintf("%v", i)))
	}

	userToUserMap := make(map[string]string)
	for _, u1 := range usersUuids {
		for _, u2 := range usersUuids {
			userToUserMap[u1] = u2
		}
	}

	fmt.Printf("userToUserMap: %T, %d\n", userToUserMap, unsafe.Sizeof(userToUserMap))

	i := 1000000000
	j := 0
	m := 0

	channelUuid1, _ := msg.SendMessage("9e3eefda-b56e-56bd-8a3a-0b8009d4a536", "9e3eefda-b56e-56bd-8a3a-0b8009d4a536", "12121")
	channelUuid2, _ := msg.SendMessage("9e3eefda-b56e-56bd-8a3a-0b8009d4a536", "9e3eefda-b56e-56bd-8a3a-0b8009d4a536", "12121")

	println(channelUuid1, channelUuid2)

	count, err := msg.CountChannelMessages("87c1cebd-27b8-4236-9e44-399ec3beba7f")
	if err != nil {
		log.Fatal(err)
	}

	log.Info(count)

	for k := 0; k < 100; k++ {
		go func() {
			for {
				ok, fromUuid := helpers.RandomFromSlice(usersUuids)
				if !ok {
					return
				}

				ok, toUuid := helpers.RandomFromSlice(usersUuids)
				if !ok {
					return
				}

				ts := []string{}
				for t := 0; t < 1; t++ {

					i--
					if i < 0 {
						return
					}

					ts = append(ts, helpers.FakeSentence(20))
				}

				if _, err = msg.SendMessage(
					fromUuid,
					toUuid,
					ts...,
				); err != nil {
					log.Fatal(err)
				}

				//println(toUuid, len(messages))
			}
		}()
	}

	// for k := 0; k < 10; k++ {
	// 	go func() {
	// 		for {
	// 			ok, fromUuid := helpers.RandomFromSlice(usersUuids)
	// 			if !ok {
	// 				return
	// 			}

	// 			messages, err := msg.FetchMessagesByUser(fromUuid)
	// 			if err != nil {
	// 				log.Fatal(err)
	// 			}

	// 			j++
	// 			m += len(messages)
	// 		}
	// 	}()
	// }

	{
		diff := 1
		go func() {
			for {
				progress := 1000000000 - i
				println("inserts", progress-diff)
				time.Sleep(1 * time.Second)
				diff = progress
			}
		}()
	}

	{
		diffJ := 1
		go func() {
			for {
				progress := j
				println("reads", progress-diffJ)
				time.Sleep(1 * time.Second)
				diffJ = progress
			}
		}()
	}

	{
		diff := 1
		go func() {
			for {
				progress := m
				println("messages fetched", progress-diff)
				time.Sleep(1 * time.Second)
				diff = progress
			}
		}()
	}

	time.Sleep(50 * time.Second)

	println(i)

	log.Debug("--------------------")
}

func addAction(opt *configs.Configs) {
	p, err := postgres.NewPDB(postgres.Config{
		Host:     opt.DB_HOST,
		Port:     opt.DB_PORT,
		User:     opt.DB_USER,
		Password: opt.DB_PASSWORD,
		DBName:   opt.DB_NAME,
		SSLMode:  "disable",
	})
	if err != nil {
		log.Fatal(err)
	}

	todoRepository := todo.NewTodoRepository(p)
	todoService := todo.NewTodoService(todoRepository)

	id, err := todoService.InsertTodo(todo.Todo{
		OwnerUuid: "dec8dc3a-52f8-44aa-b57b-d3275cd50fdc",
		Title:     "test",
	})
	if err != nil {
		log.Fatal(err)
	}

	todoService.GetUserTodos("dec8dc3a-52f8-44aa-b57b-d3275cd50fdc", todo.TodoFilter{})

	if err != nil {
		log.Fatal(err)
	} else {
		log.WithField("id", id).Debug("todo was inserted")
	}

	todoService.UpdateTitle("71fe2e7d-f421-4aca-b8ba-a24beac8108f", "????")

	todoService.DoneTodo("71fe2e7d-f421-4aca-b8ba-a24beac8108f", true)

	if err != nil {
		log.Fatal(err)
	}

	log.Debug("--------------------")
}
