import React, { FC, useState, useEffect, useContext } from "react";
import { MdArrowForward, MdSearch } from "react-icons/md";
import { Button } from "@material-ui/core";
import { TxsData } from "modules/tracker";
import { Pagination, Language, Popup } from "components";
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
        };
    };
    getTxs: (payload) => any;
    txs: TxsData;
    txsLoading: boolean;
}

const TxList: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [count, setCount] = useState(10);
    const [page, setPage] = useState(1);
    const { channelId: channel, name } = props.match.params;
    const [searchData, setSearchData] = useState({});

    useEffect(() => {
        const payload = {
            channel,
            limit: count,
            offset: 0
        };
        props.getTxs(payload);
        setPage(1);
    }, [count]);

    useEffect(() => {
        props.getTxs({
            channel,
            limit: count,
            offset: page === 1 ? 0 : count * (page - 1),
            ...searchData
        });
    }, [page]);

    const onClickSearch = payload => {
        setSearchData(payload);
        setPage(1);
        setCount(10);
        props.getTxs({
            channel,
            limit: count,
            offset: 0,
            ...payload
        });
    };

    const onClickSearchMore = () => {
        Popup.txListSearchPopup({
            className: "search",
            label: "Time Stamp",
            onClickSearch
        });
    };

    const onClickBlock = id => {
        props.history.push(`/tracker/block_info/${name}/${channel}/${id}`);
    };

    const onClickTx = txHash => {
        props.history.push(`/tracker/tx_info/${name}/${channel}/${txHash}`);
    };

    return (
        <div className="content tx-list">
            <Search channel={channel} name={name} />
            <div className="table-box type-b">
                <p className="title">
                    {I18n.transactionList}
                    <Button className="search-btn" onClick={onClickSearchMore}>
                        {I18n.searchMore}
                        <MdSearch />
                    </Button>
                </p>
                <div className="table-content">
                    <table>
                        <thead>
                            <tr>
                                <th>{I18n.tableTxHash}</th>
                                <th>{I18n.tableBlockHeight}</th>
                                <th>{I18n.tableTimestamp}</th>
                                <th>
                                    {I18n.tableFrom}
                                    <MdArrowForward />
                                </th>
                                <th>{I18n.tableTo}</th>
                            </tr>
                        </thead>
                        <tbody>
                            {props.txs.data.map((el, key) => {
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
                                        <th>
                                            <button
                                                className="block-link-btn"
                                                onClick={() =>
                                                    onClickBlock(el.blockHeight)
                                                }
                                            >
                                                {el.blockHeight}
                                            </button>
                                        </th>
                                        <th>{formatMoment(el.timeStamp)}</th>
                                        <th className="ellipsis">
                                            {el.from}
                                            <MdArrowForward />
                                        </th>
                                        <th className="ellipsis">{el.to}</th>
                                    </tr>
                                );
                            })}
                        </tbody>
                    </table>
                </div>
                <Pagination
                    setCount={setCount}
                    setPage={setPage}
                    count={count}
                    page={page}
                    total={props.txs.total}
                />
            </div>
        </div>
    );
};

export default TxList;
