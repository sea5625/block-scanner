import React, { FC, useState, useEffect } from "react";
import {
    Input,
    IconButton,
    InputLabel,
    Select,
    MenuItem
} from "@material-ui/core";
import { MdKeyboardArrowLeft, MdKeyboardArrowRight } from "react-icons/md";

interface Props {
    setCount: (payload) => any;
    setPage: (payload) => any;
    count: number;
    page: number;
    total: number;
}

const Pagination: FC<Props> = props => {
    const [totalPage, setTotalPage] = useState(1);
    useEffect(() => {
        const total =
            Math.ceil(props.total / props.count) < 1
                ? 1
                : Math.ceil(props.total / props.count);

        setTotalPage(total);
    }, [props.count, props.page, props.total, totalPage]);

    const onClickCount = e => {
        if (e.target.value) {
            props.setCount(e.target.value);
        }
    };

    return (
        <div className="pagination">
            <div className="page-count">
                <InputLabel htmlFor="page-helper">per page: </InputLabel>
                <Select
                    value={props.count}
                    onClick={onClickCount}
                    input={<Input name="page-count" />}
                    displayEmpty
                    defaultValue={10}
                    name="page-count"
                >
                    <MenuItem value={10} className="count-item">
                        10
                    </MenuItem>
                    <MenuItem value={25} className="count-item">
                        25
                    </MenuItem>
                    <MenuItem value={50} className="count-item">
                        50
                    </MenuItem>
                    <MenuItem value={100} className="count-item">
                        100
                    </MenuItem>
                </Select>
            </div>
            <div className="current-page">
                <p>
                    {props.page === 1 ? 1 : props.count * (props.page - 1) + 1}{" "}
                    -{" "}
                    {props.total > props.count * props.page
                        ? props.count * props.page
                        : props.total}{" "}
                    of {props.total}
                </p>
            </div>
            <IconButton
                className="prev-btn"
                disabled={props.page === 1}
                onClick={() => props.setPage(props.page - 1)}
            >
                <MdKeyboardArrowLeft />
            </IconButton>
            <IconButton
                className="next-btn"
                disabled={
                    props.page * props.count >= props.total || totalPage === 1
                }
                onClick={() => props.setPage(props.page + 1)}
            >
                <MdKeyboardArrowRight />
            </IconButton>
        </div>
    );
};

export default Pagination;
