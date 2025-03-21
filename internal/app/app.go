package app

import "github.com/google/uuid"

// App represents the application
type App struct {
	signedURLGenerator SignedURLGenerator
}

// SignedURLGenerator ...
type SignedURLGenerator interface {
	Generate(objectName string) (string, error)
}

// New creates a new instance of the application.
func New(generator SignedURLGenerator) *App {
	return &App{generator}
}

// Run starts the application.
func (a *App) Run() (string, error) {
	objectName := uuid.New().String()
	url, err := a.signedURLGenerator.Generate(objectName)
	if err != nil {
		return "", err
	}
	return url, nil
}
