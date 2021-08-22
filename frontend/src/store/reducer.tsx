// tslint:disable-next-line
import { History } from "history";
import { combineReducers } from "redux";
// import { connectRouter } from "connected-react-router";
import history from "./history";
import { authReducer as auth } from "modules/auth";
import { usersReducer as users } from "modules/users";
import { channelsReducer as channels } from "modules/channels";
import { symptomReducer as symptom } from "modules/symptom";
import { prometheusReducer as prometheus } from "modules/prometheus";
import { nodesReducer as nodes } from "modules/nodes";
import { trackerReducer as tracker } from "modules/tracker";
import { settingReducer as setting } from "modules/setting";
import { alertingReducer as alerting } from "modules/alerting";
import { storageReducer as storage } from "modules/storage";
import { resourceReducer as resource } from "modules/resource";

const createRootReducer = (_history: History) =>
    combineReducers({
        // router: connectRouter(_history),
        auth,
        users,
        channels,
        symptom,
        nodes,
        tracker,
        setting,
        alerting,
        storage,
        prometheus,
        resource
    });

const rootReducer = createRootReducer(history);

export default rootReducer;
