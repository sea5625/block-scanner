import React, { useEffect, useState } from "react";
import { Line } from "react-chartjs-2";
import moment from "moment";
import { Loading } from "../../index";

const TxChart = props => {
    const [currentTx, setCurrentTx] = useState(0);
    const [chartData, setChartData] = useState([]);

    useEffect(() => {
        makeTxCountData();
    }, [props.todayTxCountData]);

    const makeTxCountData = () => {
        const { txCountData, todayTxCountData } = props;
        let arr = [];
        if (typeof txCountData === "undefined") {
            return <Loading />;
        } else {
            Object.keys(txCountData).forEach(key => {
                arr.push(txCountData[key]);
            });
            arr.unshift(todayTxCountData);
            let result = [];
            for (let i = 1; i < arr.length; i++) {
                result.push({
                    x: moment(arr[i - 1].x * 1000).format("YYYY-MM-DD"),
                    y: arr[i - 1].y - arr[i].y
                });
            }
            setChartData(result);
            setCurrentTx(result[0].y);
        }
    };

    const data = {
        datasets: [
            {
                label: "Daily Transactions(1day)",
                fill: true,
                backgroundColor:
                    props.isError === 4
                        ? "rgba(94, 114, 228, 0.8)"
                        : "rgba(234, 84, 85, 0.8)",
                lineTension: 0,
                borderColor:
                    props.isError === 4 ? "rgba(94, 114, 228, 1)" : "#f55a4e",
                borderWidth: 1,
                borderCapStyle: "butt",
                borderDash: [],
                borderDashOffset: 0.0,
                borderJoinStyle: "miter",
                // pointBorderColor: props.error ? "#d32f2f" : "#ffa726",
                pointBackgroundColor:
                    props.isError === 4 ? "rgba(94, 114, 228, 0.8)" : "#f55a4e",
                pointBorderWidth: 2,
                pointHoverRadius: 2,
                // pointHoverBackgroundColor: props.error ? "#d32f2f" : "#ffa726",
                // pointHoverBorderColor: props.error ? "#d32f2f" : "#ffa726",
                pointRadius: 0,
                pointHitRadius: 10,
                data: chartData
            }
        ]
    };
    const options = {
        responsive: true,
        maintainAspectRatio: props.length === 3,
        legend: {
            display: false
        },
        scales: {
            xAxes: [
                {
                    type: "time",
                    labelString: "probability",
                    display: true,
                    time: {
                        unit: "day"
                    },
                    ticks: {
                        fontSize: 8
                    }
                }
            ],
            yAxes: [
                {
                    ticks: {
                        min: 0,
                        fontSize: 8
                    }
                }
            ]
        }
    };

    return (
        <div className="tx-chart chart">
            <div className="chart-title">
                Daily Transactions
                <p className="chart-current-tx-value">{currentTx}</p>
            </div>

            <Line data={data} options={options} />
        </div>
    );
};

export default TxChart;
