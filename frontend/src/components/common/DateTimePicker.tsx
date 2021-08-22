import React, { FC, useState } from "react";
import moment from "moment";
import { DATE_FORMAT } from "utils/const";
import DatetimeRangePicker from "lib/datetimepicker";

interface Props {
    label: string;
    start: any;
    end: any;
    setStart: (start) => any;
    setEnd: (end) => any;
    onClickApply: (e, payload) => any;
    placeholder?: string;
}

const DateTimePicker: FC<Props> = props => {
    const [visible, setVisible] = useState(false);
    const [start, setStart] = useState(null);
    const [end, setEnd] = useState(null);

    const locale = {
        format: DATE_FORMAT,
        separator: " - ",
        applyLabel: "Apply",
        cancelLabel: "Cancel",
        weekLabel: "W",
        daysOfWeek: moment.weekdaysMin(),
        monthNames: moment.monthsShort(),
        firstDay: moment.localeData().firstDayOfWeek()
    };

    const onClick = () => {
        setVisible(true);
    };

    const onClickApply = (e, picker) => {
        setStart(picker.startDate);
        setEnd(picker.setEnd);

        props.onClickApply(e, picker);
    };

    return (
        <div
            className={`input-box datepicker ${visible && "visible"}`}
            onClick={onClick}
        >
            <label>{props.label}</label>
            <DatetimeRangePicker
                startDate={start}
                endDate={end}
                onApply={onClickApply}
                timePicker
                timePicker24Hour
                showDropdowns
                timePickerSeconds
                locale={locale}
                placeholder={props.placeholder}
            >
                <div className="input-group">
                    <input
                        type="text"
                        className="form-control"
                        value={
                            props.start && props.end
                                ? `${props.start.format(
                                      locale.format
                                  )} - ${props.end.format(locale.format)}`
                                : ""
                        }
                    />
                    <span className="input-group-btn">
                        <button className="default date-range-toggle">
                            <i className="fa fa-calendar" />
                        </button>
                    </span>
                </div>
            </DatetimeRangePicker>
        </div>
    );
};

export default DateTimePicker;
