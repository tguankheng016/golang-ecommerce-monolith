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
        "/api/v1/accounts/app-permissions": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get All App Permissions",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Get All App Permissions",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetAllPermissionResult"
                        }
                    }
                }
            }
        },
        "/api/v1/accounts/authenticate": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Authenticate",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Authenticate",
                "parameters": [
                    {
                        "description": "AuthenticateRequest",
                        "name": "AuthenticateRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/AuthenticateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/AuthenticateResult"
                        }
                    }
                }
            }
        },
        "/api/v1/accounts/current-session": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get Current User Session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Get Current User Session",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetCurrentSessionResult"
                        }
                    }
                }
            }
        },
        "/api/v1/accounts/refresh-token": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Refresh access token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Refresh access token",
                "parameters": [
                    {
                        "description": "RefreshTokenRequest",
                        "name": "RefreshTokenRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RefreshTokenRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RefreshTokenResult"
                        }
                    }
                }
            }
        },
        "/api/v1/accounts/sign-out": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Sign out",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Accounts"
                ],
                "summary": "Sign out",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/role": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update role",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Update role",
                "parameters": [
                    {
                        "description": "EditRoleDto",
                        "name": "EditRoleDto",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/EditRoleDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RoleDto"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create new role",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Create new role",
                "parameters": [
                    {
                        "description": "CreateRoleDto",
                        "name": "CreateRoleDto",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/CreateRoleDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/RoleDto"
                        }
                    }
                }
            }
        },
        "/api/v1/role/{roleId}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get role by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Get role by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Role Id",
                        "name": "roleId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetRoleByIdResult"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete role",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Delete role",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Role Id",
                        "name": "roleId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/roles": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all roles",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Roles"
                ],
                "summary": "Get all roles",
                "parameters": [
                    {
                        "type": "string",
                        "name": "filters",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "maxResultCount",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "skipCount",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "sorting",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetRolesResult"
                        }
                    }
                }
            }
        },
        "/api/v1/user": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update user",
                "parameters": [
                    {
                        "description": "EditUserDto",
                        "name": "EditUserDto",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/EditUserDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/UserDto"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Create new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Create new user",
                "parameters": [
                    {
                        "description": "CreateUserDto",
                        "name": "CreateUserDto",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/CreateUserDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/UserDto"
                        }
                    }
                }
            }
        },
        "/api/v1/user/{userId}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get user by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get user by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User Id",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetUserByIdResult"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Delete user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Delete user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User Id",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/user/{userId}/permissions": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get user permissions",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get user permissions",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User Id",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/UserPermissionsResult"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Update user permissions",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update user permissions",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User Id",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "UpdateUserPermissionDto",
                        "name": "UpdateUserPermissionDto",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/UpdateUserPermissionDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/user/{userId}/reset-permissions": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Reset user permissions",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Reset user permissions",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User Id",
                        "name": "userId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/v1/users": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Get all users",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get all users",
                "parameters": [
                    {
                        "type": "string",
                        "name": "filters",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "maxResultCount",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "name": "skipCount",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "name": "sorting",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetUsersResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "AuthenticateRequest": {
            "type": "object",
            "required": [
                "password",
                "usernameOrEmailAddress"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "usernameOrEmailAddress": {
                    "type": "string"
                }
            }
        },
        "AuthenticateResult": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "expireInSeconds": {
                    "type": "integer"
                },
                "refreshToken": {
                    "type": "string"
                },
                "refreshTokenExpireInSeconds": {
                    "type": "integer"
                }
            }
        },
        "CreateOrEditRoleDto": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "grantedPermissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "CreateOrEditUserDto": {
            "type": "object",
            "required": [
                "email",
                "firstName",
                "lastName",
                "userName"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 256
                },
                "firstName": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 3
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 3
                },
                "password": {
                    "type": "string"
                },
                "roleIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "userName": {
                    "type": "string",
                    "maxLength": 256,
                    "minLength": 8
                }
            }
        },
        "CreateRoleDto": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "grantedPermissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "CreateUserDto": {
            "type": "object",
            "required": [
                "email",
                "firstName",
                "lastName",
                "userName"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 256
                },
                "firstName": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 3
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 3
                },
                "password": {
                    "type": "string"
                },
                "roleIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "userName": {
                    "type": "string",
                    "maxLength": 256,
                    "minLength": 8
                }
            }
        },
        "EditRoleDto": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "grantedPermissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "EditUserDto": {
            "type": "object",
            "required": [
                "email",
                "firstName",
                "lastName",
                "userName"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "maxLength": 256
                },
                "firstName": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 3
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 3
                },
                "password": {
                    "type": "string"
                },
                "roleIds": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "userName": {
                    "type": "string",
                    "maxLength": 256,
                    "minLength": 8
                }
            }
        },
        "GetAllPermissionResult": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/PermissionGroupDto"
                    }
                }
            }
        },
        "GetCurrentSessionResult": {
            "type": "object",
            "properties": {
                "allPermissions": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "boolean"
                    }
                },
                "grantedPermissions": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "boolean"
                    }
                },
                "user": {
                    "$ref": "#/definitions/UserLoginInfoDto"
                }
            }
        },
        "GetRoleByIdResult": {
            "type": "object",
            "properties": {
                "role": {
                    "$ref": "#/definitions/CreateOrEditRoleDto"
                }
            }
        },
        "GetRolesResult": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/RoleDto"
                    }
                },
                "totalCount": {
                    "type": "integer"
                }
            }
        },
        "GetUserByIdResult": {
            "type": "object",
            "properties": {
                "user": {
                    "$ref": "#/definitions/CreateOrEditUserDto"
                }
            }
        },
        "GetUsersResult": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/UserDto"
                    }
                },
                "totalCount": {
                    "type": "integer"
                }
            }
        },
        "PermissionDto": {
            "type": "object",
            "properties": {
                "displayName": {
                    "type": "string"
                },
                "isGranted": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "PermissionGroupDto": {
            "type": "object",
            "properties": {
                "groupName": {
                    "type": "string"
                },
                "permissions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/PermissionDto"
                    }
                }
            }
        },
        "RefreshTokenRequest": {
            "type": "object",
            "required": [
                "token"
            ],
            "properties": {
                "token": {
                    "type": "string",
                    "minLength": 10
                }
            }
        },
        "RefreshTokenResult": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "expireInSeconds": {
                    "type": "integer"
                }
            }
        },
        "RoleDto": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "UpdateUserPermissionDto": {
            "type": "object",
            "properties": {
                "grantedPermissions": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "UserDto": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                },
                "userName": {
                    "type": "string"
                }
            }
        },
        "UserLoginInfoDto": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                },
                "userName": {
                    "type": "string"
                }
            }
        },
        "UserPermissionsResult": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
