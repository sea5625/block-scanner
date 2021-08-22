import React, { FC, useState, useEffect, useContext } from "react";
import { MdSearch, MdArrowDropDown } from "react-icons/md";
import { Pagination, Language, DateTimePicker } from "components";
import { SymptomData } from "modules/symptom";
import { formatMoment, encodeURIMoment } from "utils/utils";

interface Props {
    match: {
        params: {
            name: string;
        };
    };
    getSymptom: (payload) => any;
    symptom: SymptomData;
}
const MonitoringLog: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [count, setCount] = useState(10);
    const [page, setPage] = useState(1);
    const [logList, setLogList] = useState(props.symptom.data);
    const [start, setStart] = useState(null);
    const [end, setEnd] = useState(null);
    const [search, setSearch] = useState(false);
    const [filter, setFilter] = useState({
        symptom: "",
        channel: ""
    });

    useEffect(() => {
        props.getSymptom({ limit: count, offset: 0 });
        setStart(null);
        setEnd(null);
        setPage(1);
    }, [count]);

    useEffect(() => {
        props.getSymptom({
            limit: count,
            offset: 0,
            from: encodeURIMoment(start),
            to: encodeURIMoment(end)
        });
        setPage(1);
    }, [search]);

    useEffect(() => {
        props.getSymptom({
            limit: count,
            offset: page === 1 ? 0 : count * (page - 1),
            from: encodeURIMoment(start),
            to: encodeURIMoment(end)
        });
    }, [page]);

    const onClickApply = (event, picker) => {
        setStart(picker.startDate);
        setEnd(picker.endDate);
    };
    const onClickSearch = () => {
        setSearch(!search);
    };

    useEffect(() => {
        setLogList(props.symptom.data);
    }, [props.symptom.data]);

    const makeFilterList = () => {
        let symptoms = [];
        let channels = [];
        props.symptom.data.forEach(el => {
            if (el.symptom === "Slow response") {
                symptoms.push("slowRes");
            } else if (el.symptom === "Unsync block") {
                symptoms.push("unsyncBlock");
            }
            channels.push(el.channel);
        });
        return {
            symptomFilter: [...new Set(symptoms.sort())],
            channelsFilter: [...new Set(channels.sort())]
        };
    };

    const onClickSymptomFilter = value => {
        let selectSymptom;
        if (value === "slowRes") {
            selectSymptom = "Slow response";
        } else {
            selectSymptom = "Unsync block";
        }
        const newLogList = props.symptom.data.filter(el => {
            return filter.channel
                ? el.channel === filter.channel && el.symptom === selectSymptom
                : el.symptom === selectSymptom;
        });
        setFilter({ ...filter, symptom: selectSymptom });
        setLogList(newLogList);
    };

    const onClickSymptomRefresh = () => {
        const newLogList = props.symptom.data.filter(el => {
            return filter.channel ? el.channel === filter.channel : el;
        });
        setFilter({ ...filter, symptom: "" });
        setLogList(newLogList);
    };

    const onClickChannelFilter = value => {
        const newLogList = props.symptom.data.filter(el => {
            return filter.symptom
                ? el.channel === value && el.symptom === filter.symptom
                : el.channel === value;
        });
        setFilter({ ...filter, channel: value });
        setLogList(newLogList);
    };

    const onClickChannelRefresh = () => {
        const newLogList = props.symptom.data.filter(el => {
            return filter.symptom ? el.symptom === filter.symptom : el;
        });
        setFilter({ ...filter, channel: "" });
        setLogList(newLogList);
    };

    const filterList = makeFilterList();
    return (
        <div className="content monitoring-log">
            <div className="table-box type-b">
                <div className="title">
                    {I18n.monitoringLog}
                    <div className="time-search-box">
                        <DateTimePicker
                            label={""}
                            onClickApply={onClickApply}
                            start={start}
                            end={end}
                            setStart={setStart}
                            setEnd={setEnd}
                            placeholder={I18n.selectDate}
                        />
                        <button
                            className="time-search-btn"
                            onClick={onClickSearch}
                        >
                            <span className="btn-txt">
                                {I18n.search}
                                <i>
                                    <MdSearch />
                                </i>
                            </span>
                        </button>
                    </div>
                </div>
                <div className="table-content">
                    <table>
                        <thead>
                            <tr>
                                <th>{I18n.timestamp}</th>
                                <th className="dropdown">
                                    {I18n.symptom}
                                    <MdArrowDropDown />
                                    <ul className="dropdown-menu">
                                        <li
                                            className="dropdown-item"
                                            onClick={onClickSymptomRefresh}
                                        >
                                            {I18n.viewAll}
                                        </li>
                                        {filterList.symptomFilter.map(
                                            (el, key) => {
                                                return (
                                                    <li
                                                        className="dropdown-item"
                                                        key={key}
                                                        onClick={() =>
                                                            onClickSymptomFilter(
                                                                el
                                                            )
                                                        }
                                                    >
                                                        {I18n[el]}
                                                    </li>
                                                );
                                            }
                                        )}
                                    </ul>
                                </th>
                                <th className="dropdown">
                                    {I18n.tableChannelName}
                                    <MdArrowDropDown />
                                    <ul className="dropdown-menu">
                                        <li
                                            className="dropdown-item"
                                            onClick={onClickChannelRefresh}
                                        >
                                            {I18n.viewAll}
                                        </li>
                                        {filterList.channelsFilter.map(
                                            (el, key) => {
                                                return (
                                                    <li
                                                        className="dropdown-item"
                                                        key={key}
                                                        onClick={() =>
                                                            onClickChannelFilter(
                                                                el
                                                            )
                                                        }
                                                    >
                                                        {el}
                                                    </li>
                                                );
                                            }
                                        )}
                                    </ul>
                                </th>
                                <th>{I18n.information}</th>
                            </tr>
                        </thead>
                        <tbody>
                            {logList.map((el, key) => {
                                return (
                                    <tr key={key}>
                                        <th>{formatMoment(el.timeStamp)}</th>
                                        <th>{el.symptom}</th>
                                        <th>{el.channel}</th>
                                        <th>{el.msg}</th>
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
                    total={props.symptom.total}
                />
            </div>
        </div>
    );
};

export default MonitoringLog;
