// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/changePassword": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "changes profile password",
                "operationId": "ChangePassword",
                "parameters": [
                    {
                        "description": "Old and new passwords",
                        "name": "Password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/delivery.changePasswordStruct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-int"
                        }
                    },
                    "400": {
                        "description": "passwords are empty",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "401": {
                        "description": "Person not authorized",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        },
        "/checkAuth": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "checks that user is authenticated",
                "operationId": "checkAuth",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-int"
                        }
                    },
                    "401": {
                        "description": "Person not authorized",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        },
        "/getChats": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "gets Chats previews for user",
                "operationId": "GetChats",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Chats"
                        }
                    },
                    "400": {
                        "description": "Person not authorized",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        },
        "/getContacts": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "uploads or changes avatar",
                "operationId": "GetContacts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-delivery_docsContacts"
                        }
                    },
                    "400": {
                        "description": "Описание ошибки",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        },
        "/getProfileInfo": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "gets profile info",
                "operationId": "GetProfileInfo",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-delivery_docsUserForGetProfile"
                        }
                    },
                    "400": {
                        "description": "Person not authorized",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "logs user in",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "Person",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Person"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-int"
                        }
                    },
                    "400": {
                        "description": "wrong json structure | user not found | wrong password",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "405": {
                        "description": "use POST",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        },
        "/logout": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "logs user out",
                "operationId": "logout",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-int"
                        }
                    },
                    "400": {
                        "description": "no session to logout",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        },
        "/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "registers user",
                "operationId": "register",
                "parameters": [
                    {
                        "description": "Person",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.Person"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-int"
                        }
                    },
                    "400": {
                        "description": "user already exists | required field empty | wrong json structure",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "405": {
                        "description": "use POST",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        },
        "/updateProfileInfo": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "updates profile info",
                "operationId": "UpdateProfileInfo",
                "parameters": [
                    {
                        "description": "Send only the updated fields, and number of them",
                        "name": "userAndNumOfUpdatedFields",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/delivery.updateUserStruct-delivery_docsUserForGetProfile"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-int"
                        }
                    },
                    "401": {
                        "description": "Person not authorized",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        },
        "/uploadAvatar": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "uploads or changes avatar",
                "operationId": "UploadAvatar",
                "parameters": [
                    {
                        "type": "file",
                        "description": "avatar image",
                        "name": "avatar",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-int"
                        }
                    },
                    "400": {
                        "description": "Описание ошибки",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/domain.Response-domain_Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "delivery.changePasswordStruct": {
            "type": "object",
            "properties": {
                "newPassword": {
                    "type": "string"
                },
                "oldPassword": {
                    "type": "string"
                }
            }
        },
        "delivery.docsContacts": {
            "type": "object",
            "properties": {
                "contacts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/delivery.docsUserForGetProfile"
                    }
                }
            }
        },
        "delivery.docsUserForGetProfile": {
            "type": "object",
            "properties": {
                "about": {
                    "type": "string"
                },
                "avatar": {
                    "type": "string"
                },
                "create_date": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "last_seen_date": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "delivery.updateUserStruct-delivery_docsUserForGetProfile": {
            "type": "object",
            "properties": {
                "numOfUpdatedFields": {
                    "type": "integer"
                },
                "user": {
                    "$ref": "#/definitions/delivery.docsUserForGetProfile"
                }
            }
        },
        "domain.Chat": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "creator": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "last_message": {
                    "$ref": "#/definitions/domain.Message"
                },
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Message"
                    }
                },
                "name": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                },
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.ChatUser"
                    }
                }
            }
        },
        "domain.ChatUser": {
            "type": "object",
            "properties": {
                "chat_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "domain.Chats": {
            "type": "object",
            "properties": {
                "chats": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Chat"
                    }
                }
            }
        },
        "domain.Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "error description"
                }
            }
        },
        "domain.Message": {
            "type": "object",
            "properties": {
                "chat_id": {
                    "type": "integer"
                },
                "edited": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "message_text": {
                    "type": "string"
                },
                "sent_at": {
                    "type": "string"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "domain.Person": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "domain.Response-delivery_docsContacts": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/delivery.docsContacts"
                },
                "status": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "domain.Response-delivery_docsUserForGetProfile": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/delivery.docsUserForGetProfile"
                },
                "status": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "domain.Response-domain_Chats": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/domain.Chats"
                },
                "status": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "domain.Response-domain_Error": {
            "type": "object",
            "properties": {
                "body": {
                    "$ref": "#/definitions/domain.Error"
                },
                "status": {
                    "type": "integer",
                    "example": 200
                }
            }
        },
        "domain.Response-int": {
            "type": "object",
            "properties": {
                "body": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer",
                    "example": 200
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Messenger authorization API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
