package adapter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/TechLoCo/env-generator/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/pkg/errors"
)

// Env .
type Env struct{}

// NewEnv 生成メソッド
func NewEnv() *Env {
	return &Env{}
}

// Load secrets managerから取得
func (e *Env) Load(args model.Args) (model.Env, error) {
	// aws configに記述されているユーザーを利用する
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Profile:                 args.Profile,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(args.Region),
		},
	}))

	svc := secretsmanager.New(sess, aws.NewConfig().WithRegion(args.Region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(args.Secret),
		VersionStage: aws.String(args.Version),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get secret value")
	}

	var res model.Env
	if err := json.Unmarshal([]byte(*result.SecretString), &res); err != nil {
		return nil, errors.Wrap(err, "failed to secret value to json")
	}
	return res, nil
}

// write 標準出力にenvを出力
func (e *Env) Write(env model.Env) {
	// prefixごとにまとめる
	prefixMap := make(map[string]map[string]string)
	for k, v := range env {
		prefix := strings.Split(k, "_")[0]
		if valueEnv, ok := prefixMap[prefix]; ok {
			valueEnv[k] = v
			prefixMap[prefix] = valueEnv
		} else {
			valueEnv := make(map[string]string)
			valueEnv[k] = v
			prefixMap[prefix] = valueEnv
		}
	}

	// prefixでソートするためにprefixList作成
	prefixList := make([]string, len(prefixMap))
	for prefix, _ := range prefixMap {
		prefixList = append(prefixList, prefix)
	}
	sort.Strings(prefixList)

	// 出力
	for _, prefix := range prefixList {
		fmt.Printf("# %s\n", prefix)
		for k, v := range prefixMap[prefix] {
			fmt.Printf("%s=%s\n", k, v)
		}
		fmt.Println()
	}
}
