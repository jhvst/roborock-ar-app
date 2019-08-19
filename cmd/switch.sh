curl -X POST -H "x-ha-access: TelephonePatConsideration" \
       -H "Content-Type: application/json" \
       -d '{"entity_id": "light.office"}' \
       http://192.168.4.1:8123/api/services/light/turn_off