import json
import sys

# При генерации сваггера добавляет в него запрос принятия ТК.
custom_path = {
    "/receiver/cms/v1": {
          "post": {
            "summary": "Receive a file",
            "description": "Receive a file in the request body.",
            "consumes": [
              "multipart/form-data"
            ],
            "produces": [
              "text/plain"
            ],
            "operationId": "PackageReceiver_ReceivePackageV1",
            "parameters": [
              {
                "name": "file",
                "in": "formData",
                "description": "The file to upload",
                "required": True,
                "type": "file"
              },
              {
                "name": "Send-Receipt-To",
                "description": "URL for TK",
                "in": "header",
                "required": True,
                "type": "string"
              },
              {
                  "name": "Content-Disposition",
                  "description": "Name of TP",
                  "in": "header",
                  "required": False,
                  "type": "string"
              },
              {
                "name": "X-Sender-Operator",
                "in": "header",
                "required": False,
                "type": "string"
                },
              {
                  "name": "X-Recipient-Operator",
                  "in": "header",
                  "required": False,
                  "type": "string"
              },
            ],
            "responses": {
              "200": {
                "description": "A successful response."
              },
              "default": {
                "description": "An unexpected error response.",
                "schema": {
                    "type": "string",
                    "format": "byte"
                  }
              }
            },
            "tags": [
              "PackageReceiver"
            ]
          }
        }
}

def add_custom_path(swagger_file_path, custom_path):
    try:
        # Read existing Swagger JSON file
        with open(swagger_file_path, 'r') as file:
            swagger_data = json.load(file)

        # Update paths with the custom path
        custom_path.update(swagger_data['paths'])
        swagger_data['paths'] = custom_path

        # Write the updated Swagger JSON back to the file
        with open(swagger_file_path, 'w') as file:
            json.dump(swagger_data, file, indent=2)

        print("Custom path added successfully.")

    except FileNotFoundError:
        print("Swagger JSON file not found.")
    except Exception as e:
        print("Error:", e)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 script.py <swagger_file_path>")
        sys.exit(1)

    swagger_file_path = sys.argv[1]
    add_custom_path(swagger_file_path, custom_path)
