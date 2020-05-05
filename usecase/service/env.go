package service

import (
	"env-generator/model"
	"env-generator/usecase/repository"
)

// Env service interface
type Env interface {
	Exec(args model.Args) error
}

type envImpl struct {
	envRepo repository.Env
}

// NewEnv 生成メソッド
func NewEnv(envRepo repository.Env) Env {
	return &envImpl{envRepo: envRepo}
}

// Exec 実行
func (e envImpl) Exec(args model.Args) error {
	env, err := e.envRepo.Load(args)
	if err != nil {
		return err
	}
	e.envRepo.Write(env)
	return nil
}
