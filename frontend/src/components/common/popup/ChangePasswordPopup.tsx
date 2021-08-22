import React, { FC, useContext, useState, useEffect } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import useReactRouter from "use-react-router";
import { Button } from "@material-ui/core";
import { MdLockOutline } from "react-icons/md";
import { LayerPopup } from "lib/popup";
import Language from "components/common/Language";
import Loading from "components/common/Loading";
import {
    actions as usersActions,
    api as userApi,
    UserData
} from "modules/users";

interface Props {
    id: string;
    getUser: (payload) => any;
    user: UserData;
    putUser: (payload) => any;
    layerKey: string;
    loading: boolean;
}

const ChangePasswordPopup: FC<Props> = props => {
    const [errorMessage, setErrorMessage] = useState("");
    const { history } = useReactRouter();

    const LanguageContext = useContext(Language());
    const [updatePasswordData, setUpdateData] = useState({
        password: "",
        confirmPassword: ""
    });

    const { I18n } = LanguageContext;

    useEffect(() => {
        props.getUser({ id: props.id });
    }, []);

    const handleChange = e => {
        const { name, value } = e.target;

        setUpdateData({
            ...updatePasswordData,
            [name]: value
        });
    };

    const onClickCancel = () => {
        LayerPopup.hide(props.layerKey);
    };
    const handleSubmit = async () => {
        const { password, confirmPassword } = updatePasswordData;
        if (password !== confirmPassword) {
            setErrorMessage(I18n.AlertMessageWarningConfirmPassword);
            return;
        }
        if (
            !password.match(
                /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[~!@#$%^&*()_+|<>?:{}]).{8,}$/
            )
        ) {
            setErrorMessage(I18n.AlertMessageWarningInvalidPassword);
            return;
        }
        try {
            let _channelIds = [];
            props.user.channels.forEach(el => {
                _channelIds.push(el.id);
            });
            const payload = {
                ...props.user,
                newPassword: updatePasswordData.password,
                channels: _channelIds
            };
            const res = await userApi.putUser(payload);
            if (res) {
                LayerPopup.hide(props.layerKey);
                history.push("/");
            }
        } catch (e) {
            console.log(e);
        }
    };
    if (props.loading) {
        return <Loading />;
    } else {
        return (
            <div className="popup-box3 popup-type1">
                <MdLockOutline />
                <p className="change-password-title">{I18n.firstLoginUsers}</p>
                <p className="message">{I18n.firstLoginUsersSubtitle}</p>
                <div className="input-id-box">
                    <label className="User-Id">{I18n.editUserId}</label>
                    <input
                        className="popup-type1 input-box-user-id"
                        type="text"
                        name="userId"
                        value={props.user.userId}
                        readOnly
                    />
                </div>
                <div className="input-password-box">
                    <div className="input-box">
                        <label className="password">
                            {I18n.editUserPassword}
                        </label>
                        <input
                            className="popup-type1 input-box-user-password"
                            type="password"
                            name="password"
                            value={updatePasswordData.password}
                            onChange={handleChange}
                            id="addUserPassword"
                            placeholder={I18n.editUserPassword}
                        />
                    </div>
                    <div className="input-box">
                        <label className="password">
                            {I18n.editUserPasswordConfirm}
                        </label>
                        <input
                            className="popup-type1 input-box-user-cpassword"
                            type="password"
                            name="confirmPassword"
                            value={updatePasswordData.confirmPassword}
                            onChange={handleChange}
                            id="addUserPasswordConfirmInput"
                            placeholder={I18n.editUserPasswordConfirmInput}
                        />
                    </div>
                </div>
                <p className="error-message">{errorMessage}</p>
                <div className="btn-box">
                    <Button className="cancel-btn" onClick={onClickCancel}>
                        {I18n.cancel}
                    </Button>
                    <Button
                        className="agree-btn"
                        onClick={handleSubmit}
                        disabled={
                            !updatePasswordData.confirmPassword ||
                            !updatePasswordData.password
                        }
                    >
                        {I18n.agree}
                    </Button>
                </div>
            </div>
        );
    }
};

const mapStateToProps = state => ({
    user: state.users.user,
    loading: state.users.userLoading
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getUser: payload => dispatch(usersActions.getUser(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(ChangePasswordPopup);
