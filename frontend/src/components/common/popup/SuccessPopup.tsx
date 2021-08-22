import React, { FC, useContext } from "react";
import { FaRegCheckCircle } from "react-icons/fa";
import { Button } from "@material-ui/core";
import { LayerPopup } from "lib/popup";
import { Language } from "components";

interface Props {
    message: string;
    layerKey: string;
    callbackFunc?: () => any;
}
const SuccessPopup: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const onClickCallback = () => {
        LayerPopup.hide(props.layerKey);
        if (props.callbackFunc) {
            props.callbackFunc();
        }
    };

    return (
        <div className="success-content">
            <i className="success-ic">
                <FaRegCheckCircle />
            </i>
            <p className="success-title">{I18n.success}</p>
            <p className="message">{I18n[props.message]}</p>
            <Button className="ok-btn" onClick={onClickCallback}>
                {I18n.ok}
            </Button>
        </div>
    );
};

export default SuccessPopup;
