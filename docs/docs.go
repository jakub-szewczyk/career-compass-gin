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
        "/health-check": {
            "get": {
                "description": "Returns the health status of the service",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health check"
                ],
                "summary": "Health check",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.HealthCheckResBody"
                        }
                    }
                }
            }
        },
        "/job-applications": {
            "post": {
                "description": "Processes and creates a new job application with the provided data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Job application"
                ],
                "summary": "Submit a new job application",
                "parameters": [
                    {
                        "description": "Job application details",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.CreateJobApplicationReqBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.CreateJobApplicationResBody"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/password/reset": {
            "put": {
                "description": "Allows a user to set a new password using a valid reset token. This endpoint is typically used in the \"forgot password\" flow.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Password"
                ],
                "summary": "Reset password",
                "parameters": [
                    {
                        "description": "New user credentials",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.ResetPasswordReqBody"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            },
            "post": {
                "description": "Generates and sends a password reset token to the user's email address",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Password"
                ],
                "summary": "Initiate password reset",
                "parameters": [
                    {
                        "description": "User's email address",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.InitPasswordResetReqBody"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/profile": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Retrieves and returns the profile information of the currently authenticated user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Profile"
                ],
                "summary": "Get user profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ProfileResBody"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/profile/verify-email": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Sends a verification email to the user. This endpoint can be used to resend the email if needed.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Profile"
                ],
                "summary": "Send user verification email",
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            },
            "patch": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Confirms the email address of the currently authenticated user. This endpoint requires an email verification token sent to the user's registered email.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Profile"
                ],
                "summary": "Verify user email",
                "parameters": [
                    {
                        "description": "Email verification data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.VerifyEmailReqBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.ProfileResBody"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/sign-in": {
            "post": {
                "description": "Authenticates a user and returns a JWT token for session management. Valid credentials are required to access the system.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User sign in",
                "parameters": [
                    {
                        "description": "User sign in data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SignInReqBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SignInResBody"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        },
        "/sign-up": {
            "post": {
                "description": "Registers a new user account with the provided details, including email, password, and other relevant information. Verification email will be sent.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User sign up",
                "parameters": [
                    {
                        "description": "User sign up data",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SignUpReqBody"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.SignUpResBody"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.Status": {
            "type": "string",
            "enum": [
                "IN_PROGRESS",
                "REJECTED",
                "ACCEPTED"
            ],
            "x-enum-varnames": [
                "StatusINPROGRESS",
                "StatusREJECTED",
                "StatusACCEPTED"
            ]
        },
        "models.CreateJobApplicationReqBody": {
            "type": "object",
            "required": [
                "companyName",
                "dateApplied",
                "jobTitle",
                "status"
            ],
            "properties": {
                "companyName": {
                    "type": "string",
                    "example": "Evil Corp Inc."
                },
                "dateApplied": {
                    "type": "string",
                    "example": "2025-03-14T12:34:56Z"
                },
                "jobPostingURL": {
                    "type": "string",
                    "example": "https://glassbore.com/jobs/swe420692137"
                },
                "jobTitle": {
                    "type": "string",
                    "example": "Software Engineer"
                },
                "maxSalary": {
                    "type": "number",
                    "example": 70000
                },
                "minSalary": {
                    "type": "number",
                    "example": 50000
                },
                "notes": {
                    "type": "string",
                    "example": "Follow up in two weeks"
                },
                "status": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/db.Status"
                        }
                    ],
                    "example": "IN_PROGRESS"
                }
            }
        },
        "models.CreateJobApplicationResBody": {
            "type": "object",
            "properties": {
                "companyName": {
                    "type": "string",
                    "example": "Evil Corp Inc."
                },
                "dateApplied": {
                    "type": "string",
                    "example": "2025-03-14T12:34:56Z"
                },
                "id": {
                    "type": "string",
                    "example": "f4d15edc-e780-42b5-957d-c4352401d9ca"
                },
                "jobPostingURL": {
                    "type": "string",
                    "example": "https://glassbore.com/jobs/swe420692137"
                },
                "jobTitle": {
                    "type": "string",
                    "example": "Software Engineer"
                },
                "maxSalary": {
                    "type": "number",
                    "example": 70000
                },
                "minSalary": {
                    "type": "number",
                    "example": 50000
                },
                "notes": {
                    "type": "string",
                    "example": "Follow up in two weeks"
                },
                "status": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/db.Status"
                        }
                    ],
                    "example": "IN_PROGRESS"
                }
            }
        },
        "models.Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "models.HealthCheckResBody": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "healthy"
                }
            }
        },
        "models.InitPasswordResetReqBody": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john.doe@example.com"
                }
            }
        },
        "models.ProfileResBody": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john.doe@example.com"
                },
                "firstName": {
                    "type": "string",
                    "example": "John"
                },
                "id": {
                    "type": "string",
                    "example": "f4d15edc-e780-42b5-957d-c4352401d9ca"
                },
                "isEmailVerified": {
                    "type": "boolean",
                    "example": true
                },
                "lastName": {
                    "type": "string",
                    "example": "Doe"
                }
            }
        },
        "models.ResetPasswordReqBody": {
            "type": "object",
            "required": [
                "confirmPassword",
                "password",
                "passwordResetToken"
            ],
            "properties": {
                "confirmPassword": {
                    "type": "string",
                    "example": "qwerty!123456789"
                },
                "password": {
                    "description": "TODO: Improve password strength",
                    "type": "string",
                    "minLength": 16,
                    "example": "qwerty!123456789"
                },
                "passwordResetToken": {
                    "type": "string",
                    "example": "ec6c66fbd3d92b1ad44f21613c5ee2e82c3dd65e8c918945308087ce77b5fe47"
                }
            }
        },
        "models.SignInReqBody": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "john.doe@example.com"
                },
                "password": {
                    "description": "TODO: Improve password strength",
                    "type": "string",
                    "minLength": 16,
                    "example": "qwerty!123456789"
                }
            }
        },
        "models.SignInResBody": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/models.ProfileResBody"
                }
            }
        },
        "models.SignUpReqBody": {
            "type": "object",
            "required": [
                "confirmPassword",
                "email",
                "firstName",
                "lastName",
                "password"
            ],
            "properties": {
                "confirmPassword": {
                    "type": "string",
                    "example": "qwerty!123456789"
                },
                "email": {
                    "type": "string",
                    "example": "john.doe@example.com"
                },
                "firstName": {
                    "type": "string",
                    "example": "John"
                },
                "lastName": {
                    "type": "string",
                    "example": "Doe"
                },
                "password": {
                    "description": "TODO: Improve password strength",
                    "type": "string",
                    "minLength": 16,
                    "example": "qwerty!123456789"
                }
            }
        },
        "models.SignUpResBody": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string",
                    "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mzk4MDQ1NTEsInN1YiI6ImpvaG4uZG9lQGdtYWlsLmNvbSIsInVpZCI6IjZiZTA1YTcyLTc5OGQtNGI3Ny1iOGQzLTc3MjNhN2JmM2FkYSJ9.5sj2fHB3pky3N6-mDgaPQCQA0gkEz4oQsdtVEC9BLqE"
                },
                "user": {
                    "$ref": "#/definitions/models.ProfileResBody"
                }
            }
        },
        "models.VerifyEmailReqBody": {
            "type": "object",
            "required": [
                "verificationToken"
            ],
            "properties": {
                "verificationToken": {
                    "type": "string",
                    "example": "2cc313c8b72f8e5b725e07130d0b851811f2e60c8b19f085b3aa58d1516ef767"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
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
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Career Compass REST API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
