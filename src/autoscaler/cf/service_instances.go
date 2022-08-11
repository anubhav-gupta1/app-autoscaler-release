package cf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"code.cloudfoundry.org/lager"
)

const (
	ServiceInstancesPath = "/v2/service_instances"
	ServicePlansPath     = "v2/service_plans"
	ResultsPerPageParam  = "results-per-page"
)

type (
	ServiceInstance struct {
		Guid string `json:"guid"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
)

/*GetServiceInstancesInOrg
 * get service instances in org uses
 * https://v3-apidocs.cloudfoundry.org/version/3.122.0/index.html#list-service-instances
 */
func (c *Client) GetServiceInstancesInOrg(orgGUID, brokerServicePlanGuid string) (int, error) {
	parameters := url.Values{}
	parameters.Add("organization_guids", orgGUID)
	parameters.Add("service_plan_guids", brokerServicePlanGuid)

	theUrl := fmt.Sprintf("/v3/service_instances?%s", parameters.Encode())
	instances, err := PagedResourceRetriever[ServiceInstance]{c}.GetAllPages(theUrl)
	if err != nil {
		return 0, fmt.Errorf("failed GetServiceInstancesInOrg org(%s), servicePlan(%s): %w", orgGUID, brokerServicePlanGuid, err)
	}
	return len(instances), err
}

type Metadata struct {
	Guid string `json:"guid"`
}

type Resource struct {
	Metadata Metadata `json:"metadata"`
}

type Result struct {
	TotalResults int        `json:"total_results"`
	Resources    []Resource `json:"resources"`
}

type ServiceInstanceEntity struct {
	ServicePlanGuid string `json:"service_plan_guid"`
}

type ServiceInstanceResource struct {
	Entity ServiceInstanceEntity `json:"entity"`
}

type ServicePlanEntity struct {
	UniqueId string `json:"unique_id"`
}

type ServicePlanResource struct {
	Entity ServicePlanEntity `json:"entity"`
}

//TODO we need to call this first before
func (c *Client) GetCCServicePlanGuid(brokerPlanGuid string) (string, error) {
	logger := c.logger.Session("cf-client-get-service-plan-guid", lager.Data{"brokerPlanGuid": brokerPlanGuid})
	logger.Debug("start")
	defer logger.Debug("end")

	c.planMapsLock.Lock()
	defer c.planMapsLock.Unlock()

	if g, ok := c.brokerPlanGuidToCCServicePlanGuid[brokerPlanGuid]; ok {
		return g, nil
	}

	servicePlansUrl, err := url.Parse(c.conf.API)
	if err != nil {
		return "", fmt.Errorf("cf-client-get-service-plan-guid: failed to parse CF API URL: %w", err)
	}
	servicePlansUrl.Path = servicePlansUrl.Path + ServicePlansPath

	parameters := url.Values{}
	parameters.Add("q", "unique_id:"+brokerPlanGuid)
	servicePlansUrl.RawQuery = parameters.Encode()

	logger.Info("created-url", lager.Data{"url": servicePlansUrl.String()})

	req, err := http.NewRequest("GET", servicePlansUrl.String(), nil)
	if err != nil {
		logger.Error("new-request", err)
		return "", fmt.Errorf("cf-client-get-service-plan-guid: failed to create request to CF API: %w", err)
	}
	tokens, _ := c.GetTokens()
	req.Header.Set("Authorization", TokenTypeBearer+" "+tokens.AccessToken)

	var resp *http.Response
	resp, err = c.httpClient.Do(req)

	if err != nil {
		logger.Error("do-request", err)
		return "", fmt.Errorf("cf-client-get-service-plan-guid: failed to execute request to CF API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("cf-client-get-service-plan-guid: failed to get service plan: %s [%d] %s", servicePlansUrl.String(), resp.StatusCode, resp.Status)
		logger.Error("get-response", err)
		return "", err
	}

	result := &Result{}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		logger.Error("decode", err)
		return "", fmt.Errorf("cf-client-get-service-plan-guid: failed to decode response from CF API: %w", err)
	}

	if result.TotalResults != 1 && len(result.Resources) != 1 {
		err = fmt.Errorf("cf-client-get-service-plan-guid: failed to find service plan: %s found %d plans", servicePlansUrl.String(), result.TotalResults)
		logger.Error("did-not-find-plan", err)
		return "", err
	}

	ccServicePlanGuid := result.Resources[0].Metadata.Guid
	logger.Info("found-guid", lager.Data{"ccServicePlanGuid": ccServicePlanGuid})
	c.brokerPlanGuidToCCServicePlanGuid[brokerPlanGuid] = ccServicePlanGuid
	c.ccServicePlanToBrokerPlanGuid[ccServicePlanGuid] = brokerPlanGuid

	return ccServicePlanGuid, nil
}

func (c *Client) getBrokerPlanGuid(ccServicePlanGuid string) (string, error) {
	logger := c.logger.Session("cf-client-get-broker-plan-guid", lager.Data{"ccServicePlanGuid": ccServicePlanGuid})
	logger.Debug("start")
	defer logger.Debug("end")

	c.planMapsLock.Lock()
	defer c.planMapsLock.Unlock()

	if g, ok := c.ccServicePlanToBrokerPlanGuid[ccServicePlanGuid]; ok {
		return g, nil
	}

	servicePlansUrl, err := url.Parse(c.conf.API)
	if err != nil {
		return "", fmt.Errorf("cf-client-get-broker-plan-guid: failed to parse CF API URL: %w", err)
	}
	servicePlansUrl.Path = servicePlansUrl.Path + path.Join(ServicePlansPath, ccServicePlanGuid)

	logger.Info("created-url", lager.Data{"url": servicePlansUrl.String()})

	req, err := http.NewRequest("GET", servicePlansUrl.String(), nil)
	if err != nil {
		logger.Error("new-request", err)
		return "", fmt.Errorf("cf-client-get-broker-plan-guid: failed to create request to CF API: %w", err)
	}

	tokens, _ := c.GetTokens()
	req.Header.Set("Authorization", TokenTypeBearer+" "+tokens.AccessToken)

	var resp *http.Response
	resp, err = c.httpClient.Do(req)

	if err != nil {
		logger.Error("do-request", err)
		return "", fmt.Errorf("cf-client-get-broker-plan-guid: failed to execute request to CF API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("cf-client-get-broker-plan-guid: failed to get service plan: %s [%d] %s", servicePlansUrl.String(), resp.StatusCode, resp.Status)
		logger.Error("get-response", err)
		return "", err
	}

	result := &ServicePlanResource{}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		logger.Error("decode", err)
		return "", fmt.Errorf("cf-client-get-broker-plan-guid: failed to decode response from CF API: %w", err)
	}

	brokerPlanGuid := result.Entity.UniqueId
	logger.Info("found-guid", lager.Data{"brokerPlanGuid": brokerPlanGuid})
	c.ccServicePlanToBrokerPlanGuid[ccServicePlanGuid] = brokerPlanGuid
	c.brokerPlanGuidToCCServicePlanGuid[brokerPlanGuid] = ccServicePlanGuid

	return brokerPlanGuid, nil
}

func (c *Client) GetServicePlan(serviceInstanceGuid string) (string, error) {
	logger := c.logger.Session("cf-client-get-service-plan", lager.Data{"serviceInstanceGuid": serviceInstanceGuid})
	logger.Debug("start")
	defer logger.Debug("end")

	c.instanceMapLock.Lock()
	defer c.instanceMapLock.Unlock()

	if g, ok := c.serviceInstanceGuidToBrokerPlanGuid[serviceInstanceGuid]; ok {
		return g, nil
	}

	serviceInstancesUrl, err := url.Parse(c.conf.API)
	if err != nil {
		return "", fmt.Errorf("cf-client-get-service-plan: failed to parse CF API URL: %w", err)
	}
	serviceInstancesUrl.Path = path.Join(ServiceInstancesPath, serviceInstanceGuid)

	logger.Debug("created-url", lager.Data{"url": serviceInstancesUrl.String()})

	req, err := http.NewRequest("GET", serviceInstancesUrl.String(), nil)
	if err != nil {
		logger.Error("new-request", err)
		return "", fmt.Errorf("cf-client-get-service-plan: failed to create request to CF API: %w", err)
	}

	tokens, _ := c.GetTokens()
	req.Header.Set("Authorization", TokenTypeBearer+" "+tokens.AccessToken)

	var resp *http.Response
	resp, err = c.httpClient.Do(req)

	if err != nil {
		logger.Error("do-request", err)
		return "", fmt.Errorf("cf-client-get-service-plan: failed to execute request to CF API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("cf-client-get-service-plan: failed to get service plan: %s [%d] %s", serviceInstancesUrl.String(), resp.StatusCode, resp.Status)
		logger.Error("get-response", err)
		return "", err
	}

	result := &ServiceInstanceResource{}

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		logger.Error("decode", err)
		return "", fmt.Errorf("cf-client-get-service-plan: failed to decode response from CF API: %w", err)
	}

	servicePlanGuid := result.Entity.ServicePlanGuid
	logger.Info("found-guid", lager.Data{"servicePlanGuid": servicePlanGuid})
	brokerPlanGuid, err := c.getBrokerPlanGuid(servicePlanGuid)
	if err != nil {
		logger.Error("cc-plan-to-broker-plan", err)
		return "", fmt.Errorf("cf-client-get-service-plan: failed to translate Cloud Controller service plan to broker service plan: %w", err)
	}

	c.serviceInstanceGuidToBrokerPlanGuid[serviceInstanceGuid] = brokerPlanGuid

	return c.serviceInstanceGuidToBrokerPlanGuid[serviceInstanceGuid], nil
}
