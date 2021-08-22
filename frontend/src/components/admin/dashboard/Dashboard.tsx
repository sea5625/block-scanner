import React, { FC, useContext, useState, useEffect, Component } from "react";
import { TX_CHART_X_LENGTH } from "utils/const";
import moment from "moment";
import { MdErrorOutline } from "react-icons/md";
import { ChannelsData } from "modules/channels";
import { AlertingData } from "modules/alerting";

import { PrometheusData } from "modules/prometheus";
import { TICK_TIME, NODE_TYPE_GOLOOP } from "utils/const";
import { Language, Loading } from "components";
import Total from "./Total";
import ResponseTimeChart from "./ResponseTimeChart";
import TxChart from "./TxChart";
import axios from "axios";

interface Props {
    refreshChannels: () => any;
    getPrometheus: () => any;
    getAlerting: () => any;
    alerting: AlertingData;
    channels: ChannelsData;
    prometheus: PrometheusData;
    storage: {
        nodeType: string;
        prometheus: string;
        jobName: string;
    };
}

const Dashboard: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;

    const [values, setValues] = useState({});
    const [today, setToday] = useState({});
    const [txCountData, setTxCountData] = useState({});
    const [todayTxCountData, setTodayTxCountData] = useState({});

    useEffect(() => {
        props.getPrometheus();
        props.getAlerting();
        const timer = setInterval(() => tick(), TICK_TIME);
        return () => {
            clearInterval(timer);
        };
    }, []);

    useEffect(() => {
        if (props.storage.prometheus) {
            getTodayTxData();
        }
    }, [props.storage.prometheus, props.channels.data]);

    const tick = () => {
        props.refreshChannels();
    };

    const makeResponseTimeChartValues = (value, name) => {
        if (!values[name]) {
            let arr = [];
            for (let i = 60; 0 < i; i--) {
                arr.push({
                    x: moment()
                        .subtract(5000 * i, "milliseconds")
                        .toDate(),
                    y: null
                });
            }
            setValues(values => ({
                ...values,
                [name]: arr
            }));
            return;
        }
        setValues(values => ({
            ...values,
            [name]: values[name]
                .concat({ x: new Date(), y: value })
                .slice(1, 61)
        }));
    };

    const renderNodeState = nodeList => {
        let check = 0;
        nodeList.forEach(function(nel) {
            if (nel.status === 1 || nel.status === 2) {
                if (check !== 0 && check !== nel.status) {
                    return (check = 3);
                } else {
                    check = nel.status;
                }
            } else if (nel.status === 3) {
                return (check = nel.status);
            }
        });
        return check;
    };

    const getPrometheusTxcount = async rfc3339Formatted => {
        try {
            const prometheusURL = props.storage.prometheus;
            let prometheusTx_Count_Query;

            if (props.storage.nodeType == NODE_TYPE_GOLOOP) {
                prometheusTx_Count_Query = `?query=${props.storage.jobName}_txpool_user_remove_sum&time=${rfc3339Formatted}`;
            } else {
                prometheusTx_Count_Query = `?query=tx_count&time=${rfc3339Formatted}`;
            }

            return await axios.get(
                `${prometheusURL}${prometheusTx_Count_Query}`
            );
        } catch (error) {
            console.error(error);
        }
    };

    const getTxData = async () => {
        let arr = [];
        let txCountData = {};
        for (let i = 0; i < TX_CHART_X_LENGTH; i++) {
            const rfc3339 = moment()
                .subtract(i, "days")
                .startOf("days")
                .toISOString();
            const res = await getPrometheusTxcount(rfc3339);
            let result = {};
            if (res && res.data) {
                const { result: data } = res.data.data;
                if (data.length < 1) {
                    props.channels.data.forEach(el => {
                        result[el.name] = {
                            x:
                                moment()
                                    .subtract(i + 1, "days")
                                    .startOf("days")
                                    .valueOf() / 1000,
                            y: 0
                        };
                    });
                } else {
                    if (props.storage.nodeType == NODE_TYPE_GOLOOP) {
                        data.forEach(el => {
                            props.channels.data.forEach(value => {
                                if ("0x" + el.metric.channel == value.id) {
                                    el.metric.channel = value.name;
                                }
                            });
                        });
                    }
                    data.forEach(el => {
                        result[el.metric.channel] = {
                            x:
                                moment()
                                    .subtract(i + 1, "days")
                                    .startOf("days")
                                    .valueOf() / 1000,
                            y: el.value[1]
                        };
                    });
                }
                arr.push(result);
                const keyArr: string[] = Object.keys(result);
                for (let j = 0; j < keyArr.length; j++) {
                    const key: string = keyArr[j];
                    txCountData[key] = {
                        ...txCountData[key],
                        [i + 1]: arr[i][key]
                    };
                }
            }
        }
        setTxCountData(txCountData);
    };

    const getTodayTxData = async () => {
        const now = moment();
        if (today !== now.day() || Object.keys(txCountData).length === 0) {
            getTxData();
            setToday(now.day());
        }
        const res = await getPrometheusTxcount(now.toISOString());
        let result = {};
        if (res && res.data) {
            const { result: data } = res.data.data;
            if (data.length < 1) {
                props.channels.data.forEach(el => {
                    result[el.name] = {
                        x:
                            moment()
                                .startOf("days")
                                .valueOf() / 1000,
                        y: 0
                    };
                });
            } else {
                if (props.storage.nodeType == NODE_TYPE_GOLOOP) {
                    data.forEach(el => {
                        props.channels.data.forEach(value => {
                            if ("0x" + el.metric.channel == value.id) {
                                el.metric.channel = value.name;
                            }
                        });
                    });
                }
                data.forEach(el => {
                    result[el.metric.channel] = {
                        x:
                            moment()
                                .startOf("days")
                                .valueOf() / 1000,
                        y: el.value[1]
                    };
                });
            }
            setTodayTxCountData(result);
        }
    };

    if (!txCountData || !todayTxCountData) {
        return <Loading />;
    } else {
        return (
            <div className="content dashboard">
                <div
                    className={`grid ${
                        props.channels.data.length < 4
                            ? `row-1 column-${props.channels.data.length}`
                            : `grid ${
                                  props.channels.data.length === 4
                                      ? "row-2 column-2"
                                      : `grid ${
                                            props.channels.data.length >= 5
                                                ? "row-2 column-3"
                                                : "row-2 column-3"
                                        }`
                              }`
                    }`}
                >
                    {/*TODO*/}
                    {props.channels.data.map((el, key) => {
                        const error = el.status === 4 ? "" : "error";
                        const channelStatus = renderNodeState(el.nodes);
                        return (
                            <div className={`channel ${error}`} key={key}>
                                <div className="top-box">
                                    <p className="title">{el.name}</p>
                                    {channelStatus === 0 ? (
                                        ""
                                    ) : channelStatus === 1 ? (
                                        <span className="error-status">
                                            <MdErrorOutline />
                                            Unsync Block
                                        </span>
                                    ) : channelStatus === 2 ? (
                                        <span className="error-status">
                                            <MdErrorOutline />
                                            Slow Response
                                        </span>
                                    ) : (
                                        <span className="error-status">
                                            <MdErrorOutline />
                                            Slow Response {"&"} Unsync Block
                                        </span>
                                    )}
                                </div>
                                <ResponseTimeChart
                                    name={el.name}
                                    values={values[el.name]}
                                    value={
                                        typeof el.responseTimeInSec === "number"
                                            ? el.responseTimeInSec.toFixed(3)
                                            : 0.0
                                    }
                                    makeResponseTimeChartValues={
                                        makeResponseTimeChartValues
                                    }
                                    isError={el.status}
                                    alertList={props.alerting.data}
                                />
                                <Total
                                    blockHeight={el.blockHeight}
                                    countOfTX={el.countOfTX}
                                    totalNodes={el.nodes.length}
                                    channelTotal={el.total}
                                    responseTimeInSec={el.responseTimeInSec}
                                    name={el.name}
                                    channelId={el.id}
                                />

                                <TxChart
                                    isError={el.status}
                                    txCountData={txCountData[el.name]}
                                    todayTxCountData={todayTxCountData[el.name]}
                                />
                            </div>
                        );
                    })}
                </div>
            </div>
        );
    }
};

export default Dashboard;
