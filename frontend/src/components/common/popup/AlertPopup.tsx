import React, { FC, useContext } from "react";
import { MdCancel } from "react-icons/md";
import { Button } from "@material-ui/core";
import { LayerPopup } from "lib/popup";
import { Language } from "components";

interface Props {
    message: string;
    btnName: string;
    layerKey: string;
    callbackFunc?: () => any;
}
const AlertPopup: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const onClickCallback = () => {
        LayerPopup.hide(props.layerKey);
        if (props.callbackFunc) {
            props.callbackFunc();
        }
    };
    return (
        <div className="alert-content">
            <i className="alert-ic">
                <MdCancel />
            </i>
            <p className="alert-title">{I18n.fail}</p>
            <p className="alert-message">{I18n[props.message]}</p>
            <Button className="alert-btn" onClick={onClickCallback}>
                {I18n[props.btnName]}
            </Button>
        </div>
    );
};

export default AlertPopup;
