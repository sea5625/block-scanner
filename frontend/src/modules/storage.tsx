import { call, put, takeLatest } from "redux-saga/effects";
import service from "helper/service";

const SET_TOKEN = "SET_TOKEN";
const REMOVE_TOKEN = "REMOVE_TOKEN";

const REFRESH_TOKEN = "REFRESH_TOKEN";
const REFRESH_TOKEN_SUCCESS = "REFRESH_TOKEN_SUCCESS";
const REFRESH_FAILED = "REFRESH_FAILED";

const SET_LANGUAGE = "SET_LANGUAGE";

const SET_PROMETHEUS = "GET_PROMETHEUS";
const SET_PROMETHEUS_SUCCESS = "GET_PROMETHEUS_SUCCESS";

interface StorageState {
    token?: string;
    language: "en" | "ko";
    nodeType: string;
    prometheus: string;
    jobName: string;
}

export const actions = {
    setToken: (token: string) => ({
        type: SET_TOKEN,
        token
    }),
    removeToken: () => ({
        type: REMOVE_TOKEN
    }),
    refreshToken: payload => ({
        type: REFRESH_TOKEN,
        payload
    }),
    setLanguage: (language: "en" | "ko") => ({
        type: SET_LANGUAGE,
        language
    })
};

const api = {
    refreshToken: async payload => {
        return await service.put("/auth/token", payload);
    }
};

export function storageReducer(
    state: StorageState = {
        language: "ko",
        token: "",
        nodeType: "",
        prometheus: "",
        jobName: ""
    },
    action
): StorageState {
    switch (action.type) {
        case SET_TOKEN:
            return {
                ...state,
                token: action.token
            };
        case REMOVE_TOKEN:
            return {
                ...state,
                token: ""
            };
        case REFRESH_TOKEN_SUCCESS:
            return {
                ...state,
                token: action.token
            };
        case SET_LANGUAGE:
            const { language } = action;
            return {
                ...state,
                language
            };
        case SET_PROMETHEUS:
            return {
                ...state
            };
        case SET_PROMETHEUS_SUCCESS:
            const { prometheus } = action;
            return {
                ...state,
                nodeType: prometheus.data.nodeType,
                prometheus: prometheus.data.prometheus,
                jobName: prometheus.data.jobName
            };

        default:
            return state;
    }
}

function* refreshTokenFunc(action) {
    try {
        const { payload } = action;
        const res = yield call(api.refreshToken, payload);
        if (res) {
            const { token } = res;
            yield put({ type: REFRESH_TOKEN_SUCCESS, token });
        }
    } catch (e) {
        console.log(e);
    }
}

export function* storageSaga() {
    yield takeLatest(REFRESH_TOKEN, refreshTokenFunc);
}
