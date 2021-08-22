import "@babel/polyfill";
import "react-app-polyfill/ie11";
import React, { FC } from "react";
import { Provider } from "react-redux";
import { Router } from "react-router-dom";
import { store, history } from "store";
import { LayerPopupContainer } from "lib/popup";
import { Admin } from "pages/admin";
import { Login } from "pages/public";
import { PublicRoute, PrivateRoute } from "utils/routes";

const App: FC = props => {
    return (
        <Provider store={store}>
            <Router history={history}>
                <PublicRoute exact path="/" component={Login} />
                <PrivateRoute path="/" component={Admin} />
                <LayerPopupContainer {...props} />
            </Router>
        </Provider>
    );
};

export default App;
