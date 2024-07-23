package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"weather_monitor/models"
	"weather_monitor/utils"
)

type OpenWeatherClient interface {
	GetCurrentWeather(ctx context.Context, latitude, longitude string) (*models.OpenWeatherResponse, *models.AppError)
}

type DefaultOpenWeatherClient struct {
	httpClient *http.Client
	env        *utils.Env
}

func NewDefaultOpenWeatherClient(httpClient *http.Client, env *utils.Env) OpenWeatherClient {
	return &DefaultOpenWeatherClient{httpClient: httpClient, env: env}
}

func (d *DefaultOpenWeatherClient) GetCurrentWeather(ctx context.Context, latitude, longitude string) (*models.OpenWeatherResponse, *models.AppError) {
	httpReq, httpReqErr := http.NewRequestWithContext(ctx, "GET", d.env.OpenWeatherUrl, nil)
	if httpReqErr != nil {
		logrus.Info("Unable to create http request for open weather.")
		return nil, &models.AppError{ErrorMessage: "Unable to create http request for open weather."}
	}
	query := httpReq.URL.Query()

	query.Add("lat", latitude)
	query.Add("lon", longitude)
	query.Add("appid", d.env.ApiKey)

	httpReq.URL.RawQuery = query.Encode()

	logrus.Infof("URL:%v", httpReq.URL.String())

	httpRes, httpResErr := d.httpClient.Do(httpReq)
	if httpResErr != nil {
		logrus.Info("Unable to send http request for open weather.")
		return nil, &models.AppError{ErrorMessage: "Unable to send http request for open weather."}
	}
	if httpRes.StatusCode != http.StatusOK {
		logrus.Info(fmt.Sprintf("http call failed with status:%v", httpRes.StatusCode))
		return nil, &models.AppError{ErrorMessage: fmt.Sprintf("http call failed with status:%v", httpRes.StatusCode)}
	}

	defer httpRes.Body.Close()
	resByte, err := io.ReadAll(httpRes.Body)
	if err != nil {
		logrus.Info(fmt.Sprintf("unable to marshall response:%v", err))
		return nil, &models.AppError{ErrorMessage: fmt.Sprintf("unable to marshall response:%v", err)}
	}

	var openWeatherResponse models.OpenWeatherResponse
	marshallErr := json.Unmarshal(resByte, &openWeatherResponse)
	if marshallErr != nil {
		logrus.Error(marshallErr.Error())
		return nil, &models.AppError{ErrorMessage: fmt.Sprintf("unable to marshall response:%v", marshallErr)}
	}

	return &openWeatherResponse, nil

}
