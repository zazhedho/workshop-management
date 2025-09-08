package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"workshop-management/pkg/logger"
	"workshop-management/utils"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func GetAppConf(key string, def interface{}, rdbCache *redis.Client) interface{} {
	var (
		err   error
		cache bool
	)

	cacheKey := utils.RedisAppConf
	appConf := make(map[string]string)
	var getNewConfig bool
	if consul := os.Getenv("CONSUL"); consul != "" {
		cache = strings.ToLower(os.Getenv("CACHE")) == "on" && rdbCache != nil

		if cache {
			if jsonAppConf, err := rdbCache.Get(context.Background(), cacheKey).Result(); err == nil {
				if err = json.Unmarshal([]byte(jsonAppConf), &appConf); err != nil {
					logger.WriteLog(logger.LogLevelError, fmt.Sprintf("utils.GetAppConf; Unmarshal conf from cache; %s; error: %+v;", jsonAppConf, err))
					getNewConfig = true
				}
			} else if errors.Is(err, redis.Nil) {
				getNewConfig = true
			}
		} else {
			getNewConfig = true
		}

		if getNewConfig {
			consulPath := fmt.Sprintf("%s/%s", os.Getenv("CONSUL_PATH"), os.Getenv("APP_ENV"))
			runtimeViper := viper.New()
			runtimeViper.AddRemoteProvider("consul", consul, consulPath)
			runtimeViper.SetConfigType("json") // Need to explicitly set this to json
			if err = runtimeViper.ReadRemoteConfig(); err != nil {
				logger.WriteLog(logger.LogLevelError, fmt.Sprintf("utils.GetAppConf; Loading config: %s/%s; error: %+v;", consul, consulPath, err))
			} else if err = runtimeViper.Unmarshal(&appConf); err != nil {
				logger.WriteLog(logger.LogLevelError, fmt.Sprintf("utils.GetAppConf; Loading congif: %s/%s; error: unable to decode into map, %+v;", consul, consulPath, err))
			}
		}
	} else {
		configName := "app"
		pathConfig := os.Getenv("APP_CONFIG")
		if pathConfig == "" {
			pathConfig = "config"
			configName = os.Getenv("APP_ENV")
		}

		viper.AddConfigPath(pathConfig)
		viper.SetConfigType("env")
		viper.SetConfigName(configName)

		err = viper.ReadInConfig()
		if err != nil {
			logger.WriteLog(logger.LogLevelError, fmt.Sprintf("utils.GetAppConf; Loading config: %s - %s.env; error:  %+v;", pathConfig, configName, err))
		} else {
			_ = viper.Unmarshal(&appConf)
		}
	}

	if len(appConf) > 0 {
		if _, ok := appConf["config_id"]; !ok {
			appConf["config_id"] = uuid.NewString()
		}

		if appConf["config_id"] != os.Getenv("CONFIG_ID") {
			for k, v := range appConf {
				//Reset all env variable
				os.Setenv(strings.ToUpper(k), v)
			}
		}

		if cache && getNewConfig {
			go func(rdbCache *redis.Client, cacheKey string, data map[string]string) {
				if cacheData, err := json.Marshal(data); err == nil {
					ttl := utils.GetEnv("TTL_CACHE_CONFIG_APP", time.Duration(60*60*24)).(time.Duration) * time.Second
					_ = rdbCache.Set(context.Background(), cacheKey, string(cacheData), ttl).Err()
				}
			}(rdbCache, cacheKey, appConf)
		}
	}

	return utils.GetEnv(key, def)
}
