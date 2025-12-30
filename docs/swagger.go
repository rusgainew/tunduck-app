package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "openapi": "3.0.0",
  "info": {
    "title": "Tunduc API System",
    "description": "Enterprise API for managing ESF documents, organizations, and users with caching, rate limiting, and health monitoring",
    "version": "1.0.0",
    "contact": {
      "name": "API Support",
      "url": "https://github.com/rusgainew/tunduck-app",
      "email": "support@example.com"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    }
  },
  "servers": [
    {"url": "http://localhost:8080", "description": "Development"},
    {"url": "https://api.example.com", "description": "Production"}
  ],
  "components": {
    "securitySchemes": {
      "BearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT"
      }
    },
    "schemas": {
      "User": {
        "type": "object",
        "properties": {
          "id": {"type": "integer"},
          "username": {"type": "string"},
          "email": {"type": "string"}
        }
      },
      "RegisterRequest": {
        "type": "object",
        "required": ["username", "email", "password"],
        "properties": {
          "username": {"type": "string", "minLength": 3},
          "email": {"type": "string", "format": "email"},
          "password": {"type": "string", "minLength": 8}
        }
      },
      "LoginRequest": {
        "type": "object",
        "required": ["email", "password"],
        "properties": {
          "email": {"type": "string", "format": "email"},
          "password": {"type": "string"}
        }
      },
      "LoginResponse": {
        "type": "object",
        "properties": {
          "token": {"type": "string"},
          "user": {"$ref": "#/components/schemas/User"}
        }
      },
      "HealthStatus": {
        "type": "object",
        "properties": {
          "status": {"type": "string", "enum": ["UP", "DOWN"]},
          "timestamp": {"type": "string", "format": "date-time"},
          "uptime": {"type": "string"},
          "components": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "name": {"type": "string"},
                "status": {"type": "string", "enum": ["UP", "DOWN"]},
                "response_time": {"type": "string"},
                "message": {"type": "string"}
              }
            }
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "error": {"type": "string"},
          "message": {"type": "string"}
        }
      },
      "RateLimitError": {
        "type": "object",
        "properties": {
          "error": {"type": "string"},
          "message": {"type": "string"},
          "reset": {"type": "integer", "format": "int64"}
        }
      }
    }
  },
  "paths": {
    "/api/auth/register": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Register a new user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/RegisterRequest"}
            }
          }
        },
        "responses": {
          "201": {
            "description": "User successfully registered",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/User"}
              }
            }
          },
          "400": {
            "description": "Invalid input",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "429": {
            "description": "Rate limit exceeded",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/RateLimitError"}
              }
            }
          }
        }
      }
    },
    "/api/auth/login": {
      "post": {
        "tags": ["Authentication"],
        "summary": "User login",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/LoginRequest"}
            }
          }
        },
        "responses": {
          "200": {
            "description": "Login successful",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/LoginResponse"}
              }
            }
          },
          "401": {
            "description": "Invalid credentials",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "429": {
            "description": "Rate limit exceeded",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/RateLimitError"}
              }
            }
          }
        }
      }
    },
    "/api/auth/logout": {
      "post": {
        "tags": ["Authentication"],
        "summary": "User logout",
        "security": [{"BearerAuth": []}],
        "responses": {
          "200": {
            "description": "Logout successful",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {"message": {"type": "string"}}
                }
              }
            }
          },
          "401": {
            "description": "Missing or invalid token",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/ErrorResponse"}
              }
            }
          },
          "429": {
            "description": "Rate limit exceeded",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/RateLimitError"}
              }
            }
          }
        }
      }
    },
    "/health": {
      "get": {
        "tags": ["System"],
        "summary": "Health check",
        "description": "Check system health (PostgreSQL, Redis)",
        "responses": {
          "200": {
            "description": "Systems operational",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/HealthStatus"}
              }
            }
          },
          "503": {
            "description": "System unavailable",
            "content": {
              "application/json": {
                "schema": {"$ref": "#/components/schemas/HealthStatus"}
              }
            }
          }
        }
      }
    },
    "/metrics": {
      "get": {
        "tags": ["System"],
        "summary": "Prometheus metrics",
        "responses": {
          "200": {
            "description": "Metrics data",
            "content": {
              "text/plain": {
                "schema": {"type": "string"}
              }
            }
          }
        }
      }
    }
  }
}`

func init() {
	swag.Register(swag.Name, &swag.Spec{
		Version:          "1.0.0",
		Host:             "localhost:8080",
		BasePath:         "",
		Schemes:          []string{"http", "https"},
		Title:            "Tunduc API System",
		Description:      "Enterprise API for managing ESF documents, organizations, and users",
		InfoInstanceName: "swagger",
		SwaggerTemplate:  docTemplate,
	})
}
