import React, { FC, useEffect, useState, useContext } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import useReactRouter from "use-react-router";
import { getToken, parsingToken } from "utils/token";
import { UserData } from "modules/users";
import { actions as authActions } from "modules/auth";
import { actions as usersActions } from "modules/users";
import { actions as storageActions } from "modules/storage";
import { Session, Loading, Language, Popup } from "components";

import {
    MdArrowDropDown,
    MdLanguage,
    MdFullscreen,
    MdFullscreenExit,
    MdMenu
} from "react-icons/md";
import { Button, Menu, MenuItem, IconButton } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";

interface Props {
    getUser: (payload) => any;
    setLanguage: (language) => any;
    language: "ko" | "en";
    loading: boolean;
    user: UserData;
    logout: () => any;
}
const CustomButton = withStyles({
    root: {
        boxShadow: "none",
        textTransform: "none",
        fontWeight: 300,
        marginTop: "0.8rem",
        border: "1px solid transparent",
        backgroundColor: "rgba(0,0,0,0.05)",
        outLine: "none",
        borderRadius: 0,
        padding: "0 1rem",
        "&:active": {
            boxShadow: "none",
            backgroundColor: "none",
            borderColor: "none"
        },
        "&:focus": {
            boxShadow: "none",
            backgroundColor: "none",
            borderColor: "none"
        }
    }
})(Button);

const Header: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const [userMenu, setUserMenu] = useState(null);
    const [languageMenu, setLanguageMenu] = useState(null);
    // const [fullScreen, setFullScreen] = useState(false);

    useEffect(() => {
        const token = getToken();
        const tokenParse = parsingToken(token);
        const { user: id } = tokenParse;
        props.getUser({ id });
    }, []);

    const { location } = useReactRouter();
    let title;
    if (location.pathname === "/") {
        title = I18n.dashboard;
    }
    if (location.pathname.includes("node_list")) {
        title = I18n.nodeList;
    }
    if (location.pathname.includes("tracker")) {
        title = I18n.tracker;
    }
    if (location.pathname.includes("monitoring_log")) {
        title = I18n.monitoringLog;
    }
    if (location.pathname.includes("setting_channel")) {
        title = I18n.settingChannel;
    }
    if (location.pathname.includes("users")) {
        title = I18n.users;
    }
    if (location.pathname.includes("setting_node")) {
        title = I18n.settingNode;
    }
    if (location.pathname.includes("setting_etc")) {
        title = I18n.settingETC;
    }

    const onClickUserMenu = e => {
        setUserMenu(e.currentTarget);
    };
    const onClickUserMenuItem = e => {
        setUserMenu(null);
        const value = e.currentTarget.getAttribute("value");
        if (value === "profile") {
        } else if (value === "sessionTimeout") {
            Popup.sessionTimeoutPopup({
                className: "session-timeout"
            });
        } else if (value === "logout") {
            props.logout();
        } else {
            return;
        }
    };

    const onClickLanguageMenu = e => {
        setLanguageMenu(e.currentTarget);
    };

    const onClickLanguageItem = e => {
        setLanguageMenu(null);
        const language = e.currentTarget.getAttribute("value");
        if (language === "en" || language === "ko") {
            props.setLanguage(language);
        }
    };
    // const onClickFullScreen = () => {
    //   const elem = document.documentElement;
    //   if (elem.requestFullscreen) {
    //     elem.requestFullscreen();
    //   }
    //   setFullScreen(true);
    // };
    // const onClickFullScreenExit = () => {
    //   if (document.exitFullscreen) {
    //     document.exitFullscreen();
    //   }
    //   setFullScreen(false);
    // };
    const onClickEdit = (id, source, type, selectUser, className) => {
        Popup.edituserPopup(id, source, type, selectUser, className);
    };

    if (props.loading) {
        return <Loading />;
    } else {
        return (
            <header className="header clearfix">
                <Session />
                <IconButton
                    className="menu-btn"
                    // onClick={fullScreen ? onClickFullScreenExit : onClickFullScreen}
                >
                    <MdMenu />
                </IconButton>
                <h2 className="title">{title}</h2>
                <div className="header-right">
                    {/* <IconButton
            className="screen-btn"
            onClick={fullScreen ? onClickFullScreenExit : onClickFullScreen}
          >
            {fullScreen ? <MdFullscreenExit /> : <MdFullscreen />}
          </IconButton> */}
                    {/* <IconButton className="language-btn">
            <MdLanguage />
          </IconButton> */}
                    <IconButton
                        className="language-btn"
                        aria-controls="language-menu"
                        aria-haspopup="true"
                        onClick={onClickLanguageMenu}
                    >
                        <MdLanguage />
                    </IconButton>
                    <Menu
                        id="language-menu"
                        anchorEl={languageMenu}
                        keepMounted
                        open={Boolean(languageMenu)}
                        onClose={onClickLanguageItem}
                    >
                        <MenuItem
                            onClick={onClickLanguageItem}
                            value="ko"
                            style={{
                                fontWeight: 400,
                                color:
                                    props.language === "ko" &&
                                    " rgb(94, 114, 228)"
                            }}
                        >
                            한국어
                        </MenuItem>
                        <MenuItem
                            onClick={onClickLanguageItem}
                            value="en"
                            style={{
                                color:
                                    props.language === "en" &&
                                    " rgb(94, 114, 228)"
                            }}
                        >
                            English
                        </MenuItem>
                    </Menu>
                    <CustomButton
                        aria-controls="user-menu"
                        aria-haspopup="true"
                        className="user-setting-btn"
                        onClick={onClickUserMenu}
                    >
                        {props.user.userId} <MdArrowDropDown />
                    </CustomButton>
                    <Menu
                        id="user-menu"
                        anchorEl={userMenu}
                        keepMounted
                        open={Boolean(userMenu)}
                        onClose={onClickUserMenuItem}
                    >
                        <MenuItem
                            onClick={() =>
                                onClickEdit(
                                    props.user.id,
                                    props.user.userId,
                                    "edit",
                                    props.user,
                                    "user-edit"
                                )
                            }
                            style={{
                                minHeight: "36px",
                                fontSize: 14,
                                fontWeight: 300,
                                padding: "4px 20px"
                            }}
                        >
                            {I18n.profile}
                        </MenuItem>
                        <MenuItem
                            onClick={onClickUserMenuItem}
                            value="sessionTimeout"
                        >
                            {I18n.setSessionTime}
                        </MenuItem>
                        <MenuItem onClick={onClickUserMenuItem} value="logout">
                            {I18n.logout}
                        </MenuItem>
                    </Menu>
                </div>
            </header>
        );
    }
};

const mapStateToProps = state => ({
    user: state.users.user,
    loading: state.users.userLoading,
    language: state.storage.language
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getUser: payload => dispatch(usersActions.getUser(payload)),
    putUser: payload => dispatch(usersActions.putUser(payload)),
    setLanguage: language => dispatch(storageActions.setLanguage(language)),
    logout: () => dispatch(authActions.logout())
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Header);
