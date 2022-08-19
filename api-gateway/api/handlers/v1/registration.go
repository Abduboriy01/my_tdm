package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/google/uuid"
	"github.com/my_tdm/api-gateway/api/auth"
	"github.com/my_tdm/api-gateway/api/model"
	pb "github.com/my_tdm/api-gateway/genproto"
	l "github.com/my_tdm/api-gateway/pkg/logger"

	//	"github.com/gin-gonic/gin/internal/json"
	"google.golang.org/protobuf/encoding/protojson"
)

func (h *handlerV1) Register(c *gin.Context) {
	var (
		body        model.RegisterModel
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("filed to bind json", l.Error(err))
		return
	}
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	exists, err := h.serviceManager.UserService().CheckUniquess(ctx, &pb.CheckUniqReq{
		Filed: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("filed to check email uniquess", l.Error(err))
		return
	}

	if exists.IsExist {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This email already in use, please use another email address",
		})
		h.log.Error("filed to check email uniquess", l.Error(err))
		return
	}

	byteUser, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("filed while marshalling user data", l.Error(err))
		return
	}

	err = h.redisStorage.SetWithTTL(body.Email, string(byteUser), int64(time.Second*300))

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		h.log.Error("filed while seting user data to redis", l.Error(err))
		return
	}
}

func (h *handlerV1) Verify(c *gin.Context) {
	var userData model.RegisterModel
	code := c.Query("code")
	email := c.Query("email")

	intCode, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		h.log.Error("filed parse query data", l.Error(err))
		return
	}

	data, err := redis.String(h.redisStorage.Get(email))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		h.log.Error("filed while getting data from redis", l.Error(err))
		return
	}

	err = json.Unmarshal([]byte(data), &userData)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		h.log.Error("filed while unmarshalling user data", l.Error(err))
		return
	}

	if intCode != userData.Code {
		if err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "your code is not match",
			})
			h.log.Error("code is invalid", l.Error(err))
			return
		}
	}

	id, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "error while generating uuid",
		})
		h.log.Error("error generating new uuid", l.Error(err))
		return
	}

	h.jwtHandler = auth.JwtHandler{
		Sub:  id.String(),
		Iss:  "client",
		Role: "authorized",
		Log:  h.log,
	}

	access, refresh, err := h.jwtHandler.GenerateJwt()
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "error while generating uuid",
		})
		h.log.Error("error generating new jwt toke", l.Error(err))
		return
	}

	// Create user logic

	c.JSON(http.StatusOK, &model.RegisterResponseModel{
		UserID:       id.String(),
		AccessToken:  access,
		RefreshToken: refresh,
	})

}

func (h *handlerV1) SendMessage(email, message string) error {
	return nil
}
