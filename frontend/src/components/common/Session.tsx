import React, { FC, useState, useEffect, Fragment } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import { actions as authActions } from "modules/auth";
import { actions as settingsActions } from "modules/setting";
import { actions as storageActions } from "modules/storage";
import { TICK_TIME } from "utils/const";
import { getToken, parsingToken } from "utils/token";

interface Props {
    getSetting: () => any;
    refreshToken: (payload) => any;
    logout: () => any;
    sessionTimeout: number;
}
export var sessionTimerID

const Session: FC<Props> = props => {
    useEffect(() => {
        const events = ["load", "click", "scroll", "keypress"];
        for (let i in events) {
            window.addEventListener(events[i], resetTimeOut);
        }
        props.getSetting();
        resetTimeOut();
        const timerID = setInterval(() => tickCheck(), TICK_TIME * 12);
        return () => {
            clearTimeoutFunc();
            clearInterval(timerID);
        };
    }, [props.sessionTimeout]);

    const tickCheck = () => {
        const token = getToken();
        if (!token) {
            return;
        }
        //check exp, iat time
        const sessionTime = parsingToken(token);
        const exp = sessionTime["exp"];
        const iat = sessionTime["iat"];
        const expTimestamp = +new Date(exp);
        const iatTimestamp = +new Date(iat);

        //check now time
        const date = new Date();
        const time = date.toISOString();
        const nowTime = time.substring(0, 19) + "Z";
        const nowTimestamp = +new Date(nowTime);

        if (nowTimestamp >= (expTimestamp + iatTimestamp) / 2) {
            // refresh token
            const user = sessionTime["user"];
            props.refreshToken({ user });
        }
    };

    const clearTimeoutFunc = () => {
        if (sessionTimerID) {
            clearTimeout(sessionTimerID);
        }
    };

    const logout = () => {
        clearTimeoutFunc();
        props.logout();
    };
    const setTimeoutFunc = () => {
        if (props.sessionTimeout) {
            sessionTimerID = setTimeout(logout, props.sessionTimeout * 1000 * 60);
        }
    };
    const resetTimeOut = () => {
        clearTimeoutFunc();
        setTimeoutFunc();
    };
    return <Fragment />;
};

const mapStateToProps = state => ({
    sessionTimeout: state.setting.sessionTimeout
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    logout: () => dispatch(authActions.logout()),
    getSetting: () => dispatch(settingsActions.getSetting()),
    refreshToken: payload => dispatch(storageActions.refreshToken(payload))
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Session);
