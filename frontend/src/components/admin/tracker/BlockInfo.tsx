import React, { FC, useEffect, useContext } from "react";
import { Link } from "react-router-dom";
import { MdArrowForward, MdChevronLeft, MdChevronRight } from "react-icons/md";
import { IconButton } from "@material-ui/core";
import { BlockData } from "modules/tracker";
import { Language, CopyButton } from "components";
import Search from "./Search";
import { formatMoment } from "utils/utils";

interface Props {
    history: {
        push: (url) => any;
    };
    match: {
        params: {
            name: string;
            channelId: string;
            id: number;
        };
    };
    getBlock: (payload) => any;
    block: BlockData;
    loading: boolean;
}

const BlockInfo: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const { name, channelId: channel, id } = props.match.params;

    useEffect(() => {
        props.getBlock({ name, channel, id });
    }, [id]);

    const onClickPrev = () => {
        const prevBlock = props.block.blockHeight - 1;
        props.history.push(
            `/tracker/block_info/${name}/${channel}/${prevBlock}`
        );
    };
    const onClickNext = () => {
        const nextBlock = props.block.blockHeight + 1;
        props.history.push(
            `/tracker/block_info/${name}/${channel}/${nextBlock}`
        );
    };
    const onClickTx = txHash => {
        props.history.push(`/tracker/tx_info/${name}/${channel}/${txHash}`);
    };
    return (
        <div className="content tracker-info block-info">
            <Search channel={channel} name={name} />
            <p className="title">
                <span className="title-txt clearfix">{I18n.blockTracker}</span>
                <span className="block-height">
                    <button
                        className="prev-btn"
                        onClick={onClickPrev}
                        disabled={props.block.blockHeight === 1}
                    >
                        <MdChevronLeft />
                    </button>
                    {`# ${props.block.blockHeight}`}
                    <button className="next-btn" onClick={onClickNext}>
                        <MdChevronRight />
                    </button>
                </span>
            </p>
            <div className="overview clearfix">
                <p className="overview-title">{I18n.blockInfo}</p>
                <ul>
                    <li>
                        <span className="item-title">
                            {I18n.tableBlockHeight}
                        </span>
                        {props.block.blockHeight}
                    </li>
                    <li>
                        <span className="item-title">{I18n.tablePeerID}</span>
                        {props.block.peerID}
                    </li>
                    <li>
                        <span className="item-title">
                            {I18n.tableTimestamp}
                        </span>
                        {formatMoment(props.block.timeStamp)}
                    </li>
                    <li>
                        <span className="item-title">
                            {I18n.tableBlockHash}
                        </span>
                        {props.block.blockHash}
                        <CopyButton value={props.block.blockHash} />
                    </li>
                </ul>
            </div>
            <div className="table-box type-a">
                <p className="title">{I18n.txTracker}</p>
                <div className="table-content">
                    <table>
                        <thead>
                            <tr>
                                <th>{I18n.tableTxHash}</th>
                                <th>
                                    {I18n.tableFrom}
                                    <MdArrowForward />
                                </th>
                                <th>{I18n.tableTo}</th>
                                <th>{I18n.tableData}</th>
                            </tr>
                        </thead>
                        <tbody>
                            {props.block.confirmedTx.map((el, key) => {
                                return (
                                    <tr key={key}>
                                        <th>
                                            <button
                                                className="hash-link-btn ellipsis"
                                                onClick={() =>
                                                    onClickTx(el.txHash)
                                                }
                                            >
                                                {el.txHash}
                                            </button>
                                        </th>
                                        <th className="ellipsis">
                                            {el.from}
                                            <MdArrowForward />
                                        </th>
                                        <th className="ellipsis">{el.to}</th>
                                        <th className="ellipsis">{el.data}</th>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};

export default BlockInfo;
