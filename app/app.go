package app

import (
	"appdoki-be/app/repositories"
	"appdoki-be/config"
	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Application struct {
	conf            *config.Config
	firebaseApp     *firebase.App
	usersRepository repositories.UsersRepositoryInterface
	beersRepository repositories.BeersRepositoryInterface
	notifier        notifier
}

func NewApplication(conf *config.Config, db *sqlx.DB, firebaseApp *firebase.App) *Application {
	notifierSrv, err := newNotifier(firebaseApp, conf.AppConfig.NotifierTestMode)
	if err != nil {
		log.Fatal("could not instantiate a notifier")
	}

	return &Application{
		conf:            conf,
		firebaseApp:     firebaseApp,
		usersRepository: repositories.NewUsersRepository(db),
		beersRepository: repositories.NewBeersRepository(db),
		notifier:        notifierSrv,
	}
}

func (a *Application) Routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)

	a.AuthRouter(router)
	a.UsersRouter(router)
	a.BeersRouter(router)

	fs := http.FileServer(http.Dir("./swaggerui/"))
	router.
		PathPrefix("/docs").
		Handler(http.StripPrefix("/docs", fs))

	return middlewareChain([]middleware{
		trimSuffixMiddleware,
		loggingMiddleware,
	}, router)
}

type TopicInfo struct {
	Topic       string
	Description string
}
type HomeResponse struct {
	Version         string
	DocsEndpoint    string
	MessagingTopics []TopicInfo
}

func homeHandler(w http.ResponseWriter, _ *http.Request) {
	res := HomeResponse{
		Version:      "1.0.0",
		DocsEndpoint: "/docs/",
		MessagingTopics: []TopicInfo{
			{
				Topic:       "beers",
				Description: "Global beer transfer notifications",
			},
			{
				Topic:       "users",
				Description: "Global notification for joined users",
			},
		},
	}

	respondJSON(w, res, http.StatusOK)
}
