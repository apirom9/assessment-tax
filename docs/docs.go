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
        "/admin/deductions/personal": {
            "post": {
                "description": "Update personal deduction",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tax"
                ],
                "summary": "Update personal deduction",
                "parameters": [
                    {
                        "description": "Body for update personal deduction",
                        "name": "UpdatePersonalDeductionRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tax.UpdatePersonalDeductionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tax.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    }
                }
            }
        },
        "/tax/calculations": {
            "post": {
                "description": "Calculate Tax",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tax"
                ],
                "summary": "Calculate Tax",
                "parameters": [
                    {
                        "description": "Body for calculation request",
                        "name": "CalculationRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/tax.CalculationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tax.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    }
                }
            }
        },
        "/tax/calculations/upload-csv": {
            "post": {
                "description": "Calculate Tax for upload CSV file",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tax"
                ],
                "summary": "Calculate Tax for upload CSV file",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Uploaded CSV for tax calculation",
                        "name": "taxes.csv",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/tax.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/tax.Err"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "tax.AllowanceRequest": {
            "type": "object",
            "properties": {
                "allowanceType": {
                    "type": "string",
                    "example": "donation"
                },
                "amount": {
                    "type": "number",
                    "example": 0
                }
            }
        },
        "tax.CalculationRequest": {
            "type": "object",
            "properties": {
                "allowances": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tax.AllowanceRequest"
                    }
                },
                "totalIncome": {
                    "type": "number",
                    "example": 500000
                },
                "wht": {
                    "type": "number",
                    "example": 0
                }
            }
        },
        "tax.Err": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "tax.Response": {
            "type": "object",
            "properties": {
                "tax": {
                    "type": "number",
                    "example": 29000
                },
                "taxLevel": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/tax.TaxLevelResponse"
                    }
                },
                "taxRefund": {
                    "type": "number",
                    "example": 29000
                }
            }
        },
        "tax.TaxLevelResponse": {
            "type": "object",
            "properties": {
                "level": {
                    "type": "string",
                    "example": "0-150,000"
                },
                "tax": {
                    "type": "number",
                    "example": 0
                }
            }
        },
        "tax.UpdatePersonalDeductionRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number",
                    "example": 29000
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:1323",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Tax API",
	Description:      "Tax API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
