package cf

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"code.cloudfoundry.org/app-autoscaler/src/autoscaler/models"
)

const (
	TokenTypeBearer = "Bearer"
)

type (
	//App the app information from cf for full version look at https://v3-apidocs.cloudfoundry.org/version/3.122.0/index.html#apps
	App struct {
		Guid          string        `json:"guid"`
		Name          string        `json:"name"`
		State         string        `json:"state"`
		CreatedAt     time.Time     `json:"created_at"`
		UpdatedAt     time.Time     `json:"updated_at"`
		Relationships Relationships `json:"relationships"`
	}

	AppAndProcesses struct {
		App       *App
		Processes Processes
	}

	Relationships struct {
		Space *Space `json:"space"`
	}
	SpaceData struct {
		Guid SpaceId `json:"guid"`
	}
	Space struct {
		Data SpaceData `json:"data"`
	}

	//Processes the processes information for an App from cf for full version look at https://v3-apidocs.cloudfoundry.org/version/3.122.0/index.html#processes
	Process struct {
		Guid       string    `json:"guid"`
		Type       string    `json:"type"`
		Instances  int       `json:"instances"`
		MemoryInMb int       `json:"memory_in_mb"`
		DiskInMb   int       `json:"disk_in_mb"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}
	Processes []Process
)

func (p Processes) GetInstances() int {
	instances := 0
	for _, process := range p {
		instances += process.Instances
	}
	return instances
}

/*GetAppAndProcesses
 * A utility function that gets the app and processes for the app in one call in parallel
 */
func (c *Client) GetAppAndProcesses(appID string) (*AppAndProcesses, error) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var app *App
	var processes Processes
	var errApp, errProc error
	go func() {
		app, errApp = c.GetApp(appID)
		wg.Done()
	}()
	go func() {
		processes, errProc = c.GetAppProcesses(appID)
		wg.Done()
	}()
	wg.Wait()
	if errApp != nil {
		return nil, fmt.Errorf("get state&instances getApp failed: %w", errApp)
	}
	if errProc != nil {
		return nil, fmt.Errorf("get state&instances GetAppProcesses failed: %w", errProc)
	}
	return &AppAndProcesses{App: app, Processes: processes}, nil
}

/*GetApp
 * Get the information for a specific app
 * from the v3 api https://v3-apidocs.cloudfoundry.org/version/3.122.0/index.html#apps
 */
func (c *Client) GetApp(appID string) (*App, error) {
	url := fmt.Sprintf("%s/v3/apps/%s", c.conf.API, appID)

	resp, err := c.get(url)
	if err != nil {
		return nil, fmt.Errorf("failed getting app '%s': %w", appID, err)
	}
	defer func() { _ = resp.Body.Close() }()

	app := &App{}
	err = json.NewDecoder(resp.Body).Decode(app)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshalling app information for '%s': %w", appID, err)
	}
	return app, nil
}

/*GetAppProcesses
 * Get the processes information for a specific app
 * from the v3 api https://v3-apidocs.cloudfoundry.org/version/3.122.0/index.html#apps
 */
func (c *Client) GetAppProcesses(appID string) (Processes, error) {
	url := fmt.Sprintf("%s/v3/apps/%s/processes?per_page=%d", c.conf.API, appID, c.conf.PerPage)

	pages, err := ResourceRetriever[Process]{c}.getAllPages(url)
	if err != nil {
		return nil, fmt.Errorf("failed GetAppProcesses '%s': %w", appID, err)
	}
	return pages, nil
}

/*ScaleAppWebProcess
 * Scale the given application Web processes to the given amount
 * https://v3-apidocs.cloudfoundry.org/version/3.122.0/index.html#scale-a-process
 */
func (c *Client) ScaleAppWebProcess(appID string, num int) error {
	url := fmt.Sprintf("%s/v3/apps/%s/processes/web/actions/scale", c.conf.API, appID)
	type scaleApp struct {
		Instances int `json:"instances"`
	}
	_, err := c.post(url, scaleApp{Instances: num})
	if err != nil {
		return fmt.Errorf("failed scaling app '%s' to %d: %w", appID, num, err)
	}
	return err
}
