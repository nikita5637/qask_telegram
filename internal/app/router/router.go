package router

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"qask_telegram/internal/app/model"
	"strings"
)

type RouterHandler func(*model.User, *tgbotapi.Update)

type Route struct {
	isPublic     bool
	routeHandler RouterHandler
}

type Router struct {
	r      map[string]*Route
	logger *logrus.Logger
}

func NewRouter(logger *logrus.Logger) *Router {
	return &Router{
		r:      make(map[string]*Route),
		logger: logger,
	}
}

func (r *Router) NewRoute(path string, isPublic bool, handler RouterHandler) *Router {
	r.logger.Debugf("Creating new route '%s'", path)
	if path == "" {
		return nil
	}

	newRoute := &Route{
		isPublic:     isPublic,
		routeHandler: handler,
	}

	r.r[path] = newRoute
	return r
}

func (r *Router) GetHandler(path string) RouterHandler {
	r.logger.Debugf("Looking for handler '%s'", path)

	route, ok := r.r[path]
	if !ok {
		return nil
	}

	return route.routeHandler
}

func (r *Router) CommandIsRegistered(command string) bool {
	tmpCommand := strings.TrimRight(command, " ")
	_, ok := r.r[tmpCommand]
	return ok
}

func (r *Router) CommandIsPublic(command string) bool {
	tmpCommand := strings.TrimRight(command, " ")
	route, ok := r.r[tmpCommand]
	if !ok {
		return false
	}

	return route.isPublic
}
