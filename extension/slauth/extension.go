package slauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"
)

var (
	errNotAuthenticated       = errors.New("authentication didn't succeed")
	errNoAuthenticationheader = errors.New("authentication header not found")
)

type SlauthExtension struct {
	cfg       *ExtensionConfig
	logger    *zap.Logger
	authCache *AuthCache
}

func NewSlauthExtension() *SlauthExtension {
	return &SlauthExtension{}
}
func (e *SlauthExtension) init(cfg *ExtensionConfig, logger *zap.Logger) error {
	if err := cfg.Validate(); err != nil {
		return err
	}
	e.cfg = cfg
	authCache := NewAuthCache()

	if err := authCache.init(cfg); err != nil {
		return err
	}
	e.authCache = authCache
	e.logger = logger
	return nil
}
func (e *SlauthExtension) start(context.Context, component.Host) error {
	if e.cfg.LogFileEnabled {
		logFileLocation, _ := os.OpenFile("log-archive.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
		log.SetOutput(logFileLocation)
	}
	return nil
}

func (e *SlauthExtension) findAuthKey(headers map[string][]string) (string, error) {
	if len(headers["Authorization"]) > 0 {
		return "Authorization", nil
	}

	if len(headers["authorization"]) > 0 {
		return "authorization", nil
	}

	return "", errNoAuthenticationheader
}

func (e *SlauthExtension) findAgentInstanceHeader(headers map[string][]string) string {
	if len(headers["X-Agent-Instance-Id"]) > 0 {
		return "X-Agent-Instance-Id"
	}

	if len(headers["x-agent-instance-id"]) > 0 {
		return "x-agent-instance-id"
	}

	return ""
}

func (e *SlauthExtension) findIgnoreSpanHeader(headers map[string][]string) string {
	if len(headers["x-span-ignore"]) > 0 {
		return "x-span-ignore"
	}

	return ""
}

func (e *SlauthExtension) validateToken(token string) bool {
	return strings.HasPrefix(token, "Bearer ")
}

func (e *SlauthExtension) logHeaders(headers map[string][]string) {
	for key, value := range headers {
		if e.cfg.LogFileEnabled {
			log.Println(fmt.Sprintf("All - Headers: %s, value: %s", key, value))
		}
	}
}

func (e *SlauthExtension) getAuthResponse(token string) (authResponse *AuthenticationResponse, err error) {
	client := &http.Client{}
	authResponse = NewAuthenticationResponse()

	req, err := http.NewRequest(e.cfg.IssuerUrlVerb, e.cfg.IssuerURL, nil)
	if err != nil {
		return authResponse, err
	}

	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		return authResponse, err
	}

	authResponse.ResponseInfo.StatusCode = resp.StatusCode
	authResponse.ResponseInfo.Body, err = ioutil.ReadAll(resp.Body)

	if e.cfg.LogFileEnabled {
		e.logger.Info(fmt.Sprintf("AUTH STATUS CODE. Auth server status code: %+v", authResponse.ResponseInfo.StatusCode))
		body := string(authResponse.ResponseInfo.Body[:])
		e.logger.Info(fmt.Sprintf("AUTH RESPONSE BODY: %s", body))
	}

	if err != nil {
		return authResponse, err
	}

	defer func() {
		cerr := resp.Body.Close()
		if err == nil {
			err = cerr
		}
	}()

	if err := json.Unmarshal(authResponse.ResponseInfo.Body, &authResponse); err != nil {
		return authResponse, err
	}

	return authResponse, nil
}

func (e *SlauthExtension) httpAuthenticate(token string) (*AuthenticationResponse, error) {
	if authCacheData, found := e.authCache.get(token); found {
		authenticationResponse := NewAuthenticationResponse()
		authenticationResponse.SetCustomerId(authCacheData.customerId)

		return authenticationResponse, nil
	}

	authenticationResponse, err := e.getAuthResponse(token)

	if authenticationResponse.ResponseInfo.StatusCode != http.StatusOK {
		authError := authenticationResponse.constructError()
		return nil, authError
	}

	var expirationAfter time.Duration

	if authenticationResponse.ValidUntil == 0 {
		expirationAfter = e.cfg.TimeExpirationMin
	} else {
		expirationAfter = time.Unix(authenticationResponse.ValidUntil/1000, 0).Sub(time.Now())
	}

	if _, err = e.authCache.add(NewAuthCacheData().
		setToken(token).
		setCustomerId(authenticationResponse.CustomerId),
		expirationAfter); err != nil {
		return nil, err
	}

	return authenticationResponse, nil
}

func (e *SlauthExtension) authenticate(ctx context.Context, headers map[string][]string) (context.Context, error) {
	e.logHeaders(headers)

	authKey, err := e.findAuthKey(headers)

	if err != nil {
		e.logHeaders(headers)
		return nil, errNotAuthenticated
	}

	if e.cfg.LogFileEnabled {
		log.Println(fmt.Sprintf("Auth headers: Key: %s , Value: %s", authKey, headers[authKey]))
	}

	token := headers[authKey][0]

	isValid := e.validateToken(token)

	if !isValid {
		e.logger.Error(fmt.Sprint(fmt.Errorf("invalid token format")))
		return nil, errNotAuthenticated
	}

	authenticationResponse, err := e.httpAuthenticate(token)

	if err != nil {
		e.logger.Error(fmt.Sprint(err))
		e.logger.Error(fmt.Sprintf("Authentication response: %+v", authenticationResponse))
		return nil, errNotAuthenticated
	}

	e.setCustomerId(&ctx, authenticationResponse.CustomerId)
	e.setAgentInstanceId(&ctx, headers, authenticationResponse.CustomerId)
	e.setIgnoreSpanFlagIfPresent(&ctx, headers)

	return ctx, nil
}

func (e *SlauthExtension) setCustomerId(ctx *context.Context, customerId string) {
	*ctx = context.WithValue(*ctx, e.cfg.ConsumerIdFieldName, customerId)
}

func (e *SlauthExtension) setAgentInstanceId(ctx *context.Context, headers map[string][]string, customerId string) {
	agentInstanceIdHeader := e.findAgentInstanceHeader(headers)

	if agentInstanceIdHeader == "" {
		return
	}

	//key := ContextKeyConsumerIdFieldName(e.cfg.HttpHeaderNameAgentId)
	//*ctx = context.WithValue(*ctx, key, headers[agentInstanceIdHeader][0])

	*ctx = context.WithValue(*ctx, e.cfg.HttpHeaderNameAgentId, headers[agentInstanceIdHeader][0])
}

func (e *SlauthExtension) setIgnoreSpanFlagIfPresent(ctx *context.Context, headers map[string][]string) {
	ignoreSpanHeader := e.findIgnoreSpanHeader(headers)

	if ignoreSpanHeader == "" {
		return
	}

	//key := ContextKeyIgnoreSpan("IgnoreSpan")
	*ctx = context.WithValue(*ctx, "IgnoreSpan", "true")
}
