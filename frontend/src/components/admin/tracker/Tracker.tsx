import React, { FC, useEffect, useContext } from "react";
import { Link } from "react-router-dom";
import { MdChevronRight, MdErrorOutline } from "react-icons/md";
import { FaRegCheckCircle } from "react-icons/fa";
import Button from "@material-ui/core/Button";
import { BlocksData, TxsData } from "modules/tracker";
import { ChannelData } from "modules/channels";
import { Loading, Language } from "components";
import Search from "./Search";
import Total from "./Total";
import { formatMoment } from "utils/utils";

interface Props {
    history: {
        push: (url) => any;
    };
    match: {
        params: {
            name: string;
            channelId: string;
        };
    };
    location: {
        state: {
            channelId: string;
        };
    };
    getChannel: (payload) => any;
    getBlocks: (payload) => any;
    getTxs: (payload) => any;
    channel: ChannelData;
    blocks: BlocksData;
    txs: TxsData;
    blocksLoading: boolean;
    txsLoading: boolean;
}
const Tracker: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const { channelId: channel, name } = props.match.params;
    useEffect(() => {
        const payload = {
            channel,
            limit: 10,
            offset: 0
        };
        props.getBlocks(payload);
        props.getTxs(payload);
    }, [channel]);

    useEffect(() => {
        props.getChannel({ id: channel });
    }, [channel]);

    const onClickBlock = id => {
        props.history.push(`/tracker/block_info/${name}/${channel}/${id}`);
    };
    const onClickTx = txHash => {
        props.history.push(`/tracker/tx_info/${name}/${channel}/${txHash}`);
    };

    if (props.blocksLoading || props.txsLoading) {
        return <Loading />;
    }
    return (
        <div className="content tracker">
            <Search channel={channel} name={name} />
            <Total channel={props.channel} />
            <div className="recent-list-box">
                <div className="recent-list recent-block-list">
                    <div className="table-box type-b">
                        <p className="title">
                            {I18n.recentBlocks}
                            <span>
                                <Link
                                    to={{
                                        pathname: `/tracker/block_list/${name}/${channel}`
                                    }}
                                >
                                    <Button className="view-btn">
                                        {I18n.viewAll}
                                        <MdChevronRight />
                                    </Button>
                                </Link>
                            </span>
                        </p>
                        <div className="table-content">
                            <table>
                                <thead>
                                    <tr>
                                        <th>{I18n.tableBlockHeight}</th>
                                        <th>{I18n.tableTimestamp}</th>
                                        <th>{I18n.tableTotalTxCount}</th>
                                        <th>{I18n.tableBlockHash}</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {props.blocks.data.map((el, key) => {
                                        return (
                                            <tr key={key}>
                                                <th>
                                                    <button
                                                        className="block-link-btn"
                                                        onClick={() =>
                                                            onClickBlock(
                                                                el.blockHeight
                                                            )
                                                        }
                                                    >
                                                        {el.blockHeight}
                                                    </button>
                                                </th>
                                                <th>
                                                    {formatMoment(el.timeStamp)}
                                                </th>
                                                <th>{el.confirmedTx.length}</th>
                                                <th>
                                                    <button
                                                        className="hash-link-btn ellipsis"
                                                        onClick={() =>
                                                            onClickBlock(
                                                                el.blockHeight
                                                            )
                                                        }
                                                    >
                                                        {el.blockHash}
                                                    </button>
                                                </th>
                                            </tr>
                                        );
                                    })}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
                <div className="recent-list recent-tx-list">
                    <div className="table-box type-b">
                        <p className="title">
                            {I18n.recentTransactions}
                            <span>
                                <Link
                                    to={{
                                        pathname: `/tracker/tx_list/${name}/${channel}`,
                                        state: {
                                            channelId: channel
                                        }
                                    }}
                                >
                                    <Button className="view-btn">
                                        {I18n.viewAll}
                                        <MdChevronRight />
                                    </Button>
                                </Link>
                            </span>
                        </p>
                        <div className="table-content">
                            <table>
                                <thead>
                                    <tr>
                                        <th>{I18n.tableStatus}</th>
                                        <th>{I18n.tableTimestamp}</th>
                                        <th>{I18n.tableTxHash}</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    {props.txs.data.map((el, key) => {
                                        return (
                                            <tr key={key}>
                                                <th
                                                    className={
                                                        el.status === "Success"
                                                            ? "success"
                                                            : "fail"
                                                    }
                                                >
                                                    {el.status}
                                                    {el.status === "Success" ? (
                                                        <FaRegCheckCircle />
                                                    ) : (
                                                        <MdErrorOutline />
                                                    )}
                                                </th>
                                                <th>
                                                    {formatMoment(el.timeStamp)}
                                                </th>
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
                                            </tr>
                                        );
                                    })}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Tracker;
