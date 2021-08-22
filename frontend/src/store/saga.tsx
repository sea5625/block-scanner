import { all, fork } from "redux-saga/effects";

import { authSaga } from "modules/auth";
import { usersSaga } from "modules/users";
import { channelsSaga } from "modules/channels";
import { nodesSaga } from "modules/nodes";
import { symptomSaga } from "modules/symptom";
import { trackerSaga } from "modules/tracker";
import { prometheusSaga } from "modules/prometheus";
import { settingSaga } from "modules/setting";
import { alertingSaga } from "modules/alerting";
import { storageSaga } from "modules/storage";
import { resourceSaga } from "modules/resource";

export default function* rootSaga() {
    yield all([
        fork(authSaga),
        fork(usersSaga),
        fork(channelsSaga),
        fork(nodesSaga),
        fork(symptomSaga),
        fork(trackerSaga),
        fork(settingSaga),
        fork(alertingSaga),
        fork(storageSaga),
        fork(prometheusSaga),
        fork(resourceSaga)
    ]);
}
