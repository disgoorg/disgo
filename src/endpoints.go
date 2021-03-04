package src

const APIVersion = "8"

const (
	BaseEndpoint = "https://discord.com/"
	CDNEndpoint  = "https://cdn.discordapp.com/"

	EndpointAPI        = BaseEndpoint + "api/v/" + APIVersion + "/"
	EndpointGateway    = EndpointAPI + "gateway"
	EndpointGatewayBot = EndpointGateway + "/bot"
)
