package geetest

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/parnurzeal/gorequest"
)

const apiHost = "api.geetest.com"

// Geetest is used for captcha registration and validation
type Geetest struct {
	captchaID  string
	privateKey string

	registerURL, validateURL string

	registerTimeout, validateTimeout time.Duration

	pool chan *gorequest.SuperAgent
}

// New constructs and returns a Geetest
func New(captchaID, privateKey string, enableHTTPS bool, registerTimeout, validateTimeout time.Duration, poolSize int) *Geetest {
	scheme := "http"
	if enableHTTPS {
		scheme = "https"
	}
	apiServer := fmt.Sprintf("%s://%s", scheme, apiHost)

	return &Geetest{
		captchaID:  captchaID,
		privateKey: privateKey,

		registerURL: fmt.Sprintf("%s/register.php", apiServer),
		validateURL: fmt.Sprintf("%s/validate.php", apiServer),

		registerTimeout: registerTimeout,
		validateTimeout: validateTimeout,

		pool: make(chan *gorequest.SuperAgent, poolSize),
	}
}

// CaptchaID returns captchaID
func (g *Geetest) CaptchaID() string { return g.captchaID }

// Register returns challenge get from api server
func (g *Geetest) Register() (string, error) {
	agent := g.getSuperAgent()
	defer g.putSuperAgent(agent)

	query := struct {
		CaptchaID string `json:"gt"`
	}{g.captchaID}
	_, body, errs := agent.Get(g.registerURL).Query(query).Timeout(g.registerTimeout).End()
	if errs != nil {
		return "", &multierror.Error{Errors: errs}
	}

	return hexmd5(body + g.privateKey), nil
}

// Validate validates challenge
func (g *Geetest) Validate(challenge, validate, seccode string) (bool, error) {
	hash := g.privateKey + "geetest" + challenge
	if validate != hexmd5(hash) {
		return false, nil
	}

	agent := g.getSuperAgent()
	defer g.putSuperAgent(agent)

	query := struct {
		Seccode string `json:"seccode"`
	}{seccode}
	_, body, errs := agent.Post(g.validateURL).Query(query).Timeout(g.validateTimeout).End()
	if errs != nil {
		return false, &multierror.Error{Errors: errs}
	}

	return hexmd5(seccode) == body, nil
}

func (g *Geetest) getSuperAgent() *gorequest.SuperAgent {
	var agent *gorequest.SuperAgent
	select {
	case agent = <-g.pool:
	default:
		agent = gorequest.New()
	}
	return agent
}

func (g *Geetest) putSuperAgent(agent *gorequest.SuperAgent) {
	select {
	case g.pool <- agent:
	default:
	}
}

func hexmd5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
