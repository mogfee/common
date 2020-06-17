package bind

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Bind(s interface{}, c *gin.Context) (interface{}, error) {
	b := binding.Default(c.Request.Method, c.ContentType())
	if err := c.ShouldBindWith(s, b); err != nil {
		return nil, err
	}

	// 参数验证
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		return nil, err
	}
	return s, nil
}
