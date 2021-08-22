import React, { FC, useEffect, useContext } from "react";
import { FaRegCheckCircle } from "react-icons/fa";
import { MdErrorOutline } from "react-icons/md";
import { TxData } from "modules/tracker";
import { Language, CopyButton } from "components";
import Search from "./Search";
import { stringifyJSON, formatMoment } from "utils/utils";

interface Props {
    history: {
        push: (url) => any;
    };
    match: {
        params: {
            name: string;
            txHash: string;
            channelId: string;
        };
    };
    getTx: (payload) => any;
    txHash: string;
    tx: TxData;
    loading: boolean;
}

const TxInfo: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const { txHash, name, channelId: channel } = props.match.params;
    useEffect(() => {
        props.getTx({ name, channel, txHash });
    }, [txHash]);
    return (
        <div className="content tracker-info tx-info">
            <Search channel={channel} name={name} />
            <p className="title">{I18n.txTracker}</p>
            <span className="hash">
                <span className="txt">{txHash}</span>
                <CopyButton value={txHash} />
            </span>

            <div className="overview clearfix">
                <p className="overview-title">
                    {I18n.txInfo}
                    <span
                        className={`status ${
                            props.tx.status === "Success" ? "success" : "fail"
                        }`}
                    >
                        {props.tx.status}
                        {props.tx.status === "Success" ? (
                            <FaRegCheckCircle />
                        ) : (
                            <MdErrorOutline />
                        )}
                    </span>
                </p>
                <ul>
                    <li>
                        <span className="item-title">
                            {I18n.tableBlockHeight}
                        </span>
                        {props.tx.blockHeight}
                    </li>
                    <li>
                        <span className="item-title">
                            {I18n.tableTimestamp}
                        </span>
                        {formatMoment(props.tx.timeStamp)}
                    </li>
                    <li>
                        <span className="item-title">{I18n.tableFrom}</span>
                        {props.tx.from}
                        <CopyButton value={props.tx.from} />
                    </li>
                    <li>
                        <span className="item-title">{I18n.tableTo}</span>
                        {props.tx.to}
                        <CopyButton value={props.tx.to} />
                    </li>
                </ul>
            </div>
            <div className="data">
                <p className="data-title">{I18n.tableData}</p>
                <div className="data-content">
                    {stringifyJSON(props.tx.data) !== '""' &&
                        JSON.stringify(
                            JSON.parse(stringifyJSON(props.tx.data)),
                            null,
                            "\t"
                        )}
                </div>
            </div>
        </div>
    );
};

export default TxInfo;
