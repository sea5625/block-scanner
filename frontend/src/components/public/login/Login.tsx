import React, { FC, useContext, useState, useEffect } from "react";
import Language from "components/common/Language";
import { Button } from "@material-ui/core";
import Loading from "components/common/Loading";

interface Props {
    login: (paylaod) => any;
    setToken: (paylaod) => any;
    getResource: (payload) => any;
    authStatus: AuthStatus;
    resource: {
        data: string;
    };
    loading: boolean;
    errorMessage: string;
}

interface AuthStatus {
    token: string;
    passwordStatus: number;
}

const Login: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [loginData, setLoginData] = useState({
        id: "",
        password: ""
    });
    const [errorMessage, setErrorMessage] = useState("");

    useEffect(() => {
        const payload = { id: "loginLogoImage" };
        props.getResource(payload);
    }, []);

    useEffect(() => {
        if (props.errorMessage === "No user in DB.") {
            setErrorMessage(I18n.ErrorNoUserInDB);
        } else if (props.errorMessage === "Fail to valid password.") {
            setErrorMessage(I18n.AlertMessagePleaseCheckYourPassword);
        }
    }, [props.errorMessage]);

    const onChange = e => {
        const { name, value } = e.target;
        setLoginData({
            ...loginData,
            [name]: value
        });
    };

    const onSubmit = () => {
        // setRender(false);
        props.login(loginData);
    };

    const onKeyPress = e => {
        if (e.charCode === 13) {
            onSubmit();
        }
    };
    if (props.loading) {
        return <Loading />;
    }
    return (
        <div className="auth-container">
            <div className="login-box">
                <div className="login-logo">
                    <img src={props.resource.data} alt="" />
                </div>
                <div className="login-form">
                    <div className="input-box">
                        <label htmlFor="id">{I18n.userID}</label>
                        <input
                            type="text"
                            placeholder={I18n.enterUserID}
                            name="id"
                            value={loginData.id}
                            onChange={onChange}
                            autoFocus={true}
                        />
                    </div>
                    <div className="input-box">
                        <label>{I18n.password}</label>
                        <input
                            type="password"
                            placeholder={I18n.enterPassword}
                            name="password"
                            value={loginData.password}
                            onChange={onChange}
                            onKeyPress={onKeyPress}
                        />
                    </div>
                    <p className="error-message">{errorMessage}</p>

                    <Button
                        className="login-btn"
                        onClick={onSubmit}
                        disabled={!loginData.id || !loginData.password}
                    >
                        {I18n.login}
                    </Button>
                </div>
            </div>
        </div>
    );
};

export default Login;
