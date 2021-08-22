import { call, put, takeLatest } from "redux-saga/effects";
import service from "helper/service";
import { LOGOUT_SUCCESS } from "modules/auth";

const GET_PROMETHEUS = "GET_PROMETHEUS";
const GET_PROMETHEUS_SUCCESS = "GET_PROMETHEUS_SUCCESS";
const GET_PROMETHEUS_FAILURE = "GET_PROMETHEUS_FAILURE";

//prometheusTypes
interface PrometheusState {
    prometheus: PrometheusData;
    loading: boolean;
}

export interface PrometheusData {
    data: {
        nodeType: string;
        prometheus: string;
        jobName: string;
    }
}

//prometheusActions
export const actions = {
    getPrometheus: () => ({
        type: GET_PROMETHEUS
    }),
    getPrometheusSuccess: (payload: PrometheusState) => ({
        type: GET_PROMETHEUS_SUCCESS,
        payload
    }),
    getPrometheusFailure: (error: string) => ({
        type: GET_PROMETHEUS_FAILURE,
        error
    }),
};

//prometheusReducer
export function prometheusReducer(
    state: PrometheusState = {
        prometheus: { data : { nodeType: "", prometheus: "", jobName: "" } },
        loading: false
    },
    action
): PrometheusState {
    switch (action.type) {
        case GET_PROMETHEUS:
            return {
                ...state,
                loading: true
            };
        case GET_PROMETHEUS_SUCCESS:
            const { prometheus } = action;
            return {
                ...state,
                prometheus,
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

//prometheusAPI
export const api = {
    getPrometheus: async () => {
        return await service.get(`/prometheus`);
    }
};

//channelsSaga
function* getPrometheusFunc() {
    try {
        const res: PrometheusData = yield call(api.getPrometheus);
        if (res) {
            yield put({ type: GET_PROMETHEUS_SUCCESS, prometheus: res });
        }
    } catch (e) {
        console.log(e);
    }
}

export function* prometheusSaga() {
    yield takeLatest(GET_PROMETHEUS, getPrometheusFunc);
}
