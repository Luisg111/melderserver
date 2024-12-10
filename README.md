# tiny "smart home"

## features
### pager integration
This project integrates the pager issued by my volunteer fire departement into my tiny "smart home".
When an alarm is triggered, a telegram message is sent to my smartphone so I get notified even when im not carrying my pager.

Additionally, the information is relayed to a ESP8266 microcontroller connected via usb serial.
It turns on the light in my room (ceiling fan with remote control) and turns then off again after 5 minutes.

the repo for the ESP8266 program can be found here: [Pager Project](https://github.com/Luisg111/melder_esp8266)

### garage door status indication
an additional feature is a small red LED that illuminates when our garage door is open.
The garage door cannot be seen from the house so we sometimes forgot to close it at night.

The status of the door is checked by an ESP8266 microcontroller and sent via HTTP over WiFi.
You can find the repo for the ESP8266 program here: [Garage Door Status Project](https://github.com/Luisg111/garagentor_status)

### weather warnings
current weather warnings for my hometown are loaded periodically from the [German Meteorological Service (DWD)](https://www.dwd.de/DE/Home/home_node.html) which offers a download option as json.

A telegram bot allows them to be send to my smartphone when I write a message to the bot.
The warnings will also be included in the telegram message sent when an alarm is triggered by my pager.

### remote control for light and fan
using the telegram bot that sends the alarm messages, the light and ceiling fan in my room can be remote controlled by my smartphone by sending specific messages to the bot.

## usage
this project uses three environment variables for credentials:
- ALARM_BOT_TOKEN for the alarm telegram bot
- WEATHER_BOT_TOKEN for the weather telegram bot
- ALARM_CHAT_ID for the chat where the alarm message should be posted