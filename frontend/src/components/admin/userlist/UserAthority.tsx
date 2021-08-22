import React, {FC, useContext, useEffect, useState} from "react";
import {Checkbox, FormControlLabel, Switch, withStyles} from "@material-ui/core";
import {ChannelsData} from "../../../modules/channels";
import Language from "../../common/Language";

interface Props {
    channels?: ChannelsData;
    getAthorityData?: (payload) => any;
    getFormData?: (payload) => any;
}

const CustomCheckbox = withStyles({
    root: {
        color: "rgb(94, 114, 228)",
        "&$checked": {
            color: "rgb(94, 114, 228)"
        },
        width: 10,
        height: 10
    },
    checked: {}
})(props => <Checkbox color="default" {...props} />);

const CustomSwitch = withStyles({
    switchBase: {
        color: "default",
        '&$checked': {
            color: "rgb(94, 114, 228)",
            "&:hover": {
                backgroundColor: "rgba(208,208,208,0.3)",
            }
        },
        '&$checked + $track': {
            backgroundColor: "rgb(94, 114, 228)",
        },
    },
    checked: {},
    track: {},
})(Switch);


const UserAthority: FC<Props> = props => {

    const LanguageContext = useContext(Language());
    const {I18n} = LanguageContext;

    const [checkListChannel, setCheckListChannel] = useState([]);
    const [checkListPermission, setCheckListPermission] = useState([]);
    const [userType, setUserType] = useState("");
    const [checkUserType, setCheckUserType] = useState({
        checkUserTypeValue: true
    });

    const SwitchHandleChange = name => event => {
        setCheckUserType({...checkUserType, [name]: event.target.checked});
    };

    useEffect(() => {
        if (checkUserType.checkUserTypeValue) {
            setUserType("USER_COMMON");
        }
        if (!checkUserType.checkUserTypeValue) {
            setUserType("THIRD_PARTY");
            setCheckListPermission([]);
            setCheckListChannel([]);
        }
    }, [checkUserType]);

    const onClickPermissionCheck = id => {
        if (!checkListPermission.includes(id)) {
            setCheckListPermission(checkListPermission.concat(id));
        } else {
            const newCheckListPermission = checkListPermission.filter(el => {
                return el !== id;
            });
            setCheckListPermission(newCheckListPermission);
        }
    };

    const onClickChannelCheck = id => {
        if (!checkListChannel.includes(id)) {
            setCheckListChannel(checkListChannel.concat(id));
        } else {
            const newCheckListChannel = checkListChannel.filter(el => {
                return el !== id;
            });
            setCheckListChannel(newCheckListChannel);
        }
    };

    props.getAthorityData({
        permissionToAccess: checkListPermission,
        channels: checkListChannel,
        userType: userType
    });

    return (
        <div className="user-athority-form">
            <p className="form-title">
                <div className="switch-box">
                    {I18n.tableUserType}
                <span className="third-title">{I18n.THIRD_PARTY}</span>

                <FormControlLabel
                    control={<CustomSwitch className="authority-switch"
                                           checked={checkUserType.checkUserTypeValue}
                                           onChange={SwitchHandleChange("checkUserTypeValue")}
                                           value="checkUserTypeValue"
                    />}
                    label=""
                />
                <span className="common-title">{I18n.USER_COMMON}</span>
                </div>
            </p>

            <div className="access">
                <p className="ahtority-title">{I18n.addPerAccess}</p>
                <ul className="clearfix">
                    <li>
                        <FormControlLabel
                            control={<CustomCheckbox/>}
                            label="Node List"
                            checked={!!checkListPermission.includes("Node")}
                            onChange={() => onClickPermissionCheck("Node")}
                            disabled={!checkUserType.checkUserTypeValue}
                        />
                    </li>
                    <li>
                        <FormControlLabel
                            control={<CustomCheckbox/>}
                            label="Monitoring Log"
                            checked={
                                !!checkListPermission.includes("MonitoringLog")
                            }
                            onChange={() =>
                                onClickPermissionCheck("MonitoringLog")
                            }
                            disabled={!checkUserType.checkUserTypeValue}
                        />
                    </li>
                </ul>
            </div>

            <div className="channel">
                <p className="ahtority-title">{I18n.addChannelselect}</p>
                <ul className="clearfix">
                    {props.channels.data.map(el => {
                        return (
                            <li>
                                <FormControlLabel
                                    control={<CustomCheckbox/>}
                                    checked={!!checkListChannel.includes(el.id)}
                                    onChange={() => onClickChannelCheck(el.id)}
                                    disabled={!checkUserType.checkUserTypeValue}
                                    label={el.name}
                                />
                            </li>
                        );
                    })}
                </ul>
            </div>
        </div>
    );
};

export default UserAthority;
