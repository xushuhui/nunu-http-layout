// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"goal-advanced-layout/internal/conf"
	"goal-advanced-layout/internal/data"
	"goal-advanced-layout/internal/handler"
	"goal-advanced-layout/internal/server"
	"goal-advanced-layout/internal/service"
	"goal-advanced-layout/pkg/app"
	"goal-advanced-layout/pkg/helper/sid"
	"goal-advanced-layout/pkg/jwt"
	"goal-advanced-layout/pkg/log"
)

// Injectors from wire.go:

func NewWire(confServer *conf.Server, confData *conf.Data, logger *log.Logger) (*app.App, func(), error) {
	jwtJWT := jwt.NewJwt(confServer)
	handlerHandler := handler.NewHandler(logger)
	sidSid := sid.NewSid()
	serviceService := service.NewService(logger, sidSid, jwtJWT)
	db := data.NewDB(confData, logger)
	client := data.NewRedis(confData)
	dataData := data.NewData(db, client, logger)
	userRepo := data.NewUserRepo(dataData)
	userService := service.NewUserService(serviceService, userRepo)
	userHandler := handler.NewUserHandler(handlerHandler, userService)
	httpServer := server.NewHTTPServer(logger, confServer, jwtJWT, userHandler)
	job := server.NewJob(logger)
	appApp := newApp(httpServer, job)
	return appApp, func() {
	}, nil
}