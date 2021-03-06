package core

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"main/middlewares"
	"main/models"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
)

type Server struct {
	Echo        *echo.Echo
	Firebase    *firebase.App
	Auth        *auth.Client
	Config      *models.Config
	Firestore   *firestore.Client
	Middlewares struct {
		Auth *middlewares.Auth
	}
}

func NewServer() (svr *Server, err error) {
	e := echo.New()

	// Read configuration file

	config := &models.Config{}

	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	// Load Firebase config file

	file, err = ioutil.ReadFile(config.FirebaseConfig)

	if err != nil {
		return nil, err
	}

	firebaseConfig := &firebase.Config{}

	err = json.Unmarshal(file, firebaseConfig)

	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsFile(config.FirebaseConfig)

	//Init Firebase admin

	firebaseApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	auth, err := firebaseApp.Auth(context.Background())
	if err != nil {
		return nil, err
	}

	firestoreApp, err := firebaseApp.Firestore(context.Background())

	svr = &Server{
		Echo:      e,
		Firebase:  firebaseApp,
		Auth:      auth,
		Config:    config,
		Firestore: firestoreApp,
		Middlewares: struct{ Auth *middlewares.Auth }{
			Auth: middlewares.NewAuth(auth),
		},
	}
	return svr, nil
}

func (s *Server) Start() (err error) {
	return s.Echo.Start(s.Config.HostIP)
}
