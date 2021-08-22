import React, { Fragment } from "react";
import { Switch } from "react-router-dom";
import { PrivateRoute } from "utils/routes";
import {
    Dashboard,
    NodeList,
    MonitoringLog,
    Tracker,
    BlockList,
    TxList,
    BlockInfo,
    TxInfo,
    UserList,
    SettingChannel,
    SettingNode,
    SettingETC
} from "pages/admin";

const AdminRoutes = () => {
    return (
        <Fragment>
            <PrivateRoute exact path="/" component={Dashboard} />
            <PrivateRoute path="/node_list/:name" component={NodeList} />
            <PrivateRoute path="/monitoring_log" component={MonitoringLog} />

            {/*[TO DO] Tracker Router 따로 만들고 content 모듈화*/}
            <Switch>
                <PrivateRoute
                    exact
                    path="/tracker/:name/:channelId"
                    component={Tracker}
                />
                <PrivateRoute
                    exact
                    path="/tracker/block_list/:name/:channelId"
                    component={BlockList}
                />
                <PrivateRoute
                    exact
                    path="/tracker/tx_list/:name/:channelId"
                    component={TxList}
                />
                <PrivateRoute
                    exact
                    path="/tracker/block_info/:name/:channelId/:id"
                    component={BlockInfo}
                />
                <PrivateRoute
                    exact
                    path="/tracker/tx_info/:name/:channelId/:txHash"
                    component={TxInfo}
                />
            </Switch>
            <PrivateRoute path="/users" component={UserList} />
            <PrivateRoute path="/setting_channel" component={SettingChannel} />
            <PrivateRoute path="/setting_node" component={SettingNode} />
            <PrivateRoute path="/setting_etc" component={SettingETC} />
        </Fragment>
    );
};

export default AdminRoutes;
