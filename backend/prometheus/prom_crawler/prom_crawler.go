package prom_crawler

import (
	"github.com/jasonlvhit/gocron"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/prometheus"
	"motherbear/backend/prometheus/goloop_prom_crawler"
	"motherbear/backend/prometheus/loopchain_prom_crawler"
)

type IPromCrawler interface {
	Crawler() (*prometheus.PrometheusData, error)
}

var scheduler *gocron.Scheduler

func BeginToCrawl() (*gocron.Scheduler, error) {
	var crawler IPromCrawler

	switch configuration.QueryNodeType() {
	case constants.NodeType1:
		crawler = loopchain_prom_crawler.LoopChainPrometheus{}
	case constants.NodeType2:
		crawler = goloop_prom_crawler.GoloopPrometheus{}
	default:
		err := isaacerror.SysErrFailToGetNodeType
		logger.Error(err.Error())
		return nil, err
	}

	job := func() {
		var prometheusData *prometheus.PrometheusData
		prometheusData, err := crawler.Crawler()
		if err != nil {
			logger.Error(isaacerror.SysErrFailToGetPrometheusData.Error())
			prometheusData = &prometheus.PrometheusData{
				TimeStamp: 0,
				Status:    prometheus.CrawlingFail,
			}
		} else {
			if prometheusData == nil {
				prometheusData = &prometheus.PrometheusData{
					TimeStamp: 0,
					Status:    prometheus.WarmingUpCrawling,
				}
				logger.Warn("Waiting for prometheus ready. ")
			} else {
				logger.Debug("Succeed to crawl the data from the prometheus.")
			}
		}

		prometheus.SetPrometheusData(prometheusData)
	}

	var interval uint64
	interval = uint64(configuration.Conf().Prometheus.CrawlingInterval)

	scheduler = gocron.NewScheduler()
	scheduler.Every(interval).Seconds().Do(job)
	scheduler.Start()

	return scheduler, nil
}

func StopToCrawl() {
	logger.Info("Stop crawling data from prometheus.")
	scheduler.Clear()
}
