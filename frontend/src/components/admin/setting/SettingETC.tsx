import React, { FC, useState, useEffect, useContext } from "react";
import { FormControl, Select, Input, MenuItem } from "@material-ui/core";
import { FaRegCheckCircle } from "react-icons/fa";
import { AlertingData } from "modules/alerting";
import { Language, Loading, Popup } from "components";
import { makeSortList } from "utils/utils";

interface Props {
    alerting: AlertingData;
    getAlerting: () => any;
    updateAlerting: (payload) => any;
    alertingLoading: boolean;
}

const SettingETC: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [alerting, setAlerting] = useState([]);

    useEffect(() => {
        props.getAlerting();
    }, []);

    useEffect(() => {
        setAlerting(makeSortList(props.alerting.data));
    }, [props.alerting]);

    const onChangeAlerting = (e, key, idx) => {
        const _alerting = JSON.parse(JSON.stringify(alerting));
        _alerting[idx][key] = e.target.value;
        setAlerting(_alerting);
    };

    const onClickConfirm = () => {
        if (
            JSON.stringify(makeSortList(props.alerting.data)) ===
            JSON.stringify(alerting)
        ) {
            Popup.alertPopup({
                className: "alert",
                message: "AlertMessageNotChangeSetting",
                btnName: "ok"
            });
            return;
        }
        props.updateAlerting(alerting);
    };

    if (props.alertingLoading) {
        return <Loading />;
    }
    return (
        <div className="content setting-etc">
            <div className="title">
                {I18n.settingAlert}
                <button className="confirm-btn" onClick={onClickConfirm}>
                    <span className="btn-txt">
                        <p>{I18n.save}</p>
                        <i>
                            <FaRegCheckCircle />
                        </i>
                    </span>
                </button>
            </div>
            <div className="table-box type-a">
                <div className="table-content">
                    <table>
                        <thead>
                            <tr>
                                <th>{I18n.channelName}</th>
                                <th>{I18n.tableTimeLimitForDataSync}</th>
                                <th>{I18n.tableTimeLimitForResponseTime}</th>
                            </tr>
                        </thead>
                        <tbody>
                            {alerting.map((el, key) => {
                                return (
                                    <tr key={key}>
                                        <th>{el.name}</th>
                                        <th>
                                            <FormControl>
                                                <Select
                                                    value={
                                                        el.unsyncBlockToleranceTime
                                                    }
                                                    onChange={e =>
                                                        onChangeAlerting(
                                                            e,
                                                            "unsyncBlockToleranceTime",
                                                            key
                                                        )
                                                    }
                                                    input={
                                                        <Input
                                                            name="unsync"
                                                            id="unsync-select"
                                                        />
                                                    }
                                                >
                                                    <MenuItem value={360}>
                                                        6 {I18n.min} (
                                                        {I18n.default})
                                                    </MenuItem>
                                                    <MenuItem value={420}>
                                                        7 {I18n.min}
                                                    </MenuItem>
                                                    <MenuItem value={480}>
                                                        8 {I18n.min}
                                                    </MenuItem>
                                                    <MenuItem value={540}>
                                                        9 {I18n.min}
                                                    </MenuItem>
                                                    <MenuItem value={600}>
                                                        10 {I18n.min}
                                                    </MenuItem>
                                                </Select>
                                            </FormControl>
                                            {/* {el.unsyncBlockToleranceTime} */}
                                        </th>
                                        <th>
                                            <FormControl>
                                                <Select
                                                    value={el.slowResponseTime}
                                                    onChange={e =>
                                                        onChangeAlerting(
                                                            e,
                                                            "slowResponseTime",
                                                            key
                                                        )
                                                    }
                                                    input={
                                                        <Input
                                                            name="response"
                                                            id="response-select"
                                                        />
                                                    }
                                                >
                                                    <MenuItem value={2}>
                                                        2 {I18n.sec}
                                                    </MenuItem>
                                                    <MenuItem value={5}>
                                                        5 {I18n.sec} (
                                                        {I18n.default})
                                                    </MenuItem>
                                                    <MenuItem value={8}>
                                                        8 {I18n.sec}
                                                    </MenuItem>
                                                    <MenuItem value={10}>
                                                        10 {I18n.sec}
                                                    </MenuItem>
                                                    <MenuItem value={20}>
                                                        20 {I18n.sec}
                                                    </MenuItem>
                                                </Select>
                                            </FormControl>
                                            {/* {el.slowResponseTime} */}
                                        </th>
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

export default SettingETC;
