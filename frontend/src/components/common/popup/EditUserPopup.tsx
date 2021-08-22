import React, {FC, useContext, useState, useEffect} from "react";
import {connect} from "react-redux";
import {Dispatch, Action} from "redux";
import {Button, Checkbox, FormControlLabel, Switch, TextField} from "@material-ui/core";
import {actions as usersActions} from "modules/users";
import {actions as channelsActions} from "modules/channels"
import {ChannelsData} from "modules/channels";
import {Language} from "components";
import {LayerPopup} from 'lib/popup'

interface Props {
    // user: UserData; => update
    putUser:(payload) => any;
    title: string;
    loading: boolean;
    getChannels: () => any;
    channels: ChannelsData;
    errorMessage: string;
    source:string;
    selectUser:{
        id:string;
        userId:string;
        email:string;
        phoneNumber:string;
        firstName:string;
        lastName:string;
        channels:[{
            id:string;
            name:string;
        }],
        permissionToAccess:[],
        userType:string;
    };
    layerKey:number;
}

const admin_pk = "PKID_0000000000000000";

const EditUserPopup: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const {I18n} = LanguageContext;
    const [errorMessage, setErrorMessage] = useState("");

    const [userId, setID] = useState(props.selectUser.userId);
    const [email, setEmail] = useState(props.selectUser.email);
    const [firstName, setFirstName] = useState(props.selectUser.firstName);
    const [userType] = useState(props.selectUser.userType);
    const [lastName, setLastName] = useState(props.selectUser.lastName);
    const [newPassword, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [phoneNumber, setPhoneNumber] = useState(props.selectUser.phoneNumber);
    const [channels, setChannels] = useState([]);
    const [permissionToAccess, setPermissionToAccess] = useState([]);

    useEffect(() => {
        props.getChannels();
        let newCheckListPermission = [];
        let newCheckListChannel = [];

        props.selectUser.channels.map(el => {
            return (
                newCheckListChannel.push(el.id)
            )
        });
        props.selectUser.permissionToAccess.map(el => {
            return(
                newCheckListPermission.push(el)
            )
        })

        setChannels(newCheckListChannel);
        setPermissionToAccess(newCheckListPermission);
    }, []);

    const onClickPermissionCheck=(id)=>{
        if (props.source === "admin") {
            return
        }
        if(!permissionToAccess.includes(id)){
            setPermissionToAccess(permissionToAccess.concat(id));
        }else{
            const newCheckListPermission = permissionToAccess.filter((el)=>{
                return el !== id

            });
            setPermissionToAccess(newCheckListPermission)
        }
    };

    const onClickChannelCheck=(id)=>{
        if (props.source === "admin") {
            return
        }
        if(!channels.includes(id)){
            setChannels(channels.concat(id));
        }else{
            const newCheckListChannel = channels.filter((el)=>{
                return el !== id

            });
            setChannels(newCheckListChannel)
        }
    };

    const onClickSubmit = () => {
        const id = props.selectUser.id;
        if (handleValidation(id) === false) {
            return
        }
        const payload = {
            id,
            userId,
            firstName,
            lastName,
            email,
            phoneNumber,
            channels,
            permissionToAccess,
            newPassword,
            userType
        };
        props.putUser(payload)
        LayerPopup.hide(props.layerKey)
    };

    const handleValidation = (id) => {
        // Blank value check
        if ( userId === "" ) {
            setErrorMessage(I18n.AlertMessagePleaseEnterUserID);
            return false
        } else if ( firstName === "" ) {
            setErrorMessage(I18n.AlertMessagePleaseEnterFirstName);
            return false
        } else if ( lastName === "" ) {
            setErrorMessage(I18n.AlertMessagePleaseEnterLastName);
            return false
        } else if ( email === "" ) {
            setErrorMessage(I18n.AlertMessagePleaseEnterEMail);
            return false
        } else if ( phoneNumber === "" ) {
            setErrorMessage(I18n.AlertMessagePleaseEnterTelNum);
            return false
        } else if ( channels.length === 0 ) {
            setErrorMessage(I18n.AlertMessageWarningNoSelectChannel);
            return false
        }

        // ID value check
        if ( id !== admin_pk ) {
            if(!userId.match(/^[a-zA-Z0-9]+$/)){
                setErrorMessage(I18n.AlertMessageWarningInvalidUserID);
                return false
            }
            if(userId.length <= 7){
                setErrorMessage(I18n.AlertMessageWarningUserIDLength);
                return false
            }
        }

        // Password value check
        if (newPassword !== "") {
            if(!newPassword.match(/^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[~!@#$%^&*()_+|<>?:{}]).{8,}$/))
            {
                setErrorMessage(I18n.AlertMessageWarningInvalidPassword);
                return false
            }
        }
        if (newPassword !== confirmPassword ) {
            setErrorMessage(I18n.AlertMessageWarningConfirmPassword);
            return false
        }

        //Email value check
        if(!email.match(/^[0-9a-zA-Z]([-_.]?[0-9a-zA-Z])*@[0-9a-zA-Z]([-_.]?[0-9a-zA-Z])*.[a-zA-Z]{2,3}$/))
        {
            setErrorMessage(I18n.AlertMessageWarningInvalidEMail);
            return false
        }

        //Phone num value check
        if(!phoneNumber.match(/^[0-9-]+$/))
        {
            setErrorMessage(I18n.AlertMessageWarningInvalidTelNum);
            return false
        }
        return true
    };

    const onClickCancel = () => {
        LayerPopup.hide(props.layerKey)
    };

    return (
        <div className="popup-content clearfix">
            <p className="title">{I18n.editCreateUsers}</p>
            <div className="flex-box between">
                <div className="user-info-form">
                    <div className="id full">
                        <div className="input-box">
                            <label htmlFor="id">{I18n.editUserId}</label>
                            <input
                                id="id"
                                name="id"
                                type="text"
                                value={userId}
                                onChange={value => setID(value.target.value)}
                                disabled={userId === "admin"}
                            />
                        </div>
                    </div>
                    <div className="password half">
                        <div className="input-box">
                            <label htmlFor="id">{I18n.editUserPassword}</label>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                onChange={value => setPassword(value.target.value)}
                            />
                        </div>
                        <div className="input-box">
                            <label htmlFor="password-confirm">{I18n.editUserPasswordConfirm}</label>
                            <input
                                id="password-confirm"
                                name="password-confirm"
                                type="password"
                                onChange={value => setConfirmPassword(value.target.value)}
                            />
                        </div>
                    </div>
                    <div className="name half">
                        <div className="input-box">
                            <label htmlFor="first-name">{I18n.editFirstName}</label>
                            <input
                                id="First Name"
                                name="First Name"
                                type="text"
                                value={firstName}
                                onChange={value => setFirstName(value.target.value)}
                            />
                        </div>
                        <div className="input-box">
                            <label htmlFor="last-name">{I18n.editLastName}</label>
                            <input
                                id="Last Name"
                                name="Last Name"
                                type="text"
                                value={lastName}
                                onChange={value => setLastName(value.target.value)}
                            />
                        </div>
                    </div>
                    <div className="email full">
                        <div className="input-box">
                            <label htmlFor="Email">{I18n.editEmail}</label>
                            <input
                                id="Email"
                                name="Email"
                                type="text"
                                value={email}
                                onChange={value => setEmail(value.target.value)}
                            />
                        </div>
                    </div>
                    <div className="phone full">
                        <div className="input-box">
                            <label htmlFor="Phone">{I18n.editTelNum}</label>
                            <input
                                id="Phone"
                                name="Phone"
                                type="text"
                                value={phoneNumber}
                                onChange={value => setPhoneNumber(value.target.value)}
                            />
                        </div>
                    </div>
                </div>
            </div>
            <div className="user-athority-form">
                <div className="access">
                    <p className="ahtority-title">{I18n.editPerAccess}</p>
                    <ul className="clearfix">
                        <li>
                            <FormControlLabel
                                control={<Checkbox value="checkedA"/>}
                                label="Node List"
                                checked={!!permissionToAccess.includes("Node")}
                                onChange={() => onClickPermissionCheck("Node")}
                            />
                        </li>
                        <li>
                            <FormControlLabel
                                control={<Checkbox value="checkedA"/>}
                                label="Monitoring Log"
                                checked={!!permissionToAccess.includes("MonitoringLog")}
                                onChange={() => onClickPermissionCheck("MonitoringLog")}
                            />
                        </li>
                    </ul>
                </div>

                <div className="channel">
                    <p className="ahtority-title">{I18n.editChannelselect}</p>
                    <ul className="clearfix">
                        {props.channels.data.map((el) => {
                            return (
                                <li>
                                    <FormControlLabel
                                        control={<Checkbox value="checkedA"/>}
                                        checked={!!channels.includes(el.id)}
                                        onChange={() => onClickChannelCheck(el.id)}
                                        label={el.name}
                                    />
                                </li>
                            )
                        })}
                    </ul>
                </div>
                <p className="error-message">{errorMessage}</p>
            </div>
            <div className="btn-box">
                <Button className="cancle-btn" onClick={onClickCancel}>{I18n.cancel}</Button>
                <Button className="submit-btn"
                        onClick={onClickSubmit}
                        disabled={!userId || !firstName || !lastName || !email || !phoneNumber}
                >{I18n.confirm}</Button>
            </div>
        </div>
    );
};

const mapStateToProps = state => ({
    // user: state.users.user, => update
    channels: state.channels.channels

});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    putUser : payload => dispatch(usersActions.putUser(payload)),
    getChannels: () => dispatch(channelsActions.getChannels()),
});


export default connect(
    mapStateToProps,
    mapDispatchToProps
)(EditUserPopup);
