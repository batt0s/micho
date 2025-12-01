package helm

import (
	"fmt"
	"log"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

var settings = cli.New()

type ReleaseConfig struct {
	Namespace   string
	ReleaseName string
	ChartPath   string
	Values      map[string]interface{}
}

func debug(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func getActionConfig(namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secret", debug)
	return actionConfig, err
}

// Install the Chart
func InstallRelease(cfg ReleaseConfig) error {
	actionConfig, err := getActionConfig(cfg.Namespace)
	if err != nil {
		return err
	}

	client := action.NewInstall(actionConfig)
	client.ReleaseName = cfg.ReleaseName
	client.Namespace = cfg.Namespace

	client.CreateNamespace = true

	chart, err := loader.Load(cfg.ChartPath)
	if err != nil {
		return fmt.Errorf("could not load chart %s: %w", cfg.ChartPath, err)
	}

	_, err = client.Run(chart, cfg.Values)
	if err != nil {
		return fmt.Errorf("deployment error: %w", err)
	}

	return nil
}

func UpgradeRelease(cfg ReleaseConfig) error {
	actionConfig, err := getActionConfig(cfg.Namespace)
	if err != nil {
		return err
	}

	client := action.NewUpgrade(actionConfig)
	client.Namespace = cfg.Namespace

	chart, err := loader.Load(cfg.ChartPath)
	if err != nil {
		return fmt.Errorf("could not load chart %s: %w", cfg.ChartPath, err)
	}

	_, err = client.Run(cfg.ReleaseName, chart, cfg.Values)
	if err != nil {
		return fmt.Errorf("deployment error: %w", err)
	}

	return nil
}

func UninstallRelease(namespace, releaseName string) error {
	actionConfig, err := getActionConfig(namespace)
	if err != nil {
		return err
	}

	client := action.NewUninstall(actionConfig)

	_, err = client.Run(releaseName)
	return err
}

func GetReleaseStatus(namespace, releaseName string) (string, error) {
	actionConfig, err := getActionConfig(namespace)
	if err != nil {
		return "", err
	}
	client := action.NewStatus(actionConfig)
	res, err := client.Run(releaseName)
	return res.Info.Status.String(), err
}
