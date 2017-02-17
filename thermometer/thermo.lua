temp = require("ds18b20")


-- Common variables
temp_pin = 3
led_pin = 4


-- Port setup
gpio.mode(led_pin, gpio.OUTPUT)
gpio.write(led_pin, gpio.LOW)
temp.setup(temp_pin)
temp.read()


-- Schedule temp data sending
function schedule_mqtt_send(interval)
    tmr.stop(0)
    tmr.alarm(0, interval, tmr.ALARM_AUTO, function()
        --gpio.write(led_pin, gpio.LOW)
        local val = temp.read()
        print("Temperature: " .. val)
        -- publish a message with QoS = 0, retain = 0
        mqtt:publish("sensors/temp/internal", val, 0, 0, function(conn) 
            print("MQTT data sent") 
            --gpio.write(led_pin, gpio.HIGH)
        end)
    end)
    print("Temperature update inretval set to: " .. interval)
end


-- Main entry point
function main()
    print("Configuring MQTT")
    --tls.cert.verify(true)
    mqtt = mqtt.Client("ThermometerThing", 120)
    
    mqtt:on("offline", function(conn) 
        print ("MQTT offline")
        tmr.stop(0)
    end)
    
    mqtt:on("message", function(conn, topic, data)
        if topic == "sensors/setting/light" then
            if data == "on" then
                gpio.write(led_pin, gpio.LOW)
                mqtt:publish("sensors/light/internal", "true", 0, 0)
            elseif data == "off" then
                gpio.write(led_pin, gpio.HIGH)
                mqtt:publish("sensors/light/internal", "false", 0, 0)
            end
        elseif topic == "sensors/setting/interval" then
            local val = tonumber(data)
            if val ~= nil then
                schedule_mqtt_send(val)
            end
        else
            print("Unknown MQTT message: " .. topic .. ":" .. data)
        end
    end)

    -- Enable SSL by 3rd parameter set to 1
    mqtt:connect("192.168.1.103", 1883, 0, 0,
        function(conn) 
            print("MQTT connected")
            -- subscribe topic with qos = 0
            mqtt:subscribe("sensors/setting/light", 0)
            mqtt:subscribe("sensors/setting/interval", 0)
    
            schedule_mqtt_send(30*1000)
        end,
        function(conn, reason) 
            print("MQTT NOT connected: " .. reason)
        end
    )
        
    print("Initialized")
    gpio.write(led_pin, gpio.HIGH)
end


-- Execute main
main()
