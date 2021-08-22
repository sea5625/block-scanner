import React, {FC, useContext, useState, useEffect} from "react";
import {connect} from "react-redux";
import {Dispatch, Action} from "redux";
import {Button, Checkbox, withStyles} from "@material-ui/core";
import {actions as usersActions} from "modules/users";
import {actions as channelsActions} from "modules/channels"
import {ChannelsData} from "modules/channels";
import Language from "components/common/Language";
import {LayerPopup} from 'lib/popup';
import {UserForm, UserAthority} from "components";

interface Props {
    createUser: (payload) => any;
    title: string;
    loading: boolean;
    getChannels: () => any;
    channels: ChannelsData;
    errorMessage: string;
    layerKey:number;
}

const AddUserPopup: FC<Props> = props => {
    ;

    const LanguageContext = useContext(Language());
    const {I18n} = LanguageContext;
    const [errorMessage, setErrorMessage] = useState("");

    useEffect(() => {
        props.getChannels();
    }, []);

    //UserForm
    const [userId, setID] = useState("");
    const [email, setEmail] = useState("");
    const [firstName, setFirstName] = useState("");
    const [lastName, setLastName] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [phoneNumber, setPhoneNumber] = useState("");

    //UserAthority
    const [permissionToAccess, setPermissionToAccess] = useState([]);
    const [channels, setChannels] = useState([]);
    const [userType, setUserType] = useState("");

    const getFormData = ({id, email, firstName, lastName, confirmPassword, password, phoneNumber}) => {
        setID(id);
        setEmail(email);
        setFirstName(firstName);
        setLastName(lastName);
        setPassword(password);
        setConfirmPassword(confirmPassword);
        setPhoneNumber(phoneNumber);
    };
    const getAthorityData = ({permissionToAccess, channels, userType}) => {
        setPermissionToAccess(permissionToAccess);
        setChannels(channels);
        setUserType(userType);
    };
    const onClickSubmit = () => {

        if (handleValidation() === false) {
            return;
        }
        const payload = {
            userId,
            firstName,
            lastName,
            email,
            phoneNumber,
            channels,
            permissionToAccess,
            password,
            userType
        }
        props.createUser(payload);
        LayerPopup.hide(props.layerKey)

    };

    const handleValidation = () => {

        // ID value check
        if (!userId.match(/^[a-zA-Z0-9]+$/)) {
            setErrorMessage(I18n.AlertMessageWarningInvalidUserID)
            // Popup.alertPopup("AlertMessageWarningInvalidUserID");
            return false;
        }
        if (userId.length <= 7) {
            setErrorMessage(I18n.AlertMessageWarningUserIDLength)
            // Popup.alertPopup("AlertMessageWarningUserIDLength");
            return false;
        }

        // Password value check
        if (
            !password.match(
                /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[~!@#$%^&*()_+|<>?:{}]).{8,}$/
            )
        ) {
            setErrorMessage(I18n.AlertMessageWarningInvalidPassword)
            // Popup.alertPopup("AlertMessageWarningInvalidPassword");
            return false;
        }
        if (password !== confirmPassword) {
            setErrorMessage(I18n.AlertMessageWarningConfirmPassword)
            // Popup.alertPopup("AlertMessageWarningConfirmPassword");
            return false;
        }

        //Email value check
        if (
            !email.match(
                /^(([^<>()[\].,;:\s@"]+(\.[^<>()[\].,;:\s@"]+)*)|(".+"))@(([^<>()[\].,;:\s@"]+\.)+[^<>()[\].,;:\s@"]{2,})$/i
            )
        ) {
            setErrorMessage(I18n.AlertMessageWarningInvalidEMail)
            // Popup.alertPopup("AlertMessageWarningInvalidEMail");
            return false;
        }

        //Phone num value check
        if (!phoneNumber.match(/^[0-9-]+$/)) {
            setErrorMessage(I18n.AlertMessageWarningInvalidTelNum)
            // Popup.alertPopup("AlertMessageWarningInvalidTelNum");
            return false;
        }

        //channel value check
        if(userType === "USER_COMMON") {
            if (channels.length === 0) {
                setErrorMessage(I18n.AlertMessageWarningNoSelectChannel)
                // Popup.alertPopup("AlertMessageWarningNoSelectChannel");
                return false;
            }

        }
        return true;
    };

    const onClickCancel = () => {
        LayerPopup.hide(props.layerKey)
    };

    return (
        <div className="popup-content clearfix">
            <p className="title">{I18n.addCreateUsers}</p>
            <div className="flex-box between">
                <UserForm
                    getFormData={getFormData}
                    userType={userType}
                    errorMessage={errorMessage}
                />
                <UserAthority
                    getAthorityData={getAthorityData}
                    getFormData={getFormData}
                    channels={props.channels}/>
            </div>
            <div className="btn-box">
                <Button className="cancle-btn" onClick={onClickCancel}>{I18n.cancel}</Button>
                <Button className="submit-btn"
                        onClick={onClickSubmit}
                        disabled={userType === "USER_COMMON" ?
                            (!userId || !password || !firstName || !lastName || !email || !phoneNumber || !confirmPassword)
                            : (!userId || !password || !lastName || !email || !phoneNumber || !confirmPassword)}
                >{I18n.confirm}</Button>
            </div>
        </div>
    );
};

const mapStateToProps = state => ({
    channels: state.channels.channels

});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    createUser: payload => dispatch(usersActions.createUser(payload)),
    getChannels: () => dispatch(channelsActions.getChannels())
});


export default connect(
    mapStateToProps,
    mapDispatchToProps
)(AddUserPopup);
