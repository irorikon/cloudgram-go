package utils

// verification rules
var (
	LoginVerify   = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}}
	ChannelVerify = Rules{"ChannelID": {NotEmpty()}, "ChannelName": {NotEmpty()}}
)
