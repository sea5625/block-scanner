import React, { FC, useContext } from "react";
import { FaTrashAlt } from "react-icons/fa";
import { Button } from "@material-ui/core";
import { LayerPopup } from "lib/popup";
import { Language } from "components";

interface Props {
    name: string;
    layerKey: string;
    callbackFunc?: () => any;
}
const DeletePopup: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;

    const onClickCancel = () => {
        LayerPopup.hide(props.layerKey);
    };
    const onClickCallback = () => {
        LayerPopup.hide(props.layerKey);
        if (props.callbackFunc) {
            props.callbackFunc();
        }
    };

    return (
        <div className="delete-content">
            <i className="delete-ic">
                <FaTrashAlt />
            </i>
            <p className="delete-title">{props.name}</p>
            <p className="message">{I18n.areYouSure}</p>
            <div className="btn-box">
                <Button className="cancel-btn" onClick={onClickCancel}>
                    {I18n.cancel}
                </Button>
                <Button className="delete-btn" onClick={onClickCallback}>
                    {I18n.delete}
                </Button>
            </div>
        </div>
    );
};

export default DeletePopup;
