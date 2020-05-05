package repository

import (
	"github.com/TechLoCo/env-generator/model"
)

// Env interface
type Env interface {
	Load(args model.Args) (model.Env, error)
	Write(env model.Env)
}
