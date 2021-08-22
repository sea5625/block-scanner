import React, { useState, Fragment, useEffect } from "react";
import { CopyToClipboard } from "react-copy-to-clipboard";
import { FaRegCopy } from "react-icons/fa";

const CopyButton = props => {
    const [text, setText] = useState("COPY");

    useEffect(() => {}, [props.value]);

    const onClick = () => {
        setTimeout(() => {
            setText("COPY");
        }, 1500);
        setText("COPIED");
    };
    if (!props.value) {
        return <Fragment />;
    }
    return (
        <CopyToClipboard text={props.value} onCopy={onClick}>
            <button className={`copy-btn ${text === "COPIED" && "active"}`}>
                {text}
                <FaRegCopy />
            </button>
        </CopyToClipboard>
    );
};

export default CopyButton;
