// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "get clusters state",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "get clusters state",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authentication header",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/models_icos.Infrastructure"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models_icos.API": {
            "type": "object",
            "properties": {
                "authentication": {
                    "type": "string"
                },
                "authorization": {
                    "type": "string"
                },
                "commProtocol": {
                    "type": "string"
                },
                "dataFormat": {
                    "type": "string"
                },
                "protocolVersion": {
                    "type": "string"
                }
            }
        },
        "models_icos.AvailableStorage": {
            "type": "object",
            "properties": {
                "free": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models_icos.Cluster": {
            "type": "object",
            "properties": {
                "API": {
                    "$ref": "#/definitions/models_icos.API"
                },
                "any": {},
                "name": {
                    "type": "string"
                },
                "node": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/models_icos.Node"
                    }
                },
                "pod": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/models_icos.Pod"
                    }
                },
                "serviceLevelAgreement": {
                    "$ref": "#/definitions/models_icos.ServiceLevelAgreement"
                },
                "type": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "models_icos.Container": {
            "type": "object",
            "properties": {
                "containerMemory": {
                    "type": "string"
                },
                "cpuUsage": {
                    "type": "number"
                },
                "ip": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "node": {
                    "type": "string"
                },
                "port": {
                    "type": "string"
                }
            }
        },
        "models_icos.Controller": {
            "type": "object",
            "properties": {
                "API": {
                    "$ref": "#/definitions/models_icos.API"
                },
                "any": {},
                "location": {
                    "$ref": "#/definitions/models_icos.Location"
                },
                "name": {
                    "type": "string"
                },
                "serviceLevelAgreement": {
                    "$ref": "#/definitions/models_icos.ServiceLevelAgreement"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "models_icos.Device": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "models_icos.DynamicMetrics": {
            "type": "object",
            "properties": {
                "availableStorage": {
                    "$ref": "#/definitions/models_icos.AvailableStorage"
                },
                "cpuEnergyConsumption": {
                    "type": "number"
                },
                "cpuFrequency": {
                    "type": "string"
                },
                "cpuTemperature": {
                    "type": "number"
                },
                "freeRAM": {
                    "type": "integer"
                },
                "gpuEnergyConsumption": {
                    "type": "number"
                },
                "gpuFrequency": {
                    "type": "string"
                },
                "gpuTemperature": {
                    "type": "number"
                },
                "upTime": {
                    "type": "number"
                }
            }
        },
        "models_icos.Infrastructure": {
            "type": "object",
            "properties": {
                "cluster": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/models_icos.Cluster"
                    }
                },
                "controller": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/models_icos.Controller"
                    }
                },
                "timestamp": {
                    "$ref": "#/definitions/models_icos.Timestamp"
                }
            }
        },
        "models_icos.Interface": {
            "type": "object",
            "properties": {
                "engressUssage": {
                    "type": "string"
                },
                "ingressUssage": {
                    "type": "string"
                },
                "ip": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "speed": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "subnetMask": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "models_icos.Location": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "continent": {
                    "type": "string"
                },
                "country": {
                    "type": "string"
                },
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "models_icos.Node": {
            "type": "object",
            "properties": {
                "ScaScore": {
                    "type": "integer"
                },
                "devices": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/models_icos.Device"
                    }
                },
                "dynamicMetrics": {
                    "$ref": "#/definitions/models_icos.DynamicMetrics"
                },
                "location": {
                    "$ref": "#/definitions/models_icos.Location"
                },
                "name": {
                    "type": "string"
                },
                "networkInterfaces": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/models_icos.Interface"
                    }
                },
                "staticMetrics": {
                    "$ref": "#/definitions/models_icos.StaticMetrics"
                },
                "type": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                },
                "vulnerabilities": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                }
            }
        },
        "models_icos.Pod": {
            "type": "object",
            "properties": {
                "container": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/models_icos.Container"
                    }
                },
                "ip": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "numberOfApps": {
                    "type": "integer"
                },
                "numberOfContainers": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "models_icos.ServiceLevelAgreement": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "models_icos.StaticMetrics": {
            "type": "object",
            "properties": {
                "RAMMemory": {
                    "type": "integer"
                },
                "cpuArchitecture": {
                    "type": "string"
                },
                "cpuCores": {
                    "type": "integer"
                },
                "cpuMaxFrequency": {
                    "type": "integer"
                },
                "gpuCores": {
                    "type": "number"
                },
                "gpuMaxFrequency": {
                    "type": "string"
                },
                "gpuRAMMemory": {
                    "type": "string"
                },
                "storage": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models_icos.Storage"
                    }
                }
            }
        },
        "models_icos.Storage": {
            "type": "object",
            "properties": {
                "capacity": {
                    "type": "number"
                },
                "name": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "models_icos.Timestamp": {
            "type": "object",
            "properties": {
                "oldestTimestamp": {
                    "type": "number"
                },
                "timeSinceOldest": {
                    "type": "number"
                }
            }
        }
    },
    "securityDefinitions": {
        "OAuth 2.0": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Swagger Aggregator API",
	Description:      "Aggregator Microservice.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}