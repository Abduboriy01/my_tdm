// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/v1/users": {
            "get": {
                "description": "This api is using for getting user list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user list summary",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "page",
                        "name": "page",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "This api is using for creating new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Create user summary",
                "parameters": [
                    {
                        "description": "user body",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.CreateUserRequestBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/v1/users/login/{email}/{password}": {
            "get": {
                "description": "This api using for logging registered user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email",
                        "name": "email",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Password",
                        "name": "password",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/v1/users/register": {
            "post": {
                "description": "This api is using for registering user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Register user summary",
                "parameters": [
                    {
                        "description": "user_body",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.RegisterUserAuthReqBody"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/users/verfication": {
            "post": {
                "description": "This api using for verifying registered user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "parameters": [
                    {
                        "description": "user body",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.Emailver"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/users/{id}": {
            "get": {
                "description": "This api is using for getting user by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get user summary",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "v1.CreateUserRequestBody": {
            "type": "object",
            "properties": {
                "first_name": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "last_name": {
                    "type": "string"
                },
                "posts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.Post"
                    }
                }
            }
        },
        "v1.Emailver": {
            "type": "object",
            "properties": {
                "Code": {
                    "type": "string"
                },
                "Email": {
                    "type": "string"
                }
            }
        },
        "v1.Media": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "v1.Post": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "medias": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.Media"
                    }
                },
                "name": {
                    "type": "string"
                },
                "user_id": {
                    "type": "string"
                }
            }
        },
        "v1.RegisterUserAuthReqBody": {
            "type": "object",
            "properties": {
                "Code": {
                    "type": "string"
                },
                "Email": {
                    "type": "string"
                },
                "FirstName": {
                    "description": "Id          string ` + "`" + `protobuf:\"bytes,1,opt,name=Id,proto3\" json:\"Id\"` + "`" + `",
                    "type": "string"
                },
                "Password": {
                    "type": "string"
                },
                "PhoneNumber": {
                    "type": "string"
                },
                "Username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}