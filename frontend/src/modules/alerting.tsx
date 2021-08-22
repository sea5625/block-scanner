import { call, put, takeLatest } from "redux-saga/effects";
import service from "helper/service";
import { LOGOUT_SUCCESS } from "modules/auth";
import { Popup } from "components";

const GET_ALERTING = "GET_ALERTING";
const GET_ALERTING_SUCCESS = "GET_ALERTING_SUCCESS";
const GET_ALERTING_FAILURE = "GET_ALERTING_FAILURE";

const UPDATE_ALERTING = "UPDATE_ALERTING";
const UPDATE_ALERTING_SUCCESS = "UPDATE_ALERTING_SUCCESS";
const UPDATE_ALERTING_FAILURE = "UPDATE_ALERTING_FAILURE";

//alertingParamsTypes
interface AlertingState {
    alerting: AlertingData;
    loading: boolean;
}

export interface AlertingData {
    data: Alerting[];
    total: number;
}

export interface Alerting {
    id: string;
    name: string;
    slowResponseTime: number;
    unsyncBlockToleranceTime: number;
}

//alertingActions
export const actions = {
    getAlerting: () => ({
        type: GET_ALERTING
    }),
    getAlertingSuccess: payload => ({
        type: GET_ALERTING_SUCCESS,
        payload
    }),
    getAlertingFailure: (error: string) => ({
        type: GET_ALERTING_FAILURE,
        error
    }),
    updateAlerting: payload => ({
        type: UPDATE_ALERTING,
        payload
    }),
    updateAlertingSuccess: payload => ({
        type: UPDATE_ALERTING_SUCCESS,
        payload
    }),
    updateAlertingFailure: (error: string) => ({
        type: UPDATE_ALERTING_FAILURE,
        error
    })
};

//alertingReducer
export function alertingReducer(
    state: AlertingState = {
        alerting: {
            data: [],
            total: 0
        },
        loading: false
    },
    action
): AlertingState {
    switch (action.type) {
        case GET_ALERTING:
            return {
                ...state,
                loading: true
            };
        case GET_ALERTING_SUCCESS:
            const { payload } = action;
            return {
                ...state,
                alerting: payload,
                loading: false
            };
        case LOGOUT_SUCCESS:
            return {
                ...state
            };
        default:
            return state;
    }
}

//alertingAPI
export const api = {
    getAlerting: async () => {
        return await service.get(`/alerting`);
    },
    updateAlerting: async payload => {
        const data = {
            data: payload
        };
        return await service.put(`/alerting`, data);
    }
};

//alertingSaga
function* getAlertingFunc() {
    try {
        const res = yield call(api.getAlerting);
        if (res) {
            yield put({ type: GET_ALERTING_SUCCESS, payload: res });
        }
    } catch (e) {}
}

function* updateAlertingFunc(action) {
    const { payload } = action;
    try {
        const res: AlertingData = yield call(api.updateAlerting, payload);
        if (res) {
            yield put({ type: UPDATE_ALERTING_SUCCESS, payload: res });
            Popup.successPopup({
                className: "success",
                message: "AlertMessageSuccessfullyUpdated",
                btnName: "ok"
            });
        }
    } catch (e) {
        Popup.alertPopup({
            className: "alert",
            message: "AlertMessageFailedToUpdated",
            btnName: "ok"
        });
    }
}

export function* alertingSaga() {
    yield takeLatest(GET_ALERTING, getAlertingFunc);
    yield takeLatest(UPDATE_ALERTING, updateAlertingFunc);
}
