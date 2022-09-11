package mlib

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
	"strconv"
)

var Verson = "1.0.0"

////////////////////////////////////////////////////////////////////////////
// 基本
////////////////////////////////////////////////////////////////////////////

type BzError struct {
	Code    string
	Message string
}

func (m BzError) Error() string {
	return fmt.Sprintf("[%s] %s", m.Code, m.Message)
}

////////////////////////////////////////////////////////////////////////////
// Http Server
////////////////////////////////////////////////////////////////////////////

func StartHttpServer(port int) {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			var bzError BzError
			if errors.As(err, &bzError) {
				return ctx.Status(fiber.StatusBadRequest).JSON(bzError)
			}

			code := fiber.StatusInternalServerError

			// 必须，比如实现无路由 404
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			log.Error(err)

			return ctx.Status(code).SendString(err.Error())
		},
	})

	app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	if port <= 0 {
		port = 9090
	}

	if err := app.Listen(":" + strconv.Itoa(port)); err != nil {
		panic(err)
	}

	log.Info("Http server end.")
}

////////////////////////////////////////////////////////////////////////////
// Meta
////////////////////////////////////////////////////////////////////////////

type NsEntityMeta struct {
	Name     string        `json:"name,omitempty"`
	Label    string        `json:"label,omitempty"`
	Type     string        `json:"type,omitempty"`
	Digest   string        `json:"digest,omitempty"`
	ListSort string        `json:"listSort,omitempty"`
	Fields   []NsFieldMeta `json:"fields,omitempty"`
}

type NsFieldMeta struct {
	Name            string       `json:"name,omitempty"`
	Label           string       `json:"label,omitempty"`
	Type            NsFieldType  `json:"type,omitempty"`
	Scale           NsFieldScale `json:"scale,omitempty"`
	OptionBill      string       `json:"optionBill,omitempty"`
	DefaultValueStr string       `json:"defaultValueStr,omitempty"`
	RefEntity       string       `json:"refEntity,omitempty"`
}

type NsFieldType string

const (
	FtString    NsFieldType = "String"
	FtInt       NsFieldType = "Int"
	FtLong      NsFieldType = "Long"
	FtDouble    NsFieldType = "Double"
	FtBoolean   NsFieldType = "Boolean"
	FtDate      NsFieldType = "Date"
	FtComponent NsFieldType = "Component"
	FtImage     NsFieldType = "Image"
	FtReference NsFieldType = "Reference"
)

type NsFieldScale string

const (
	ScaleSingle NsFieldScale = "Single"
	ScaleList   NsFieldScale = "List"
)

////////////////////////////////////////////////////////////////////////////
// 实体增删改查
////////////////////////////////////////////////////////////////////////////
