package models

// events
const (
	EventHTTPRequest                  = "http request"
	EventCreateWithdrawals            = "create withdrawals"
	EventProcessWithdrawals           = "process withdrawals"
	EventLogBalanceAndAddress         = "log balance and address"
	EventValidateCaptcha              = "validate captcha"
	EventRegisterCaptcha              = "register captcha"
	EventReward                       = "reward"
	EventGetGeoFromIP                 = "get geo from ip"
	EventUserSignup                   = "user signup"
	EventOfferwowCallback             = "offerwow callback"
	EventOfferwowInvalidSignature     = "offerwow invalid signature"
	EventSuperrewardsCallback         = "superrewards callback"
	EventSuperrewardsInvalidSignature = "superrewards invalid signature"
	EventPTCWallCallback              = "ptcwall callback"
	EventClixwallCallback             = "clixwall callback"
	EventPersonalyCallback            = "personaly callback"
	EventPersonalyInvalidSignature    = "personaly invalid signature"
	EventKiwiwallCallback             = "kiwiwall callback"
	EventKiwiwallInvalidSignature     = "kiwiwall invalid signature"
	EventTrialpayCallback             = "trialpay callback"
	EventTrialpayInvalidSignature     = "trialpay invalid signature"
	EventAdscendMediaCallback         = "adscend media callback"
)
