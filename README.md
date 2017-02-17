#AWS IoT @ Go

####Project uses next stack of technologies:
- Golang
- LUA
- JS
- AWS IoT
- DynamoDB
- MQTT
- Mosquitto
- NodeMCU
- Angular

####TODO:
- ~~setup AWS IoT (things, certificates & policies)~~
- ~~implement & test AWS IoT MQTT interaction source~~
- ~~implement gateway with REST API & MQTT connectivity~~
- ~~implement forecast temperature sensor thing (retrieving local outdoor temp from weather forecast srv)~~
- ~~implement thermometer temperature sensor thing (retrieving local indoor temp from thermo sensor)
  (uses local net MQTT broker with bridge to AWS as far as NodeMCU does not support TLS)~~
- ~~implement UI (hosted on gateway side, simple Angular app)~~
- implement unit tests for Go functionality
- add DynamoDB support for logging data
    
####MQTT topics:
- sensors/temp/external
- sensors/temp/internal
- sensors/light/internal
- sensors/setting/interval
- sensors/setting/light