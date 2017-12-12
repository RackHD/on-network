[//]: # ( Copyright 2017, Dell EMC, Inc.)


# on-network 

Copyright © 2017 Dell Inc. or its subsidiaries.  All Rights Reserved. 

`on-network` is the network service for RackHD. This service performs a variety of network configuration functions, including the following: 

* Gets the running configuration for a switch. 
* Retrieves the short version for the firmware currently running on a Cisco Nexus 3132, Cisco Nexus 3164, Cisco Nexus 3172, Cisco Nexus 9332, or Cisco Nexus 93180 switch. 
* Retrieves the full version for the firmware currently running on a Cisco Nexus 3132, Cisco Nexus 3164, Cisco Nexus 3172, Cisco Nexus 9332, or Cisco Nexus 93180 switch. 
* Updates switch firmware based on a specified switch type and firmware image. For Cisco 3164, Cisco 9332, and Cisco 93180 switches, the update is performed by a non-disruptive "install all" command. For Cisco 3132 and Cisco 3172 switches, the update is performed by a disruptive "install all" command. 

These functions are exposed through a REST API. Before making these calls, a client must first use the `/login` REST call to authenticate to the service. The `/login` call returns a token that must be passed in the header for each subsequent call.

The `on-network` service handles network configuration requests in the following manner:

1. A Symphony PAQX (such as RCM Fitness or Dell Node Expansion) initiates a network configuration request. 
2. The Symphony RackHD Adapter passes along the request to the RackHD `on-http` service.
3. The `on-http` service passes the request along to the `on-taskgraph` service, which executes a workflow to satisfy the request.
4. The workflow initiated by the `on-taskgraph` service makes a REST API call to the `on-network` service.
5. The `on-network` service communicates directly with the Cisco switch through the NX-API to perform the network configuration.

The `on-network` service is written in Golang.

## Before you begin

Verify that the following tools are installed:

* Go Programming Language release 1.8 (go1.8) or higher
* Docker 1.12+
* Docker Compose 1.8.0+
* Java Development Kit (version 8)

## Building and running

You can use any of the following techniques to build the project and run the service:

* Use the go command to run the main.go file
* Use the make command to build a binary
* Use Docker to build a Docker image

If you want to build the project by using the go command or the make command, you need to set some environment variables first. Here are the environment variables and some suggested values:

```
SERVICE_USERNAME=<username>;                                       
SERVICE_PASSWORD=<password>;
CISCO_BOOT_TIME_IN_SECONDS=20;                            // time it takes for the switch to reboot after firmware update
CISCO_INSTALL_TIME_IN_MINUTES=4;                           // time it takes to install firmware
CISCO_RECONNECTION_TIMEOUT_IN_SECONDS=240;  //time it takes  for cisco to reboot 
SWITCH_MODELS_FILE_PATH=switchModels.yml           //this  is the file path which specifies supported switches and whether they are distuptive/non-disruptive

```
To set the environment variables on Linux, use the export command, as shown in the following example:
```
export SERVICE_USERNAME=<username>
export SERVICE_PASSWORD=<password>
export CISCO_BOOT_TIME_IN_SECONDS=20
export CISCO_INSTALL_TIME_IN_MINUTES=4
export CISCO_RECONNECTION_TIMEOUT_IN_SECONDS=240
export SWITCH_MODELS_FILE_PATH=switchModels.yml
```

### Using the go command to run the main.go file

To build the project and run the main.go file:

1. Set the required environment variables.

2. Run the make command to resolve dependencies:

``` 
make link clean deps 
```


3. Run the main.go to start the service. For example:

``` 
go run cmd/on-network-server-impl/main.go --port 8081 --host 0.0.0.0 --write-timeout 10m
```

 The parameters for main.go are described below:

| Parameter | Description |
| --- | --- |
| `port` | Specifies the port on which the service is running. | 
| `host` | Specifies the IP address on which the service is listening. You can use 0.0.0.0 for any IP. |
| `write-timeout` | Defines the maximum duration before the write of the response will time out. You would typically want to set this parameter to 10 minutes because a switch update could take around 7 minutes. |



### Using the make command to build a binary

To use the make command to build a binary:

1. Set the required environment variables.

2. Execute the make command that builds a binary for your operating system:

On Linux:

``` 
make link clean deps linux
```

On a Mac:

``` 
make link clean deps darwin 
```

On Windows:

``` 
make link clean deps windows 
```

3. Start the binary for your operating system. For example:

```
cmd/on-network-server/on-network-linux-amd64  --port 8081 --host 0.0.0.0 --write-timeout 10m 
```


### Using Docker to build a Docker image

To use Docker to build and run a Docker image:

1. Run the make command with the Docker target to build the image:

``` 
make docker
```

2. Execute the Docker run command to start the image:

```
sudo docker run -i -p 8081:8080 -e "SERVICE_USERNAME=<username>" -e "SERVICE_PASSWORD=<password>" -e "SWITCH_MODELS_FILE_PATH=switchModels.yml" -e "CISCO_BOOT_TIME_IN_SECONDS=20" -e "CISCO_RECONNECTION_TIMEOUT_IN_SECONDS=240" -e "CISCO_INSTALL_TIME_IN_MINUTES=4" on-network:dev
```
## Adding on-network to RackHD

To add on-network to RackHD: 

1. Edit `/opt/monorail/config.json` on a RackHD deployment.

2. Add the on-network configuration, as shown in the following example:

```
{
...
   "onNetwork": {
       "url": "http://localhost:8081",
       "username": "<username>",
       "password": "<password>"
    }

...
}
```

3. Restart the RackHD service.

**Note:** When RackHD is deployed within Symphony, you do not need to add the on-network configuration to the `config.json` file.

## Generating the server code with go-swagger

To generate the server code: 

1. Install go-swagger:

```
go get -u github.com/go-swagger/go-swagger/cmd/swagger
```

2. Generate the code:
```
swagger generate server -A "on-network"
```

## Generating the Swagger API documentation with swagger-codegen

To generate the Swagger API documentation:

1. Download swagger-codegen:

```
wget http://central.maven.org/maven2/io/swagger/swagger-codegen-cli/2.2.3/swagger-codegen-cli-2.2.3.jar -O swagger-codegen-cli.jar
```

2. Generate the docs and place the output in the /client folder:

```
java -jar swagger-codegen-cli.jar generate -i swagger.yaml -l go -o client/
```

For more information on the Go API client for Swagger, see https://github.com/RackHD/on-network/tree/master/client. This folder contains an automatically generated README.md file that provides an overview of the supported API calls as well as links to additional markdown files that contain details on the API paths and model files. 


## Starting the Swagger user interface 

To start the Swagger UI, open this page in a browser: https://localhost:8080/docs

## Sample API calls

Here are some examples that show how to make API calls to the `on-network` service. 

### Logging in to the service

#### Path
```
POST localhost:8081/login
```

#### Header
```
Content-Type: application/json
```


#### Body 

```
{
     "username": "<username>",
     "password": "<password>"
}
```

#### Return

```
{
    "token": "<token>"
}
```

### Retrieving the switch version

#### Path

```
POST localhost:8081/switchVersion
```

#### Header

```
Content-Type: application/json
Authorization: "Bearer <token-from-login>"
```


#### Body

```
{
    "endpoint": {
        "ip": "<ip-address>",
        "username": "<username>",
        "password": "<password>",
        "switchType": "cisco"    
  }
}
```

### Updating the switch firmware

#### Path

```
POST localhost:8081/updateSwitch

```

#### Header

```
Content-Type: application/json
Authorization: "Bearer <token-from-login>"
```

#### Body

```
{   
	"endpoint": {
	"ip":"<ip-address>",
    	"username":"<username>",
    	"password":"<password>",
    	"switchType":"cisco"
	 },
	  "imageURL":"http://<ip-address>:8080/nxos.7.0.3.I5.2.bin",
	  "switchModel": "93180"
}
```

## Licensing

Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the License. You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

RackHD is a Trademark of Dell EMC.

