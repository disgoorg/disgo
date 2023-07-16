# Proxy

This example shows how to use disgo with a gateway-proxy & rest-proxy such as https://github.com/Gelbpunkt/gateway-proxy and https://github.com/twilight-rs/http-proxy

For configuring those proxies, please refer to their documentation.

## Environment Variables

```env
disgo_token=discord bot token

disgo_guild_id=guild id to register commmands to

disgo_gateway_url=url to your gateway proxy. example: ws://gateway-proxy:7878

disgo_rest_url=url to your rest proxy. example: http://rest-proxy:7979/api/v10
```