curl -X POST -H "x-ha-access: TelephonePatConsideration" \
       -H "Content-Type: application/json" \
       -d '{"entity_id": "media_player.apple_tv"}' \
       http://192.168.4.1:8123/api/services/media_player/media_pause