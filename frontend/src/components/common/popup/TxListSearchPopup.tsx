import React, { FC, useContext, useState, useEffect } from "react";
import { Button } from "@material-ui/core";
import { LayerPopup } from "lib/popup";
import { DateTimePicker } from "components";
import { encodeURIMoment } from "utils/utils";

interface Props {
    label: string;
    layerKey: string;
    onClickSearch: (payload) => any;
}

const TxListSearchPopup: FC<Props> = props => {
    const [searchData, setSearchData] = useState({
        status: "",
        fromAddress: "",
        toAddress: "",
        blockHeight: "",
        data: ""
    });
    const [start, setStart] = useState(null);
    const [end, setEnd] = useState(null);

    const onChange = e => {
        const { name, value } = e.target;
        setSearchData(prevSearchData => ({
            ...prevSearchData,
            [name]: value
        }));
    };

    const onClickApply = (event, picker) => {
        setStart(picker.startDate);
        setEnd(picker.endDate);
    };

    const onClickInit = () => {
        setSearchData({
            status: "",
            blockHeight: "",
            fromAddress: "",
            toAddress: "",
            data: ""
        });
        setStart(null);
        setEnd(null);
    };

    const onClickSearch = () => {
        props.onClickSearch({
            ...searchData,
            from: encodeURIMoment(start),
            to: encodeURIMoment(end)
        });
        LayerPopup.hide(props.layerKey);
    };

    return (
        <div className="popup-content clearfix">
            <div className="flex-box between">
                <div className="input-box">
                    <label>Status</label>
                    <select
                        value={searchData.status}
                        name="status"
                        onChange={onChange}
                    >
                        {!searchData.status && <option disabled></option>}
                        <option value="success">Success</option>
                        <option value="failure">Failure</option>
                    </select>
                </div>
                <DateTimePicker
                    label={props.label}
                    start={start}
                    end={end}
                    setStart={setStart}
                    setEnd={setEnd}
                    onClickApply={onClickApply}
                />
                <div className="input-box">
                    <label>From Wallet</label>
                    <input
                        type="text"
                        name="fromAddress"
                        value={searchData.fromAddress}
                        onChange={onChange}
                    />
                </div>
                <div className="input-box">
                    <label>To Wallet</label>
                    <input
                        type="text"
                        name="toAddress"
                        value={searchData.toAddress}
                        onChange={onChange}
                    />
                </div>
                <div className="input-box">
                    <label>Block Height</label>
                    <input
                        type="text"
                        name="blockHeight"
                        value={searchData.blockHeight}
                        onChange={onChange}
                    />
                </div>
                <div className="input-box">
                    <label>Data</label>
                    <input
                        type="text"
                        name="data"
                        value={searchData.data}
                        onChange={onChange}
                    />
                </div>
            </div>
            <div className="btn-box">
                <Button className="init-btn" onClick={onClickInit}>
                    Initialize
                </Button>
                <Button className="search-btn" onClick={onClickSearch}>
                    Search
                </Button>
            </div>
        </div>
    );
};

export default TxListSearchPopup;
