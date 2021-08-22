import React, { useState, useEffect, useRef } from "react";
import { Line } from "react-chartjs-2";
import "chartjs-plugin-annotation";
import { TICK_TIME } from "utils/const";

const ResponseTimeChart = props => {
    const responseTimeChart = useRef(null);
    const [values] = useState(props.values);

    useEffect(() => {
        advance();
        const interval = setInterval(() => advance(), TICK_TIME);
        return () => clearInterval(interval);
    }, [props.value]);
    const annotationValue = props.alertList.filter(el => {
        return props.name === el.name;
    });

    const data = {
        datasets: [
            {
                data: props.values,
                dataPoints: props.values,
                // fill: props.isError ? true : false,
                fill: false,
                borderColor:
                    props.isError === 4
                        ? "rgb(94, 114, 228)"
                        : "rgb(234, 84, 85)",
                borderWidth: 1.5,
                lineTension: 0.25,
                pointRadius: 0
            }
        ]
    };
    const options = {
        maintainAspectRatio: false,
        responsive: true,
        animation: {
            duration: 0
        },
        legend: false,
        scales: {
            xAxes: [
                {
                    type: "time",
                    display: true,
                    time: {
                        unit: "minute"
                    },
                    ticks: {
                        fontSize: 8
                    }
                }
            ],
            yAxes: [
                {
                    ticks: {
                        max: annotationValue[0]
                            ? annotationValue[0].slowResponseTime +
                              (annotationValue[0].slowResponseTime % 2 === 0
                                  ? 4
                                  : 3)
                            : 0,
                        min: 0,
                        fontSize: 8
                    }
                }
            ]
        },
        annotation: {
            annotations: [
                {
                    type: "line",
                    mode: "horizontal",
                    scaleID: "y-axis-0",
                    value: annotationValue[0]
                        ? annotationValue[0].slowResponseTime
                        : 0,
                    background: "rgba(200,60,60,0.25)",
                    // borderColor: "rgba(200,60,60,0.25)"
                    borderColor:
                        props.isError === 4 || props.isError === 5
                            ? " rgb(234, 84, 85)"
                            : "rgba(94, 114, 228,0.8)",
                    borderWidth: 0.8
                }
            ],
            drawTime: "afterDraw"
        }
    };

    const updateCharts = () => {
        if (responseTimeChart.current) {
            responseTimeChart.current.chartInstance.update();
        }
    };

    const progress = () => {
        const { makeResponseTimeChartValues, value, name } = props;
        makeResponseTimeChartValues(value, name);
    };

    const advance = () => {
        progress();
        updateCharts();
    };

    // if (!props.values) {
    //   return <Loader />;
    // }
    return (
        <div className="response-time-chart chart">
            <p className="chart-title">
                <span className="chart-last-res-time">Response Time(sec)</span>
                <span className="chart-current-res-time">{props.value}</span>
            </p>
            <Line ref={responseTimeChart} data={data} options={options} />
        </div>
    );
};

export default ResponseTimeChart;
