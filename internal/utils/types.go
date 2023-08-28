package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

type ContextWithCache struct {
	echo.Context
	Cache *cache.Cache
}
