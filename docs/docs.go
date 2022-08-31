// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2022-08-31 11:13:24.778614 +0800 CST m=+0.112263770

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "Iobscan Ibc Explorer Support API document",
        "title": "Iobscan Ibc Explorer Support API",
        "contact": {},
        "license": {},
        "version": "visit /version"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/data/accounts_daily/api_support": {
            "get": {
                "description": "get daily accounts of chains",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api_support"
                ],
                "summary": "list",
                "operationId": "list",
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/vo.AccountsDailyResp"
                        }
                    }
                }
            }
        },
        "/data/chainList/api_support": {
            "get": {
                "description": "get IBC BUSD of chains",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api_support"
                ],
                "summary": "list",
                "operationId": "list",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "page num",
                        "name": "page_num",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "page size",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "enum": [
                            true,
                            false
                        ],
                        "type": "boolean",
                        "description": "if used count",
                        "name": "use_count",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/vo.ChainListResp"
                        }
                    }
                }
            }
        },
        "/data/fail_txs/api_support": {
            "get": {
                "description": "get  fail txs list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api_support"
                ],
                "summary": "list",
                "operationId": "list",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "page num",
                        "name": "page_num",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "page size",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/vo.FailTxsListResp"
                        }
                    }
                }
            }
        },
        "/data/relayers_fee/api_support": {
            "get": {
                "description": "get relayers fee list of chains",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api_support"
                ],
                "summary": "list",
                "operationId": "list",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "page num",
                        "name": "page_num",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "page size",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "tx_hash",
                        "name": "tx_hash",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "chain_id",
                        "name": "chain_id",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/vo.RelayerTxFeesResp"
                        }
                    }
                }
            }
        },
        "/data/statistics/api_support": {
            "get": {
                "description": "get ibc txs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "api_support"
                ],
                "summary": "list",
                "operationId": "list",
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/vo.StatisticInfoResp"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Coin": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "denom": {
                    "type": "string"
                }
            }
        },
        "model.Fee": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Coin"
                    }
                },
                "gas": {
                    "type": "integer"
                }
            }
        },
        "vo.AccountsDailyDto": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "chain_name": {
                    "type": "string"
                }
            }
        },
        "vo.AccountsDailyResp": {
            "type": "object",
            "properties": {
                "date_time": {
                    "type": "string"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/vo.AccountsDailyDto"
                    }
                },
                "time_stamp": {
                    "type": "integer"
                }
            }
        },
        "vo.ChainDto": {
            "type": "object",
            "properties": {
                "chain_id": {
                    "type": "string"
                },
                "channels": {
                    "type": "integer"
                },
                "connected_chains": {
                    "type": "integer"
                },
                "currency": {
                    "type": "string"
                },
                "ibc_tokens": {
                    "type": "integer"
                },
                "ibc_tokens_value": {
                    "type": "string"
                },
                "relayers": {
                    "type": "integer"
                },
                "transfer_txs": {
                    "type": "integer"
                },
                "transfer_txs_value": {
                    "type": "string"
                }
            }
        },
        "vo.ChainListResp": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/vo.ChainDto"
                    }
                },
                "page_info": {
                    "type": "object",
                    "$ref": "#/definitions/vo.PageInfo"
                },
                "time_stamp": {
                    "type": "integer"
                }
            }
        },
        "vo.FailTxsListDto": {
            "type": "object",
            "properties": {
                "chain_id": {
                    "type": "string"
                },
                "recv_chain": {
                    "type": "string"
                },
                "send_chain": {
                    "type": "string"
                },
                "tx_error_log": {
                    "type": "string"
                },
                "tx_hash": {
                    "type": "string"
                }
            }
        },
        "vo.FailTxsListResp": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/vo.FailTxsListDto"
                    }
                },
                "page_info": {
                    "type": "object",
                    "$ref": "#/definitions/vo.PageInfo"
                },
                "time_stamp": {
                    "type": "integer"
                }
            }
        },
        "vo.IbcStatisticDto": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "statistics_name": {
                    "type": "string"
                }
            }
        },
        "vo.PageInfo": {
            "type": "object",
            "properties": {
                "page_num": {
                    "type": "integer"
                },
                "page_size": {
                    "type": "integer"
                },
                "total_item": {
                    "type": "integer"
                },
                "total_page": {
                    "type": "integer"
                }
            }
        },
        "vo.RelayerTxFeeDto": {
            "type": "object",
            "properties": {
                "chain_id": {
                    "type": "string"
                },
                "fee": {
                    "type": "object",
                    "$ref": "#/definitions/model.Fee"
                },
                "relayer_addr": {
                    "type": "string"
                },
                "tx_hash": {
                    "type": "string"
                }
            }
        },
        "vo.RelayerTxFeesResp": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/vo.RelayerTxFeeDto"
                    }
                },
                "page_info": {
                    "type": "object",
                    "$ref": "#/definitions/vo.PageInfo"
                },
                "time_stamp": {
                    "type": "integer"
                }
            }
        },
        "vo.StatisticInfoResp": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/vo.IbcStatisticDto"
                    }
                },
                "time_stamp": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
