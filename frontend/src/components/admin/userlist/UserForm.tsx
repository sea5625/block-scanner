import React, {FC, useContext, useEffect, useState} from "react";
import {Theme, createStyles, withStyles, TextField} from "@material-ui/core";
import {TICK_TIME} from "../../../utils/const";
import Language from "../../common/Language";

interface Props {
    classes?: {};
    getFormData?: (payload) => any;
    errorMessage?: string;
    userType?: string;
}

const styles = (theme: Theme) =>
    createStyles({
        root: {
            fullInput: {
                width: "100%",
                height: 44
            }
        }
    });

const UserForm: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const {I18n} = LanguageContext;

    const [userID, setUserID] = useState("");
    const [addUserPassword, setAdduserPassword] = useState("");
    const [addUserPasswordConfirm, setAdduserPasswordConfirm] = useState("");
    const [addFirstName, setAddFirstName] = useState("(주)");
    const [addLastName, setAddLastName] = useState("");
    const [addEmail, setAddEmail] = useState("");
    const [addTelNum, setAddTelNum] = useState("");

    props.getFormData({
        id: userID,
        email: addEmail,
        firstName: addFirstName,
        lastName: addLastName,
        confirmPassword: addUserPassword,
        password: addUserPasswordConfirm,
        phoneNumber: addTelNum
    });

    return (
        <div className="user-info-form">
            <p className="form-title">{I18n.addUserInfo}</p>
            <div className="id full">
                <div className="input-box">
                    <label htmlFor="id">{I18n.addUserId}</label>
                    <input
                        id="id"
                        name="id"
                        type="text"
                        onChange={value => setUserID(value.target.value)}
                    />
                </div>
            </div>
            <div className="password half">
                <div className="input-box">
                    <label htmlFor="password">{I18n.addUserPassword}</label>
                    <input
                        id="password"
                        name="password"
                        type="password"
                        onChange={value => setAdduserPassword(value.target.value)}
                    />
                </div>
                <div className="input-box">
                    <label htmlFor="password-confirm">{I18n.addUserPasswordConfirm}</label>
                    <input
                        id="password-confirm"
                        name="password-confirm"
                        type="password"
                        onChange={value => setAdduserPasswordConfirm(value.target.value)}
                    />
                </div>
            </div>
            {props.userType === "USER_COMMON" ? (
                    <div className="name half">
                        <div className="input-box">
                            <label htmlFor="first-name">{I18n.addFirstName}</label>
                            <input
                                id="First Name"
                                name="First Name"
                                type="text"
                                onChange={value => setAddFirstName(value.target.value)}
                            />
                        </div>
                        <div className="input-box">
                            <label htmlFor="last-name">{I18n.addLastName}</label>
                            <input
                                id="Last Name"
                                name="Last Name"
                                type="text"
                                onChange={value => setAddLastName(value.target.value)}
                            />
                        </div>
                    </div>
                ) :
                (<div className="name full">
                        <div className="input-box">
                            <label htmlFor="last-name">{I18n.addTpLastName}</label>
                            <input
                                id="TP Last Name"
                                name="TP Last Name"
                                type="text"
                                onChange={value => setAddLastName(value.target.value)}
                            />
                        </div>
                    </div>
                )}

            <div className="email full">
                <div className="input-box">
                    <label htmlFor="Email">{I18n.addEmail}</label>
                    <input
                        id="Email"
                        name="Email"
                        type="text"
                        onChange={value => setAddEmail(value.target.value)}
                    />
                </div>
            </div>
            <div className="phone full">
                <div className="input-box">
                    <label htmlFor="Phone">{I18n.addTelNum}</label>
                    <input
                        id="Phone"
                        name="Phone"
                        type="text"
                        onChange={value => setAddTelNum(value.target.value)}
                    />
                </div>
            </div>
            <p className="error-message">{props.errorMessage}</p>
        </div>
    );
};

export default withStyles(styles)(UserForm);
