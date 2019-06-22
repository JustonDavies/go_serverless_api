Go Serverless API
==========

This project exists as a simple service-oriented server-less showcase API to CRUD basic `Task` items. 

The code in the this project enables users:
  - Deploy and manage infrastructure
  - Deploy and manage multiple applications
  - Simple command line scripts to test the API
  - Create/Read/Update/Delete/List Tasks  
  
NOTE: No code generation or templates were leveraged in this project. All ~6,000 lines were entered by hand to help demonstrate in depth knowledge of each moving part of this application stack.  

Development
===========

This project is developed using Docker ( and by extension Terraform and Serverless.js ) 
 
The project assumes you have a basic understanding of how to run these tools but provides a small set of scripts to help setup runtime environments and version states securely.

### Project Structure
The project is follows the established paradigms in [golang standards project-layout](https://github.com/golang-standards/project-layout) project and is broken down into the following directory structure:

```
├── .state/
│   ├── .serverless
│   ├── terraform.tfstate.d
│   ├── backup
│   │   ├── ...
├── build/
│   ├── ...
├── cmd/
│   ├── shared
│   ├── endpoint
│   ├── ...
├── configs/
│   ├── secrets
│   ├── ...
├── dockerfiles/
│   ├── ...
├── pkg/
│   ├── services
│   ├── ...
├── scripts/
│   ├── api
│   ├── ...
...
```

* `.state/` A directory which helps Terraform and Serverless maintain a versioned and secure state
* `build/` A directory used to store a collection of compiled executables
* `cmd/` A collection of `main` executables (and closely coupled helpers), each representing one endpoint
* `configs/` Configuration files, templates, secrets and default configs
* `dockerfiles/` Docker files describing build and run containers
* `pkg/` Library code that's ok to use by external applications
* `scripts/` Scripts to perform various build, deploy and test a deployed collection of endpoints

### Idioms
* `go fmt`
* Whole-word variable names
* Strong, predictable commenting and organization of functions and declarations 
* Try not to log errors, return them and expect consumer handling

Build
===========
To build binaries for all supported endpoints just run the docker-based script below:

```
    $ scripts/build.sh
```

Provision
===========
To build out managed infrastructure for all supported endpoints you just need to run the docker-based `./scripts/provision.sh` script and then run the following commands:

```
    $ terraform init
    $ terraform apply -var-file='configs/secrets/production.tfvars'
```

The following `best practices` are also highly recommended:
 
 - Terraform Workspaces: Take advantage of `workspaces` in to support multiple stages / environments by running the following before getting started

```
    $ terraform init
    $ terraform workspace new YOUR_ENVIRONMENT_NAME_HERE
```

  - AWS Credentials: It is recommended you have your AWS credentials configured in the default manner documented [here](https://docs.aws.amazon.com/cli/latest/userguide/cli-config-files.html) however the AWS CLI is not required.

Deploy
===========
To deploy the application (after building) for all supported endpoints you just need to run the docker-based `./scripts/provision.sh` script and then run the following commands:

```
    $ serverless deploy
```

Secrets
===========
  The project expects a secrets file to be present for each environment in `./configs/secrets/` for each of the following items:
  
  - Terraform `ENVIRONMENT.tfvars`

  ```
    //-- General -----------------------------------------------------------------------------------------------------------
    stage = "..."
    
    aws_profile = "..."
    aws_region = "..."
    
    //-- Secrets -----------------------------------------------------------------------------------------------------------
    go_serverless_api_database_password = "..."
  ```
  
  - Serverless.js `ENVIRONMENT.yml` (NOTE: The contents of this file should be taken from the output of `terraform apply` of `terraform show`)

  ```
    #-- BEGIN MINIMUM VIABLE SERVERLESS SECRETS ----------
    aws:
      stage: ...
      region: ...
      profile: ...
    
      rds:
        engine: ...
        url: ...
        name: ...
        username: ...
        password: ...
        ssl_mode: ...
    
      vpc:
        subnet_ids: ...
        security_group_ids: ...
    
      schedule:
        warming: ...
    
    #-- END MINIMUM VIABLE SERVERLESS SECRETS ----------

  ```
  
Documentation
===========

### Testing
This application can be tested provided you have:
  
  - An available instance of Postgres 11.1 available with a valid connection string in the environment variables
  - Environment variables which adhere to the following:    
    ```
      DATABASE_CONNECTION_PARAMETERS=postgres://task_service_user:task_service_password@127.0.0.1/task_service_database?sslmode=disable&timezone=UTC
      DATABASE_MIGRATION_PATH=file://migrations
      ENVIRONMENT=test
    ```
  - You can run `go test` on the directory on one file/package at a time  

### Endpoints
This application has just one set of HTTPS endpoints

`POST /tasks`
  - Parameters:
    - URL: This endpoint will not acknowledge URL encoded parameters
    - Body: This endpoint expects a request with the following format where:
      - `name`: A string which represents the name of the task, it must be present
      - `details`: A string which represents the details of the task
      - `resolved_at`: A string which represents the resolution date of the task (RFC3339)
      - Example:    
        ```
        {
          "name": "Create an example task",
          "details": "Here is an example task",
          "resolved_at": "2019-01-01T00:00:01+00:00"
        }
        ```
  - Exceptions:
    - StatusBadRequest: If the request body is malformed or cannot be parsed the application will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 400 
    - Internal Error: If the endpoint hits a critical error while encoding the results, connecting to providers, or infrastructure it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 500
    - Unprocessable Entry Error: If the endpoint is unable to validate or sanitize the provided data it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 422
  - Return:
    - If no errors are encountered the endpoint will return a JSON encoded Task item and a status 200
      - `id`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `name`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `details`: A string which represents the details of the task
      - `resolved_at`: A string which represents the resolution date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - `created_at`: A string which represents the create date of the task (RFC3339) It will always be present (NOTE: All timestamps will be within the UTC timezone)
      - `updated_at`: A string which represents the create date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - Example:            
        ```
          {
            "id": 1,
            "name": "Create an example task",
            "details": "Here is an example task",
            "resolved_at": "2019-01-01T00:00:01Z",
            "created_at": "2019-03-25T13:49:03.171049643Z"
          }
        ```
  - A usable example can also be found in this repository in  `./scripts/api/task_create.sh`
  
`DELETE /tasks/{id}`
  - Parameters:
    - URL: This endpoint expects an ID of a valid Task in the system
    - Body: This endpoint will not acknowledge body parameters
  - Exceptions:
    - BadPathParameterErr: If the url encoded ID is malformed or cannot be parsed the application will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 400 
    - Internal Error: If the endpoint hits a critical error while encoding the results, connecting to providers, or infrastructure it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 500
    - Not Found: If the endpoint is unable to find a valid record based on the provided data it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 404
  - Return:
    - If no errors are encountered the endpoint will return a JSON encoded Task item and a status 200
      - `id`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `name`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `details`: A string which represents the details of the task
      - `resolved_at`: A string which represents the resolution date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - `created_at`: A string which represents the create date of the task (RFC3339) It will always be present (NOTE: All timestamps will be within the UTC timezone)
      - `updated_at`: A string which represents the create date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - Example:           
        ```
          {
            "id": 1,
            "name": "Deleted example task",
            "details": "Here is an example task",
            "resolved_at": "2019-01-01T00:00:01Z",
            "created_at": "2019-03-25T13:49:03.171049643Z"
          }
          ```
  
`GET /tasks`
  - Parameters:
    - URL: This endpoint will not acknowledge URL encoded parameters
    - Body: This endpoint expects a request with the following format where:
      - `limit`: An integer which represents a maximum number of items to fetch, this value must be present
      - `offset`: An integer which represents the offset on a limited amount of items, this value must be present
      - Example:     
        ```
        {
          "limit": 100,
          "offset": 0,
        }
        ```
  - Exceptions:
    - StatusBadRequest: If the request body is malformed or cannot be parsed the application will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 400 
    - Internal Error: If the endpoint hits a critical error while encoding the results, connecting to providers, or infrastructure it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 500
    - Not Found: If the endpoint is unable to find a valid records based on the provided data it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 404
  - Return:
    - If no errors are encountered the endpoint will return a JSON encoded collection of Task items and a status 200
      - `id`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `name`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `details`: A string which represents the details of the task
      - `resolved_at`: A string which represents the resolution date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - `created_at`: A string which represents the create date of the task (RFC3339) It will always be present (NOTE: All timestamps will be within the UTC timezone)
      - `updated_at`: A string which represents the create date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - Example:            
        ```
          {
            "tasks":[          
              {
                "id": 1,
                "name": "Create an example task",
                "details": "Here is an example task",
                "resolved_at": "2019-01-01T00:00:01Z",
                "created_at": "2019-03-25T13:49:03.171049643Z"
              }
            }
          }
        ```
  - A usable example can also be found in this repository in  `./scripts/api/task_create.sh`
  
`GET /tasks/{id}`
  - Parameters:
    - URL: This endpoint expects an ID of a valid Task in the system
    - Body: This endpoint will not acknowledge body parameters
  - Exceptions:
    - BadPathParameterErr: If the url encoded ID is malformed or cannot be parsed the application will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 400 
    - Internal Error: If the endpoint hits a critical error while encoding the results, connecting to providers, or infrastructure it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 500
    - Not Found: If the endpoint is unable to find a valid record based on the provided data it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 404
  - Return:
    - If no errors are encountered the endpoint will return a JSON encoded Task item and a status 200
      - `id`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `name`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `details`: A string which represents the details of the task
      - `resolved_at`: A string which represents the resolution date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - `created_at`: A string which represents the create date of the task (RFC3339) It will always be present (NOTE: All timestamps will be within the UTC timezone)
      - `updated_at`: A string which represents the create date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - Example:   
        ```
          {
            "id": 1,
            "name": "Read example task",
            "details": "Here is an example task",
            "resolved_at": "2019-01-01T00:00:01Z",
            "created_at": "2019-03-25T13:49:03.171049643Z"
          }
        ```
          
`PUT /tasks/{id}`
  - Parameters:
    - URL: his endpoint expects an ID of a valid Task in the system
    - Body: This endpoint expects a request with the following format where:
      - `name`: A string which represents the name of the task, it must be present
      - `details`: A string which represents the details of the task
      - `resolved_at`: A string which represents the resolution date of the task (RFC3339)
      - Example: 
      ```
        {
          "id": 1
          "name": "Update an example task",
          "details": "Here is an example task",
          "resolved_at": "2019-01-01T00:00:01+00:00"
        }
    ```
  - Exceptions:
    - StatusBadRequest: If the request body and url encoded id do not match the application will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 406
    - StatusBadRequest: If the request body is malformed or cannot be parsed the application will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 400 
    - Internal Error: If the endpoint hits a critical error while encoding the results, connecting to providers, or infrastructure it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 500
    - Unprocessable Entry Error: If the endpoint is unable to validate or sanitize the provided data it will return a [JSON API encoded exception](https://jsonapi.org/format/) and response code of 422
  - Return:
    - If no errors are encountered the endpoint will return a JSON encoded Task item and a status 200
      - `id`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `name`: An unsigned integer which represents the unique ID of the new record, it will always be present
      - `details`: A string which represents the details of the task
      - `resolved_at`: A string which represents the resolution date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - `created_at`: A string which represents the create date of the task (RFC3339) It will always be present (NOTE: All timestamps will be within the UTC timezone)
      - `updated_at`: A string which represents the create date of the task (RFC3339) (NOTE: All timestamps will be within the UTC timezone)
      - Example:  
        ```
          {
            "id": 1,
            "name": "Create an example task",
            "details": "Here is an example task",
            "resolved_at": "2019-01-01T00:00:01Z",
            "created_at": "2019-03-25T13:49:03.171049643Z"
            "updated_at": "2019-03-25T13:49:03.171049643Z"
          }
          ```
          
  ### Current deployment
  This API is currently deployed at: `https://me78vc7i2c.execute-api.us-west-2.amazonaws.com/production`   
            