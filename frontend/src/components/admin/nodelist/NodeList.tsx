import React, { FC, useState, useEffect, useContext } from "react";
import { MdErrorOutline } from "react-icons/md";
import { NodesData } from "modules/nodes";
import { Loading, Language } from "components";
import { makeSortList, formatNumberWithComma } from "utils/utils";
import { TICK_TIME } from "utils/const";

interface Props {
    match: {
        params: {
            name: string;
        };
    };
    location: {
        state: {
            channelId: string;
        };
    };
    loading: boolean;
    getNodes: (payload) => any;
    nodes: NodesData;
    language: "ko" | "en";
}

const NodeList: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;

    useEffect(() => {
        const { channelId } = props.location.state;
        const timer = setInterval(() => tick(), TICK_TIME);
        props.getNodes(channelId);
        return () => {
            clearInterval(timer);
        };
    }, [props.location.state.channelId]);

    const tick = () => {
        const { channelId } = props.location.state;
        props.getNodes(channelId);
    };

    if (props.loading && !props.nodes.data.length) {
        return <Loading />;
    } else {
        return (
            <div className="content node-list">
                <div className="table-box type-b">
                    <p className="title">{I18n.nodeList}</p>
                    <div className="table-content">
                        <table>
                            <thead>
                                <tr>
                                    <th>{I18n.nodeName}</th>
                                    <th>{I18n.nodeIP}</th>
                                    <th>{I18n.tableBlockHeight}</th>
                                    <th>{I18n.tableConfirmedTx}</th>
                                    <th>{I18n.tableUnconfirmedTx}</th>
                                    <th>{I18n.tableResponseTime}</th>
                                </tr>
                            </thead>
                            <tbody>
                                {makeSortList(props.nodes.data).map(
                                    (el, key) => {
                                        return (
                                            <tr
                                                key={key}
                                                className={`${el.isLeader ===
                                                    1 &&
                                                    "leader"} ${el.status ===
                                                    1 ||
                                                    (el.status === 2 &&
                                                        "error")}`}
                                            >
                                                {/* <th>{status}</th> */}
                                                <th>{el.name}</th>
                                                <th>{el.ip}</th>
                                                <th>
                                                    {el.blockHeight}
                                                    {el.status === 1 && (
                                                        <span className="error-status">
                                                            <MdErrorOutline />
                                                            Unsync
                                                        </span>
                                                    )}
                                                </th>
                                                <th>
                                                    {formatNumberWithComma(
                                                        el.countOfTX
                                                    )}
                                                </th>
                                                <th>
                                                    {formatNumberWithComma(
                                                        el.countOfUnconfirmedTX
                                                    )}
                                                </th>
                                                <th>
                                                    {typeof el.responseTimeInSec ===
                                                        "number" &&
                                                        el.responseTimeInSec.toFixed(
                                                            3
                                                        )}
                                                    {el.status === 2 && (
                                                        <span className="error-status">
                                                            <MdErrorOutline />
                                                            Delay
                                                        </span>
                                                    )}
                                                </th>
                                            </tr>
                                        );
                                    }
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        );
    }
};

export default NodeList;
