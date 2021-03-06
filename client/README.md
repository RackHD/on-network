[//]: # ( Copyright 2017, Dell EMC, Inc.)

# Go API client for swagger

No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)

## Overview
This API client was generated by the [swagger-codegen](https://github.com/swagger-api/swagger-codegen) project.  By using the [swagger-spec](https://github.com/swagger-api/swagger-spec) from a remote server, you can easily generate an API client.

- API version: 0.0.1
- Package version: 1.0.0
- Build package: io.swagger.codegen.languages.GoClientCodegen

## Installation
Put the package under your project folder and add the following in import:
```
    "./swagger"
```

## Documentation for API Endpoints

All URIs are relative to *http://localhost*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AuthApi* | [**LoginPost**](docs/AuthApi.md#loginpost) | **Post** /login | 
*AboutApi* | [**AboutGet**](docs/AboutApi.md#aboutget) | **Get** /about | Get about
*SwitchConfigApi* | [**SwitchConfig**](docs/SwitchConfigApi.md#switchconfig) | **Post** /switchConfig | Get switch running config
*SwitchFirmwareApi* | [**SwitchFirmware**](docs/SwitchFirmwareApi.md#switchfirmware) | **Post** /switchFirmware | Get switch Firmware Version
*SwitchVersionApi* | [**SwitchVersion**](docs/SwitchVersionApi.md#switchversion) | **Post** /switchVersion | Get switch Firmware Version
*UpdateSwitchApi* | [**UpdateSwitch**](docs/UpdateSwitchApi.md#updateswitch) | **Post** /updateSwitch | Update switch firmware


## Documentation For Models

 - [About](docs/About.md)
 - [ErrorResponse](docs/ErrorResponse.md)
 - [Login](docs/Login.md)
 - [LoginError](docs/LoginError.md)
 - [ModelSwitch](docs/ModelSwitch.md)
 - [Status](docs/Status.md)
 - [SwitchConfigResponse](docs/SwitchConfigResponse.md)
 - [SwitchEndpoint](docs/SwitchEndpoint.md)
 - [SwitchVersionResponse](docs/SwitchVersionResponse.md)
 - [Token](docs/Token.md)
 - [UpdateSwitch](docs/UpdateSwitch.md)


## Documentation For Authorization


## Bearer

- **Type**: API key 
- **API key parameter name**: authorization
- **Location**: HTTP header


## Author



