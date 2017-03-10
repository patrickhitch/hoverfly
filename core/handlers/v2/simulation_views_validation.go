package v2

var requestResponsePairDefinition = map[string]interface{}{
	"type": "object",
	"required": []string{
		"request",
		"response",
	},
	"properties": map[string]interface{}{
		"request": map[string]interface{}{
			"$ref": "#/definitions/request",
		},
		"response": map[string]interface{}{
			"$ref": "#/definitions/response",
		},
	},
}

var requestV1Definition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"scheme": map[string]interface{}{
			"type": "string",
		},
		"destination": map[string]interface{}{
			"type": "string",
		},
		"path": map[string]interface{}{
			"type": "string",
		},
		"query": map[string]interface{}{
			"type": "string",
		},
		"body": map[string]interface{}{
			"type": "string",
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/headers",
		},
	},
}

var requestV2Definition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"scheme": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"destination": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"path": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"query": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"body": map[string]interface{}{
			"$ref": "#/definitions/field-matchers",
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/headers",
		},
	},
}

var responseDefinition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"body": map[string]interface{}{
			"type": "string",
		},
		"encodedBody": map[string]interface{}{
			"type": "boolean",
		},
		"headers": map[string]interface{}{
			"$ref": "#/definitions/headers",
		},
		"status": map[string]interface{}{
			"type": "integer",
		},
	},
}

var requestFieldMatchersV2Definition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"exactMatch": map[string]interface{}{
			"type": "string",
		},
		"globMatch": map[string]interface{}{
			"type": "string",
		},
		"regexMatch": map[string]interface{}{
			"type": "string",
		},
		"xpathMatch": map[string]interface{}{
			"type": "string",
		},
		"jsonMatch": map[string]interface{}{
			"type": "string",
		},
	},
}

var headersDefinition = map[string]interface{}{
	"type": "object",
	"additionalProperties": map[string]interface{}{
		"type": "array",
		"items": map[string]interface{}{
			"type": "string",
		},
	},
}

var delaysDefinition = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"urlPattern": map[string]interface{}{
			"type": "string",
		},
		"httpMethod": map[string]interface{}{
			"type": "string",
		},
		"delay": map[string]interface{}{
			"type": "integer",
		},
	},
}

var metaDefinition = map[string]interface{}{
	"type": "object",
	"required": []string{
		"schemaVersion",
	},
	"properties": map[string]interface{}{
		"schemaVersion": map[string]interface{}{
			"type": "string",
		},
		"hoverflyVersion": map[string]interface{}{
			"type": "string",
		},
		"timeExported": map[string]interface{}{
			"type": "string",
		},
	},
}

var SimulationViewV2Schema = map[string]interface{}{
	"description": "Hoverfly simulation schema",
	"type":        "object",
	"required": []string{
		"data", "meta",
	},
	"additionalProperties": false,
	"properties": map[string]interface{}{
		"data": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pairs": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/request-response-pair",
					},
				},
				"globalActions": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"delays": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/delay",
							},
						},
					},
				},
			},
		},
		"meta": map[string]interface{}{
			"$ref": "#/definitions/meta",
		},
	},
	"definitions": map[string]interface{}{
		"request-response-pair": requestResponsePairDefinition,
		"request":               requestV2Definition,
		"response":              responseDefinition,
		"field-matchers":        requestFieldMatchersV2Definition,
		"headers":               headersDefinition,
		"delay":                 delaysDefinition,
		"meta":                  metaDefinition,
	},
}

var SimulationViewV1Schema = map[string]interface{}{
	"description": "Hoverfly simulation schema",
	"type":        "object",
	"required": []string{
		"data", "meta",
	},
	"additionalProperties": false,
	"properties": map[string]interface{}{
		"data": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"pairs": map[string]interface{}{
					"type": "array",
					"items": map[string]interface{}{
						"$ref": "#/definitions/request-response-pair",
					},
				},
				"globalActions": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"delays": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"$ref": "#/definitions/delay",
							},
						},
					},
				},
			},
		},
		"meta": map[string]interface{}{
			"$ref": "#/definitions/meta",
		},
	},
	"definitions": map[string]interface{}{
		"request-response-pair": requestResponsePairDefinition,
		"request":               requestV1Definition,
		"response":              responseDefinition,
		"headers":               headersDefinition,
		"delay":                 delaysDefinition,
		"meta":                  metaDefinition,
	},
}
