import React, { FC, useEffect, useState, useContext, Fragment } from "react";
import { connect } from "react-redux";
import { Dispatch, Action } from "redux";
import useReactRouter from "use-react-router";
import { Link } from "react-router-dom";
import { ChannelsData } from "modules/channels";
import { actions as channelsActions } from "modules/channels";
import Button from "@material-ui/core/Button";
import {
    MdDashboard,
    MdViewList,
    MdLaptop,
    MdSettings,
    MdViewModule,
    MdArrowDropDown
} from "react-icons/md";
import { FaUsers } from "react-icons/fa";
import { Language } from "components";
import { getUserInfo } from "utils/token";

interface Props {
    getChannels: () => any;
    channels: ChannelsData;
}

// [TO DO] className명 변경
const Sidebar: FC<Props> = props => {
    const LanguageContext = useContext(Language());
    const { I18n } = LanguageContext;
    const userInfo = getUserInfo();
    const { location } = useReactRouter();

    const [dropdown, setDropdown] = useState({
        nodeList: location.pathname.includes("/node_list"),
        monitoringLog: false,
        tracker: location.pathname.includes("/tracker"),
        setting: false
    });

    useEffect(() => {
        props.getChannels();
    }, []);

    const onClickDropdown = target => {
        setDropdown({ ...dropdown, [target]: !dropdown[target] });
    };
    return (
        <div className="sidebar">
            {/* <h1 className="logo">
        <img src={logo} alt="" />
      </h1> */}
            <ul className="sidebar-menu">
                <li>
                    <Link to="/">
                        <Button
                            variant="contained"
                            className={`nav-item ${location.pathname === "/" &&
                                "current"}`}
                        >
                            <MdDashboard />
                            {I18n.dashboard}
                        </Button>
                    </Link>
                </li>

                {userInfo.permission.includes("Node") && (
                    <Fragment>
                        <li>
                            <Button
                                variant="contained"
                                className="nav-item"
                                onClick={() => onClickDropdown("nodeList")}
                            >
                                <MdLaptop />
                                {I18n.nodeList}
                                <MdArrowDropDown
                                    className={`arrow-ic ${dropdown.nodeList &&
                                        "active"}`}
                                />
                            </Button>
                        </li>

                        <div
                            className={`sub-menu ${dropdown.nodeList &&
                                "active"}`}
                        >
                            <ul>
                                {props.channels.data.map((el, key) => {
                                    return (
                                        <li key={key}>
                                            <Link
                                                to={{
                                                    pathname: `/node_list/${el.name}`,
                                                    state: {
                                                        channelId: el.id
                                                    }
                                                }}
                                            >
                                                <Button
                                                    className={`nav-item ${location.pathname.includes(
                                                        "/node_list"
                                                    ) &&
                                                        location.pathname.includes(
                                                            el.name
                                                        ) &&
                                                        "current"}`}
                                                >
                                                    {el.name}
                                                </Button>
                                            </Link>
                                        </li>
                                    );
                                })}
                            </ul>
                        </div>
                    </Fragment>
                )}
                {userInfo.permission.includes("MonitoringLog") && (
                    <li>
                        <Link
                            to={{
                                pathname: "/monitoring_log"
                            }}
                        >
                            <Button
                                variant="contained"
                                className={`nav-item ${location.pathname.includes(
                                    "/monitoring_log"
                                ) && "current"}`}
                            >
                                <MdViewList />
                                {I18n.monitoringLog}
                            </Button>
                        </Link>
                    </li>
                )}
                <li>
                    <Button
                        variant="contained"
                        className="nav-item"
                        onClick={() => onClickDropdown("tracker")}
                    >
                        <MdViewModule />
                        {I18n.blockTx}
                        <MdArrowDropDown
                            className={`arrow-ic ${dropdown.tracker &&
                                "active"}`}
                            style={{
                                float: "right",
                                marginRight: 0,
                                width: "22px"
                            }}
                        />
                    </Button>
                </li>
                <div className={`sub-menu ${dropdown.tracker && "active"}`}>
                    {props.channels.data.map((el, key) => {
                        return (
                            <li key={key}>
                                <Link
                                    to={{
                                        pathname: `/tracker/${el.name}/${el.id}`
                                    }}
                                >
                                    <Button
                                        className={`nav-item ${location.pathname.includes(
                                            "/tracker"
                                        ) &&
                                            location.pathname.includes(
                                                el.name
                                            ) &&
                                            "current"}
                `}
                                    >
                                        {el.name}
                                    </Button>
                                </Link>
                            </li>
                        );
                    })}
                </div>
                {userInfo.admin && (
                    <Fragment>
                        <li>
                            <Link to="/users">
                                <Button
                                    variant="contained"
                                    className={`nav-item ${location.pathname.includes(
                                        "/user"
                                    ) && "current"}`}
                                >
                                    <FaUsers />
                                    {I18n.user}
                                </Button>
                            </Link>
                        </li>
                        <li>
                            <Button
                                variant="contained"
                                className="nav-item"
                                onClick={() => onClickDropdown("setting")}
                            >
                                <MdSettings />
                                {I18n.setting}
                                <MdArrowDropDown
                                    className={`arrow-ic ${dropdown.setting &&
                                        "active"}`}
                                />
                            </Button>
                        </li>
                        <div
                            className={`sub-menu ${dropdown.setting &&
                                "active"}`}
                        >
                            <ul>
                                <li>
                                    <Link
                                        to={{
                                            pathname: `/setting_channel`
                                        }}
                                    >
                                        <Button
                                            className={`nav-item ${location.pathname.includes(
                                                "/setting_channel"
                                            ) && "current"}`}
                                        >
                                            {I18n.settingChannel}
                                        </Button>
                                    </Link>
                                </li>
                                <li>
                                    <Link
                                        to={{
                                            pathname: `/setting_node`
                                        }}
                                    >
                                        <Button
                                            className={`nav-item ${location.pathname.includes(
                                                "/setting_node"
                                            ) && "current"}`}
                                        >
                                            {I18n.settingNode}
                                        </Button>
                                    </Link>
                                </li>
                                <li>
                                    <Link
                                        to={{
                                            pathname: `/setting_etc`
                                        }}
                                    >
                                        <Button
                                            className={`nav-item ${location.pathname.includes(
                                                "/setting_etc"
                                            ) && "current"}`}
                                        >
                                            {I18n.settingETC}
                                        </Button>
                                    </Link>
                                </li>
                            </ul>
                        </div>
                    </Fragment>
                )}
            </ul>
        </div>
    );
};

const mapStateToProps = state => ({
    channels: state.channels.channels,
    loading: state.channels.loading,
    language: state.storage.language
});
const mapDispatchToProps = (dispatch: Dispatch<Action>) => ({
    getChannels: () => dispatch(channelsActions.getChannels())
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(Sidebar);
