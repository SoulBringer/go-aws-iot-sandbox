#AWS IoT @ Go

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

####AWS IoT settings:
    arn:aws:iot:eu-central-1:648323319497:thing/ForecastThing
    a1x6d6ym1e2e7r.iot.eu-central-1.amazonaws.com
    
    arn:aws:iot:eu-central-1:648323319497:thing/ThermometerThing

####Weather forecast settings:
    https://home.openweathermap.org/api_keys
    78ac22a76a45e175dbb87e0fb0b38bd6
    http://api.openweathermap.org/data/2.5/weather?appid=78ac22a76a45e175dbb87e0fb0b38bd6&units=metric&q=Vinnitsya,ua
    
####MQTT topics:
- sensors/temp/external
- sensors/temp/internal
- sensors/light/internal
- sensors/setting/interval
- sensors/setting/light