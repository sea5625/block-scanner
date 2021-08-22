import React, { FC, useState, useEffect, useContext } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import {
    FormControl,
    Select,
    Input,
    MenuItem,
    Button
} from "@material-ui/core";
import { actions as settingActions } from "modules/setting";
import { Language } from "components";
import LayerPopup from "lib/popup";

interface Props {
    layerKey: string;
    getSetting: () => any;
    updateSetting: (payload) => any;
    sessionTimeout: number;
}

const SessionTimeOutPopup: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [sessionTimeout, setSessionTimeout] = useState(0);
    useEffect(() => {
        props.getSetting();
    }, []);
    useEffect(() => {
        setSessionTimeout(props.sessionTimeout);
    }, [props.sessionTimeout]);

    const onChangeSessionTimeout = e => {
        setSessionTimeout(e.target.value);
    };

    const onClickCancel = () => {
        LayerPopup.hide(props.layerKey);
    };
    const onClickSave = () => {
        props.updateSetting({ sessionTimeout });
        LayerPopup.hide(props.layerKey);
    };

    return (
        <div className="popup-content">
            <p className="title">{I18n.setSessionTime}</p>
            <div className="setting-form">
                <FormControl>
                    <Select
                        value={sessionTimeout}
                        onChange={onChangeSessionTimeout}
                        input={
                            <Input
                                name="sessionTimeout"
                                id="session-timeout-select"
                            />
                        }
                    >
                        <MenuItem value={5}>5 {I18n.min}</MenuItem>
                        <MenuItem value={15}>
                            15 {I18n.min} ({I18n.default})
                        </MenuItem>
                        <MenuItem value={30}>30 {I18n.min}</MenuItem>
                        <MenuItem value={60}>60 {I18n.min}</MenuItem>
                        <MenuItem value={0}>{I18n.sessionNotExpired}</MenuItem>
                    </Select>
                </FormControl>
            </div>
            <div className="btn-box">
                <Button className="cancel-btn" onClick={onClickCancel}>
                    {I18n.cancel}
                </Button>
                <Button className="save-btn" onClick={onClickSave}>
                    {I18n.save}
                </Button>
            </div>
        </div>
    );
};

const mapStateToProps = state => ({
    sessionTimeout: state.setting.sessionTimeout
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getSetting: () => dispatch(settingActions.getSetting()),
    updateSetting: payload => dispatch(settingActions.updateSetting(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(SessionTimeOutPopup);
