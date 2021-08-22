import React, { FC, useState, useEffect, useContext } from "react";
import { BlocksData } from "modules/tracker";
import { Pagination, Language } from "components";
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
    getBlocks: (payload) => any;
    blocks: BlocksData;
    blocksLoading: boolean;
}

const BlockList: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [count, setCount] = useState(10);
    const [page, setPage] = useState(1);
    const { channelId: channel, name } = props.match.params;

    useEffect(() => {
        const payload = {
            channel,
            limit: count,
            offset: 0
        };
        props.getBlocks(payload);
        setPage(1);
    }, [count]);

    useEffect(() => {
        props.getBlocks({
            channel,
            limit: count,
            offset: page === 1 ? 0 : count * (page - 1)
        });
    }, [page, count]);

    const onClickBlock = id => {
        props.history.push(`/tracker/block_info/${name}/${channel}/${id}`);
    };
    return (
        <div className="content block-list">
            <Search channel={channel} name={name} />
            <div className="table-box type-b">
                <p className="title">{I18n.blockList}</p>
                <div className="table-content">
                    <table>
                        <thead>
                            <tr>
                                <th>{I18n.tableBlockHeight}</th>
                                <th>{I18n.tableTimestamp}</th>
                                <th>{I18n.tableTxCountInBlock}</th>
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
                                                    onClickBlock(el.blockHeight)
                                                }
                                            >
                                                {el.blockHeight}
                                            </button>
                                        </th>
                                        <th> {formatMoment(el.timeStamp)}</th>
                                        <th>{el.confirmedTx.length}</th>
                                        <th>
                                            <button
                                                className="hash-link-btn ellipsis"
                                                onClick={() =>
                                                    onClickBlock(el.blockHeight)
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
                <Pagination
                    setCount={setCount}
                    setPage={setPage}
                    count={count}
                    page={page}
                    total={props.blocks.total}
                />
            </div>
        </div>
    );
};

export default BlockList;
